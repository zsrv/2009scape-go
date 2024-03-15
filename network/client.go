package network

import (
	"math"
	"math/rand"
	"net"
	"strconv"

	"github.com/zsrv/rt5-server-go/util"
	"github.com/zsrv/rt5-server-go/util/packet"
)

const (
	ClientStateClosed = -1
	ClientStateNew    = 0
	ClientStateJS5    = 1
	ClientStateWL     = 2
	ClientStateLogin  = 3
	ClientStateGame   = 4
)

type Client struct {
	Server *Server
	Socket net.Conn
	State  int

	NetOut []NetOutData

	RandomIn  *util.IsaacRandom
	RandomOut *util.IsaacRandom

	Player          *Player
	BufferStart     int
	BufferInOffset  int
	BufferOutOffset int

	PacketCount []uint8

	BufferInRaw packet.Packet
}

func NewClient(socket net.Conn, server *Server) *Client {
	// TODO
	return &Client{
		Server: server,
		Socket: socket,
		State:  ClientStateNew,

		PacketCount: make([]uint8, 256),
	}
}

func (c *Client) handleData() {
	c.Server.Logger.Debug("entered handleData()")
	switch c.State {
	case ClientStateClosed:
		c.Server.Logger.Debug("closing connection")
		c.Socket.Close()
	case ClientStateNew:
		c.handleNew()
	case ClientStateJS5:
		c.handleJS5()
	case ClientStateWL:
		c.handleWL()
	case ClientStateLogin:
		c.handleLogin()
	case ClientStateGame:
		c.handleGame()
	default:
		c.Socket.Close()
		c.State = ClientStateClosed
		c.Server.Logger.Error("invalid client state", "state", strconv.Itoa(c.State))
	}
}

