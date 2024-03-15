package engine

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net"
	"sync"
	"time"

	"github.com/zsrv/rt5-server-go/util"
)

var (
	ErrServerClosed = errors.New("server: server already closed")
)

type Server struct {
	Addr string

	Logger       slog.Logger
	ReadTimeout  time.Duration
	WriteTimeout time.Duration

	wg sync.WaitGroup

	done chan struct{}

	locker    sync.Mutex
	listeners []net.Listener
	Clients   map[*Client]struct{}

	World *World

	BufferIn  []uint8
	BufferOut []uint8
}

func NewServer() *Server {
	return &Server{
		// TODO: init buffers here?
		Logger:  *util.NewLogger(),
		Clients: make(map[*Client]struct{}),

		World: NewWorld(),

		BufferIn:  make([]uint8, 2048*30000), // pre-allocate 61MB for incoming packets, reduces GC pressure
		BufferOut: make([]uint8, 2048*30000), // pre-allocate 61MB for outgoing packets, reduces GC pressure
	}
}

// ListenAndServe listens on the network address s.Addr and then calls Serve
// to handle requests on incoming connections.
func (s *Server) ListenAndServe() error {
	l, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return err
	}

	return s.Serve(l)
}

// Serve accepts incoming connections on the Listener l.
func (s *Server) Serve(l net.Listener) error {
	s.locker.Lock()
	s.listeners = append(s.listeners, l)
	s.locker.Unlock()

	for {
		socket, err := l.Accept()
		if err != nil {
			select {
			case <-s.done:
				// we called Close()
				s.Logger.Info("connection closed by server", "remoteAddr", socket.RemoteAddr())
				return nil
			default:
				return err
			}
		}

		s.wg.Add(1)
		go func() {
			defer s.wg.Done()

			err := s.handleConn(NewClient(socket, s))
			if err != nil {
				s.Logger.Error("handler error", "error", err)
			}
		}()
	}
}

// Close immediately closes all active listeners and connections.
//
// Close returns any error returned from closing the server's underlying
// listener(s).
func (s *Server) Close() error {
	select {
	case <-s.done:
		return ErrServerClosed
	default:
		close(s.done)
	}

	var err error
	s.locker.Lock()
	for _, l := range s.listeners {
		if lerr := l.Close(); lerr != nil && err == nil {
			err = lerr
		}
	}

	for conn := range s.Clients {
		conn.Socket.Close()
	}
	s.locker.Unlock()

	return err
}

// Shutdown gracefully shuts down the server without interrupting any
// active connections. Shutdown works by first closing all open
// listeners and then waiting indefinitely for connections to return to
// idle and then shut down.
// If the provided context expires before the shutdown is complete,
// Shutdown returns the context's error, otherwise it returns any
// error returned from closing the Server's underlying Listener(s).
func (s *Server) Shutdown(ctx context.Context) error {
	select {
	case <-s.done:
		return ErrServerClosed
	default:
		close(s.done)
	}

	var err error
	s.locker.Lock()
	for _, l := range s.listeners {
		if lerr := l.Close(); lerr != nil && err == nil {
			err = lerr
		}
	}
	s.locker.Unlock()

	connDone := make(chan struct{})
	go func() {
		defer close(connDone)
		s.wg.Wait()
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-connDone:
		return err
	}
}

func (s *Server) handleConn(c *Client) error {
	s.locker.Lock()
	s.Clients[c] = struct{}{}
	s.locker.Unlock()

	defer func() {
		c.Socket.Close()
		s.Logger.Info("connection closed", "remoteAddr", c.Socket.RemoteAddr())

		if c.Player != nil {
			s.World.RemovePlayer(*c)
		}

		s.locker.Lock()
		delete(s.Clients, c)
		s.locker.Unlock()
	}()

	if d := s.ReadTimeout; d != 0 {
		c.Socket.SetReadDeadline(time.Now().Add(d))
	}
	if d := s.WriteTimeout; d != 0 {
		c.Socket.SetWriteDeadline(time.Now().Add(d))
	}

	buf := make([]byte, 65536)

	for {
		s.Logger.Debug("waiting for new data")
		n, err := c.Socket.Read(buf)
		if err == io.EOF {
			// Connection closed
			return nil
		}
		if err != nil {
			s.Logger.Error("conn read error", "error", err)
			return err
		}

		msg := buf[:n]
		util.DebugfBytes(&s.Logger, "conn read", msg)

		c.BufferInRaw.Reset()
		_, err = c.BufferInRaw.Write(msg)
		if err != nil {
			s.Logger.Error("conn write error", "error", err)
		}

		c.handleData()
	}
}
