package object

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/encode/zjson"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/object-pool/pool"
)

func TestResetObjects(t *testing.T) {
	list := make(Objects, 14)
	for ix := range list {
		n := testObjManager.Acquire().(*testObject1)
		n.i1 = rand.Intn(100000)
		n.f1 = rand.Float64() * 100000
		n.b1 = rand.Intn(100) < 50
		list[ix] = n
	}

	testCases := []struct {
		name string
		objs Objects
		min  int
		max  int
	}{
		{"nil", nil, 4, 8},
		{"empty", Objects{}, 4, 8},
		{"single", Objects{list[0]}, 4, 8},
		{"few", append(Objects{}, list[1:4]...), 4, 8},
		{"many", append(Objects{}, list[4:]...), 4, 8},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			n := ResetObjects(tc.objs, tc.min, tc.max)
			require.NotNil(t, n)
			assert.GreaterOrEqual(t, cap(n), tc.min)
			assert.LessOrEqual(t, cap(n), tc.max)

			l := cap(n)
			n = n[:l]
			for ix := range n {
				assert.Nil(t, n[ix])
			}
		})
	}

	for ix := range list {
		n := list[ix].(*testObject1)
		assert.Zero(t, n.i1)
		assert.Zero(t, n.f1)
		assert.False(t, n.b1)
	}
	list = list[:0]
}

func TestObjects_Replace(t *testing.T) {
	list1 := make(Objects, 20)
	for ix := range list1 {
		n := testObjManager.Acquire().(*testObject1)
		n.i1 = rand.Intn(100000)
		n.f1 = rand.Float64() * 100000
		n.b1 = rand.Intn(100) < 50
		list1[ix] = n
	}

	list2 := make(Objects, 28)
	for ix := range list2 {
		n := testObjManager.Acquire().(*testObject1)
		n.i1 = rand.Intn(100000)
		n.f1 = rand.Float64() * 100000
		n.b1 = rand.Intn(100) < 50
		list2[ix] = n
	}

	testCases := []struct {
		name  string
		objs  Objects
		input Objects
		want  Objects
	}{
		{"nil - nil", nil, nil, nil},
		{"empty - nil", Objects{}, nil, Objects{}},
		{"nil - empty", nil, Objects{}, Objects{}},
		{"empty - empty", Objects{}, Objects{}, Objects{}},
		{"one - nil", Objects{list1[0]}, nil, Objects{}},
		{"one - empty", Objects{list1[1]}, Objects{}, Objects{}},
		{"one - one", Objects{list1[2]}, Objects{list2[0]}, Objects{list2[0]}},
		{"one - few", Objects{list1[3]}, list2[1:5], list2[1:5]},
		{"one - many", Objects{list1[4]}, list2[5:14], list2[5:14]},
		{"few - nil", Objects{list1[5], list1[6], list1[7]}, nil, Objects{}},
		{"few - empty", Objects{list1[8], list1[9], list1[10]}, Objects{}, Objects{}},
		{"few - one", Objects{list1[11], list1[12], list1[13]}, Objects{list2[14]}, Objects{list2[14]}},
		{"few - few", Objects{list1[14], list1[15], list1[16]}, list2[15:19], list2[15:19]},
		{"few - many", Objects{list1[17], list1[18], list1[19]}, list2[19:], list2[19:]},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_ = list1[0]
			got := Replace(tc.objs, tc.input)
			if tc.want == nil {
				require.Nil(t, got)
			} else {
				require.NotNil(t, got)
				assert.Equal(t, tc.input.IsEmpty(), got.IsEmpty())
				assert.EqualValues(t, tc.want, got)
			}
		})
	}

	list1 = list1[:0]
	list2 = list2[:0]
}

