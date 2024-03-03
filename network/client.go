package network

import (
	"fmt"
	"math"
	"math/rand"
	"net"
	"path"
	"runtime"
	"strconv"

	"github.com/zsrv/rt5-server-go/util"
	"github.com/zsrv/rt5-server-go/util/bytebuffer"
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

	BufferInRaw bytebuffer.ByteBuffer
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
	fmt.Println("entered handleData()")
	// switch on the client State here, then pass to the more specific handlers?
	switch c.State {
	case ClientStateClosed:
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
		fmt.Printf("Invalid client state: %v\n", strconv.Itoa(c.State))
	}
}

func (c *Client) handleNew() {
	fmt.Println("entered handleNew()")
	opcode, err := c.BufferInRaw.G1()
	if err != nil {
		fmt.Println(err)
		return
	}

	switch opcode {
	case util.LoginProtJS5Open:
		fmt.Println("IN LoginProtJS5Open")
		clientVersion, err := c.BufferInRaw.G4()
		if err != nil {
			fmt.Println(err)
		}

		if clientVersion == 578 {
			fmt.Println("client version is 578")
			c.WriteRawSocket([]byte{util.JS5ProtOutSuccess})
			c.State = ClientStateJS5
		} else {
			fmt.Printf("not 578! received clientVersion %v\n", clientVersion)
			c.WriteRawSocket([]byte{util.JS5ProtOutOutOfDate})
			c.Socket.Close()
			c.State = ClientStateClosed
		}

	case util.LoginProtWorldListFetch:
		fmt.Println("IN LoginProtWorldListFetch")
		checksum, err := c.BufferInRaw.G4B()
		if err != nil {
			fmt.Println(err)
		}

		c.WriteRawSocket([]byte{util.WlProtOutSuccess})
		c.State = ClientStateWL

		var response bytebuffer.ByteBuffer

		response.P2(0)
		start := response.Len() // offset

		response.PBool(true) // encoding a world list update

		if uint32(checksum) != util.WorldListChecksum {
			response.PBool(true) // encoding all information about the world list (countries, size of list, etc.)
			response.PData(&util.WorldListRaw)
			response.P4(util.WorldListChecksum)
		} else {
			response.PBool(false) // not encoding any world list information, just updating the player counts
		}

		for _, world := range util.WorldList {
			response.PSmart(uint16(world.ID - util.MinID))
			response.P2(uint16(world.Players))
		}

		response.PSize2(response.Len() - start)
		c.WriteRawSocket(response.Bytes())

	case util.LoginProtWorldHandshake: // login
		fmt.Println("IN LoginProtWorldHandshake")
		var response bytebuffer.ByteBuffer
		response.P1(0)
		response.P8(uint64(math.Floor(rand.Float64()*0xFFFF_FFFF))<<32 | uint64(math.Floor(rand.Float64()*0xFFFF_FFFF)))

		c.WriteRawSocket(response.Bytes())
		c.State = ClientStateLogin

	case util.LoginProtCreateLogProgress:
		fmt.Println("IN LoginProtCreateLogProgress")
		day, err := c.BufferInRaw.G1()
		if err != nil {
			fmt.Println(err)
		}

		month, err := c.BufferInRaw.G1()
		if err != nil {
			fmt.Println(err)
		}

		year, err := c.BufferInRaw.G2()
		if err != nil {
			fmt.Println(err)
		}

		country, err := c.BufferInRaw.G2()
		if err != nil {
			fmt.Println(err)
		}

		fmt.Printf("LOG PROGRESS: %v-%v-%v %v\n", year, month, day, country)

		c.WriteRawSocket([]byte{2})

	case util.LoginProtCreateCheckName:
		fmt.Println("IN LoginProtCreateCheckName")
		usernameBase37, err := c.BufferInRaw.G8()
		if err != nil {
			fmt.Println(err)
		}
		username := util.FromBase37(usernameBase37)
		fmt.Printf("CREATE CHECK NAME: %v\n", username)

		// success:
		c.WriteRawSocket([]byte{2})

		// suggested names:
		var response bytebuffer.ByteBuffer
		response.P1(21)

		names := []string{"test", "test2"}
		response.P1(uint8(len(names)))
		for _, v := range names {
			response.P8(util.ToBase37(v))
		}

		c.WriteRawSocket(response.Bytes())

	case util.LoginProtCreateAccount:
		fmt.Println("IN LoginProtCreateAccount")
		length, err := c.BufferInRaw.G2()
		if err != nil {
			fmt.Println(err)
		}

		newBuf, err := c.BufferInRaw.GData(int(length))
		if err != nil {
			fmt.Println(err)
		}
		c.BufferInRaw = *bytebuffer.NewBuffer(newBuf)

		revision, err := c.BufferInRaw.G2()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("Revision: %v\n", revision)

		decrypted, err := c.BufferInRaw.RSADec()
		if err != nil {
			fmt.Println(err)
		}

		rsaMagic, err := decrypted.G1()
		if err != nil {
			fmt.Println(err)
		}
		if rsaMagic != 10 {
			// TODO: read failure
			fmt.Println("rsaMagic read failure!", rsaMagic)
		}

		key := make([]uint32, 4)
		optIn, err := decrypted.G2()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("OptIn: %v\n", optIn)

		usernameBase37, err := decrypted.G8()
		if err != nil {
			fmt.Println(err)
		}
		username := util.FromBase37(usernameBase37)
		fmt.Printf("Username: %v\n", username)

		key1, err := decrypted.G4()
		if err != nil {
			fmt.Println(err)
		}
		key[0] = key1

		password := decrypted.GJStr()
		fmt.Printf("Password: %v\n", password)

		key2, err := decrypted.G4()
		if err != nil {
			fmt.Println(err)
		}
		key[1] = key2

		affiliate, err := decrypted.G2()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("Affiliate: %v\n", affiliate)

		day, err := decrypted.G1()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("Day: %v\n", day)

		month, err := decrypted.G1()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("Month: %v\n", month)

		key3, err := decrypted.G4()
		if err != nil {
			fmt.Println(err)
		}
		key[2] = key3

		year, err := decrypted.G2()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("Year: %v\n", year)

		country, err := decrypted.G2()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("Country: %v\n", country)

		key4, err := decrypted.G4()
		if err != nil {
			fmt.Println(err)
		}
		key[3] = key4

		extra := bytebuffer.NewBuffer(c.BufferInRaw.Bytes())
		err = extra.TinyDec(key, extra.Len(), 0)
		if err != nil {
			fmt.Println(err)
		}

		email := extra.GJStr()
		fmt.Printf("Email: %v\n", email)

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
		fmt.Printf("unknown opcode %v\n", opcode)
		c.State = ClientStateClosed
	}
}

