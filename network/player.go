package network

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
	"time"

	"github.com/zsrv/rt5-server-go/util"
	"github.com/zsrv/rt5-server-go/util/bytebuffer"
)

// TODO: Could not put this file in the engine package because there is
// a circular dependency between player.go and client.go!

const (
	MessageTypeGame        = 0
	MessageTypePublic      = 1
	MessageTypeTrade       = 4
	MessageTypePrivateTo   = 6
	MessageTypePrivateFrom = 7
	MessageTypeAssist      = 10
	MessageTypeDevConsole  = 99
)

type Player struct {
	Client *Client

	FirstLoad    bool
	Reconnecting bool
	Loaded       bool
	Loading      bool
	Appearance   *bytebuffer.ByteBuffer
	Placement    bool
	VerifyID     int

	ID         int
	Username   string
	WindowMode uint8

	World *World

	LastPos *util.Position

	Pos *util.Position
}

func NewPlayer(client *Client) *Player {
	return &Player{
		Client: client,

		FirstLoad:    true,
		Reconnecting: false,
		Loaded:       false,
		Loading:      false,
		Appearance:   nil,
		Placement:    false,
		VerifyID:     1,

		ID:         1,
		Username:   "",
		WindowMode: 0,

		LastPos: util.NewPosition(0, 0, 0),

		// make-over mage: 2925, 3323, 0
		// varrock square: 3213, 3443
		Pos: util.NewPosition(3162, 3490, 0),
	}
}

func (p *Player) Tick() {
	//fmt.Print(" PTick ")
	if !p.Loaded && !p.Loading {
		p.Loading = true

		if p.Reconnecting {
			var response bytebuffer.ByteBuffer
			response.P2(0)
			start := response.Len() // offset

			// INIT_GPI

			response.AccessBits()
			response.PBit(30, p.Pos.HighRes())
			p.LastPos.Clone(p.Pos)

			for i := 1; i < 2048; i++ {
				if p.ID == i {
					continue
				}
				response.PBit(18, 0)
			}
			response.AccessBytes()

			response.PSize2(response.Len() - start)
			debugBytes := response.Bytes()
			//p.Client.Queue(response.Bytes(), false)
			p.Client.Queue(debugBytes, false)
			fmt.Printf("Player.tick reconnecting response %v: %v\n", len(debugBytes), debugBytes)
		} else if p.FirstLoad {
			var response bytebuffer.ByteBuffer
			response.P1(98)
			response.P2(0)
			start := response.Len()

			// INIT_GPI

			response.AccessBits()
			response.PBit(30, p.Pos.HighRes())
			p.LastPos.Clone(p.Pos)

			for i := 1; i < 2048; i++ {
				if p.ID == i {
					continue
				}
				response.PBit(18, 0)
			}
			response.AccessBytes()

			// REBUILD_NORMAL

			response.IP2(uint16(p.Pos.ZoneX()))
			response.P2(uint16(p.Pos.ZoneZ()))
			response.P1(uint8(p.Pos.BAIndex))
			response.P1Neg(0)

			for mapsquareX := (p.Pos.ZoneX() - (p.Pos.BASizeX >> 4)) >> 3; mapsquareX <= (p.Pos.ZoneX()+(p.Pos.BASizeX>>4))>>3; mapsquareX++ {
				for mapsquareZ := (p.Pos.ZoneZ() - (p.Pos.BASizeZ >> 4)) >> 3; mapsquareZ <= (p.Pos.ZoneZ()+(p.Pos.BASizeZ>>4))>>3; mapsquareZ++ {
					xtea, found := util.GetXTEA(mapsquareX, mapsquareZ)
					if found {
						for i := 0; i < len(xtea.Key); i++ {
							// TODO: converting signed to unsigned!!
							response.P4(uint32(xtea.Key[i]))
						}
					} else {
						for i := 0; i < 4; i++ {
							response.P4(0)
						}
					}
				}
			}

			response.PSize2(response.Len() - start)
			debugBytes := response.Bytes()
			//p.Client.Queue(response.Bytes(), true)
			p.Client.Queue(debugBytes, true)
			fmt.Printf("Player.tick first load response %v: %v\n", len(debugBytes), debugBytes)
		}

		if p.FirstLoad {
			if p.IsClientResizable() {
				p.OpenGameFrame(746) // fixed?
			} else {
				p.OpenGameFrame(548) // resizable?
			}
			p.OpenChatBox(752)

			p.OpenTab(0, 884)
			p.OpenTab(1, 320)
			p.OpenTab(2, 190)
			p.OpenTab(3, 259)
			p.OpenTab(4, 149)
			p.OpenTab(5, 387)
			p.OpenTab(6, 271)
			p.OpenTab(7, 192)
			p.OpenTab(8, 891)
			p.OpenTab(9, 550)
			p.OpenTab(10, 551)
			p.OpenTab(11, 589)
			p.OpenTab(12, 261)
			p.OpenTab(13, 464)
			p.OpenTab(14, 187)
			p.OpenTab(15, 34)
			p.OpenTab(16, 182)

			if !p.Reconnecting {
				p.MessageGame("Welcome to RuneScape.", MessageTypeGame, "", "")
			}
		}

		p.FirstLoad = false
		p.Loading = false
		p.Loaded = true
	}

	// player info
	if p.Loaded {
		var response bytebuffer.ByteBuffer
		var updateBlock bytebuffer.ByteBuffer

		response.P1(72)
		response.P2(0)
		start := response.Len() // offset

		p.ProcessActivePlayers(&response, &updateBlock, true)
		p.ProcessActivePlayers(&response, &updateBlock, false)
		p.ProcessInactivePlayers(&response, &updateBlock, true)
		p.ProcessInactivePlayers(&response, &updateBlock, false)
		response.PData(&updateBlock)

		response.PSize2(response.Len() - start) // offset
		//p.Client.Queue(response.Bytes(), true)
		debugBytes := response.Bytes()
		p.Client.Queue(debugBytes, true)
		fmt.Printf("Player.tick player info response %v: %v\n", len(debugBytes), debugBytes)
	}
}

