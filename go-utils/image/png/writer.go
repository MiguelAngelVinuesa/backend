// Code taken from standard image/png library to add support for transparent grayscale and 1/2/4-bit grayscale images.

package png

import (
	"bufio"
	"compress/zlib"
	"encoding/binary"
	"hash/crc32"
	"image"
	"image/color"
	"io"
	"math"
	"strconv"
)

type Encoder struct {
	CompressionLevel CompressionLevel
	BufferPool       EncoderBufferPool
}

type EncoderBufferPool interface {
	Get() *EncoderBuffer
	Put(*EncoderBuffer)
}

type EncoderBuffer encoder

type encoder struct {
	enc     *Encoder
	w       io.Writer
	m       image.Image
	cb      int
	err     error
	header  [8]byte
	footer  [4]byte
	tmp     [4 * 256]byte
	cr      [nFilter][]uint8
	pr      []uint8
	zw      *zlib.Writer
	zwLevel int
	bw      *bufio.Writer
}

type CompressionLevel int

const (
	DefaultCompression CompressionLevel = 0
	NoCompression      CompressionLevel = -1
	BestSpeed          CompressionLevel = -2
	BestCompression    CompressionLevel = -3
)

type opaquer interface {
	Opaque() bool
}

func opaque(m image.Image) bool {
	if o, ok := m.(opaquer); ok {
		return o.Opaque()
	}
	b := m.Bounds()
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			_, _, _, a := m.At(x, y).RGBA()
			if a != 0xffff {
				return false
			}
		}
	}
	return true
}

func abs8(d uint8) int {
	if d < 128 {
		return int(d)
	}
	return 256 - int(d)
}

func (e *encoder) writeChunk(b []byte, name string) {
	if e.err != nil {
		return
	}
	n := uint32(len(b))
	if int(n) != len(b) {
		e.err = UnsupportedError(name + " chunk is too large: " + strconv.Itoa(len(b)))
		return
	}
	binary.BigEndian.PutUint32(e.header[:4], n)

	e.header[4] = name[0]
	e.header[5] = name[1]
	e.header[6] = name[2]
	e.header[7] = name[3]

	crc := crc32.NewIEEE()
	crc.Write(e.header[4:8])
	crc.Write(b)
	binary.BigEndian.PutUint32(e.footer[:4], crc.Sum32())

	_, e.err = e.w.Write(e.header[:8])
	if e.err != nil {
		return
	}
	_, e.err = e.w.Write(b)
	if e.err != nil {
		return
	}
	_, e.err = e.w.Write(e.footer[:4])
}

func (e *encoder) writeIHDR() {
	b := e.m.Bounds()
	binary.BigEndian.PutUint32(e.tmp[0:4], uint32(b.Dx()))
	binary.BigEndian.PutUint32(e.tmp[4:8], uint32(b.Dy()))

	switch e.cb {
	case cbG1:
		e.tmp[8] = 1
		e.tmp[9] = ctGrayscale
	case cbG2:
		e.tmp[8] = 2
		e.tmp[9] = ctGrayscale
	case cbG4:
		e.tmp[8] = 4
		e.tmp[9] = ctGrayscale
	case cbG8:
		e.tmp[8] = 8
		e.tmp[9] = ctGrayscale
	case cbGA8:
		e.tmp[8] = 8
		e.tmp[9] = ctGrayscaleAlpha
	case cbTC8:
		e.tmp[8] = 8
		e.tmp[9] = ctTrueColor
	case cbP8:
		e.tmp[8] = 8
		e.tmp[9] = ctPaletted
	case cbP4:
		e.tmp[8] = 4
		e.tmp[9] = ctPaletted
	case cbP2:
		e.tmp[8] = 2
		e.tmp[9] = ctPaletted
	case cbP1:
		e.tmp[8] = 1
		e.tmp[9] = ctPaletted
	case cbTCA8:
		e.tmp[8] = 8
		e.tmp[9] = ctTrueColorAlpha
	case cbG16:
		e.tmp[8] = 16
		e.tmp[9] = ctGrayscale
	case cbGA16:
		e.tmp[8] = 16
		e.tmp[9] = ctGrayscaleAlpha
	case cbTC16:
		e.tmp[8] = 16
		e.tmp[9] = ctTrueColor
	case cbTCA16:
		e.tmp[8] = 16
		e.tmp[9] = ctTrueColorAlpha
	}

	e.tmp[10] = 0
	e.tmp[11] = 0
	e.tmp[12] = 0
	e.writeChunk(e.tmp[:13], "IHDR")
}