func (c *Client) handleJS5() {
	fmt.Println("entered handleJS5()")
	type QueueData struct {
		Type    uint8
		Archive uint8
		Group   uint16
	}

	var queue []QueueData

	for c.BufferInRaw.Len() != 0 {
		xType, err := c.BufferInRaw.G1()
		if err != nil {
			fmt.Println(err)
			return
		}

		switch xType {
		case util.JS5ProtInRequest, util.JS5ProtInPriorityRequest:
			archive, err := c.BufferInRaw.G1()
			if err != nil {
				fmt.Println(err)
				return
			}
			group, err := c.BufferInRaw.G2()
			if err != nil {
				fmt.Println(err)
				return
			}

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
	fmt.Printf("%+v\n", queue) // DEBUG
	for _, v := range queue {
		file, err := util.GetGroup(v.Archive, v.Group)
		if err != nil {
			fmt.Println(err)
			return
		}

		if v.Archive == 255 && v.Group == 255 {
			// checksum table for all archives
			//var response bytes.Buffer
			var response bytebuffer.ByteBuffer

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

			var response bytebuffer.ByteBuffer
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
	fmt.Println("entered handleWL()")
}

func (c *Client) handleLogin() {
	fmt.Println("entered handleLogin()")
	opcode, err := c.BufferInRaw.G1()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Opcode: %v\n", opcode)

	length, err := c.BufferInRaw.G2()
	if err != nil {
		fmt.Println(err)
	}

	data, err := c.BufferInRaw.GData(int(length))
	if err != nil {
		fmt.Println(err)
	}
	c.BufferInRaw = *bytebuffer.NewBuffer(data)

	revision, err := c.BufferInRaw.G4()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Revision: %v\n", revision)

	byte1, err := c.BufferInRaw.G1()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Byte1: %v\n", byte1)

	windowMode, err := c.BufferInRaw.G1()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Window Mode: %v\n", windowMode)

	canvasWidth, err := c.BufferInRaw.G2()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Canvas Width: %v\n", canvasWidth)

	canvasHeight, err := c.BufferInRaw.G2()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Canvas Height: %v\n", canvasHeight)

	prefInt, err := c.BufferInRaw.G1()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("PrefInt: %v\n", prefInt)

	uid, err := c.BufferInRaw.GData(24)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("UID: %v\n", uid)

	settings := c.BufferInRaw.GJStr()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Settings: %v\n", settings)

	affiliate, err := c.BufferInRaw.G4()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Affiliate: %v\n", affiliate)

	preferencesLen, err := c.BufferInRaw.G1()
	if err != nil {
		fmt.Println(err)
	}

	preferences, err := c.BufferInRaw.GData(int(preferencesLen))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Preferences: %v\n", preferences)

	verifyId, err := c.BufferInRaw.G2()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("VerifyID: %v\n", verifyId)

	checksums := make([]uint32, 29)
	for i := 0; i < 29; i++ {
		checksums[i], err = c.BufferInRaw.G4()
		if err != nil {
			fmt.Println(err)
		}
	}

	decrypted, err := c.BufferInRaw.RSADec()
	if err != nil {
		fmt.Println(err)
	}
	rsaMagic, err := decrypted.G1()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("rsaMagic: %v\n", rsaMagic)
	key := make([]uint32, 4)
	for i := 0; i < 4; i++ {
		key[i], err = decrypted.G4()
		if err != nil {
			fmt.Println(err)
		}
	}

	username1, err := decrypted.G8()
	if err != nil {
		fmt.Println(err)
	}
	username := util.FromBase37(username1)
	fmt.Printf("Username: %v\n", username)

	password := decrypted.GJStr()
	fmt.Printf("Password: %v\n", password)

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

	var response bytebuffer.ByteBuffer
	if opcode == util.LoginProtWorldReconnect {
		response.P1(15)
	} else {
		response.P1(2)
	}

	if opcode == util.LoginProtWorldConnect {
		response.P1(0)                 // staff mod level
		response.P1(0)                 // player mod level
		response.PBool(false)          // player underage
		response.PBool(false)          // parentalChatConsent
		response.PBool(false)          // parentalAdvertConsent
		response.PBool(false)          // mapQuickChat
		response.P2(uint16(player.ID)) // selfId
		response.PBool(false)          // MouseRecorder
		response.PBool(true)           // mapMembers
	}

	c.WriteRawSocket(response.Bytes())

	c.State = ClientStateGame
	c.Server.World.AddPlayer(player)

	fmt.Println("LOGIN COMPLETE")
}

func (c *Client) handleGame() {
	fmt.Println("entered handleGame()")
	data := c.BufferInRaw.Bytes()

	offset := 0
	for offset < len(data) {
		start := offset

		if c.RandomIn != nil {
			data[offset] -= byte(c.RandomIn.NextInt())
		}

		opcode := data[offset]
		offset += 1

		length := util.ClientProtLengths[opcode]

		if length == 255 {
			length = data[offset]
			offset += 1
		} else if length == 254 {
			length = data[offset]<<8 | data[offset+1]
			offset += 2
		}

		if int(length) > 30000-c.BufferInOffset {
			fmt.Println("Packet overflow for this tick")
			return
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
	Data bytebuffer.ByteBuffer
}

func (c *Client) DecodeIn() []DecodedData {
	fmt.Println("Client.DecodeIn() ENTERED")
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
			Data: *bytebuffer.NewBuffer(c.Server.BufferIn[c.BufferStart+offset : c.BufferStart+offset+int(length)]),
		})

		offset += int(length)
	}

	return decoded
}

func (c *Client) Write(data []byte) {
	fmt.Printf("Client.Write() DATA TO WRITE (%v) %v: %v\n", getCallerInfo(2), len(data), data)
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
	//fmt.Println("entered Client.Flush()")
	if c.BufferOutOffset > 0 {
		//fmt.Printf("Client.Flush DATA: %v\n", c.Server.BufferOut[c.BufferStart:c.BufferStart+c.BufferOutOffset])
		c.WriteRawSocket(c.Server.BufferOut[c.BufferStart : c.BufferStart+c.BufferOutOffset])
		c.BufferOutOffset = 0
	}
}

type NetOutData struct {
	Data    []byte
	Encrypt bool
}

func (c *Client) Queue(data []byte, encrypt bool) {
	//fmt.Println("Client.Queue() ENTERED")
	//c.NetOut = append(c.NetOut)
	c.NetOut = append(c.NetOut, NetOutData{
		Data:    data,
		Encrypt: encrypt,
	})
	fmt.Printf("Client.Queue() (%v): Encrypt %v, Data %v\n", getCallerInfo(2), encrypt, data)
}

func (c *Client) EncodeOut() {
	//fmt.Println("EncodeOut() ENTERED")
	for i := 0; i < len(c.NetOut); i++ {
		packet := c.NetOut[i]

		if /*c.RandomOut != nil &&*/ packet.Encrypt {
			// TODO: uint32 converted to byte!!
			packet.Data[0] += byte(c.RandomOut.NextInt())
		}

		c.Write(packet.Data)
	}
}

func (c *Client) WriteRawSocket(data []byte) {
	fmt.Printf("NET WRITE %v (%v): %v\n", len(data), getCallerInfo(2), data)
	n, err := c.Socket.Write(data)
	if err != nil {
		fmt.Printf("NET WRITE ERROR: %v\n", err)
	}
	fmt.Printf("NET WRITE %v DONE\n", n)
}

func getCallerInfo(skip int) (info string) {
	pc, file, lineNo, ok := runtime.Caller(skip)
	if !ok {
		info = "runtime.Caller() failed"
		return
	}
	funcName := runtime.FuncForPC(pc).Name()
	fileName := path.Base(file) // The Base function returns the last element of the path
	return fmt.Sprintf("FuncName:%s, file:%s, line:%d ", funcName, fileName, lineNo)
}