func TestObjects_IsEmpty(t *testing.T) {
	list := make(Objects, 3)
	for ix := range list {
		n := testObjManager.Acquire().(*testObject1)
		n.i1 = rand.Intn(100000)
		n.f1 = rand.Float64() * 100000
		n.b1 = rand.Intn(100) < 50
		list[ix] = n
	}
	testCases := []struct {
		name string
		objs Objects
		want bool
	}{
		{"nil", nil, true},
		{"empty", Objects{}, true},
		{"not empty", list, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.objs.IsEmpty()
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestObjects_Encode(t *testing.T) {
	o1 := &testObject1{i1: 1, f1: 11.1, b1: true}
	o2 := &testObject1{i1: 1, f1: 11.1}
	o3 := &testObject1{b1: true}
	o4 := &testObject1{f1: 11.1}
	o5 := &testObject1{f1: 11.1, b1: true}
	o6 := &testObject1{i1: 1}

	testCases := []struct {
		name string
		ints Objects
		want string
	}{
		{
			name: "nil",
			want: "[]",
		},
		{
			name: "empty",
			ints: Objects{},
			want: "[]",
		},
		{
			name: "one",
			ints: Objects{o4},
			want: `[{"f1":11.1}]`,
		},
		{
			name: "few",
			ints: Objects{o1, o2},
			want: `[{"i1":1,"f1":11.1,"b1":true},{"i1":1,"f1":11.1}]`,
		},
		{
			name: "many",
			ints: Objects{o3, o4, o2, o1, o5, o6},
			want: `[{"b1":true},{"f1":11.1},{"i1":1,"f1":11.1},{"i1":1,"f1":11.1,"b1":true},{"f1":11.1,"b1":true},{"i1":1}]`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			enc := zjson.AcquireEncoder(512)
			defer enc.Release()

			tc.ints.Encode(enc)
			j := enc.Bytes()
			require.NotNil(t, j)
			assert.Equal(t, tc.want, string(j))
		})
	}
}

func TestNewObjectsProducer(t *testing.T) {
	t.Run("new object producer", func(t *testing.T) {
		p1 := pool.NewProducer(func() (pool.Objecter, func()) { o := &testObject1{}; return o, o.reset })
		require.NotNil(t, p1)

		p2 := NewObjectsProducer(8, 16, p1)
		require.NotNil(t, p2)

		n1 := p2.Acquire()
		require.NotNil(t, n1)

		o1, ok1 := n1.(*ObjectsManager)
		require.True(t, ok1)
		require.NotNil(t, o1)
		assert.Equal(t, 8, o1.minSize)
		assert.Equal(t, 16, o1.maxSize)
		assert.NotNil(t, o1.Items)
		assert.Zero(t, len(o1.Items))
		assert.True(t, o1.IsEmpty())

		for ix := 0; ix < 19; ix++ {
			o1.Append(p1.Acquire().(pool.Objecter))
		}
		assert.Equal(t, 19, len(o1.Items))
		assert.False(t, o1.IsEmpty())

		o1.Release()

		o2 := p2.Acquire().(*ObjectsManager)
		require.NotNil(t, o2)

		assert.NotNil(t, o2.Items)
		assert.Zero(t, len(o2.Items))

		o2.Items = o2.Items[:cap(o2.Items)]
		for ix := range o2.Items {
			assert.Nil(t, o2.Items[ix])
		}

		o2.Release()
	})
}

func TestObjectsManager_Clone(t *testing.T) {
	testCases := []struct {
		name  string
		count int
	}{
		{"empty", 0},
		{"1 value", 1},
		{"few values", 3},
		{"many values", 15},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p1 := pool.NewProducer(func() (pool.Objecter, func()) { o := &testObject1{}; return o, o.reset })
			require.NotNil(t, p1)

			p2 := NewObjectsProducer(4, 32, p1)
			require.NotNil(t, p2)

			n1 := p2.Acquire()
			require.NotNil(t, n1)

			o1, ok1 := n1.(*ObjectsManager)
			require.True(t, ok1)
			require.NotNil(t, o1)
			assert.True(t, o1.IsEmpty())

			for ix := 0; ix < tc.count; ix++ {
				n := p1.Acquire().(*testObject1)
				n.i1 = ix + 1
				o1.Append(n)
			}

			assert.Equal(t, tc.count, len(o1.Items))
			assert.Equal(t, tc.count == 0, o1.IsEmpty())

			n2 := o1.Clone()
			require.NotNil(t, n2)

			o2, ok2 := n2.(*ObjectsManager)
			require.True(t, ok2)
			require.NotNil(t, o2)

			o1.Release()

			assert.Equal(t, tc.count, len(o2.Items))
			assert.Equal(t, tc.count == 0, o2.IsEmpty())

			for ix := 0; ix < tc.count; ix++ {
				o := o2.Items[ix].(*testObject1)
				assert.Equal(t, ix+1, o.i1)
			}

			enc := zjson.AcquireEncoder(512)
			defer enc.Release()

			enc.Array(o2)
			j := enc.Bytes()
			assert.NotZero(t, len(j))

			n3, err := p2.AcquireFromJSON(j)
			require.NoError(t, err)
			require.NotNil(t, n3)

			o3 := n3.(*ObjectsManager)
			require.NotNil(t, o3)
			require.Equal(t, len(o2.Items), len(o3.Items))

			for ix := range o2.Items {
				obj2 := o2.Items[ix].(*testObject1)
				obj3 := o3.Items[ix].(*testObject1)
				require.NotNil(t, obj2)
				require.NotNil(t, obj3)
				assert.Equal(t, obj2.b1, obj3.b1)
				assert.Equal(t, obj2.i1, obj3.i1)
				assert.Equal(t, obj2.f1, obj3.f1)
			}

			o3.Release()
			o2.Release()
		})
	}
}

type testObject1 struct {
	b1 bool
	i1 int
	f1 float64
	pool.Object
}

func (o *testObject1) reset() {
	o.i1 = 0
	o.f1 = 0.0
	o.b1 = false
}

func (o *testObject1) IsEmpty() bool {
	return o == nil
}

func (o *testObject1) EncodeFields(enc *zjson.Encoder) {
	enc.IntFieldOpt("i1", o.i1)
	enc.FloatFieldOpt("f1", o.f1, 'g', 4)
	enc.BoolFieldOpt("b1", o.b1)
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

var testObjManager = pool.NewProducer(func() (pool.Objecter, func()) {
	o := &testObject1{}
	return o, o.reset
})

func init() {
	rand.Seed(time.Now().UnixNano())
}