func (e *encoder) writePLTEAndTRNS(p color.Palette) {
	if len(p) < 1 || len(p) > 256 {
		e.err = FormatError("bad palette length: " + strconv.Itoa(len(p)))
		return
	}

	last := -1
	for i, c := range p {
		c1 := color.NRGBAModel.Convert(c).(color.NRGBA)
		e.tmp[3*i+0] = c1.R
		e.tmp[3*i+1] = c1.G
		e.tmp[3*i+2] = c1.B
		if c1.A != 0xff {
			last = i
		}
		e.tmp[3*256+i] = c1.A
	}

	e.writeChunk(e.tmp[:3*len(p)], "PLTE")
	if last != -1 {
		e.writeChunk(e.tmp[3*256:3*256+1+last], "tRNS")
	}
}

func (e *encoder) Write(b []byte) (int, error) {
	e.writeChunk(b, "IDAT")
	if e.err != nil {
		return 0, e.err
	}
	return len(b), nil
}

func filter(cr *[nFilter][]byte, pr []byte, bpp int) int {
	cdat0 := cr[0][1:]
	cdat1 := cr[1][1:]
	cdat2 := cr[2][1:]
	cdat3 := cr[3][1:]
	cdat4 := cr[4][1:]
	pdat := pr[1:]
	n := len(cdat0)

	sum := 0
	for i := 0; i < n; i++ {
		cdat2[i] = cdat0[i] - pdat[i]
		sum += abs8(cdat2[i])
	}
	best := sum
	filter := ftUp

	sum = 0
	for i := 0; i < bpp; i++ {
		cdat4[i] = cdat0[i] - pdat[i]
		sum += abs8(cdat4[i])
	}
	for i := bpp; i < n; i++ {
		cdat4[i] = cdat0[i] - paeth(cdat0[i-bpp], pdat[i], pdat[i-bpp])
		sum += abs8(cdat4[i])
		if sum >= best {
			break
		}
	}
	if sum < best {
		best = sum
		filter = ftPaeth
	}

	sum = 0
	for i := 0; i < n; i++ {
		sum += abs8(cdat0[i])
		if sum >= best {
			break
		}
	}
	if sum < best {
		best = sum
		filter = ftNone
	}

	sum = 0
	for i := 0; i < bpp; i++ {
		cdat1[i] = cdat0[i]
		sum += abs8(cdat1[i])
	}
	for i := bpp; i < n; i++ {
		cdat1[i] = cdat0[i] - cdat0[i-bpp]
		sum += abs8(cdat1[i])
		if sum >= best {
			break
		}
	}
	if sum < best {
		best = sum
		filter = ftSub
	}

	sum = 0
	for i := 0; i < bpp; i++ {
		cdat3[i] = cdat0[i] - pdat[i]/2
		sum += abs8(cdat3[i])
	}
	for i := bpp; i < n; i++ {
		cdat3[i] = cdat0[i] - uint8((int(cdat0[i-bpp])+int(pdat[i]))/2)
		sum += abs8(cdat3[i])
		if sum >= best {
			break
		}
	}
	if sum < best {
		filter = ftAverage
	}

	return filter
}

func zeroMemory(v []uint8) {
	for i := range v {
		v[i] = 0
	}
}

