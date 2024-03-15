package engine

import (
	"time"
)

type World struct {
	Players []*Player
}

func NewWorld() *World {
	// the client index starts at 1
	w := &World{
		Players: make([]*Player, 2046),
	}
	w.Tick()
	return w
}

func (w *World) RegisterPlayer(player *Player) {
	for i := range w.Players {
		if w.Players[i] == nil {
			player.ID = i + 1
			break
		}
	}
}

func (w *World) AddPlayer(player *Player) {
	player.World = w
	w.Players[player.ID-1] = player
}

func (w *World) RemovePlayer(client Client) {
	w.Players[client.Player.ID-1] = nil
}

func (w *World) Tick() {
	start := time.Now().UnixMilli()

	// read packets
	for _, v := range w.Players {
		if v == nil {
			//return
			continue
		}

		v.ProcessIn()
	}

	// npc processing
	// player processing
	for _, v := range w.Players {
		if v == nil {
			//return
			continue
		}

		v.Tick()
	}

	// game tasks
	// flushing packets
	for _, v := range w.Players {
		if v == nil {
			//return
			continue
		}

		if len(v.Client.NetOut) > 0 {
			v.Client.EncodeOut()
			v.Client.NetOut = make([]NetOutData, 0)
		}

		v.Client.Flush()
		v.Client.ResetIn()

		v.Placement = false
	}

	// npc aggro etc
	end := time.Now().UnixMilli()

	delta := 600 - (end - start)
	go func() {
		time.Sleep(time.Duration(delta) * time.Millisecond)
		w.Tick()
	}()
}
