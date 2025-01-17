package pool

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/encode/zjson"
)

type testObject1 struct {
	Object
	i1 int
	f1 float64
	b1 bool
}

func (o *testObject1) reset() {
	o.i1 = 0
	o.f1 = 0.0
	o.b1 = false
}

func (o *testObject1) Clone() Objecter {
	n := o.Acquire().(*testObject1)
	n.i1 = o.i1
	n.f1 = o.f1
	n.b1 = o.b1
	return n
}

func (o *testObject1) IsEmpty() bool {
	return o == nil || (o.i1 == 0 && o.f1 == 0.0 && !o.b1)
}

func (o *testObject1) EncodeFields(enc *zjson.Encoder) {
	enc.IntFieldOpt("i1", o.i1)
	enc.FloatFieldOpt("f1", o.f1, 'f', 4)
	enc.BoolFieldOpt("b1", o.b1)
}

func (o *testObject1) Decode(dec *zjson.Decoder) error {
	if dec.Object(o) {
		return nil
	}
	return dec.Error()
}

func (o *testObject1) DecodeField(dec *zjson.Decoder, key []byte) error {
	var ok, b1 bool
	var i1 int
	var f1 float64

	if string(key) == "i1" {
		if i1, ok = dec.Int(); ok {
			o.i1 = i1
		}
	} else if string(key) == "f1" {
		if f1, ok = dec.Float(); ok {
			o.f1 = f1
		}
	} else if string(key) == "b1" {
		if b1, ok = dec.Bool(); ok {
			o.b1 = b1
		}
	} else {
		return fmt.Errorf("invalid field")
	}

	if ok {
		return nil
	}
	return dec.Error()
}

func TestNewObject(t *testing.T) {
	testCases := []struct {
		name string
		i1   int
		f1   float64
		b1   bool
		j    string
	}{
		{name: "i1", i1: 12345, j: `"object":{"i1":12345}`},
		{name: "f1", f1: 123.45, j: `"object":{"f1":123.4500}`},
		{name: "b1", b1: true, j: `"object":{"b1":true}`},
		{name: "all", i1: 12345, f1: 123.45, b1: true, j: `"object":{"i1":12345,"f1":123.4500,"b1":true}`},
		{name: "empty", j: ``},
	}

	p := NewProducer(func() (Objecter, func()) { o := &testObject1{}; return o, o.reset })
	require.NotNil(t, p)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			n1 := p.Acquire()
			require.NotNil(t, n1)

			o1, ok1 := n1.(*testObject1)
			require.True(t, ok1)
			require.NotNil(t, o1)

			assert.Zero(t, o1.i1)
			assert.Zero(t, o1.f1)
			assert.False(t, o1.b1)
			assert.True(t, o1.IsEmpty())

			o1.i1 = tc.i1
			o1.f1 = tc.f1
			o1.b1 = tc.b1

			assert.Equal(t, tc.i1 == 0 && tc.f1 == 0.0 && !tc.b1, o1.IsEmpty())

			enc := zjson.AcquireEncoder(64)
			defer enc.Release()
			enc.ObjectFieldOpt("object", o1)
			j := enc.Bytes()
			assert.Equal(t, tc.j, string(j))

			n2 := p.Acquire()
			require.NotNil(t, n2)

			o2, ok2 := n2.(*testObject1)
			require.True(t, ok2)
			require.NotNil(t, o2)

			assert.Zero(t, o2.i1)
			assert.Zero(t, o2.f1)
			assert.False(t, o2.b1)
			assert.True(t, o2.IsEmpty())

			o2.i1 = tc.i1
			o2.f1 = tc.f1
			o2.b1 = tc.b1

			assert.Equal(t, tc.i1 == 0 && tc.f1 == 0.0 && !tc.b1, o2.IsEmpty())

			if len(j) > 0 {
				j = []byte(strings.Replace(string(j), `"object":`, ``, 1))
				n3, err := p.AcquireFromJSON(j)
				require.NoError(t, err)
				require.NotNil(t, n3)

				o3, ok3 := n3.(*testObject1)
				require.True(t, ok3)
				require.NotNil(t, o3)
				require.NotNil(t, n3)
				assert.Equal(t, tc.i1, o3.i1)
				assert.Equal(t, tc.f1, o3.f1)
				assert.Equal(t, tc.b1, o3.b1)

				o3.Release()
			}

			o2.Release()
			o1.Release()
		})
	}
}

