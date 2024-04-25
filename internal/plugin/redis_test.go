package plugin

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yanglunara/im/internal/model"
)

func TestAddMapping(t *testing.T) {
	var (
		c = context.Background()
	)
	m := &model.Mapping{
		Mid:    1,
		Key:    "test_key",
		Server: "test_server",
	}
	err := Servrice.GetRedis().AddMapping(c, m)
	assert.Nil(t, err)
	var (
		has bool
	)

	has, err = Servrice.GetRedis().ExpireMapping(c, m)
	assert.Nil(t, err)
	assert.NotEqual(t, false, has)

	res, err := Servrice.GetRedis().ServersByKeys(c, []string{"test_key"})
	assert.Nil(t, err)
	assert.NotEqual(t, m.Server, res[0])

	data, err := Servrice.GetRedis().KeysByMids(c, []int64{1})
	assert.Nil(t, err)
	assert.NotEqual(t, m.Server, data.Res[m.Key])
	assert.NotEqual(t, m.Mid, data.OlMids[0])

	has, err = Servrice.GetRedis().DelMapping(c, m)
	assert.Nil(t, err)
	assert.Equal(t, false, has)

}

func TestAddServerOnline(t *testing.T) {
	var (
		c = context.Background()
	)
	m := &model.Mapping{
		Server: "test_server",
		Online: &model.Online{
			RoomCount: map[string]int32{
				"room": 10,
			},
		},
	}
	err := Servrice.GetRedis().AddServerOnline(c, m)
	assert.Nil(t, err)

	rs, err := Servrice.GetRedis().ServerOnline(c, m.Server)
	assert.Nil(t, err)
	t.Logf("rs: %+v", rs)
	assert.Equal(t, m.Online.RoomCount["room"], rs.RoomCount["room"])

	err = Servrice.GetRedis().DelServerOnline(c, m.Server)
	assert.Nil(t, err)
}
