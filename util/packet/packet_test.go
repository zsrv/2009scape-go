package packet

import (
	"math/big"
	"slices"
	"testing"
)

func TestGetCRC(t *testing.T) {
	type args struct {
		length int
		offset int
		src    []uint8
	}
	tests := []struct {
		name string
		args args
		want uint32
	}{
		{
			name: "valid",
			args: args{
				length: 8,
				offset: 0,
				src:    []byte{1, 2, 3, 4, 5, 6, 7, 8},
			},
			want: 1070237893,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetCRC(tt.args.length, tt.args.offset, tt.args.src); got != tt.want {
				t.Errorf("GetCRC() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPJStrLen(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "valid",
			args: args{
				str: "testing",
			},
			want: 8,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PJStrLen(tt.args.str); got != tt.want {
				t.Errorf("PJStrLen() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPacket_AddCRC(t *testing.T) {
	type fields struct {
		Buf      []byte
		Pos      int
		lastRead readOp
	}
	type args struct {
		off int
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		want      uint32
		wantBytes []byte
	}{
		{
			name: "valid, offset 0",
			fields: fields{
				Buf:      []byte{1, 2, 3, 4, 5, 6, 7, 8},
				Pos:      0,
				lastRead: 0,
			},
			args: args{
				off: 0,
			},
			want:      1070237893,
			wantBytes: []byte{1, 2, 3, 4, 5, 6, 7, 8, 63, 202, 136, 197},
		},
		{
			name: "valid, offset 4",
			fields: fields{
				Buf:      []byte{1, 2, 3, 4, 5, 6, 7, 8},
				Pos:      0,
				lastRead: 0,
			},
			args: args{
				off: 4,
			},
			want:      1401769321,
			wantBytes: []byte{1, 2, 3, 4, 5, 6, 7, 8, 83, 141, 77, 105},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Packet{
				Buf:      tt.fields.Buf,
				Pos:      tt.fields.Pos,
				lastRead: tt.fields.lastRead,
			}
			if got := p.AddCRC(tt.args.off); got != tt.want {
				t.Errorf("AddCRC() = %v, want %v", got, tt.want)
			}
			if gotBytes := p.Bytes(); !slices.Equal(gotBytes, tt.wantBytes) {
				t.Errorf("got %v, want %v", gotBytes, tt.wantBytes)
			}
		})
	}
}

func TestPacket_CheckCRC(t *testing.T) {
	type fields struct {
		Buf      []byte
		Pos      int
		lastRead readOp
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "valid",
			fields: fields{
				Buf:      []byte{1, 2, 3, 4, 5, 6, 7, 8, 63, 202, 136, 197},
				Pos:      0,
				lastRead: 0,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Packet{
				Buf:      tt.fields.Buf,
				Pos:      tt.fields.Pos,
				lastRead: tt.fields.lastRead,
			}
			if got := p.CheckCRC(); got != tt.want {
				t.Errorf("CheckCRC() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPacket_FastGJStr(t *testing.T) {
	type fields struct {
		Buf      []byte
		Pos      int
		lastRead readOp
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "password",
			fields: fields{
				Buf:      []byte{112, 97, 115, 115, 119, 111, 114, 100, 0},
				Pos:      0,
				lastRead: 0,
			},
			want: "password",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Packet{
				Buf:      tt.fields.Buf,
				Pos:      tt.fields.Pos,
				lastRead: tt.fields.lastRead,
			}
			if got := p.FastGJStr(); got != tt.want {
				t.Errorf("FastGJStr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPacket_G1(t *testing.T) {
	type fields struct {
		Buf      []byte
		Pos      int
		lastRead readOp
	}
	tests := []struct {
		name   string
		fields fields
		want   uint8
	}{
		{
			name: "valid",
			fields: fields{
				Buf:      []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
				Pos:      0,
				lastRead: 0,
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Packet{
				Buf:      tt.fields.Buf,
				Pos:      tt.fields.Pos,
				lastRead: tt.fields.lastRead,
			}
			if got := p.G1(); got != tt.want {
				t.Errorf("G1() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPacket_G1Alt1(t *testing.T) {
	type fields struct {
		Buf      []byte
		Pos      int
		lastRead readOp
	}
	tests := []struct {
		name   string
		fields fields
		want   uint8
	}{
		{
			name: "valid",
			fields: fields{
				Buf:      []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
				Pos:      0,
				lastRead: 0,
			},
			want: 0x81,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Packet{
				Buf:      tt.fields.Buf,
				Pos:      tt.fields.Pos,
				lastRead: tt.fields.lastRead,
			}
			if got := p.G1Alt1(); got != tt.want {
				t.Errorf("G1Alt1() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPacket_G1Alt2(t *testing.T) {
	type fields struct {
		Buf      []byte
		Pos      int
		lastRead readOp
	}
	tests := []struct {
		name   string
		fields fields
		want   uint8
	}{
		{
			name: "valid",
			fields: fields{
				Buf:      []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
				Pos:      0,
				lastRead: 0,
			},
			want: 0xFF,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Packet{
				Buf:      tt.fields.Buf,
				Pos:      tt.fields.Pos,
				lastRead: tt.fields.lastRead,
			}
			if got := p.G1Alt2(); got != tt.want {
				t.Errorf("G1Alt2() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPacket_G1Alt3(t *testing.T) {
	type fields struct {
		Buf      []byte
		Pos      int
		lastRead readOp
	}
	tests := []struct {
		name   string
		fields fields
		want   uint8
	}{
		{
			name: "valid",
			fields: fields{
				Buf:      []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
				Pos:      0,
				lastRead: 0,
			},
			want: 0x7F,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Packet{
				Buf:      tt.fields.Buf,
				Pos:      tt.fields.Pos,
				lastRead: tt.fields.lastRead,
			}
			if got := p.G1Alt3(); got != tt.want {
				t.Errorf("G1Alt3() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPacket_G1B(t *testing.T) {
	type fields struct {
		Buf      []byte
		Pos      int
		lastRead readOp
	}
	tests := []struct {
		name   string
		fields fields
		want   int8
	}{
		{
			name: "valid",
			fields: fields{
				Buf:      []byte{150, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
				Pos:      0,
				lastRead: 0,
			},
			want: -106,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Packet{
				Buf:      tt.fields.Buf,
				Pos:      tt.fields.Pos,
				lastRead: tt.fields.lastRead,
			}
			if got := p.G1B(); got != tt.want {
				t.Errorf("G1B() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPacket_G1BAlt1(t *testing.T) {
	type fields struct {
		Buf      []byte
		Pos      int
		lastRead readOp
	}
	tests := []struct {
		name   string
		fields fields
		want   int8
	}{
		{
			name: "valid",
			fields: fields{
				Buf:      []byte{150, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
				Pos:      0,
				lastRead: 0,
			},
			want: 22,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Packet{
				Buf:      tt.fields.Buf,
				Pos:      tt.fields.Pos,
				lastRead: tt.fields.lastRead,
			}
			if got := p.G1BAlt1(); got != tt.want {
				t.Errorf("G1BAlt1() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPacket_G1BAlt2(t *testing.T) {
	type fields struct {
		Buf      []byte
		Pos      int
		lastRead readOp
	}
	tests := []struct {
		name   string
		fields fields
		want   int8
	}{
		{
			name: "valid",
			fields: fields{
				Buf:      []byte{150, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
				Pos:      0,
				lastRead: 0,
			},
			want: 106,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Packet{
				Buf:      tt.fields.Buf,
				Pos:      tt.fields.Pos,
				lastRead: tt.fields.lastRead,
			}
			if got := p.G1BAlt2(); got != tt.want {
				t.Errorf("G1BAlt2() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPacket_G1BAlt3(t *testing.T) {
	type fields struct {
		Buf      []byte
		Pos      int
		lastRead readOp
	}
	tests := []struct {
		name   string
		fields fields
		want   int8
	}{
		{
			name: "valid",
			fields: fields{
				Buf:      []byte{150, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
				Pos:      0,
				lastRead: 0,
			},
			want: -0x16,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Packet{
				Buf:      tt.fields.Buf,
				Pos:      tt.fields.Pos,
				lastRead: tt.fields.lastRead,
			}
			if got := p.G1BAlt3(); got != tt.want {
				t.Errorf("G1BAlt3() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPacket_G2(t *testing.T) {
	type fields struct {
		Buf      []byte
		Pos      int
		lastRead readOp
	}
	tests := []struct {
		name   string
		fields fields
		want   uint16
	}{
		{
			name: "valid",
			fields: fields{
				Buf:      []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
				Pos:      0,
				lastRead: 0,
			},
			want: 0x0102,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Packet{
				Buf:      tt.fields.Buf,
				Pos:      tt.fields.Pos,
				lastRead: tt.fields.lastRead,
			}
			if got := p.G2(); got != tt.want {
				t.Errorf("G2() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPacket_G2Alt2(t *testing.T) {
	type fields struct {
		Buf      []byte
		Pos      int
		lastRead readOp
	}
	tests := []struct {
		name   string
		fields fields
		want   uint16
	}{
		{
			name: "valid",
			fields: fields{
				Buf:      []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
				Pos:      0,
				lastRead: 0,
			},
			want: 0x0182,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Packet{
				Buf:      tt.fields.Buf,
				Pos:      tt.fields.Pos,
				lastRead: tt.fields.lastRead,
			}
			if got := p.G2Alt2(); got != tt.want {
				t.Errorf("G2Alt2() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPacket_G2Alt3(t *testing.T) {
	type fields struct {
		Buf      []byte
		Pos      int
		lastRead readOp
	}
	tests := []struct {
		name   string
		fields fields
		want   uint16
	}{
		{
			name: "valid",
			fields: fields{
				Buf:      []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
				Pos:      0,
				lastRead: 0,
			},
			want: 0x0281,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Packet{
				Buf:      tt.fields.Buf,
				Pos:      tt.fields.Pos,
				lastRead: tt.fields.lastRead,
			}
			if got := p.G2Alt3(); got != tt.want {
				t.Errorf("G2Alt3() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPacket_G2S(t *testing.T) {
	type fields struct {
		Buf      []byte
		Pos      int
		lastRead readOp
	}
	tests := []struct {
		name   string
		fields fields
		want   int16
	}{
		{
			name: "valid",
			fields: fields{
				Buf:      []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
				Pos:      0,
				lastRead: 0,
			},
			want: 0x0102,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Packet{
				Buf:      tt.fields.Buf,
				Pos:      tt.fields.Pos,
				lastRead: tt.fields.lastRead,
			}
			if got := p.G2S(); got != tt.want {
				t.Errorf("G2S() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPacket_G2SAlt1(t *testing.T) {
	type fields struct {
		Buf      []byte
		Pos      int
		lastRead readOp
	}
	tests := []struct {
		name   string
		fields fields
		want   int16
	}{
		{
			name: "valid",
			fields: fields{
				Buf:      []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
				Pos:      0,
				lastRead: 0,
			},
			want: 0x0201,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Packet{
				Buf:      tt.fields.Buf,
				Pos:      tt.fields.Pos,
				lastRead: tt.fields.lastRead,
			}
			if got := p.G2SAlt1(); got != tt.want {
				t.Errorf("G2SAlt1() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPacket_G2SAlt2(t *testing.T) {
	type fields struct {
		Buf      []byte
		Pos      int
		lastRead readOp
	}
	tests := []struct {
		name   string
		fields fields
		want   int16
	}{
		{
			name: "valid",
			fields: fields{
				Buf:      []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
				Pos:      0,
				lastRead: 0,
			},
			want: 0x0182,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Packet{
				Buf:      tt.fields.Buf,
				Pos:      tt.fields.Pos,
				lastRead: tt.fields.lastRead,
			}
			if got := p.G2SAlt2(); got != tt.want {
				t.Errorf("G2SAlt2() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPacket_G2SAlt3(t *testing.T) {
	type fields struct {
		Buf      []byte
		Pos      int
		lastRead readOp
	}
	tests := []struct {
		name   string
		fields fields
		want   int16
	}{
		{
			name: "valid",
			fields: fields{
				Buf:      []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
				Pos:      0,
				lastRead: 0,
			},
			want: 0x0281,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Packet{
				Buf:      tt.fields.Buf,
				Pos:      tt.fields.Pos,
				lastRead: tt.fields.lastRead,
			}
			if got := p.G2SAlt3(); got != tt.want {
				t.Errorf("G2SAlt3() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPacket_G3(t *testing.T) {
	type fields struct {
		Buf      []byte
		Pos      int
		lastRead readOp
	}
	tests := []struct {
		name   string
		fields fields
		want   uint32
	}{
		{
			name: "valid",
			fields: fields{
				Buf:      []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
				Pos:      0,
				lastRead: 0,
			},
			want: 0x010203,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Packet{
				Buf:      tt.fields.Buf,
				Pos:      tt.fields.Pos,
				lastRead: tt.fields.lastRead,
			}
			if got := p.G3(); got != tt.want {
				t.Errorf("G3() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPacket_G4(t *testing.T) {
	type fields struct {
		Buf      []byte
		Pos      int
		lastRead readOp
	}
	tests := []struct {
		name   string
		fields fields
		want   uint32
	}{
		{
			name: "valid",
			fields: fields{
				Buf:      []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
				Pos:      0,
				lastRead: 0,
			},
			want: 0x01020304,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Packet{
				Buf:      tt.fields.Buf,
				Pos:      tt.fields.Pos,
				lastRead: tt.fields.lastRead,
			}
			if got := p.G4(); got != tt.want {
				t.Errorf("G4() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPacket_G4Alt1(t *testing.T) {
	type fields struct {
		Buf      []byte
		Pos      int
		lastRead readOp
	}
	tests := []struct {
		name   string
		fields fields
		want   uint32
	}{
		{
			name: "valid",
			fields: fields{
				Buf:      []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
				Pos:      0,
				lastRead: 0,
			},
			want: 0x04030201,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Packet{
				Buf:      tt.fields.Buf,
				Pos:      tt.fields.Pos,
				lastRead: tt.fields.lastRead,
			}
			if got := p.G4Alt1(); got != tt.want {
				t.Errorf("G4Alt1() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPacket_G4Alt2(t *testing.T) {
	type fields struct {
		Buf      []byte
		Pos      int
		lastRead readOp
	}
	tests := []struct {
		name   string
		fields fields
		want   uint32
	}{
		{
			name: "valid",
			fields: fields{
				Buf:      []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
				Pos:      0,
				lastRead: 0,
			},
			want: 0x03040102,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Packet{
				Buf:      tt.fields.Buf,
				Pos:      tt.fields.Pos,
				lastRead: tt.fields.lastRead,
			}
			if got := p.G4Alt2(); got != tt.want {
				t.Errorf("G4Alt2() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPacket_G4Alt3(t *testing.T) {
	type fields struct {
		Buf      []byte
		Pos      int
		lastRead readOp
	}
	tests := []struct {
		name   string
		fields fields
		want   uint32
	}{
		{
			name: "valid",
			fields: fields{
				Buf:      []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
				Pos:      0,
				lastRead: 0,
			},
			want: 0x02010403,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Packet{
				Buf:      tt.fields.Buf,
				Pos:      tt.fields.Pos,
				lastRead: tt.fields.lastRead,
			}
			if got := p.G4Alt3(); got != tt.want {
				t.Errorf("G4Alt3() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPacket_G8(t *testing.T) {
	type fields struct {
		Buf      []byte
		Pos      int
		lastRead readOp
	}
	tests := []struct {
		name   string
		fields fields
		want   uint64
	}{
		{
			name: "valid",
			fields: fields{
				Buf:      []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
				Pos:      0,
				lastRead: 0,
			},
			want: 0x0102030405060708,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Packet{
				Buf:      tt.fields.Buf,
				Pos:      tt.fields.Pos,
				lastRead: tt.fields.lastRead,
			}
			if got := p.G8(); got != tt.want {
				t.Errorf("G8() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPacket_GData(t *testing.T) {
	type fields struct {
		Buf      []byte
		Pos      int
		lastRead readOp
	}
	type args struct {
		length int
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
				Buf:      []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
				Pos:      0,
				lastRead: 0,
			},
			args: args{
				length: 12,
			},
			want: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Packet{
				Buf:      tt.fields.Buf,
				Pos:      tt.fields.Pos,
				lastRead: tt.fields.lastRead,
			}
			dest := make([]byte, tt.args.length)
			p.GData(dest, tt.args.length)
			if got := dest; !slices.Equal(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPacket_GDataAlt1(t *testing.T) {
	type fields struct {
		Buf      []byte
		Pos      int
		lastRead readOp
	}
	type args struct {
		length int
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
				Buf:      []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
				Pos:      0,
				lastRead: 0,
			},
			args: args{
				length: 12,
			},
			want: []byte{12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Packet{
				Buf:      tt.fields.Buf,
				Pos:      tt.fields.Pos,
				lastRead: tt.fields.lastRead,
			}
			dest := make([]byte, tt.args.length)
			p.GDataAlt1(dest, tt.args.length)
			if got := dest; !slices.Equal(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPacket_GDataAlt3(t *testing.T) {
	type fields struct {
		Buf      []byte
		Pos      int
		lastRead readOp
	}
	type args struct {
		length int
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
				Buf:      []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
				Pos:      0,
				lastRead: 0,
			},
			args: args{
				length: 12,
			},
			want: []byte{140, 139, 138, 137, 136, 135, 134, 133, 132, 131, 130, 129},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Packet{
				Buf:      tt.fields.Buf,
				Pos:      tt.fields.Pos,
				lastRead: tt.fields.lastRead,
			}
			dest := make([]byte, tt.args.length)
			p.GDataAlt3(dest, tt.args.length)
			if got := dest; !slices.Equal(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPacket_GExtended1or2(t *testing.T) {
	type fields struct {
		Buf      []byte
		Pos      int
		lastRead readOp
	}
	tests := []struct {
		name   string
		fields fields
		want   uint32
	}{
		{
			name: "one byte",
			fields: fields{
				Buf:      []byte{0x10},
				Pos:      0,
				lastRead: 0,
			},
			want: 0x10,
		},
		{
			name: "two bytes",
			fields: fields{
				Buf:      []byte{0x80, 0x81},
				Pos:      0,
				lastRead: 0,
			},
			want: 0x81,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Packet{
				Buf:      tt.fields.Buf,
				Pos:      tt.fields.Pos,
				lastRead: tt.fields.lastRead,
			}
			if got := p.GExtended1or2(); got != tt.want {
				t.Errorf("GExtended1or2() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPacket_GJStr(t *testing.T) {
	type fields struct {
		Buf      []byte
		Pos      int
		lastRead readOp
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "password",
			fields: fields{
				Buf:      []byte{112, 97, 115, 115, 119, 111, 114, 100, 0},
				Pos:      0,
				lastRead: 0,
			},
			want: "password",
		},
		{
			name: "mid-buffer",
			fields: fields{
				Buf:      []byte{10, 0, 3, 0, 0, 1, 219, 154, 95, 17, 108, 1, 155, 179, 69, 112, 97, 115, 115, 119, 111, 114, 100, 0, 0, 99, 29, 123, 0, 0, 1, 0, 3, 33, 131, 170, 7, 178, 0, 225, 0, 0, 0, 0},
				Pos:      15,
				lastRead: -1,
			},
			want: "password",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Packet{
				Buf:      tt.fields.Buf,
				Pos:      tt.fields.Pos,
				lastRead: tt.fields.lastRead,
			}
			if got := p.GJStr(); got != tt.want {
				t.Errorf("GJStr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPacket_GJStr2(t *testing.T) {
	type fields struct {
		Buf      []byte
		Pos      int
		lastRead readOp
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "password",
			fields: fields{
				Buf:      []byte{0, 112, 97, 115, 115, 119, 111, 114, 100, 0},
				Pos:      0,
				lastRead: 0,
			},
			want: "password",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Packet{
				Buf:      tt.fields.Buf,
				Pos:      tt.fields.Pos,
				lastRead: tt.fields.lastRead,
			}
			if got := p.GJStr2(); got != tt.want {
				t.Errorf("GJStr2() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPacket_GSmart(t *testing.T) {
	type fields struct {
		Buf      []byte
		Pos      int
		lastRead readOp
	}
	tests := []struct {
		name   string
		fields fields
		want   uint16
	}{
		{
			name: "64",
			fields: fields{
				Buf:      []byte{64},
				Pos:      0,
				lastRead: 0,
			},
			want: 64,
		},
		{
			name: "128, 202",
			fields: fields{
				Buf:      []byte{0x80, 0xCA},
				Pos:      0,
				lastRead: 0,
			},
			want: 0xCA,
		},
		{
			name: "150, 202",
			fields: fields{
				Buf:      []byte{0x96, 0xCA},
				Pos:      0,
				lastRead: 0,
			},
			want: 0x16CA,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Packet{
				Buf:      tt.fields.Buf,
				Pos:      tt.fields.Pos,
				lastRead: tt.fields.lastRead,
			}
			if got := p.GSmart(); got != tt.want {
				t.Errorf("GSmart() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPacket_GSmartS(t *testing.T) {
	type fields struct {
		Buf      []byte
		Pos      int
		lastRead readOp
	}
	tests := []struct {
		name   string
		fields fields
		want   int32
	}{
		{
			name: "64",
			fields: fields{
				Buf:      []byte{64},
				Pos:      0,
				lastRead: 0,
			},
			want: 0,
		},
		{
			name: "128, 202",
			fields: fields{
				Buf:      []byte{0x80, 0xCA},
				Pos:      0,
				lastRead: 0,
			},
			want: 0xC0CA,
		},
		{
			name: "150, 202",
			fields: fields{
				Buf:      []byte{0x96, 0xCA},
				Pos:      0,
				lastRead: 0,
			},
			want: 0xD6CA,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Packet{
				Buf:      tt.fields.Buf,
				Pos:      tt.fields.Pos,
				lastRead: tt.fields.lastRead,
			}
			if got := p.GSmartS(); got != tt.want {
				t.Errorf("GSmartS() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPacket_GVarInt(t *testing.T) {
	type fields struct {
		Buf      []byte
		Pos      int
		lastRead readOp
	}
	tests := []struct {
		name   string
		fields fields
		want   int32
	}{
		{
			name: "> 0",
			fields: fields{
				Buf:      []byte{1},
				Pos:      0,
				lastRead: 0,
			},
			want: 1,
		},
		{
			name: "0",
			fields: fields{
				Buf:      []byte{0},
				Pos:      0,
				lastRead: 0,
			},
			want: 0,
		},
		{
			name: "< 0 signed",
			fields: fields{
				Buf:      []byte{224, 1},
				Pos:      0,
				lastRead: 0,
			},
			want: 0x3001,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Packet{
				Buf:      tt.fields.Buf,
				Pos:      tt.fields.Pos,
				lastRead: tt.fields.lastRead,
			}
			if got := p.GVarInt(); got != tt.want {
				t.Errorf("GVarInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPacket_GVarLong(t *testing.T) {
	type fields struct {
		Buf      []byte
		Pos      int
		lastRead readOp
	}
	type args struct {
		length int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int64
	}{
		{
			name: "one byte",
			fields: fields{
				Buf:      []byte{1},
				Pos:      0,
				lastRead: 0,
			},
			args: args{
				length: 1,
			},
			want: 1,
		},
		{
			name: "two bytes",
			fields: fields{
				Buf:      []byte{1, 2},
				Pos:      0,
				lastRead: 0,
			},
			args: args{
				length: 2,
			},
			want: 0x0102,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Packet{
				Buf:      tt.fields.Buf,
				Pos:      tt.fields.Pos,
				lastRead: tt.fields.lastRead,
			}
			if got := p.GVarLong(tt.args.length); got != tt.want {
				t.Errorf("GVarLong() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPacket_IP2(t *testing.T) {
	type fields struct {
		Buf      []byte
		Pos      int
		lastRead readOp
	}
	type args struct {
		value uint16
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []byte
	}{
		{
			name:   "1",
			fields: fields{},
			args: args{
				value: 1,
			},
			want: []byte{1, 0},
		},
		{
			name:   "0x1234",
			fields: fields{},
			args: args{
				value: 0x1234,
			},
			want: []byte{0x34, 0x12},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Packet{
				Buf:      tt.fields.Buf,
				Pos:      tt.fields.Pos,
				lastRead: tt.fields.lastRead,
			}
			p.IP2(tt.args.value)
			if got := p.Bytes(); !slices.Equal(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPacket_IP4(t *testing.T) {
	type fields struct {
		Buf      []byte
		Pos      int
		lastRead readOp
	}
	type args struct {
		value uint32
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []byte
	}{
		{
			name:   "0x1234",
			fields: fields{},
			args: args{
				value: 0x1234,
			},
			want: []byte{0x34, 0x12, 0, 0},
		},
		{
			name:   "0x123456",
			fields: fields{},
			args: args{
				value: 0x123456,
			},
			want: []byte{0x56, 0x34, 0x12, 0},
		},
		{
			name:   "0x12345678",
			fields: fields{},
			args: args{
				value: 0x12345678,
			},
			want: []byte{0x78, 0x56, 0x34, 0x12},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Packet{
				Buf:      tt.fields.Buf,
				Pos:      tt.fields.Pos,
				lastRead: tt.fields.lastRead,
			}
			p.IP4(tt.args.value)
			if got := p.Bytes(); !slices.Equal(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPacket_P1(t *testing.T) {
	type fields struct {
		Buf      []byte
		Pos      int
		lastRead readOp
	}
	type args struct {
		value uint8
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []byte
	}{
		{
			name:   "0x00",
			fields: fields{},
			args: args{
				value: 0x00,
			},
			want: []byte{0x00},
		},
		{
			name:   "0x12",
			fields: fields{},
			args: args{
				value: 0x12,
			},
			want: []byte{0x12},
		},
		{
			name:   "0x80",
			fields: fields{},
			args: args{
				value: 0x80,
			},
			want: []byte{0x80},
		},
		{
			name:   "0xFF",
			fields: fields{},
			args: args{
				value: 0xFF,
			},
			want: []byte{0xFF},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Packet{
				Buf:      tt.fields.Buf,
				Pos:      tt.fields.Pos,
				lastRead: tt.fields.lastRead,
			}
			p.P1(tt.args.value)
			if got := p.Bytes(); !slices.Equal(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPacket_P1Multi(t *testing.T) {
	type fields struct {
		Buf      []byte
		Pos      int
		lastRead readOp
	}
	type args struct {
		value uint8
	}
	tests := []struct {
		name   string
		fields fields
		args   []args
		want   []byte
	}{
		{
			name:   "0x00 0x12 0x80 0xFF",
			fields: fields{},
			args: []args{
				{value: 0x00},
				{value: 0x12},
				{value: 0x80},
				{value: 0xFF},
			},
			want: []byte{0x00, 0x12, 0x80, 0xFF},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Packet{
				Buf:      tt.fields.Buf,
				Pos:      tt.fields.Pos,
				lastRead: tt.fields.lastRead,
			}
			for _, a := range tt.args {
				p.P1(a.value)
			}
			if got := p.Bytes(); !slices.Equal(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPacket_P1Alt1(t *testing.T) {
	type fields struct {
		Buf      []byte
		Pos      int
		lastRead readOp
	}
	type args struct {
		value uint8
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []byte
	}{
		{
			name:   "0x00",
			fields: fields{},
			args: args{
				value: 0x00,
			},
			want: []byte{0x80},
		},
		{
			name:   "0x12",
			fields: fields{},
			args: args{
				value: 0x12,
			},
			want: []byte{0x92},
		},
		{
			name:   "0x80",
			fields: fields{},
			args: args{
				value: 0x80,
			},
			want: []byte{0x0},
		},
		{
			name:   "0xFF",
			fields: fields{},
			args: args{
				value: 0xFF,
			},
			want: []byte{0x7F},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Packet{
				Buf:      tt.fields.Buf,
				Pos:      tt.fields.Pos,
				lastRead: tt.fields.lastRead,
			}
			p.P1Alt1(tt.args.value)
			if got := p.Bytes(); !slices.Equal(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPacket_P1Alt2(t *testing.T) {
	type fields struct {
		Buf      []byte
		Pos      int
		lastRead readOp
	}
	type args struct {
		value uint8
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []byte
	}{
		{
			name:   "0x00",
			fields: fields{},
			args: args{
				value: 0x00,
			},
			want: []byte{0x0},
		},
		{
			name:   "0x12",
			fields: fields{},
			args: args{
				value: 0x12,
			},
			want: []byte{0xEE},
		},
		{
			name:   "0x80",
			fields: fields{},
			args: args{
				value: 0x80,
			},
			want: []byte{0x80},
		},
		{
			name:   "0xFF",
			fields: fields{},
			args: args{
				value: 0xFF,
			},
			want: []byte{0x1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Packet{
				Buf:      tt.fields.Buf,
				Pos:      tt.fields.Pos,
				lastRead: tt.fields.lastRead,
			}
			p.P1Alt2(tt.args.value)
			if got := p.Bytes(); !slices.Equal(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPacket_P1Alt3(t *testing.T) {
	type fields struct {
		Buf      []byte
		Pos      int
		lastRead readOp
	}
	type args struct {
		value uint8
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []byte
	}{
		{
			name:   "0x00",
			fields: fields{},
			args: args{
				value: 0x00,
			},
			want: []byte{0x80},
		},
		{
			name:   "0x12",
			fields: fields{},
			args: args{
				value: 0x12,
			},
			want: []byte{0x6E},
		},
		{
			name:   "0x80",
			fields: fields{},
			args: args{
				value: 0x80,
			},
			want: []byte{0x0},
		},
		{
			name:   "0xFF",
			fields: fields{},
			args: args{
				value: 0xFF,
			},
			want: []byte{0x81},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Packet{
				Buf:      tt.fields.Buf,
				Pos:      tt.fields.Pos,
				lastRead: tt.fields.lastRead,
			}
			p.P1Alt3(tt.args.value)
			if got := p.Bytes(); !slices.Equal(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPacket_P2(t *testing.T) {
	type fields struct {
		Buf      []byte
		Pos      int
		lastRead readOp
	}
	type args struct {
		value uint16
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []byte
	}{
		{
			name:   "0x00",
			fields: fields{},
			args: args{
				value: 0x00,
			},
			want: []byte{0x0, 0x0},
		},
		{
			name:   "0x12",
			fields: fields{},
			args: args{
				value: 0x12,
			},
			want: []byte{0x0, 0x12},
		},
		{
			name:   "0x80",
			fields: fields{},
			args: args{
				value: 0x80,
			},
			want: []byte{0x0, 0x80},
		},
		{
			name:   "0xFF",
			fields: fields{},
			args: args{
				value: 0xFF,
			},
			want: []byte{0x0, 0xFF},
		},
		{
			name:   "0x1234",
			fields: fields{},
			args: args{
				value: 0x1234,
			},
			want: []byte{0x12, 0x34},
		},
		{
			name:   "0x8000",
			fields: fields{},
			args: args{
				value: 0x8000,
			},
			want: []byte{0x80, 0x0},
		},
		{
			name:   "0x7FFF",
			fields: fields{},
			args: args{
				value: 0x7FFF,
			},
			want: []byte{0x7F, 0xFF},
		},
		{
			name:   "0xFFFF",
			fields: fields{},
			args: args{
				value: 0xFFFF,
			},
			want: []byte{0xFF, 0xFF},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Packet{
				Buf:      tt.fields.Buf,
				Pos:      tt.fields.Pos,
				lastRead: tt.fields.lastRead,
			}
			p.P2(tt.args.value)
			if got := p.Bytes(); !slices.Equal(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPacket_P2Alt1(t *testing.T) {
	type fields struct {
		Buf      []byte
		Pos      int
		lastRead readOp
	}
	type args struct {
		value uint16
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []byte
	}{
		{
			name:   "0x00",
			fields: fields{},
			args: args{
				value: 0x00,
			},
			want: []byte{0x0, 0x0},
		},
		{
			name:   "0x12",
			fields: fields{},
			args: args{
				value: 0x12,
			},
			want: []byte{0x12, 0x0},
		},
		{
			name:   "0x80",
			fields: fields{},
			args: args{
				value: 0x80,
			},
			want: []byte{0x80, 0x0},
		},
		{
			name:   "0xFF",
			fields: fields{},
			args: args{
				value: 0xFF,
			},
			want: []byte{0xFF, 0x0},
		},
		{
			name:   "0x1234",
			fields: fields{},
			args: args{
				value: 0x1234,
			},
			want: []byte{0x34, 0x12},
		},
		{
			name:   "0x8000",
			fields: fields{},
			args: args{
				value: 0x8000,
			},
			want: []byte{0x0, 0x80},
		},
		{
			name:   "0x7FFF",
			fields: fields{},
			args: args{
				value: 0x7FFF,
			},
			want: []byte{0xFF, 0x7F},
		},
		{
			name:   "0xFFFF",
			fields: fields{},
			args: args{
				value: 0xFFFF,
			},
			want: []byte{0xFF, 0xFF},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Packet{
				Buf:      tt.fields.Buf,
				Pos:      tt.fields.Pos,
				lastRead: tt.fields.lastRead,
			}
			p.P2Alt1(tt.args.value)
			if got := p.Bytes(); !slices.Equal(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPacket_P2Alt2(t *testing.T) {
	type fields struct {
		Buf      []byte
		Pos      int
		lastRead readOp
	}
	type args struct {
		value uint16
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []byte
	}{
		{
			name:   "0x00",
			fields: fields{},
			args: args{
				value: 0x00,
			},
			want: []byte{0x0, 0x80},
		},
		{
			name:   "0x12",
			fields: fields{},
			args: args{
				value: 0x12,
			},
			want: []byte{0x0, 0x92},
		},
		{
			name:   "0x80",
			fields: fields{},
			args: args{
				value: 0x80,
			},
			want: []byte{0x0, 0x0},
		},
		{
			name:   "0xFF",
			fields: fields{},
			args: args{
				value: 0xFF,
			},
			want: []byte{0x0, 0x7F},
		},
		{
			name:   "0x1234",
			fields: fields{},
			args: args{
				value: 0x1234,
			},
			want: []byte{0x12, 0xB4},
		},
		{
			name:   "0x8000",
			fields: fields{},
			args: args{
				value: 0x8000,
			},
			want: []byte{0x80, 0x80},
		},
		{
			name:   "0x7FFF",
			fields: fields{},
			args: args{
				value: 0x7FFF,
			},
			want: []byte{0x7F, 0x7F},
		},
		{
			name:   "0xFFFF",
			fields: fields{},
			args: args{
				value: 0xFFFF,
			},
			want: []byte{0xFF, 0x7F},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Packet{
				Buf:      tt.fields.Buf,
				Pos:      tt.fields.Pos,
				lastRead: tt.fields.lastRead,
			}
			p.P2Alt2(tt.args.value)
			if got := p.Bytes(); !slices.Equal(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPacket_P2Alt3(t *testing.T) {
	type fields struct {
		Buf      []byte
		Pos      int
		lastRead readOp
	}
	type args struct {
		value uint16
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []byte
	}{
		{
			name:   "0x00",
			fields: fields{},
			args: args{
				value: 0x00,
			},
			want: []byte{0x80, 0x0},
		},
		{
			name:   "0x12",
			fields: fields{},
			args: args{
				value: 0x12,
			},
			want: []byte{0x92, 0x0},
		},
		{
			name:   "0x80",
			fields: fields{},
			args: args{
				value: 0x80,
			},
			want: []byte{0x0, 0x0},
		},
		{
			name:   "0xFF",
			fields: fields{},
			args: args{
				value: 0xFF,
			},
			want: []byte{0x7F, 0x0},
		},
		{
			name:   "0x1234",
			fields: fields{},
			args: args{
				value: 0x1234,
			},
			want: []byte{0xB4, 0x12},
		},
		{
			name:   "0x8000",
			fields: fields{},
			args: args{
				value: 0x8000,
			},
			want: []byte{0x80, 0x80},
		},
		{
			name:   "0x7FFF",
			fields: fields{},
			args: args{
				value: 0x7FFF,
			},
			want: []byte{0x7F, 0x7F},
		},
		{
			name:   "0xFFFF",
			fields: fields{},
			args: args{
				value: 0xFFFF,
			},
			want: []byte{0x7F, 0xFF},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Packet{
				Buf:      tt.fields.Buf,
				Pos:      tt.fields.Pos,
				lastRead: tt.fields.lastRead,
			}
			p.P2Alt3(tt.args.value)
			if got := p.Bytes(); !slices.Equal(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPacket_P3(t *testing.T) {
	type fields struct {
		Buf      []byte
		Pos      int
		lastRead readOp
	}
	type args struct {
		value uint32
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []byte
	}{
		{
			name:   "0x00",
			fields: fields{},
			args: args{
				value: 0x00,
			},
			want: []byte{0x0, 0x0, 0x0},
		},
		{
			name:   "0x12",
			fields: fields{},
			args: args{
				value: 0x12,
			},
			want: []byte{0x0, 0x0, 0x12},
		},
		{
			name:   "0x80",
			fields: fields{},
			args: args{
				value: 0x80,
			},
			want: []byte{0x0, 0x0, 0x80},
		},
		{
			name:   "0xFF",
			fields: fields{},
			args: args{
				value: 0xFF,
			},
			want: []byte{0x0, 0x0, 0xFF},
		},
		{
			name:   "0x1234",
			fields: fields{},
			args: args{
				value: 0x1234,
			},
			want: []byte{0x0, 0x12, 0x34},
		},
		{
			name:   "0x8000",
			fields: fields{},
			args: args{
				value: 0x8000,
			},
			want: []byte{0x0, 0x80, 0x0},
		},
		{
			name:   "0x7FFF",
			fields: fields{},
			args: args{
				value: 0x7FFF,
			},
			want: []byte{0x0, 0x7F, 0xFF},
		},
		{
			name:   "0xFFFF",
			fields: fields{},
			args: args{
				value: 0xFFFF,
			},
			want: []byte{0x0, 0xFF, 0xFF},
		},
		{
			name:   "0x123456",
			fields: fields{},
			args: args{
				value: 0x123456,
			},
			want: []byte{0x12, 0x34, 0x56},
		},
		{
			name:   "0x7FFFFF",
			fields: fields{},
			args: args{
				value: 0x7FFFFF,
			},
			want: []byte{0x7F, 0xFF, 0xFF},
		},
		{
			name:   "0xFFFFFF",
			fields: fields{},
			args: args{
				value: 0xFFFFFF,
			},
			want: []byte{0xFF, 0xFF, 0xFF},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Packet{
				Buf:      tt.fields.Buf,
				Pos:      tt.fields.Pos,
				lastRead: tt.fields.lastRead,
			}
			p.P3(tt.args.value)
			if got := p.Bytes(); !slices.Equal(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPacket_P4(t *testing.T) {
	type fields struct {
		Buf      []byte
		Pos      int
		lastRead readOp
	}
	type args struct {
		value uint32
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []byte
	}{
		{
			name:   "0x00",
			fields: fields{},
			args: args{
				value: 0x00,
			},
			want: []byte{0x0, 0x0, 0x0, 0x0},
		},
		{
			name:   "0x12",
			fields: fields{},
			args: args{
				value: 0x12,
			},
			want: []byte{0x0, 0x0, 0x0, 0x12},
		},
		{
			name:   "0x80",
			fields: fields{},
			args: args{
				value: 0x80,
			},
			want: []byte{0x0, 0x0, 0x0, 0x80},
		},
		{
			name:   "0xFF",
			fields: fields{},
			args: args{
				value: 0xFF,
			},
			want: []byte{0x0, 0x0, 0x0, 0xFF},
		},
		{
			name:   "0x1234",
			fields: fields{},
			args: args{
				value: 0x1234,
			},
			want: []byte{0x0, 0x0, 0x12, 0x34},
		},
		{
			name:   "0x8000",
			fields: fields{},
			args: args{
				value: 0x8000,
			},
			want: []byte{0x0, 0x0, 0x80, 0x0},
		},
		{
			name:   "0x7FFF",
			fields: fields{},
			args: args{
				value: 0x7FFF,
			},
			want: []byte{0x0, 0x0, 0x7F, 0xFF},
		},
		{
			name:   "0xFFFF",
			fields: fields{},
			args: args{
				value: 0xFFFF,
			},
			want: []byte{0x0, 0x0, 0xFF, 0xFF},
		},
		{
			name:   "0x123456",
			fields: fields{},
			args: args{
				value: 0x123456,
			},
			want: []byte{0x0, 0x12, 0x34, 0x56},
		},
		{
			name:   "0x7FFFFF",
			fields: fields{},
			args: args{
				value: 0x7FFFFF,
			},
			want: []byte{0x0, 0x7F, 0xFF, 0xFF},
		},
		{
			name:   "0xFFFFFF",
			fields: fields{},
			args: args{
				value: 0xFFFFFF,
			},
			want: []byte{0x0, 0xFF, 0xFF, 0xFF},
		},
		{
			name:   "0x12345678",
			fields: fields{},
			args: args{
				value: 0x12345678,
			},
			want: []byte{0x12, 0x34, 0x56, 0x78},
		},
		{
			name:   "0x7FFFFFFF",
			fields: fields{},
			args: args{
				value: 0x7FFFFFFF,
			},
			want: []byte{0x7F, 0xFF, 0xFF, 0xFF},
		},
		{
			name:   "0xFFFFFFFF",
			fields: fields{},
			args: args{
				value: 0xFFFFFFFF,
			},
			want: []byte{0xFF, 0xFF, 0xFF, 0xFF},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Packet{
				Buf:      tt.fields.Buf,
				Pos:      tt.fields.Pos,
				lastRead: tt.fields.lastRead,
			}
			p.P4(tt.args.value)
			if got := p.Bytes(); !slices.Equal(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPacket_P4Alt1(t *testing.T) {
	type fields struct {
		Buf      []byte
		Pos      int
		lastRead readOp
	}
	type args struct {
		value uint32
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []byte
	}{
		{
			name:   "0x00",
			fields: fields{},
			args: args{
				value: 0x00,
			},
			want: []byte{0x0, 0x0, 0x0, 0x0},
		},
		{
			name:   "0x12",
			fields: fields{},
			args: args{
				value: 0x12,
			},
			want: []byte{0x12, 0x0, 0x0, 0x0},
		},
		{
			name:   "0x80",
			fields: fields{},
			args: args{
				value: 0x80,
			},
			want: []byte{0x80, 0x0, 0x0, 0x0},
		},
		{
			name:   "0xFF",
			fields: fields{},
			args: args{
				value: 0xFF,
			},
			want: []byte{0xFF, 0x0, 0x0, 0x0},
		},
		{
			name:   "0x1234",
			fields: fields{},
			args: args{
				value: 0x1234,
			},
			want: []byte{0x34, 0x12, 0x0, 0x0},
		},
		{
			name:   "0x8000",
			fields: fields{},
			args: args{
				value: 0x8000,
			},
			want: []byte{0x0, 0x80, 0x0, 0x0},
		},
		{
			name:   "0x7FFF",
			fields: fields{},
			args: args{
				value: 0x7FFF,
			},
			want: []byte{0xFF, 0x7F, 0x0, 0x0},
		},
		{
			name:   "0xFFFF",
			fields: fields{},
			args: args{
				value: 0xFFFF,
			},
			want: []byte{0xFF, 0xFF, 0x0, 0x0},
		},
		{
			name:   "0x123456",
			fields: fields{},
			args: args{
				value: 0x123456,
			},
			want: []byte{0x56, 0x34, 0x12, 0x0},
		},
		{
			name:   "0x7FFFFF",
			fields: fields{},
			args: args{
				value: 0x7FFFFF,
			},
			want: []byte{0xFF, 0xFF, 0x7F, 0x0},
		},
		{
			name:   "0xFFFFFF",
			fields: fields{},
			args: args{
				value: 0xFFFFFF,
			},
			want: []byte{0xFF, 0xFF, 0xFF, 0x0},
		},
		{
			name:   "0x12345678",
			fields: fields{},
			args: args{
				value: 0x12345678,
			},
			want: []byte{0x78, 0x56, 0x34, 0x12},
		},
		{
			name:   "0x7FFFFFFF",
			fields: fields{},
			args: args{
				value: 0x7FFFFFFF,
			},
			want: []byte{0xFF, 0xFF, 0xFF, 0x7F},
		},
		{
			name:   "0xFFFFFFFF",
			fields: fields{},
			args: args{
				value: 0xFFFFFFFF,
			},
			want: []byte{0xFF, 0xFF, 0xFF, 0xFF},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Packet{
				Buf:      tt.fields.Buf,
				Pos:      tt.fields.Pos,
				lastRead: tt.fields.lastRead,
			}
			p.P4Alt1(tt.args.value)
			if got := p.Bytes(); !slices.Equal(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPacket_P4Alt2(t *testing.T) {
	type fields struct {
		Buf      []byte
		Pos      int
		lastRead readOp
	}
	type args struct {
		value uint32
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []byte
	}{
		{
			name:   "0x00",
			fields: fields{},
			args: args{
				value: 0x00,
			},
			want: []byte{0x0, 0x0, 0x0, 0x0},
		},
		{
			name:   "0x12",
			fields: fields{},
			args: args{
				value: 0x12,
			},
			want: []byte{0x0, 0x12, 0x0, 0x0},
		},
		{
			name:   "0x80",
			fields: fields{},
			args: args{
				value: 0x80,
			},
			want: []byte{0x0, 0x80, 0x0, 0x0},
		},
		{
			name:   "0xFF",
			fields: fields{},
			args: args{
				value: 0xFF,
			},
			want: []byte{0x0, 0xFF, 0x0, 0x0},
		},
		{
			name:   "0x1234",
			fields: fields{},
			args: args{
				value: 0x1234,
			},
			want: []byte{0x12, 0x34, 0x0, 0x0},
		},
		{
			name:   "0x8000",
			fields: fields{},
			args: args{
				value: 0x8000,
			},
			want: []byte{0x80, 0x0, 0x0, 0x0},
		},
		{
			name:   "0x7FFF",
			fields: fields{},
			args: args{
				value: 0x7FFF,
			},
			want: []byte{0x7F, 0xFF, 0x0, 0x0},
		},
		{
			name:   "0xFFFF",
			fields: fields{},
			args: args{
				value: 0xFFFF,
			},
			want: []byte{0xFF, 0xFF, 0x0, 0x0},
		},
		{
			name:   "0x123456",
			fields: fields{},
			args: args{
				value: 0x123456,
			},
			want: []byte{0x34, 0x56, 0x0, 0x12},
		},
		{
			name:   "0x7FFFFF",
			fields: fields{},
			args: args{
				value: 0x7FFFFF,
			},
			want: []byte{0xFF, 0xFF, 0x0, 0x7F},
		},
		{
			name:   "0xFFFFFF",
			fields: fields{},
			args: args{
				value: 0xFFFFFF,
			},
			want: []byte{0xFF, 0xFF, 0x0, 0xFF},
		},
		{
			name:   "0x12345678",
			fields: fields{},
			args: args{
				value: 0x12345678,
			},
			want: []byte{0x56, 0x78, 0x12, 0x34},
		},
		{
			name:   "0x7FFFFFFF",
			fields: fields{},
			args: args{
				value: 0x7FFFFFFF,
			},
			want: []byte{0xFF, 0xFF, 0x7F, 0xFF},
		},
		{
			name:   "0xFFFFFFFF",
			fields: fields{},
			args: args{
				value: 0xFFFFFFFF,
			},
			want: []byte{0xFF, 0xFF, 0xFF, 0xFF},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Packet{
				Buf:      tt.fields.Buf,
				Pos:      tt.fields.Pos,
				lastRead: tt.fields.lastRead,
			}
			p.P4Alt2(tt.args.value)
			if got := p.Bytes(); !slices.Equal(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPacket_P4Alt3(t *testing.T) {
	type fields struct {
		Buf      []byte
		Pos      int
		lastRead readOp
	}
	type args struct {
		value uint32
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []byte
	}{
		{
			name:   "0x00",
			fields: fields{},
			args: args{
				value: 0x00,
			},
			want: []byte{0x0, 0x0, 0x0, 0x0},
		},
		{
			name:   "0x12",
			fields: fields{},
			args: args{
				value: 0x12,
			},
			want: []byte{0x0, 0x0, 0x12, 0x0},
		},
		{
			name:   "0x80",
			fields: fields{},
			args: args{
				value: 0x80,
			},
			want: []byte{0x0, 0x0, 0x80, 0x0},
		},
		{
			name:   "0xFF",
			fields: fields{},
			args: args{
				value: 0xFF,
			},
			want: []byte{0x0, 0x0, 0xFF, 0x0},
		},
		{
			name:   "0x1234",
			fields: fields{},
			args: args{
				value: 0x1234,
			},
			want: []byte{0x0, 0x0, 0x34, 0x12},
		},
		{
			name:   "0x8000",
			fields: fields{},
			args: args{
				value: 0x8000,
			},
			want: []byte{0x0, 0x0, 0x0, 0x80},
		},
		{
			name:   "0x7FFF",
			fields: fields{},
			args: args{
				value: 0x7FFF,
			},
			want: []byte{0x0, 0x0, 0xFF, 0x7F},
		},
		{
			name:   "0xFFFF",
			fields: fields{},
			args: args{
				value: 0xFFFF,
			},
			want: []byte{0x0, 0x0, 0xFF, 0xFF},
		},
		{
			name:   "0x123456",
			fields: fields{},
			args: args{
				value: 0x123456,
			},
			want: []byte{0x12, 0x0, 0x56, 0x34},
		},
		{
			name:   "0x7FFFFF",
			fields: fields{},
			args: args{
				value: 0x7FFFFF,
			},
			want: []byte{0x7F, 0x0, 0xFF, 0xFF},
		},
		{
			name:   "0xFFFFFF",
			fields: fields{},
			args: args{
				value: 0xFFFFFF,
			},
			want: []byte{0xFF, 0x0, 0xFF, 0xFF},
		},
		{
			name:   "0x12345678",
			fields: fields{},
			args: args{
				value: 0x12345678,
			},
			want: []byte{0x34, 0x12, 0x78, 0x56},
		},
		{
			name:   "0x7FFFFFFF",
			fields: fields{},
			args: args{
				value: 0x7FFFFFFF,
			},
			want: []byte{0xFF, 0x7F, 0xFF, 0xFF},
		},
		{
			name:   "0xFFFFFFFF",
			fields: fields{},
			args: args{
				value: 0xFFFFFFFF,
			},
			want: []byte{0xFF, 0xFF, 0xFF, 0xFF},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Packet{
				Buf:      tt.fields.Buf,
				Pos:      tt.fields.Pos,
				lastRead: tt.fields.lastRead,
			}
			p.P4Alt3(tt.args.value)
			if got := p.Bytes(); !slices.Equal(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPacket_P8(t *testing.T) {
	type fields struct {
		Buf      []byte
		Pos      int
		lastRead readOp
	}
	type args struct {
		value uint64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []byte
	}{
		{
			name:   "0x00",
			fields: fields{},
			args: args{
				value: 0x00,
			},
			want: []byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
		},
		{
			name:   "0x12",
			fields: fields{},
			args: args{
				value: 0x12,
			},
			want: []byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x12},
		},
		{
			name:   "0x80",
			fields: fields{},
			args: args{
				value: 0x80,
			},
			want: []byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x80},
		},
		{
			name:   "0xFF",
			fields: fields{},
			args: args{
				value: 0xFF,
			},
			want: []byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xFF},
		},
		{
			name:   "0x1234",
			fields: fields{},
			args: args{
				value: 0x1234,
			},
			want: []byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x12, 0x34},
		},
		{
			name:   "0x8000",
			fields: fields{},
			args: args{
				value: 0x8000,
			},
			want: []byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x80, 0x0},
		},
		{
			name:   "0x7FFF",
			fields: fields{},
			args: args{
				value: 0x7FFF,
			},
			want: []byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x7F, 0xFF},
		},
		{
			name:   "0xFFFF",
			fields: fields{},
			args: args{
				value: 0xFFFF,
			},
			want: []byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xFF, 0xFF},
		},
		{
			name:   "0x123456",
			fields: fields{},
			args: args{
				value: 0x123456,
			},
			want: []byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x12, 0x34, 0x56},
		},
		{
			name:   "0x7FFFFF",
			fields: fields{},
			args: args{
				value: 0x7FFFFF,
			},
			want: []byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x7F, 0xFF, 0xFF},
		},
		{
			name:   "0xFFFFFF",
			fields: fields{},
			args: args{
				value: 0xFFFFFF,
			},
			want: []byte{0x0, 0x0, 0x0, 0x0, 0x0, 0xFF, 0xFF, 0xFF},
		},
		{
			name:   "0x12345678",
			fields: fields{},
			args: args{
				value: 0x12345678,
			},
			want: []byte{0x0, 0x0, 0x0, 0x0, 0x12, 0x34, 0x56, 0x78},
		},
		{
			name:   "0x7FFFFFFF",
			fields: fields{},
			args: args{
				value: 0x7FFFFFFF,
			},
			want: []byte{0x0, 0x0, 0x0, 0x0, 0x7F, 0xFF, 0xFF, 0xFF},
		},
		{
			name:   "0xFFFFFFFF",
			fields: fields{},
			args: args{
				value: 0xFFFFFFFF,
			},
			want: []byte{0x0, 0x0, 0x0, 0x0, 0xFF, 0xFF, 0xFF, 0xFF},
		},
		{
			name:   "0x123456789A",
			fields: fields{},
			args: args{
				value: 0x123456789A,
			},
			want: []byte{0x0, 0x0, 0x0, 0x12, 0x34, 0x56, 0x78, 0x9A},
		},
		{
			name:   "0x7FFFFFFFFF",
			fields: fields{},
			args: args{
				value: 0x7FFFFFFFFF,
			},
			want: []byte{0x0, 0x0, 0x0, 0x7F, 0xFF, 0xFF, 0xFF, 0xFF},
		},
		{
			name:   "0xFFFFFFFFFF",
			fields: fields{},
			args: args{
				value: 0xFFFFFFFFFF,
			},
			want: []byte{0x0, 0x0, 0x0, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Packet{
				Buf:      tt.fields.Buf,
				Pos:      tt.fields.Pos,
				lastRead: tt.fields.lastRead,
			}
			p.P8(tt.args.value)
			if got := p.Bytes(); !slices.Equal(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPacket_PData(t *testing.T) {
	type fields struct {
		Buf      []byte
		Pos      int
		lastRead readOp
	}
	type args struct {
		src    []byte
		length int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []byte
	}{
		{
			name: "buffer with existing data",
			fields: fields{
				Buf: []byte{1, 89},
			},
			args: args{
				src:    []byte{0, 33, 0, 3, 0, 116, 115, 255},
				length: 8,
			},
			want: []byte{1, 89, 0, 33, 0, 3, 0, 116, 115, 255},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Packet{
				Buf:      tt.fields.Buf,
				Pos:      tt.fields.Pos,
				lastRead: tt.fields.lastRead,
			}
			p.PData(tt.args.src, tt.args.length)
			if got := p.Bytes(); !slices.Equal(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPacket_PJStr(t *testing.T) {
	type fields struct {
		Buf      []byte
		Pos      int
		lastRead readOp
	}
	type args struct {
		str string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []byte
	}{
		{
			name:   "valid",
			fields: fields{},
			args: args{
				str: "Username",
			},
			want: []byte{85, 115, 101, 114, 110, 97, 109, 101, 0},
		},
		{
			name: "valid with non-empty buffer",
			fields: fields{
				Buf: []byte{1, 2, 3},
			},
			args: args{
				str: "Username",
			},
			want: []byte{1, 2, 3, 85, 115, 101, 114, 110, 97, 109, 101, 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Packet{
				Buf:      tt.fields.Buf,
				Pos:      tt.fields.Pos,
				lastRead: tt.fields.lastRead,
			}
			p.PJStr(tt.args.str)
			if got := p.Bytes(); !slices.Equal(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPacket_PSize1(t *testing.T) {
	type fields struct {
		Buf      []byte
		Pos      int
		lastRead readOp
	}
	type args struct {
		length int
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
				Buf: []byte{1, 2, 0, 1, 2, 3, 4, 5, 6, 7, 8},
			},
			args: args{
				length: 8,
			},
			want: []byte{1, 2, 8, 1, 2, 3, 4, 5, 6, 7, 8},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Packet{
				Buf:      tt.fields.Buf,
				Pos:      tt.fields.Pos,
				lastRead: tt.fields.lastRead,
			}
			p.PSize1(tt.args.length)
			if got := p.Bytes(); !slices.Equal(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPacket_PSize2(t *testing.T) {
	type fields struct {
		Buf      []byte
		Pos      int
		lastRead readOp
	}
	type args struct {
		length int
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
				Buf: []byte{1, 2, 0, 0, 1, 2, 3, 4, 5, 6, 7, 8},
			},
			args: args{
				length: 8,
			},
			want: []byte{1, 2, 0, 8, 1, 2, 3, 4, 5, 6, 7, 8},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Packet{
				Buf:      tt.fields.Buf,
				Pos:      tt.fields.Pos,
				lastRead: tt.fields.lastRead,
			}
			p.PSize2(tt.args.length)
			if got := p.Bytes(); !slices.Equal(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPacket_PSize4(t *testing.T) {
	type fields struct {
		Buf      []byte
		Pos      int
		lastRead readOp
	}
	type args struct {
		length int
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
				Buf: []byte{1, 2, 0, 0, 0, 0, 1, 2, 3, 4, 5, 6, 7, 8},
			},
			args: args{
				length: 8,
			},
			want: []byte{1, 2, 0, 0, 0, 8, 1, 2, 3, 4, 5, 6, 7, 8},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Packet{
				Buf:      tt.fields.Buf,
				Pos:      tt.fields.Pos,
				lastRead: tt.fields.lastRead,
			}
			p.PSize4(tt.args.length)
			if got := p.Bytes(); !slices.Equal(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPacket_PSmart(t *testing.T) {
	type fields struct {
		Buf      []byte
		Pos      int
		lastRead readOp
	}
	type args struct {
		value uint16
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []byte
	}{
		{
			name:   "value < 0x80",
			fields: fields{},
			args: args{
				value: 0x14,
			},
			want: []byte{0x14},
		},
		{
			name:   "value >= 0x80",
			fields: fields{},
			args: args{
				value: 0x98,
			},
			want: []byte{0x80, 0x98},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Packet{
				Buf:      tt.fields.Buf,
				Pos:      tt.fields.Pos,
				lastRead: tt.fields.lastRead,
			}
			p.PSmart(tt.args.value)
			if got := p.Bytes(); !slices.Equal(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPacket_PVarInt(t *testing.T) {
	type fields struct {
		Buf      []byte
		Pos      int
		lastRead readOp
	}
	type args struct {
		value uint32
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []byte
	}{
		{
			name:   "4 bits",
			fields: fields{},
			args: args{
				value: 0xF,
			},
			want: []byte{0xF},
		},
		{
			name:   "8 bits",
			fields: fields{},
			args: args{
				value: 0xFF,
			},
			want: []byte{0x81, 0x7F},
		},
		{
			name:   "16 bits",
			fields: fields{},
			args: args{
				value: 0xFFFF,
			},
			want: []byte{0x83, 0xFF, 0x7F},
		},
		{
			name:   "24 bits",
			fields: fields{},
			args: args{
				value: 0xFFFFFF,
			},
			want: []byte{0x87, 0xFF, 0xFF, 0x7F},
		},
		{
			name:   "32 bits",
			fields: fields{},
			args: args{
				value: 0xFFFFFFFF,
			},
			want: []byte{0x8F, 0xFF, 0xFF, 0xFF, 0x7F},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Packet{
				Buf:      tt.fields.Buf,
				Pos:      tt.fields.Pos,
				lastRead: tt.fields.lastRead,
			}
			p.PVarInt(tt.args.value)
			if got := p.Bytes(); !slices.Equal(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPacket_PVarLong(t *testing.T) {
	type fields struct {
		Buf      []byte
		Pos      int
		lastRead readOp
	}
	type args struct {
		length int
		value  int64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []byte
	}{
		{
			name:   "64 bits, full length",
			fields: fields{},
			args: args{
				length: 8,
				value:  0x1234567890ABCDEF,
			},
			want: []byte{0x12, 0x34, 0x56, 0x78, 0x90, 0xAB, 0xCD, 0xEF},
		},
		{
			name:   "64 bits, partial length",
			fields: fields{},
			args: args{
				length: 5,
				value:  0x1234567890ABCDEF,
			},
			want: []byte{0x78, 0x90, 0xAB, 0xCD, 0xEF},
		},
		{
			name: "64 bits, full length, non-empty buffer",
			fields: fields{
				Buf: []byte{0x1, 0x2, 0x3, 0x4},
			},
			args: args{
				length: 8,
				value:  0x1234567890ABCDEF,
			},
			want: []byte{0x1, 0x2, 0x3, 0x4, 0x12, 0x34, 0x56, 0x78, 0x90, 0xAB, 0xCD, 0xEF},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Packet{
				Buf:      tt.fields.Buf,
				Pos:      tt.fields.Pos,
				lastRead: tt.fields.lastRead,
			}
			p.PVarLong(tt.args.length, tt.args.value)
			if got := p.Bytes(); !slices.Equal(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPacket_RSAEnc(t *testing.T) {
	type fields struct {
		Buf      []byte
		Pos      int
		lastRead readOp
	}
	type args struct {
		modulus  *big.Int
		exponent *big.Int
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
				Buf: []byte{1, 2, 3, 4, 5, 6, 7, 8},
			},
			args: args{
				modulus:  new(big.Int).SetUint64(0x1234567890ABCDEF),
				exponent: new(big.Int).SetUint64(0x10001),
			},
			want: []byte{0x8, 0x11, 0xFB, 0xAC, 0x86, 0x54, 0x8B, 0x8, 0x83},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Packet{
				Buf:      tt.fields.Buf,
				Pos:      tt.fields.Pos,
				lastRead: tt.fields.lastRead,
			}
			//p.Write([]byte{1, 2, 3, 4, 5, 6, 7, 8})
			p.RSAEnc(tt.args.modulus, tt.args.exponent)
			if got := p.Bytes(); !slices.Equal(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPacket_TinyDec(t *testing.T) {
	type fields struct {
		Buf      []byte
		Pos      int
		lastRead readOp
	}
	type args struct {
		offset int
		key    []uint32
		length int
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
				Buf: []byte{9, 50, 31, 124, 116, 116, 204, 119},
			},
			args: args{
				offset: 0,
				key:    []uint32{50732917, 70165515, 4240440, 0},
				length: 8,
			},
			want: []byte{0, 210, 95, 90, 253, 100, 140, 215},
		},
		{
			name: "two blocks",
			fields: fields{
				Buf: []byte{0xEF, 0xA7, 0x56, 0xB, 0x82, 0x50, 0xDD, 0x59, 0x53, 0x2D, 0x1F, 0xC, 0xA, 0xC4, 0xCD, 0x29},
			},
			args: args{
				offset: 0,
				key:    []uint32{50732917, 70165515, 4240440, 0},
				length: 16,
			},
			want: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Packet{
				Buf:      tt.fields.Buf,
				Pos:      tt.fields.Pos,
				lastRead: tt.fields.lastRead,
			}
			p.TinyDec(tt.args.offset, tt.args.key, tt.args.length)
			if got := p.Bytes(); !slices.Equal(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPacket_TinyEnc(t *testing.T) {
	type fields struct {
		Buf      []byte
		Pos      int
		lastRead readOp
	}
	type args struct {
		key []uint32
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []byte
	}{
		{
			name: "one block",
			fields: fields{
				Buf: []byte{0, 210, 95, 90, 253, 100, 140, 215},
			},
			args: args{
				key: []uint32{50732917, 70165515, 4240440, 0},
			},
			want: []byte{9, 50, 31, 124, 116, 116, 204, 119},
		},
		{
			name: "two blocks",
			fields: fields{
				Buf: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16},
			},
			args: args{
				key: []uint32{50732917, 70165515, 4240440, 0},
			},
			want: []byte{0xEF, 0xA7, 0x56, 0xB, 0x82, 0x50, 0xDD, 0x59, 0x53, 0x2D, 0x1F, 0xC, 0xA, 0xC4, 0xCD, 0x29},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Packet{
				Buf:      tt.fields.Buf,
				Pos:      tt.fields.Pos,
				lastRead: tt.fields.lastRead,
			}
			p.TinyEnc(tt.args.key)
			if got := p.Bytes(); !slices.Equal(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}