func (e *encoder) writeImage(w io.Writer, m image.Image, cb int, level int) error {
	if e.zw == nil || e.zwLevel != level {
		zw, err := zlib.NewWriterLevel(w, level)
		if err != nil {
			return err
		}
		e.zw = zw
		e.zwLevel = level
	} else {
		e.zw.Reset(w)
	}
	defer e.zw.Close()

	bitsPerPixel := 0

	switch cb {
	case cbG8:
		bitsPerPixel = 8
	case cbG4:
		bitsPerPixel = 4
	case cbG2:
		bitsPerPixel = 2
	case cbG1:
		bitsPerPixel = 1
	case cbGA8:
		bitsPerPixel = 16
	case cbTC8:
		bitsPerPixel = 24
	case cbP8:
		bitsPerPixel = 8
	case cbP4:
		bitsPerPixel = 4
	case cbP2:
		bitsPerPixel = 2
	case cbP1:
		bitsPerPixel = 1
	case cbTCA8:
		bitsPerPixel = 32
	case cbTC16:
		bitsPerPixel = 48
	case cbTCA16:
		bitsPerPixel = 64
	case cbG16:
		bitsPerPixel = 16
	case cbGA16:
		bitsPerPixel = 32
	}

	b := m.Bounds()
	sz := 1 + (bitsPerPixel*b.Dx()+7)/8
	for i := range e.cr {
		if cap(e.cr[i]) < sz {
			e.cr[i] = make([]uint8, sz)
		} else {
			e.cr[i] = e.cr[i][:sz]
		}
		e.cr[i][0] = uint8(i)
	}
	cr := e.cr
	if cap(e.pr) < sz {
		e.pr = make([]uint8, sz)
	} else {
		e.pr = e.pr[:sz]
		zeroMemory(e.pr)
	}
	pr := e.pr

	gray, _ := m.(*image.Gray)
	rgba, _ := m.(*image.RGBA)
	paletted, _ := m.(*image.Paletted)
	nrgba, _ := m.(*image.NRGBA)

	for y := b.Min.Y; y < b.Max.Y; y++ {
		i := 1
		switch cb {
		case cbG8:
			if gray != nil {
				offset := (y - b.Min.Y) * gray.Stride
				copy(cr[0][1:], gray.Pix[offset:offset+b.Dx()])
			} else {
				for x := b.Min.X; x < b.Max.X; x++ {
					c := color.GrayModel.Convert(m.At(x, y)).(color.Gray)
					cr[0][i] = c.Y
					i++
				}
			}

		case cbG4, cbG2, cbG1:
			var a uint8
			var c int
			var s uint

			switch cb {
			case cbG4:
				s = 4
			case cbG2:
				s = 6
			case cbG1:
				s = 7
			}
			pixelsPerByte := 8 / bitsPerPixel
			for x := b.Min.X; x < b.Max.X; x++ {
				g := color.GrayModel.Convert(m.At(x, y)).(color.Gray)
				a = a<<uint(bitsPerPixel) | g.Y>>s
				c++
				if c == pixelsPerByte {
					cr[0][i] = a
					i++
					a = 0
					c = 0
				}
			}
			if c != 0 {
				for c != pixelsPerByte {
					a = a << uint(bitsPerPixel)
					c++
				}
				cr[0][i] = a
			}

		case cbGA8:
			for x := b.Min.X; x < b.Max.X; x++ {
				c := color.NRGBAModel.Convert(m.At(x, y)).(color.NRGBA)
				c2 := color.GrayModel.Convert(color.RGBA{R: c.R, G: c.G, B: c.B, A: math.MaxUint8}).(color.Gray)
				cr[0][i] = c2.Y
				cr[0][i+1] = c.A
				i += 2
			}

		case cbTC8:
			cr0 := cr[0]
			stride, pix := 0, []byte(nil)
			if rgba != nil {
				stride, pix = rgba.Stride, rgba.Pix
			} else if nrgba != nil {
				stride, pix = nrgba.Stride, nrgba.Pix
			}
			if stride != 0 {
				j0 := (y - b.Min.Y) * stride
				j1 := j0 + b.Dx()*4
				for j := j0; j < j1; j += 4 {
					cr0[i+0] = pix[j+0]
					cr0[i+1] = pix[j+1]
					cr0[i+2] = pix[j+2]
					i += 3
				}
			} else {
				for x := b.Min.X; x < b.Max.X; x++ {
					r, g, b, _ := m.At(x, y).RGBA()
					cr0[i+0] = uint8(r >> 8)
					cr0[i+1] = uint8(g >> 8)
					cr0[i+2] = uint8(b >> 8)
					i += 3
				}
			}

		case cbP8:
			if paletted != nil {
				offset := (y - b.Min.Y) * paletted.Stride
				copy(cr[0][1:], paletted.Pix[offset:offset+b.Dx()])
			} else {
				pi := m.(image.PalettedImage)
				for x := b.Min.X; x < b.Max.X; x++ {
					cr[0][i] = pi.ColorIndexAt(x, y)
					i += 1
				}
			}

		case cbP4, cbP2, cbP1:
			pi := m.(image.PalettedImage)

			var a uint8
			var c int
			pixelsPerByte := 8 / bitsPerPixel
			for x := b.Min.X; x < b.Max.X; x++ {
				a = a<<uint(bitsPerPixel) | pi.ColorIndexAt(x, y)
				c++
				if c == pixelsPerByte {
					cr[0][i] = a
					i += 1
					a = 0
					c = 0
				}
			}
			if c != 0 {
				for c != pixelsPerByte {
					a = a << uint(bitsPerPixel)
					c++
				}
				cr[0][i] = a
			}

		case cbTCA8:
			if nrgba != nil {
				offset := (y - b.Min.Y) * nrgba.Stride
				copy(cr[0][1:], nrgba.Pix[offset:offset+b.Dx()*4])
			} else {
				for x := b.Min.X; x < b.Max.X; x++ {
					c := color.NRGBAModel.Convert(m.At(x, y)).(color.NRGBA)
					cr[0][i+0] = c.R
					cr[0][i+1] = c.G
					cr[0][i+2] = c.B
					cr[0][i+3] = c.A
					i += 4
				}
			}

		case cbG16:
			for x := b.Min.X; x < b.Max.X; x++ {
				c := color.Gray16Model.Convert(m.At(x, y)).(color.Gray16)
				cr[0][i+0] = uint8(c.Y >> 8)
				cr[0][i+1] = uint8(c.Y)
				i += 2
			}

		case cbGA16:
			for x := b.Min.X; x < b.Max.X; x++ {
				c := color.NRGBA64Model.Convert(m.At(x, y)).(color.NRGBA64)
				c2 := color.Gray16Model.Convert(color.RGBA64{R: c.R, G: c.G, B: c.B, A: math.MaxUint16}).(color.Gray16)
				cr[0][i+0] = uint8(c2.Y >> 8)
				cr[0][i+1] = uint8(c2.Y)
				cr[0][i+2] = uint8(c.A >> 8)
				cr[0][i+3] = uint8(c.A)
				i += 4
			}

		case cbTC16:
			for x := b.Min.X; x < b.Max.X; x++ {
				r, g, b, _ := m.At(x, y).RGBA()
				cr[0][i+0] = uint8(r >> 8)
				cr[0][i+1] = uint8(r)
				cr[0][i+2] = uint8(g >> 8)
				cr[0][i+3] = uint8(g)
				cr[0][i+4] = uint8(b >> 8)
				cr[0][i+5] = uint8(b)
				i += 6
			}

		case cbTCA16:
			for x := b.Min.X; x < b.Max.X; x++ {
				c := color.NRGBA64Model.Convert(m.At(x, y)).(color.NRGBA64)
				cr[0][i+0] = uint8(c.R >> 8)
				cr[0][i+1] = uint8(c.R)
				cr[0][i+2] = uint8(c.G >> 8)
				cr[0][i+3] = uint8(c.G)
				cr[0][i+4] = uint8(c.B >> 8)
				cr[0][i+5] = uint8(c.B)
				cr[0][i+6] = uint8(c.A >> 8)
				cr[0][i+7] = uint8(c.A)
				i += 8
			}
		}

		f := ftNone
		if level != zlib.NoCompression && cb != cbP8 && cb != cbP4 && cb != cbP2 && cb != cbP1 && cb != cbG4 && cb != cbG2 && cb != cbG1 {
			bpp := bitsPerPixel / 8
			f = filter(&cr, pr, bpp)
		}

		if _, err := e.zw.Write(cr[f]); err != nil {
			return err
		}

		pr, cr[0] = cr[0], pr
	}

	return nil
}