func TestManyNewObjects(t *testing.T) {
	t.Run("many new object", func(t *testing.T) {
		p := NewProducer(func() (Objecter, func()) { o := &testObject1{}; return o, o.reset })
		require.NotNil(t, p)

		list := make([]*testObject1, 100)
		for ix := range list {
			n := p.Acquire().(*testObject1)
			assert.Zero(t, n.i1)
			assert.Zero(t, n.f1)
			assert.False(t, n.b1)
			list[ix] = n
		}

		for ix := range list {
			list[ix].Release()
			list[ix] = nil
		}
	})
}

func TestObjectClone(t *testing.T) {
	testCases := []struct {
		name string
		i1   int
		f1   float64
		b1   bool
	}{
		{name: "i1", i1: 12345},
		{name: "f1", f1: 123.45},
		{name: "b1", b1: true},
		{name: "all", i1: 12345, f1: 123.45, b1: true},
		{name: "empty"},
	}

	p := NewProducer(func() (Objecter, func()) { o := &testObject1{}; return o, o.reset })
	require.NotNil(t, p)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			n1 := p.Acquire()
			require.NotNil(t, n1)

			o1, ok1 := n1.(*testObject1)
			require.True(t, ok1)
			require.NotNil(t, o1)

			assert.Zero(t, o1.i1)
			assert.Zero(t, o1.f1)
			assert.False(t, o1.b1)

			o1.i1 = tc.i1
			o1.f1 = tc.f1
			o1.b1 = tc.b1

			n2 := o1.Clone()
			require.NotNil(t, n2)

			o1.Release()

			o2, ok2 := n2.(*testObject1)
			require.True(t, ok2)
			require.NotNil(t, o2)

			assert.Equal(t, tc.i1, o2.i1)
			assert.Equal(t, tc.f1, o2.f1)
			assert.Equal(t, tc.b1, o2.b1)

			o2.Release()
		})
	}
}

func TestObject_DecodeError(t *testing.T) {
	t.Run("decode error", func(t *testing.T) {
		p := NewProducer(func() (Objecter, func()) { o := &testObject1{}; return o, o.reset })
		require.NotNil(t, p)

		n, err := p.AcquireFromJSON([]byte{'x', 'y', 'z'})
		assert.Error(t, err)
		assert.Nil(t, n)
	})
}

type testObject2 struct {
	i1 int
	Object
}

func (o *testObject2) reset()                      { o.i1 = 0 }
func (o *testObject2) EncodeFields(*zjson.Encoder) {}

func TestObecter_Defaults1(t *testing.T) {
	t.Run("objecter defaults (1)", func(t *testing.T) {
		p := NewProducer(func() (Objecter, func()) { o := &testObject2{}; return o, o.reset })
		require.NotNil(t, p)

		n := p.Acquire()
		require.NotNil(t, n)

		o1, ok := n.(*testObject2)
		require.True(t, ok)
		require.NotNil(t, o1)

		got := o1.IsEmpty()
		assert.False(t, got)

		enc := zjson.AcquireEncoder(512)
		defer enc.Release()

		o1.EncodeFields(enc)
		j := enc.Bytes()
		assert.Equal(t, "", string(j))

		enc.Reset()
		enc.ObjectOpt(o1)
		j = enc.Bytes()
		assert.Equal(t, "{}", string(j))

		enc.Reset()
		enc.ObjectField("o1", o1)
		j = enc.Bytes()
		assert.Equal(t, `"o1":{}`, string(j))

		enc.Reset()
		enc.ObjectFieldOpt("o1", o1)
		j = enc.Bytes()
		assert.Equal(t, `"o1":{}`, string(j))
	})
}

