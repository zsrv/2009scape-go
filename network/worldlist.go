package network

import (
	"hash/crc32"

	"github.com/zsrv/rt5-server-go/util/packet"
)

// TODO: had to move it into network because of an import cycle
// between util/packet and util

type CountriesS struct {
	ID          int
	DisplayName string
	Flag        int
}

// TODO: is there not actually an ID, it's just written in a sequence?
var CountriesList = []CountriesS{
	{ID: 0, DisplayName: "United States", Flag: 0},
	{ID: 1, DisplayName: "Austria", Flag: 15},
	{ID: 2, DisplayName: "Australia", Flag: 16},
	{ID: 3, DisplayName: "Germany", Flag: 22},
	{ID: 4, DisplayName: "Brazil", Flag: 31},
	{ID: 5, DisplayName: "Canada", Flag: 38},
	{ID: 6, DisplayName: "Switzerland", Flag: 43},
	{ID: 7, DisplayName: "China", Flag: 48},
	{ID: 8, DisplayName: "Denmark", Flag: 58},
	{ID: 9, DisplayName: "Finland", Flag: 69},
	{ID: 10, DisplayName: "France", Flag: 74},
	{ID: 11, DisplayName: "United Kingdom", Flag: 77},
	{ID: 12, DisplayName: "Ireland", Flag: 101},
	{ID: 13, DisplayName: "India", Flag: 103},
	{ID: 14, DisplayName: "Mexico", Flag: 152},
	{ID: 15, DisplayName: "Netherlands", Flag: 161},
	{ID: 16, DisplayName: "Norway", Flag: 162},
	{ID: 17, DisplayName: "New Zealand", Flag: 166},
	{ID: 18, DisplayName: "Portugal", Flag: 179},
	{ID: 19, DisplayName: "Sweden", Flag: 191},
}

type WorldParameters struct {
	ID        int
	Hostname  string
	Port      int
	Country   int
	Activity  string
	Members   bool
	QuickChat bool
	PvP       bool
	LootShare bool
	Highlight bool
	Players   int
}

var WorldList = []WorldParameters{
	{
		ID:        1,
		Hostname:  "localhost",
		Port:      43594,
		Country:   6,
		Activity:  "",
		Members:   true,
		QuickChat: false,
		PvP:       false,
		LootShare: true,
		Highlight: false,
		Players:   5,
	},
	{
		ID:        2,
		Hostname:  "localhost",
		Port:      43594,
		Country:   6,
		Activity:  "Activity Name",
		Members:   false,
		QuickChat: false,
		PvP:       false,
		LootShare: true,
		Highlight: true,
		Players:   5,
	},
}

var WorldListRaw packet.Packet
var WorldListChecksum uint32 = 0

func init() {
	WorldListRaw.PSmart(uint16(len(CountriesList)))

	for _, v := range CountriesList {
		WorldListRaw.PSmart(uint16(v.Flag))
		WorldListRaw.PJStr2(v.DisplayName)
	}
}

var MinID int
var MaxID int

func init() {
	initialized := false
	for _, v := range WorldList {
		if !initialized {
			MinID = v.ID
			MaxID = v.ID
			initialized = true
		}
		if v.ID < MinID {
			MinID = v.ID
		}
		if v.ID > MaxID {
			MaxID = v.ID
		}
	}
}

func init() {
	WorldListRaw.PSmart(uint16(MinID))
	WorldListRaw.PSmart(uint16(MaxID))
	WorldListRaw.PSmart(uint16(len(WorldList)))
}

func init() {
	for _, world := range WorldList {
		WorldListRaw.PSmart(uint16(world.ID - MinID))
		WorldListRaw.P1(uint8(world.Country))

		var flags uint32 = 0

		if world.Members {
			flags |= 0x1
		}

		if world.QuickChat {
			flags |= 0x2
		}

		if world.PvP {
			//flags |= 0x4
		}

		if world.LootShare {
			flags |= 0x8
		}

		if world.Activity != "" && world.Highlight {
			flags |= 0x10
		}

		WorldListRaw.P4(flags)

		// if there is no activity name, client will fall back to country flag + name
		WorldListRaw.PJStr2(world.Activity)
		WorldListRaw.PJStr2(world.Hostname)
	}
}

func init() {
	b := WorldListRaw.Bytes()
	WorldListRaw = *packet.NewPacket(b) // bad workaround since we can't calc checksum without clearing buffer

	WorldListChecksum = crc32.ChecksumIEEE(b)
}

const (
	WlProtOutSuccess = 0
	WlProtOutReject  = 1
)