func (e *encoder) writeIDATs() {
	if e.err != nil {
		return
	}
	if e.bw == nil {
		e.bw = bufio.NewWriterSize(e, 1<<15)
	} else {
		e.bw.Reset(e)
	}
	e.err = e.writeImage(e.bw, e.m, e.cb, levelToZlib(e.enc.CompressionLevel))
	if e.err != nil {
		return
	}
	e.err = e.bw.Flush()
}

// Encoder.CompressionLevel to map to zlib.DefaultCompression.
func levelToZlib(l CompressionLevel) int {
	switch l {
	case DefaultCompression:
		return zlib.DefaultCompression
	case NoCompression:
		return zlib.NoCompression
	case BestSpeed:
		return zlib.BestSpeed
	case BestCompression:
		return zlib.BestCompression
	default:
		return zlib.DefaultCompression
	}
}

func (e *encoder) writeIEND() { e.writeChunk(nil, "IEND") }

func Encode(w io.Writer, m image.Image) error {
	var e Encoder
	return e.Encode(w, m)
}

func (enc *Encoder) Encode(w io.Writer, m image.Image) error {
	mw, mh := int64(m.Bounds().Dx()), int64(m.Bounds().Dy())
	if mw <= 0 || mh <= 0 || mw >= 1<<32 || mh >= 1<<32 {
		return FormatError("invalid image size: " + strconv.FormatInt(mw, 10) + "x" + strconv.FormatInt(mh, 10))
	}

	var e *encoder
	if enc.BufferPool != nil {
		buffer := enc.BufferPool.Get()
		e = (*encoder)(buffer)

	}
	if e == nil {
		e = &encoder{}
	}
	if enc.BufferPool != nil {
		defer enc.BufferPool.Put((*EncoderBuffer)(e))
	}

	e.enc = enc
	e.w = w
	e.m = m

	var pal color.Palette
	if _, ok := m.(image.PalettedImage); ok {
		pal, _ = m.ColorModel().(color.Palette)
	}
	if pal != nil {
		if len(pal) <= 2 {
			e.cb = cbP1
		} else if len(pal) <= 4 {
			e.cb = cbP2
		} else if len(pal) <= 16 {
			e.cb = cbP4
		} else {
			e.cb = cbP8
		}
	} else {
		switch m.ColorModel() {
		case color.GrayModel:
			e.cb = cbG8
		case color.Gray16Model:
			e.cb = cbG16
		case color.RGBAModel, color.NRGBAModel, color.AlphaModel:
			if opaque(m) {
				e.cb = cbTC8
			} else {
				e.cb = cbTCA8
			}
		default:
			if opaque(m) {
				e.cb = cbTC16
			} else {
				e.cb = cbTCA16
			}
		}
	}

	_, e.err = io.WriteString(w, pngHeader)
	e.writeIHDR()
	if pal != nil {
		e.writePLTEAndTRNS(pal)
	}
	e.writeIDATs()
	e.writeIEND()
	return e.err
}