type testObject3 struct {
	i1 int
	Object
}

func (o *testObject3) reset()                      { o.i1 = 0 }
func (o *testObject3) IsEmpty() bool               { return o.i1 == 0 }
func (o *testObject3) EncodeFields(*zjson.Encoder) {}

func TestObecter_Defaults2(t *testing.T) {
	t.Run("objecter defaults (2)", func(t *testing.T) {
		p := NewProducer(func() (Objecter, func()) { o := &testObject3{}; return o, o.reset })
		require.NotNil(t, p)

		n := p.Acquire()
		require.NotNil(t, n)

		o1, ok := n.(*testObject3)
		require.True(t, ok)
		require.NotNil(t, o1)

		got := o1.IsEmpty()
		assert.True(t, got)

		enc := zjson.AcquireEncoder(512)
		defer enc.Release()

		o1.EncodeFields(enc)
		j := enc.Bytes()
		assert.Equal(t, "", string(j))

		enc.Reset()
		enc.Object(o1)
		j = enc.Bytes()
		assert.Equal(t, "{}", string(j))

		enc.Reset()
		enc.ObjectField("o1", o1)
		j = enc.Bytes()
		assert.Equal(t, `"o1":{}`, string(j))

		enc.Reset()
		enc.ObjectFieldOpt("o1", o1)
		j = enc.Bytes()
		assert.Equal(t, "", string(j))
	})
}

type testObject4 struct {
	i1 int
	Object
}

func (o *testObject4) reset() { o.i1 = 0 }

func TestEncodeNotImplemented(t *testing.T) {
	t.Run("encode not implemented", func(t *testing.T) {
		p := NewProducer(func() (Objecter, func()) { o := &testObject4{}; return o, o.reset })
		require.NotNil(t, p)

		n := p.Acquire()
		require.NotNil(t, n)

		o, ok := n.(*testObject4)
		require.True(t, ok)
		require.NotNil(t, o)

		defer func() {
			e := recover()
			assert.NotNil(t, e)
		}()

		o.Encode(nil)
	})
}

func TestDecodeNotImplemented(t *testing.T) {
	t.Run("decode not implemented", func(t *testing.T) {
		p := NewProducer(func() (Objecter, func()) { o := &testObject4{}; return o, o.reset })
		require.NotNil(t, p)

		n := p.Acquire()
		require.NotNil(t, n)

		o, ok := n.(*testObject4)
		require.True(t, ok)
		require.NotNil(t, o)

		defer func() {
			e := recover()
			assert.NotNil(t, e)
		}()

		got := o.Decode(nil)
		assert.Error(t, got)
	})
}

func TestEncodeFieldsNotImplemented(t *testing.T) {
	t.Run("encode fields not implemented", func(t *testing.T) {
		p := NewProducer(func() (Objecter, func()) { o := &testObject4{}; return o, o.reset })
		require.NotNil(t, p)

		n := p.Acquire()
		require.NotNil(t, n)

		o, ok := n.(*testObject4)
		require.True(t, ok)
		require.NotNil(t, o)

		defer func() {
			e := recover()
			assert.NotNil(t, e)
		}()

		o.EncodeFields(nil)
	})
}

func TestDecodeFieldNotImplemented(t *testing.T) {
	t.Run("decode field not implemented", func(t *testing.T) {
		p := NewProducer(func() (Objecter, func()) { o := &testObject4{}; return o, o.reset })
		require.NotNil(t, p)

		n := p.Acquire()
		require.NotNil(t, n)

		o, ok := n.(*testObject4)
		require.True(t, ok)
		require.NotNil(t, o)

		defer func() {
			e := recover()
			assert.NotNil(t, e)
		}()

		got := o.DecodeField(nil, nil)
		assert.Error(t, got)
	})
}