func (p *Player) IsClientResizable() bool {
	// 1 = fixed, 2 = resizable, 3 = fullscreen
	return p.WindowMode > 1
}

func (p *Player) ProcessActivePlayers(buf *bytebuffer.ByteBuffer, updateBlock *bytebuffer.ByteBuffer, nsn0 bool) {
	buf.AccessBits()
	// TODO: this is supposed to loop, and "nsn0" is supposed to check against a player flag to skip
	if nsn0 {
		needsMaskUpdate := p.Appearance == nil
		needsUpdate := p.Placement || needsMaskUpdate

		if needsUpdate {
			buf.PBit(1, 1)
		} else {
			buf.PBit(1, 0)
		}

		if needsUpdate {
			if needsMaskUpdate {
				buf.PBit(1, 1)
			} else {
				buf.PBit(1, 0)
			}

			buf.PBit(2, 0) // no further update

			//if p.Placement {
			//	buf.PBit(2, 3) // teleport
			//	buf.PBit(1, 1) // full location update
			//	buf.PBit(30, p.Pos.Z | p.Pos.X << 14 | p.Pos.Plane << 28)
			//}
		}

		if needsMaskUpdate {
			p.AppendUpdateBlock(updateBlock)
		}
	}
	buf.AccessBytes()
}

func (p *Player) ProcessInactivePlayers(buf *bytebuffer.ByteBuffer, updateBlock *bytebuffer.ByteBuffer, nsn2 bool) {
	buf.AccessBits()
	// TODO: "nsn2" is supposed to check against a player flag to skip
	if nsn2 {
		for i := 1; i < 2048; i++ {
			if p.ID == i {
				continue
			}

			buf.PBit(1, 0)
			buf.PBit(2, 0)
		}
	}
	buf.AccessBytes()
}

func (p *Player) GenerateAppearance() {
	var buf bytebuffer.ByteBuffer

	buf.P1(0)   // flags
	buf.P1(255) // title-related (was -1)
	buf.P1(255) // pkIcon (was -1)
	buf.P1(255) // prayerIcon

	//for i := 0; i < 12; i++ {
	//	buf.P1(0) // body
	//}

	// hat, cape, amulet, weapon, chest, shield, arms, legs, hair, wrists, hands, feet, beard
	body := []uint8{255, 255, 255, 255, 18, 255, 26, 36, 0, 33, 42, 10}
	for i := 0; i < len(body); i++ {
		if body[i] == 255 {
			buf.P1(0)
		} else {
			body[i] += uint8(math.Floor(rand.Float64() * 2))
			buf.P2(uint16(body[i]) | 0x100)
		}
	}

	for i := 0; i < 5; i++ {
		buf.P1(uint8(math.Floor(rand.Float64() * 4))) // color
	}

	buf.P2(1426) // bas id
	buf.PJStr(p.Username)
	buf.P1(3)  // combat level
	buf.P2(33) // total level
	buf.P1(0)  // sound radius

	p.Appearance = new(bytebuffer.ByteBuffer)
	p.Appearance.IPData(buf.Bytes())
}

