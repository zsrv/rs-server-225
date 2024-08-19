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
			if got := GetCRC(tt.args.src, tt.args.offset, tt.args.length); got != tt.want {
				t.Errorf("GetCRC() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCheckCRC(t *testing.T) {
	type args struct {
		src      []uint8
		offset   int
		length   int
		expected uint32
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckCRC(tt.args.src, tt.args.offset, tt.args.length, tt.args.expected); got != tt.want {
				t.Errorf("CheckCRC() = %v, want %v", got, tt.want)
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

func TestPacket_IG2(t *testing.T) {
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Packet{
				Buf:      tt.fields.Buf,
				Pos:      tt.fields.Pos,
				lastRead: tt.fields.lastRead,
			}
			if got := p.IG2(); got != tt.want {
				t.Errorf("IG2() = %v, want %v", got, tt.want)
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

func TestPacket_IG4(t *testing.T) {
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
			if got := p.IG4(); got != tt.want {
				t.Errorf("IG4() = %v, want %v", got, tt.want)
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

func TestPacket_GBool(t *testing.T) {
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Packet{
				Buf:      tt.fields.Buf,
				Pos:      tt.fields.Pos,
				lastRead: tt.fields.lastRead,
			}
			if got := p.GBool(); got != tt.want {
				t.Errorf("GBool() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TODO: fix
func TestPacket_GJStr(t *testing.T) {
	type fields struct {
		Buf      []byte
		Pos      int
		lastRead readOp
	}
	type args struct {
		terminator byte
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "password",
			fields: fields{
				Buf:      []byte{112, 97, 115, 115, 119, 111, 114, 100, 0},
				Pos:      0,
				lastRead: 0,
			},
			args: args{
				terminator: 0,
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
			if got := p.GJStr(tt.args.terminator); got != tt.want {
				t.Errorf("GJStr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPacket_GJStrLF(t *testing.T) {
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Packet{
				Buf:      tt.fields.Buf,
				Pos:      tt.fields.Pos,
				lastRead: tt.fields.lastRead,
			}
			if got := p.GJStrLF(); got != tt.want {
				t.Errorf("GJStrLF() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPacket_GJStrNUL(t *testing.T) {
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Packet{
				Buf:      tt.fields.Buf,
				Pos:      tt.fields.Pos,
				lastRead: tt.fields.lastRead,
			}
			if got := p.GJStrNUL(); got != tt.want {
				t.Errorf("GJStrNUL() = %v, want %v", got, tt.want)
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

func TestPacket_PBool(t *testing.T) {
	type fields struct {
		Buf      []byte
		Pos      int
		lastRead readOp
	}
	type args struct {
		value bool
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
			p := &Packet{
				Buf:      tt.fields.Buf,
				Pos:      tt.fields.Pos,
				lastRead: tt.fields.lastRead,
			}
			p.PBool(tt.args.value)
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
		str        string
		terminator byte
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []byte
	}{
		{
			name:   "valid with NUL terminator",
			fields: fields{},
			args: args{
				str:        "Username",
				terminator: 0,
			},
			want: []byte{85, 115, 101, 114, 110, 97, 109, 101, 0},
		},
		{
			name: "valid with non-empty buffer and NUL terminator",
			fields: fields{
				Buf: []byte{1, 2, 3},
			},
			args: args{
				str:        "Username",
				terminator: 0,
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
			p.PJStr(tt.args.str, tt.args.terminator)
			if got := p.Bytes(); !slices.Equal(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPacket_PJStrLF(t *testing.T) {
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
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Packet{
				Buf:      tt.fields.Buf,
				Pos:      tt.fields.Pos,
				lastRead: tt.fields.lastRead,
			}
			p.PJStrLF(tt.args.str)
		})
	}
}

func TestPacket_PJStrNUL(t *testing.T) {
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
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Packet{
				Buf:      tt.fields.Buf,
				Pos:      tt.fields.Pos,
				lastRead: tt.fields.lastRead,
			}
			p.PJStrNUL(tt.args.str)
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
		value int32
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

func TestPacket_PSmartS(t *testing.T) {
	type fields struct {
		Buf      []byte
		Pos      int
		lastRead readOp
	}
	type args struct {
		value int32
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
			p := &Packet{
				Buf:      tt.fields.Buf,
				Pos:      tt.fields.Pos,
				lastRead: tt.fields.lastRead,
			}
			p.PSmartS(tt.args.value)
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
