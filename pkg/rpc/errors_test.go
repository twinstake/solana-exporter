package rpc

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnpackRPCErrorData(t *testing.T) {
	t.Run("node-unhealthy", func(t *testing.T) {
		// numSlotsBehind arrives as a float64 because it was decoded from JSON into map[string]any.
		rpcErr := &Error{
			Code:    NodeUnhealthyCode,
			Method:  "getHealth",
			Message: "Node is unhealthy",
			Data:    map[string]any{"numSlotsBehind": float64(42)},
		}
		var data NodeUnhealthyErrorData
		require.NoError(t, UnpackRPCErrorData(rpcErr, &data))
		assert.Equal(t, int64(42), data.NumSlotsBehind)
	})

	t.Run("nil-data", func(t *testing.T) {
		rpcErr := &Error{Code: NodeUnhealthyCode, Method: "getHealth"}
		var data NodeUnhealthyErrorData
		require.NoError(t, UnpackRPCErrorData(rpcErr, &data))
		assert.Equal(t, int64(0), data.NumSlotsBehind)
	})

	t.Run("type-mismatch", func(t *testing.T) {
		// numSlotsBehind is a string, but the target field is int64 — unmarshal must fail.
		rpcErr := &Error{
			Method: "getHealth",
			Data:   map[string]any{"numSlotsBehind": "not-a-number"},
		}
		var data NodeUnhealthyErrorData
		assert.Error(t, UnpackRPCErrorData(rpcErr, &data))
	})
}

func TestError_Error(t *testing.T) {
	err := &Error{Method: "getHealth", Code: NodeUnhealthyCode, Message: "Node is unhealthy"}
	assert.Equal(t, "getHealth rpc error (code: -32005): Node is unhealthy (data: map[])", err.Error())
}

func TestHostProduction_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    HostProduction
		wantErr bool
	}{
		{"valid", "[10, 7]", HostProduction{LeaderSlots: 10, BlocksProduced: 7}, false},
		{"too-many-elements", "[10, 7, 3]", HostProduction{}, true},
		{"empty-array", "[]", HostProduction{}, true},
		{"not-an-array", `"foo"`, HostProduction{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var hp HostProduction
			err := json.Unmarshal([]byte(tt.input), &hp)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, hp)
		})
	}
}