func (c *Client) handleNew() {
	c.Server.Logger.Debug("entered handleNew()")
	opcode := c.BufferInRaw.G1()

	switch opcode {
	case util.LoginProtJS5Open:
		c.Server.Logger.Debug("handleNew(): case LoginProtJS5Open")
		clientVersion := c.BufferInRaw.G4()

		if clientVersion == 578 {
			c.Server.Logger.Debug("client version is 578")
			c.WriteRawSocket([]byte{util.JS5ProtOutSuccess})
			c.State = ClientStateJS5
		} else {
			c.Server.Logger.Debug("client version is not 578", "clientVersion", clientVersion)
			c.WriteRawSocket([]byte{util.JS5ProtOutOutOfDate})
			c.Socket.Close()
			c.State = ClientStateClosed
		}

	case util.LoginProtWorldListFetch:
		c.Server.Logger.Debug("handleNew(): case LoginProtWorldListFetch")
		//checksum, err := c.BufferInRaw.G4B() // TODO: G4B make any diff here or no?
		checksum := c.BufferInRaw.G4()

		c.WriteRawSocket([]byte{WlProtOutSuccess})
		c.State = ClientStateWL

		var response packet.Packet

		response.P2(0)
		start := response.Len() // offset

		response.P1(1) // encoding a world list update

		if checksum != WorldListChecksum {
			response.P1(1) // encoding all information about the world list (countries, size of list, etc.)

			wlraw := WorldListRaw.Bytes()
			response.PData(wlraw, len(wlraw))

			response.P4(WorldListChecksum)
		} else {
			response.P1(0) // not encoding any world list information, just updating the player counts
		}

		for _, world := range WorldList {
			response.PSmart(uint16(world.ID - MinID))
			response.P2(uint16(world.Players))
		}

		response.PSize2(response.Len() - start)
		c.WriteRawSocket(response.Bytes())

	case util.LoginProtWorldHandshake: // login
		c.Server.Logger.Debug("handleNew(): case LoginProtWorldHandshake")
		var response packet.Packet
		response.P1(0)
		response.P8(uint64(math.Floor(rand.Float64()*0xFFFF_FFFF))<<32 | uint64(math.Floor(rand.Float64()*0xFFFF_FFFF)))

		c.WriteRawSocket(response.Bytes())
		c.State = ClientStateLogin

	case util.LoginProtCreateLogProgress:
		c.Server.Logger.Debug("handleNew(): case LoginProtCreateLogProgress")
		day := c.BufferInRaw.G1()
		month := c.BufferInRaw.G1()
		year := c.BufferInRaw.G2()
		country := c.BufferInRaw.G2()

		c.Server.Logger.Debug("progress", "year", year, "month", month, "day", day, "country", country)

		c.WriteRawSocket([]byte{2})

	case util.LoginProtCreateCheckName:
		c.Server.Logger.Debug("handleNew(): case LoginProtCreateCheckName")
		usernameBase37 := c.BufferInRaw.G8()
		username := util.FromBase37(usernameBase37)
		c.Server.Logger.Debug("decoded username", "username", username)

		// success:
		c.WriteRawSocket([]byte{2})

		// suggested names:
		var response packet.Packet
		response.P1(21)

		names := []string{"test", "test2"}
		response.P1(uint8(len(names)))
		for _, v := range names {
			response.P8(util.ToBase37(v))
		}

		c.WriteRawSocket(response.Bytes())

	case util.LoginProtCreateAccount:
		c.Server.Logger.Debug("handleNew(): case LoginProtCreateAccount")
		length := c.BufferInRaw.G2()

		newBuf := make([]byte, length)
		c.BufferInRaw.GData(newBuf, int(length))
		c.BufferInRaw = *packet.NewPacket(newBuf)

		revision := c.BufferInRaw.G2()

		decrypted, err := c.BufferInRaw.RSADec()
		if err != nil {
			// TODO: make an error wrapper that logs the error and closes the client connection?
			// check my notes for this
			c.Server.Logger.Error("rsa decryption error", "error", err)
		}

		rsaMagic := decrypted.G1()
		if rsaMagic != 10 {
			// TODO: read failure
			c.Server.Logger.Error("rsaMagic read failure", "rsaMagic", rsaMagic)
		}

		key := make([]uint32, 4)
		optIn := decrypted.G2()

		usernameBase37 := decrypted.G8()
		username := util.FromBase37(usernameBase37)

		key1 := decrypted.G4()
		key[0] = key1

		password := decrypted.GJStr()

		key2 := decrypted.G4()
		key[1] = key2

		affiliate := decrypted.G2()

		day := decrypted.G1()

		month := decrypted.G1()

		key3 := decrypted.G4()
		key[2] = key3

		year := decrypted.G2()

		country := decrypted.G2()

		key4 := decrypted.G4()
		key[3] = key4

		extra := packet.NewPacket(c.BufferInRaw.Bytes())
		extra.TinyDec(0, key, extra.Len())

		email := extra.GJStr()

		c.Server.Logger.Debug("account creation values",
			"revision", revision, "optIn", optIn, "username", username, "password", password,
			"affiliate", affiliate, "day", day, "month", month, "year", year, "country", country,
			"email", email)

		c.WriteRawSocket([]byte{2})

		// 0 - unexpected response
		// 1 - could not display video ad
		// 2 - success
		// 3 - invalid username/password
		// 4 - account is banned
		// 5 - account is logged in
		// 6 - client out of date
		// 7 - world is full
		// 8 - login server offline
		// 9 - too many connections
		// 10 - bad session id
		// 11 - weak password
		// 12 - f2p account, p2p world
		// 13 - could not login
		// 14 - server is updating
		// 15 - reconnecting
		// 16 - too many login attempts
		// 17 - p2p area, f2p world
		// 18 - account locked
		// 19 - members beta
		// 20 - invalid login server
		// 21 - moving worlds
		// 22 - malformed login packet
		// 23 - no reply from login server
		// 24 - error loading profile
		// 26 - mac address banned
		// 27 - service unavailable

	default:
		c.Server.Logger.Warn("unknown opcode", "opcode", opcode)
		c.State = ClientStateClosed
	}
}

func (c *Client) handleJS5() {
	c.Server.Logger.Debug("entered handleJS5()")
	type QueueData struct {
		Type    uint8
		Archive uint8
		Group   uint16
	}

	var queue []QueueData

	for c.BufferInRaw.Len() != 0 {
		xType := c.BufferInRaw.G1()

		switch xType {
		case util.JS5ProtInRequest, util.JS5ProtInPriorityRequest:
			archive := c.BufferInRaw.G1()
			group := c.BufferInRaw.G2()

			queue = append(queue, QueueData{
				Type:    xType,
				Archive: archive,
				Group:   group,
			})
		default:
			c.BufferInRaw.Next(3)
		}
	}

	// TODO: move this out of the network handler and into a dedicated Js5 queue loop (for all requests)
	// TODO: async?
	for _, v := range queue {
		file, err := util.GetGroup(v.Archive, v.Group)
		if err != nil {
			c.Server.Logger.Error("error getting group", "error", err)
			// TODO: close conn etc
			return
		}

		if v.Archive == 255 && v.Group == 255 {
			// checksum table for all archives
			var response packet.Packet

			response.P1(v.Archive)
			response.P2(v.Group)

			response.Write(file)
			c.WriteRawSocket(response.Bytes())
		} else {
			compression := file[0]
			var length uint32 = uint32(file[1])<<24 | uint32(file[2])<<16 | uint32(file[3])<<8 | uint32(file[4])
			realLength := 0
			if compression != 0 {
				realLength = int(length + 4)
			} else {
				realLength = int(length)
			}

			settings := compression
			if v.Type == util.JS5ProtInRequest {
				settings |= 0x80
			}

			var response packet.Packet
			response.P1(v.Archive)
			response.P2(v.Group)
			response.P1(settings)
			response.P4(length)

			for i := 5; i < realLength+5; i++ {
				if response.Len()%512 == 0 { // TODO: might not be correct equiv of offset
					response.P1(0xFF)
				}
				response.P1(file[i])
			}

			c.WriteRawSocket(response.Bytes())
		}
	}
}

