package png

import (
	"bytes"
	"image"
	"image/color"
	"math"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncodeGX(t *testing.T) {
	testCases := []struct {
		name   string
		gen    func(w, h int) *image.Gray
		width  int
		height int
		depth  int
	}{
		{
			name: "G1 - random",
			gen: func(w, h int) *image.Gray {
				img := image.NewGray(image.Rect(0, 0, w, h))
				for x := 0; x < w; x++ {
					for y := 0; y < h; y++ {
						g := uint8(rand.Int31n(2) * 0xff)
						img.Set(x, y, color.Gray{Y: g})
					}
				}
				return img
			},
			width:  255,
			height: 255,
			depth:  cbG1,
		},
		{
			name: "G1 - vertical lines",
			gen: func(w, h int) *image.Gray {
				img := image.NewGray(image.Rect(0, 0, w, h))
				for x := 0; x < w; x += 2 {
					for y := 0; y < h; y++ {
						img.Set(x, y, color.Gray{Y: 0xff})
					}
				}
				return img
			},
			width:  255,
			height: 255,
			depth:  cbG1,
		},
		{
			name: "G1 - horizontal lines",
			gen: func(w, h int) *image.Gray {
				img := image.NewGray(image.Rect(0, 0, w, h))
				for x := 0; x < w; x++ {
					for y := 0; y < h; y += 2 {
						img.Set(x, y, color.Gray{Y: 0xff})
					}
				}
				return img
			},
			width:  252,
			height: 252,
			depth:  cbG1,
		},
		{
			name: "G2 - random",
			gen: func(w, h int) *image.Gray {
				img := image.NewGray(image.Rect(0, 0, w, h))
				for x := 0; x < w; x++ {
					for y := 0; y < h; y++ {
						var g uint8
						switch rand.Int31n(4) {
						case 1:
							g = 0x55
						case 2:
							g = 0xaa
						case 3:
							g = 0xff
						}
						img.Set(x, y, color.Gray{Y: g})
					}
				}
				return img
			},
			width:  252,
			height: 252,
			depth:  cbG2,
		},
		{
			name: "G2 - vertical lines",
			gen: func(w, h int) *image.Gray {
				img := image.NewGray(image.Rect(0, 0, w, h))
				for x := 0; x < w; x += 4 {
					for y := 0; y < h; y++ {
						img.Set(x, y, color.Gray{Y: 0xff})
					}
				}
				for x := 1; x < w; x += 4 {
					for y := 0; y < h; y++ {
						img.Set(x, y, color.Gray{Y: 0xaa})
					}
				}
				for x := 2; x < w; x += 4 {
					for y := 0; y < h; y++ {
						img.Set(x, y, color.Gray{Y: 0x55})
					}
				}
				return img
			},
			width:  254,
			height: 254,
			depth:  cbG2,
		},
		{
			name: "G2 - horizontal lines",
			gen: func(w, h int) *image.Gray {
				img := image.NewGray(image.Rect(0, 0, w, h))
				for x := 0; x < w; x++ {
					for y := 0; y < h; y += 4 {
						img.Set(x, y, color.Gray{Y: 0xff})
					}
				}
				for x := 0; x < w; x++ {
					for y := 1; y < h; y += 4 {
						img.Set(x, y, color.Gray{Y: 0xaa})
					}
				}
				for x := 0; x < w; x++ {
					for y := 2; y < h; y += 4 {
						img.Set(x, y, color.Gray{Y: 0x55})
					}
				}
				return img
			},
			width:  254,
			height: 254,
			depth:  cbG2,
		},
		{
			name: "G4 - random",
			gen: func(w, h int) *image.Gray {
				img := image.NewGray(image.Rect(0, 0, w, h))
				for x := 0; x < w; x++ {
					for y := 0; y < h; y++ {
						g := uint8(rand.Int31n(16))
						g += g << 4
						img.Set(x, y, color.Gray{Y: g})
					}
				}
				return img
			},
			width:  252,
			height: 252,
			depth:  cbG4,
		},
		{
			name: "G4 - vertical lines",
			gen: func(w, h int) *image.Gray {
				img := image.NewGray(image.Rect(0, 0, w, h))
				for x := 0; x < w; x++ {
					for y := 0; y < h; y++ {
						g := uint8(x) & 0xf
						g += g << 4
						img.Set(x, y, color.Gray{Y: g})
					}
				}
				return img
			},
			width:  252,
			height: 252,
			depth:  cbG4,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g1 := tc.gen(tc.width, tc.height)

			buf := &bytes.Buffer{}
			err := Encode2(buf, g1, tc.depth)
			require.NoError(t, err)
			assert.NotZero(t, buf.Len())

			f, _ := os.CreateTemp("/tmp", "test_encode_gx_*.png")
			f.Write(buf.Bytes())
			f.Close()
			defer os.Remove(f.Name())

			var img2 image.Image
			img2, err = Decode(buf)
			require.NoError(t, err)
			require.NotNil(t, img2)

			g2, ok := img2.(*image.Gray)
			require.True(t, ok)
			require.NotNil(t, g2)

			assert.EqualValues(t, g1, g2)
		})
	}
}

