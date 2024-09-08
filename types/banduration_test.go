package types_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/supabase-community/auth-go/types"
)

func TestBanDuration(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	// Test BanDurationNone
	bd := types.BanDurationNone()
	assert.Equal("none", bd.String())

	// Test BanDurationTime
	bd = types.BanDurationTime(time.Hour)
	assert.Equal("1h0m0s", bd.String())

	bd = types.BanDurationTime(time.Minute)
	assert.Equal("1m0s", bd.String())

	bd = types.BanDurationTime(time.Second)
	assert.Equal("1s", bd.String())

	bd = types.BanDurationTime(time.Millisecond)
	assert.Equal("1ms", bd.String())

	bd = types.BanDurationTime((26 * time.Hour) + (15 * time.Minute))
	assert.Equal("26h15m0s", bd.String())

	// Test MarshalJSON
	bd = types.BanDurationNone()
	b, err := json.Marshal(bd)
	require.NoError(err)
	assert.Equal(`"none"`, string(b))

	bd = types.BanDurationTime(time.Second)
	b, err = json.Marshal(bd)
	require.NoError(err)
	assert.Equal(`"1s"`, string(b))

	// Test UnmarshalJSON
	b = []byte(`"none"`)
	err = json.Unmarshal(b, &bd)
	require.NoError(err)
	assert.Nil(bd.Value())
	assert.Equal("none", bd.String())

	b = []byte(`"1s"`)
	err = json.Unmarshal(b, &bd)
	require.NoError(err)
	assert.Equal(time.Second, *bd.Value())
	assert.Equal("1s", bd.String())

	b = []byte(`"not what it should be"`)
	err = json.Unmarshal(b, &bd)
	assert.Error(err)
}
