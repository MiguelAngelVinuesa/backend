package kafka

import (
	"testing"

	"github.com/goccy/go-json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type myData struct {
	X string  `json:"x,omitempty"`
	Y float64 `json:"y,omitempty"`
}

func TestNewEvent(t *testing.T) {
	testCases := []struct {
		name     string
		clientID string
		entity   string
		key      string
		op       OperationKind
		data     any
		want     string
	}{
		{
			name: "empty",
			want: `{"producer":"","entity":"","op":0}`,
		},
		{
			name:     "clientID",
			clientID: "bebop",
			want:     `{"producer":"bebop","entity":"","op":0}`,
		},
		{
			name:   "entity",
			entity: "SESSION",
			want:   `{"producer":"","entity":"SESSION","op":0}`,
		},
		{
			name: "key",
			key:  "abcdef",
			want: `{"producer":"","entity":"","key":"abcdef","op":0}`,
		},
		{
			name: "created",
			op:   OpCreate,
			want: `{"producer":"","entity":"","op":1}`,
		},
		{
			name: "read",
			op:   OpRead,
			want: `{"producer":"","entity":"","op":2}`,
		},
		{
			name: "updated",
			op:   OpUpdate,
			want: `{"producer":"","entity":"","op":3}`,
		},
		{
			name: "deleted",
			op:   OpDelete,
			want: `{"producer":"","entity":"","op":4}`,
		},
		{
			name: "disabled",
			op:   OpDisable,
			want: `{"producer":"","entity":"","op":5}`,
		},
		{
			name: "enabled",
			op:   OpEnable,
			want: `{"producer":"","entity":"","op":6}`,
		},
		{
			name: "blocked",
			op:   OpBlock,
			want: `{"producer":"","entity":"","op":7}`,
		},
		{
			name: "unblocked",
			op:   OpUnblock,
			want: `{"producer":"","entity":"","op":8}`,
		},
		{
			name: "data map",
			data: map[string]any{"x": "y", "z": 1.1},
			want: `{"producer":"","entity":"","op":0,"data":{"x":"y","z":1.1}}`,
		},
		{
			name: "data struct",
			data: &myData{X: "oops", Y: 3.33},
			want: `{"producer":"","entity":"","op":0,"data":{"x":"oops","y":3.33}}`,
		},
		{
			name:     "all",
			clientID: "bebop2",
			entity:   "ROUND",
			key:      "abcdefgh",
			op:       OpCreate,
			data:     map[string]any{"y": "z", "x": 6.6},
			want:     `{"producer":"bebop2","entity":"ROUND","key":"abcdefgh","op":1,"data":{"x":6.6,"y":"z"}}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			e := NewEvent(tc.clientID, tc.entity, tc.key, tc.op, tc.data)
			require.NotNil(t, e)

			assert.Equal(t, tc.clientID, e.Producer)
			assert.Equal(t, tc.entity, e.EntityCode)
			assert.Equal(t, tc.key, e.EntityKey)
			assert.Equal(t, tc.op, e.Operation)
			assert.Equal(t, tc.data, e.Data)

			got, err := json.Marshal(e)
			require.NoError(t, err)
			require.NotNil(t, got)
			assert.Equal(t, tc.want, string(got))

			var d any
			switch tc.data.(type) {
			case map[string]any:
				d = make(map[string]any)
			case *myData:
				d = &myData{}
			}

			e2, err2 := NewEventFromJSON(got, d)
			require.NoError(t, err2)
			require.NotNil(t, e2)
			assert.EqualValues(t, e, e2)
		})
	}
}