func (c *Client) handleWL() {
	// no communication
	c.Server.Logger.Debug("entered handleWL()")
}

func (c *Client) handleLogin() {
	c.Server.Logger.Debug("entered handleLogin()")
	opcode := c.BufferInRaw.G1()

	length := c.BufferInRaw.G2()

	data := make([]byte, length)
	c.BufferInRaw.GData(data, int(length))
	c.BufferInRaw = *packet.NewPacket(data)

	revision := c.BufferInRaw.G4()

	byte1 := c.BufferInRaw.G1()

	windowMode := c.BufferInRaw.G1()

	canvasWidth := c.BufferInRaw.G2()

	canvasHeight := c.BufferInRaw.G2()

	prefInt := c.BufferInRaw.G1()

	uid := make([]byte, 24)
	c.BufferInRaw.GData(uid, 24)

	settings := c.BufferInRaw.GJStr()

	affiliate := c.BufferInRaw.G4()

	preferencesLen := c.BufferInRaw.G1()

	preferences := make([]byte, preferencesLen)
	c.BufferInRaw.GData(preferences, int(preferencesLen))

	verifyId := c.BufferInRaw.G2()

	checksums := make([]uint32, 29)
	for i := 0; i < 29; i++ {
		checksums[i] = c.BufferInRaw.G4()
	}

	decrypted, err := c.BufferInRaw.RSADec()
	if err != nil {
		c.Server.Logger.Error("error decrypting buffer", "error", err)
		return // TODO: close connection etc
	}
	rsaMagic := decrypted.G1()
	key := make([]uint32, 4)
	for i := 0; i < 4; i++ {
		key[i] = decrypted.G4()
	}

	username1 := decrypted.G8()
	username := util.FromBase37(username1)

	password := decrypted.GJStr()

	c.Server.Logger.Debug("login", "opcode", opcode, "revision", revision, "byte1", byte1,
		"windowMode", windowMode, "canvasWidth", canvasWidth, "canvasHeight", canvasHeight,
		"prefInt", prefInt, "uid", uid, "settings", settings, "affiliate", affiliate,
		"preferences", preferences, "verifyId", verifyId, "rsaMagic", rsaMagic,
		"username", username, "password", password,
	)

	c.RandomIn = util.NewIsaacRandom(key)
	for i := 0; i < 4; i++ {
		key[i] += 50
	}
	c.RandomOut = util.NewIsaacRandom(key)

	player := NewPlayer(c)
	if opcode == util.LoginProtWorldReconnect {
		player.Reconnecting = true
	}
	player.WindowMode = windowMode
	player.Username = util.ToTitleCase(username)
	c.Player = player

	c.Server.World.RegisterPlayer(player)
	c.BufferStart = c.Player.ID * 30000

	var response packet.Packet
	if opcode == util.LoginProtWorldReconnect {
		response.P1(15)
	} else {
		response.P1(2)
	}

	if opcode == util.LoginProtWorldConnect {
		response.P1(0)                 // staff mod level
		response.P1(0)                 // player mod level
		response.P1(0)                 // player underage
		response.P1(0)                 // parentalChatConsent
		response.P1(0)                 // parentalAdvertConsent
		response.P1(0)                 // mapQuickChat
		response.P2(uint16(player.ID)) // selfId
		response.P1(0)                 // MouseRecorder
		response.P1(1)                 // mapMembers
	}

	c.WriteRawSocket(response.Bytes())

	c.State = ClientStateGame
	c.Server.World.AddPlayer(player)

	c.Server.Logger.Debug("login complete")
}

