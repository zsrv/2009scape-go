package util

const Size = 256
const GoldenRatio uint32 = 0x9E3779B9

type IsaacRandom struct {
	count int
	a     uint32
	b     uint32
	c     uint32
	mem   [Size]uint32
	rsl   [Size]uint32
}

func NewIsaacRandom(seed []uint32) *IsaacRandom {
	isaac := &IsaacRandom{}

	for i := 0; i < len(seed); i++ {
		isaac.rsl[i] = seed[i]
	}

	isaac.init()

	return isaac
}

func (is *IsaacRandom) isaac() {
	// TODO: unsigned right shifts
	is.c += 1
	is.b += is.c

	for i := 0; i < Size; i++ {
		x := is.mem[i]

		switch i & 3 {
		case 0:
			is.a ^= is.a << 13
		case 1:
			is.a ^= is.a >> 6
		case 2:
			is.a ^= is.a << 2
		case 3:
			is.a ^= is.a >> 16
		}

		is.a += is.mem[(i+128)&0xFF]

		y := is.mem[(x>>2)&0xFF] + is.a + is.b
		is.mem[i] = y
		is.b = is.mem[(y>>10)&0xFF] + x
		is.rsl[i] = is.b
	}
}

func (is *IsaacRandom) init() {
	// TODO: unsigned right shifts
	a := GoldenRatio
	b := GoldenRatio
	c := GoldenRatio
	d := GoldenRatio
	e := GoldenRatio
	f := GoldenRatio
	g := GoldenRatio
	h := GoldenRatio

	for i := 0; i < 4; i++ {
		a ^= b << 11
		d += a
		b += c

		b ^= c >> 2
		e += b
		c += d

		c ^= d << 8
		f += c
		d += e

		d ^= e >> 16
		g += d
		e += f

		e ^= f << 10
		h += e
		f += g

		f ^= g >> 4
		a += f
		g += h

		g ^= h << 8
		b += g
		h += a

		h ^= a >> 9
		c += h
		a += b
	}

	for i := 0; i < Size; i += 8 {
		a += is.rsl[i]
		b += is.rsl[i+1]
		c += is.rsl[i+2]
		d += is.rsl[i+3]
		e += is.rsl[i+4]
		f += is.rsl[i+5]
		g += is.rsl[i+6]
		h += is.rsl[i+7]

		a ^= b << 11
		d += a
		b += c

		b ^= c >> 2
		e += b
		c += d

		c ^= d << 8
		f += c
		d += e

		d ^= e >> 16
		g += d
		e += f

		e ^= f << 10
		h += e
		f += g

		f ^= g >> 4
		a += f
		g += h

		g ^= h << 8
		b += g
		h += a

		h ^= a >> 9
		c += h
		a += b

		is.mem[i] = a
		is.mem[i+1] = b
		is.mem[i+2] = c
		is.mem[i+3] = d
		is.mem[i+4] = e
		is.mem[i+5] = f
		is.mem[i+6] = g
		is.mem[i+7] = h
	}

	for i := 0; i < Size; i += 8 {
		a += is.mem[i]
		b += is.mem[i+1]
		c += is.mem[i+2]
		d += is.mem[i+3]
		e += is.mem[i+4]
		f += is.mem[i+5]
		g += is.mem[i+6]
		h += is.mem[i+7]

		a ^= b << 11
		d += a
		b += c

		b ^= c >> 2
		e += b
		c += d

		c ^= d << 8
		f += c
		d += e

		d ^= e >> 16
		g += d
		e += f

		e ^= f << 10
		h += e
		f += g

		f ^= g >> 4
		a += f
		g += h

		g ^= h << 8
		b += g
		h += a

		h ^= a >> 9
		c += h
		a += b

		is.mem[i] = a
		is.mem[i+1] = b
		is.mem[i+2] = c
		is.mem[i+3] = d
		is.mem[i+4] = e
		is.mem[i+5] = f
		is.mem[i+6] = g
		is.mem[i+7] = h
	}

	is.isaac()
	is.count = Size
}

func (is *IsaacRandom) NextInt() uint32 {
	if is.count == 0 {
		is.count -= 1
		is.isaac()
		is.count = Size - 1
	} else {
		is.count -= 1
	}

	if len(is.rsl)-1 >= is.count {
		return is.rsl[is.count]
	} else {
		return 0
	}
}
