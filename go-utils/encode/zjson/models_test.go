package zjson

import (
	"encoding/json"
	"fmt"
	"sync"
)

type Data1 struct {
	B1 bool     `json:"b1,omitempty"`
	B2 bool     `json:"b2,omitempty"`
	I1 uint8    `json:"i1,omitempty"`
	I2 uint16   `json:"i2,omitempty"`
	I3 int64    `json:"i3,omitempty"`
	I4 uint64   `json:"i4,omitempty"`
	F1 float64  `json:"f1,omitempty"`
	S1 string   `json:"s1,omitempty"`
	S2 string   `json:"s2,omitempty"`
	D  []*Data2 `json:"d,omitempty"`
}

type Data2 struct {
	B1 bool   `json:"b1,omitempty"`
	B2 bool   `json:"b2,omitempty"`
	I1 uint8  `json:"i1,omitempty"`
	I2 uint16 `json:"i2,omitempty"`
}

var (
	d1 = &Data1{
		B1: false,
		B2: true,
		I1: 1,
		I2: 2,
		I3: 12345678,
		I4: 1234567890,
		F1: 123.456,
		S1: "haha",
		S2: "wat een dag",
		D: []*Data2{
			{true, false, 5, 231},
			{false, false, 15, 131},
			{true, true, 25, 31},
		},
	}

	j1 []byte
)

func init() {
	j1, _ = json.Marshal(d1)
}

func (d *Data1) IsEmpty() bool { return false }

func (d *Data1) EncodeFields(e *Encoder) {
	e.BoolFieldOpt("b1", d.B1)
	e.BoolFieldOpt("b2", d.B2)
	e.Uint8FieldOpt("i1", d.I1)
	e.Uint16FieldOpt("i2", d.I2)
	e.Int64FieldOpt("i3", d.I3)
	e.Uint64FieldOpt("i4", d.I4)
	e.FloatFieldOpt("f1", d.F1, 'g', -1)
	e.StringFieldOpt("s1", d.S1)
	e.StringFieldOpt("s2", d.S2)
	if l := len(d.D); l > 0 {
		e.StartArrayField("d")
		for ix := 0; ix < l; ix++ {
			e.Object(d.D[ix])
		}
		e.EndArray()
	}
}

func (d *Data1) DecodeField(dec *Decoder, key []byte) error {
	if len(key) == 1 && key[0] == 'd' {
		if !dec.Array(d.addD) {
			return dec.Error()
		}
		return nil
	}

	if len(key) != 2 {
		return fmt.Errorf("unknown field for Data1")
	}

	var ok bool
	switch key[0] {
	case 'b':
		switch key[1] {
		case '1':
			d.B1, ok = dec.Bool()
		case '2':
			d.B2, ok = dec.Bool()
		default:
			return fmt.Errorf("unknown field for Data1")
		}
	case 'i':
		switch key[1] {
		case '1':
			d.I1, ok = dec.Uint8()
		case '2':
			d.I2, ok = dec.Uint16()
		case '3':
			d.I3, ok = dec.Int64()
		case '4':
			d.I4, ok = dec.Uint64()
		default:
			return fmt.Errorf("unknown field for Data1")
		}
	case 'f':
		switch key[1] {
		case '1':
			d.F1, ok = dec.Float()
		default:
			return fmt.Errorf("unknown field for Data1")
		}
	case 's':
		var b []byte
		var escaped bool
		switch key[1] {
		case '1':
			if b, escaped, ok = dec.String(); ok && b != nil {
				if escaped {
					d.S1 = string(dec.Unescaped(b))
				} else {
					d.S1 = string(b)
				}
			}
		case '2':
			if b, escaped, ok = dec.String(); ok && b != nil {
				if escaped {
					d.S2 = string(dec.Unescaped(b))
				} else {
					d.S2 = string(b)
				}
			}
		default:
			return fmt.Errorf("unknown field for Data1")
		}
	default:
		return fmt.Errorf("unknown field for Data1")
	}

	if ok {
		return nil
	}
	return dec.Error()
}

func (d *Data1) addD(dec *Decoder) error {
	if d.D == nil {
		d.D = make([]*Data2, 0, 4)
	}
	d2 := NewData2()
	if !dec.Object(d2) {
		return dec.Error()
	}
	d.D = append(d.D, d2)
	return nil
}

func (d *Data2) IsEmpty() bool { return false }

func (d *Data2) EncodeFields(e *Encoder) {
	e.BoolFieldOpt("b1", d.B1)
	e.BoolFieldOpt("b2", d.B2)
	e.Uint8FieldOpt("i1", d.I1)
	e.Uint16FieldOpt("i2", d.I2)
}

func (d *Data2) DecodeField(dec *Decoder, key []byte) error {
	if len(key) != 2 {
		return fmt.Errorf("unknown field for Data2")
	}

	var ok bool
	switch key[0] {
	case 'b':
		switch key[1] {
		case '1':
			d.B1, ok = dec.Bool()
		case '2':
			d.B2, ok = dec.Bool()
		default:
			return fmt.Errorf("unknown field for Data2")
		}
	case 'i':
		switch key[1] {
		case '1':
			d.I1, ok = dec.Uint8()
		case '2':
			d.I2, ok = dec.Uint16()
		default:
			return fmt.Errorf("unknown field for Data2")
		}
	}

	if ok {
		return nil
	}
	return dec.Error()
}

func NewData1() *Data1 {
	d := data1pool.Get().(*Data1)
	if d.D == nil {
		d.D = make([]*Data2, 0, 4)
	}
	return d
}

func (d *Data2) ReturnToPool() {
	if d != nil {
		d.B1 = false
		d.B2 = false
		d.I1 = 0
		d.I2 = 0
		data2pool.Put(d)
	}
}

func NewData2() *Data2 {
	return data2pool.Get().(*Data2)
}

func (d *Data1) ReturnToPool() {
	if d != nil {
		for ix := range d.D {
			d.D[ix].ReturnToPool()
			d.D[ix] = nil
		}

		d.B1 = false
		d.B2 = false
		d.I1 = 0
		d.I2 = 0
		d.I3 = 0
		d.I4 = 0
		d.F1 = 0.0
		d.S1 = ""
		d.S2 = ""
		d.D = d.D[:0]

		data1pool.Put(d)
	}
}

var (
	data1pool = sync.Pool{New: func() interface{} { return &Data1{} }}
	data2pool = sync.Pool{New: func() interface{} { return &Data2{} }}
)