func Encode2(w io.Writer, m image.Image, cb int) error {
	var e Encoder
	return e.Encode2(w, m, cb)
}

func (enc *Encoder) Encode2(w io.Writer, m image.Image, cb int) error {
	mw, mh := int64(m.Bounds().Dx()), int64(m.Bounds().Dy())
	if mw <= 0 || mh <= 0 || mw >= 1<<32 || mh >= 1<<32 {
		return FormatError("invalid image size: " + strconv.FormatInt(mw, 10) + "x" + strconv.FormatInt(mh, 10))
	}

	var e *encoder
	if enc.BufferPool != nil {
		buffer := enc.BufferPool.Get()
		e = (*encoder)(buffer)

	}
	if e == nil {
		e = &encoder{}
	}
	if enc.BufferPool != nil {
		defer enc.BufferPool.Put((*EncoderBuffer)(e))
	}

	e.enc = enc
	e.w = w
	e.m = m
	e.cb = cb

	_, e.err = io.WriteString(w, pngHeader)
	e.writeIHDR()

	if e.cb == cbP1 || e.cb == cbP2 || e.cb == cbP4 || e.cb == cbP8 {
		var pal color.Palette
		if _, ok := m.(image.PalettedImage); ok {
			pal, _ = m.ColorModel().(color.Palette)
			if pal != nil {
				e.writePLTEAndTRNS(pal)
			}
		}
	}

	e.writeIDATs()
	e.writeIEND()
	return e.err
}