func TestEncodeGA8(t *testing.T) {
	rgba1 := randomGA8(256, 256)

	buf := &bytes.Buffer{}
	err := Encode2(buf, rgba1, cbGA8)
	require.NoError(t, err)
	assert.NotZero(t, buf.Len())

	f, _ := os.CreateTemp("/tmp", "test_encode_ga8_*.png")
	f.Write(buf.Bytes())
	f.Close()
	defer os.Remove(f.Name())

	var img2 image.Image
	img2, err = Decode(buf)
	require.NoError(t, err)
	require.NotNil(t, img2)

	rgba2, ok := img2.(*image.NRGBA)
	require.True(t, ok)
	require.NotNil(t, rgba2)

	assert.EqualValues(t, rgba1, rgba2)
}

func TestEncodeGA16(t *testing.T) {
	rgba1 := randomGA16(256, 256)

	buf := &bytes.Buffer{}
	err := Encode2(buf, rgba1, cbGA16)
	require.NoError(t, err)
	assert.NotZero(t, buf.Len())

	f, _ := os.CreateTemp("/tmp", "test_encode_ga16_*.png")
	f.Write(buf.Bytes())
	f.Close()
	defer os.Remove(f.Name())

	var img2 image.Image
	img2, err = Decode(buf)
	require.NoError(t, err)
	require.NotNil(t, img2)

	rgba2, ok := img2.(*image.NRGBA64)
	require.True(t, ok)
	require.NotNil(t, rgba2)

	assert.EqualValues(t, rgba1, rgba2)
}

func randomGA8(w, h int) *image.NRGBA {
	img := image.NewNRGBA(image.Rect(0, 0, w, h))
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			g := uint8(rand.Int31n(math.MaxUint8 + 1))
			c := color.RGBAModel.Convert(color.Gray{Y: g}).(color.RGBA)
			img.Set(x, y, color.NRGBA{
				R: c.R,
				G: c.G,
				B: c.B,
				A: uint8(rand.Int31n(math.MaxUint8 + 1)),
			})
		}
	}
	return img
}

func randomGA16(w, h int) *image.NRGBA64 {
	img := image.NewNRGBA64(image.Rect(0, 0, w, h))
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			g := uint16(rand.Int31n(math.MaxUint16 + 1))
			c := color.RGBA64Model.Convert(color.Gray16{Y: g}).(color.RGBA64)
			img.Set(x, y, color.NRGBA64{
				R: c.R,
				G: c.G,
				B: c.B,
				A: uint16(rand.Int31n(math.MaxUint16 + 1)),
			})
		}
	}
	return img
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
