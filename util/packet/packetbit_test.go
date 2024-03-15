package packet

import (
	"slices"
	"testing"

	"github.com/zsrv/rt5-server-go/util/isaacrandom"
)

func TestPacketBit_AccessBits(t *testing.T) {
	type fields struct {
		Packet    Packet
		bitOffset int
		random    *isaacrandom.IsaacRandom
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PacketBit{
				Packet:    tt.fields.Packet,
				bitOffset: tt.fields.bitOffset,
				random:    tt.fields.random,
			}
			p.AccessBits()
		})
	}
}

func TestPacketBit_AccessBytes(t *testing.T) {
	type fields struct {
		Packet    Packet
		bitOffset int
		random    *isaacrandom.IsaacRandom
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PacketBit{
				Packet:    tt.fields.Packet,
				bitOffset: tt.fields.bitOffset,
				random:    tt.fields.random,
			}
			p.AccessBytes()
		})
	}
}

func TestPacketBit_AvailableBits(t *testing.T) {
	type fields struct {
		Packet    Packet
		bitOffset int
		random    *isaacrandom.IsaacRandom
	}
	type args struct {
		pos int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PacketBit{
				Packet:    tt.fields.Packet,
				bitOffset: tt.fields.bitOffset,
				random:    tt.fields.random,
			}
			if got := p.AvailableBits(tt.args.pos); got != tt.want {
				t.Errorf("AvailableBits() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPacketBit_G1Isaac(t *testing.T) {
	type fields struct {
		Packet    Packet
		bitOffset int
		random    *isaacrandom.IsaacRandom
	}
	tests := []struct {
		name   string
		fields fields
		want   uint32
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PacketBit{
				Packet:    tt.fields.Packet,
				bitOffset: tt.fields.bitOffset,
				random:    tt.fields.random,
			}
			if got := p.G1Isaac(); got != tt.want {
				t.Errorf("G1Isaac() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPacketBit_GBit(t *testing.T) {
	type fields struct {
		Packet    Packet
		bitOffset int
		random    *isaacrandom.IsaacRandom
	}
	type args struct {
		n int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   uint8
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PacketBit{
				Packet:    tt.fields.Packet,
				bitOffset: tt.fields.bitOffset,
				random:    tt.fields.random,
			}
			if got := p.GBit(tt.args.n); got != tt.want {
				t.Errorf("GBit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPacketBit_GIsaac(t *testing.T) {
	type fields struct {
		Packet    Packet
		bitOffset int
		random    *isaacrandom.IsaacRandom
	}
	type args struct {
		dest   []uint8
		length int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PacketBit{
				Packet:    tt.fields.Packet,
				bitOffset: tt.fields.bitOffset,
				random:    tt.fields.random,
			}
			p.GIsaac(tt.args.dest, tt.args.length)
		})
	}
}

func TestPacketBit_P1Isaac(t *testing.T) {
	type fields struct {
		Packet    Packet
		bitOffset int
		random    *isaacrandom.IsaacRandom
	}
	type args struct {
		op uint8
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PacketBit{
				Packet:    tt.fields.Packet,
				bitOffset: tt.fields.bitOffset,
				random:    tt.fields.random,
			}
			p.P1Isaac(tt.args.op)
		})
	}
}

func TestPacketBit_PBit(t *testing.T) {
	type fields struct {
		Packet    Packet
		bitOffset int
		random    *isaacrandom.IsaacRandom
	}
	type args struct {
		n     int
		value int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []byte
	}{
		{
			name: "valid",
			fields: fields{
				Packet: Packet{
					Buf:      []byte{98, 0, 0},
					Pos:      0,
					lastRead: 0,
				},
				bitOffset: 24,
				random:    nil,
			},
			args: args{
				n:     30,
				value: 51809698,
			},
			want: []byte{98, 0, 0, 12, 90, 54, 136},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PacketBit{
				Packet:    tt.fields.Packet,
				bitOffset: tt.fields.bitOffset,
				random:    tt.fields.random,
			}
			p.PBit(tt.args.n, tt.args.value)
			if got := p.Bytes(); !slices.Equal(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPacketBit_Peek1Isaac(t *testing.T) {
	type fields struct {
		Packet    Packet
		bitOffset int
		random    *isaacrandom.IsaacRandom
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PacketBit{
				Packet:    tt.fields.Packet,
				bitOffset: tt.fields.bitOffset,
				random:    tt.fields.random,
			}
			if got := p.Peek1Isaac(); got != tt.want {
				t.Errorf("Peek1Isaac() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPacketBit_SetKey(t *testing.T) {
	type fields struct {
		Packet    Packet
		bitOffset int
		random    *isaacrandom.IsaacRandom
	}
	type args struct {
		key []uint32
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PacketBit{
				Packet:    tt.fields.Packet,
				bitOffset: tt.fields.bitOffset,
				random:    tt.fields.random,
			}
			p.SetKey(tt.args.key)
		})
	}
}
