package main

import (
	"os"
	"sync"

	"github.com/zsrv/rt5-server-go/engine"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		s := engine.NewServer()

		s.Addr = "127.0.0.1:40001"

		s.Logger.Info("starting server", "listenAddr", s.Addr)
		err := s.ListenAndServe()
		if err != nil {
			s.Logger.Error("error", err)
			os.Exit(1)
		} else {
			s.Logger.Info("server exiting")
		}
	}()

	//wg.Add(1)
	//go func() {
	//	defer wg.Done()
	//	s := network.NewServer()
	//
	//	s.Addr = "127.0.0.1:50001"
	//
	//	log.Println("Starting server at", s.Addr)
	//	log.Fatal(s.ListenAndServe())
	//}()

	wg.Wait()

	//s := network.NewServer()
	//
	//s.Addr = "127.0.0.1:40001"
	//
	//log.Println("Starting server at", s.Addr)
	//log.Fatal(s.ListenAndServe())
}