func (c *Client) handleGame() {
	c.Server.Logger.Debug("entered handleGame()")
	// TODO: is there a Packet func that does the stuff being done to data here?
	// so then we don't have to extract the bytes from the buffer/Packet
	data := c.BufferInRaw.Bytes()

	offset := 0
	for offset < len(data) {
		start := offset

		if c.RandomIn != nil {
			data[offset] -= byte(c.RandomIn.GetNext())
		}

		opcode := data[offset]
		offset++

		length := util.ClientProtLengths[opcode]

		if length == 255 {
			length = data[offset]
			offset += 1
		} else if length == 254 {
			// TODO: adjusted this
			length = uint8(uint16(data[offset])<<8 | uint16(data[offset+1]))
			offset += 2
		}

		if int(length) > 30000-c.BufferInOffset {
			c.Server.Logger.Error("packet overflow for this tick")
			return // TODO: close conn
		}

		if c.PacketCount[opcode]+1 > 10 {
			offset += int(length)
			continue
		}

		c.PacketCount[opcode] += 1

		slicex := data[start : offset+int(length)]
		offset += int(length)

		copy(c.Server.BufferIn[c.BufferStart+c.BufferInOffset:], slicex)
		c.BufferInOffset += len(slicex)
	}
}

func (c *Client) ResetIn() {
	c.BufferInOffset = 0
	for i := 0; i < len(c.PacketCount); i++ {
		c.PacketCount[i] = 0
	}
}

type DecodedData struct {
	ID   uint8
	Data packet.Packet
}

func (c *Client) DecodeIn() []DecodedData {
	c.Server.Logger.Debug("entered DecodeIn()")
	offset := 0

	var decoded []DecodedData

	for offset < c.BufferInOffset {
		opcode := c.Server.BufferIn[c.BufferStart+offset]
		offset += 1
		length := util.ClientProtLengths[opcode]
		if length == 255 {
			length = c.Server.BufferIn[c.BufferStart+offset]
			offset += 1
		} else if length == 254 {
			length = c.Server.BufferIn[c.BufferStart+offset]<<8 | c.Server.BufferIn[c.BufferStart+offset+1]
			offset += 2
		}

		decoded = append(decoded, DecodedData{
			ID:   opcode,
			Data: *packet.NewPacket(c.Server.BufferIn[c.BufferStart+offset : c.BufferStart+offset+int(length)]),
		})

		offset += int(length)
	}

	return decoded
}

func (c *Client) Write(data []byte) {
	//util.DebugfBytes(&c.Server.Logger, "Write()", data)
	offset := 0
	remaining := len(data)

	// pack as much data as we can into a single chunk, then flush and repeat
	for remaining > 0 {
		untilNextFlush := 30000 - c.BufferOutOffset

		if remaining > untilNextFlush {
			toWrite := data[offset : offset+untilNextFlush]
			for i, v := range toWrite {
				c.Server.BufferOut[c.BufferStart+c.BufferOutOffset+i] = v
			}

			c.BufferOutOffset += untilNextFlush
			c.Flush()
			offset += untilNextFlush
			remaining -= untilNextFlush
		} else {
			toWrite := data[offset : offset+remaining]
			for i, v := range toWrite {
				c.Server.BufferOut[c.BufferStart+c.BufferOutOffset+i] = v
			}

			c.BufferOutOffset += remaining
			offset += remaining
			remaining = 0
		}
	}
}

func (c *Client) Flush() {
	if c.BufferOutOffset > 0 {
		c.WriteRawSocket(c.Server.BufferOut[c.BufferStart : c.BufferStart+c.BufferOutOffset])
		c.BufferOutOffset = 0
	}
}

type NetOutData struct {
	Data    []byte
	Encrypt bool
}

func (c *Client) Queue(data []byte, encrypt bool) {
	c.NetOut = append(c.NetOut, NetOutData{
		Data:    data,
		Encrypt: encrypt,
	})
}

func (c *Client) EncodeOut() {
	for i := 0; i < len(c.NetOut); i++ {
		xPacket := c.NetOut[i]

		if c.RandomOut != nil && xPacket.Encrypt {
			// TODO: uint32 converted to byte!!
			xPacket.Data[0] += byte(c.RandomOut.GetNext())
		}

		c.Write(xPacket.Data)
	}
}

func (c *Client) WriteRawSocket(data []byte) {
	_, err := c.Socket.Write(data)
	if err != nil {
		c.Server.Logger.Error("error writing to connection", "error", err)
	}
}