func (p *Player) AppendUpdateBlock(buf *bytebuffer.ByteBuffer) {
	var flags uint8 = 0

	if p.Appearance == nil {
		p.GenerateAppearance()
		flags |= 0x1
	}

	buf.P1(flags)

	if flags&0x1 == 1 {
		buf.P1Sub(uint8(p.Appearance.Len()))
		buf.PData(p.Appearance)
	}
}

func (p *Player) ProcessIn() {
	decoded := p.Client.DecodeIn()

	for _, v := range decoded {
		switch v.ID {
		//case 78: // MOVE_GAMECLICK
		//	ctrlClick, err := v.Data.G1() // g1add
		//	if err != nil {
		//		fmt.Println(err)
		//	}
		//
		//	x, err := v.Data.G2()
		//	if err != nil {
		//		fmt.Println(err)
		//	}
		//
		//	z := v.Data.IG2()
		//
		//	p.Pos.X = int(x)
		//	p.Pos.Z = int(z)
		//
		//	if ctrlClick != 0 {
		//		p.Placement = true
		//	}
		case util.ClientProtClientCheat:
			//tele, err := v.Data.G1()
			//if err != nil {
			//	fmt.Println(err)
			//}

			cmd := strings.ToLower(v.Data.GJStr())
			args := strings.Split(cmd, " ")

			cmd = args[0]
			args = args[1:]

			if cmd == "logout" {
				p.Logout()
			}
		default:
			fmt.Println("Unhandled packet", v.ID)
		}
	}
}

// events

func (p *Player) OpenChatBox(interfaceID uint16) {
	if interfaceID == 752 {
		if p.IsClientResizable() {
			p.OpenInterface(746, 15, 751, 3)
			p.OpenInterface(746, 18, 752, 3)
		} else {
			p.OpenInterface(548, 20, 751, 3)
			p.OpenInterface(548, 142, 752, 3)
		}

		if p.IsClientResizable() {
			p.OpenInterface(752, 9, 137, 3)
		}
	}
}

func (p *Player) OpenTab(tabID uint16, interfaceID uint16) {
	if p.IsClientResizable() {
		p.OpenInterface(746, 33+tabID, interfaceID, 3)
	} else {
		p.OpenInterface(548, 152+tabID, interfaceID, 3)
	}
}

// encoders

func (p *Player) Logout() {
	var response bytebuffer.ByteBuffer
	response.P1(58)
	p.Client.Queue(response.Bytes(), true)
}

func (p *Player) MessageGame(msg string, msgType uint8, msg2 string, msg3 string) {
	var response bytebuffer.ByteBuffer
	response.P1(99)
	response.P1(0)
	start := response.Len() // offset

	response.PSmart(uint16(msgType))
	response.P4(uint32(time.Now().UnixMilli() / 1000))

	var more uint8 = 0
	if msg2 != "" {
		more |= 0x1
	}

	if msg3 != "" {
		more |= 0x2
	}

	response.P1(more)
	if more&0x1 != 0 {
		response.PJStr(msg2)
	}

	if more&0x2 != 0 {
		response.PJStr(msg3)
	}

	response.PJStr(msg)

	response.PSize1(response.Len() - start)

	p.Client.Queue(response.Bytes(), true)
}

func (p *Player) OpenGameFrame(interfaceId uint16) {
	var response bytebuffer.ByteBuffer
	response.P1(93)

	response.P1(0)
	response.IP2(interfaceId)
	response.IP2(uint16(p.VerifyID))
	p.VerifyID += 1

	p.Client.Queue(response.Bytes(), true)
}

func (p *Player) OpenInterface(windowID uint16, componentId uint16, interfaceId uint16, flags uint8) {
	var response bytebuffer.ByteBuffer
	response.P1(52)

	response.P2Add(uint8(p.VerifyID))
	p.VerifyID += 1
	response.P1Sub(flags)
	response.IP2(componentId)
	response.IP2(windowID)
	response.P2(interfaceId)

	p.Client.Queue(response.Bytes(), true)
}
