package profile

import (
	key "github.com/0studio/storage_key"
	"github.com/dropbox/godropbox/memcache"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"time"
)

func TestMCUserStorage(t *testing.T) {
	setupTest()
	Uin := key.KeyUint64(rand.Int63())
	now := time.Now()
	user := User{}
	user.SetUin(Uin)

	store := NewMCUserStorage(memcache.NewMockClient())

	ok := store.Set(&user, now)
	assert.True(t, ok)

	playerRet, ok := store.Get(Uin, now)
	assert.True(t, ok)
	assert.Equal(t, playerRet.GetUin(), user.GetUin())

	ok = store.Delete(Uin)
	assert.True(t, ok)

	playerRet, ok = store.Get(Uin, now)
	assert.False(t, ok)
}
func TestMCUserStorageMulti(t *testing.T) {
	setupTest()
	Uin := key.KeyUint64(rand.Int63())
	Uin2 := Uin + 1
	now := time.Now()
	user := User{}
	user.SetUin(Uin)
	user2 := User{}
	user2.SetUin(Uin2)

	playerMap := NewUserMap()
	playerMap[user.GetUin()] = user
	playerMap[user2.GetUin()] = user2

	store := NewMCUserStorage(memcache.NewMockClient())

	ok := store.MultiUpdate(playerMap, now)
	assert.True(t, ok)

	playerMapRet, ok := store.MultiGet(key.KeyUint64List{user.GetUin(), user2.GetUin()}, now)
	assert.True(t, ok)
	assert.Equal(t, 2, len(playerMapRet))

	ok = store.MultiDelete(key.KeyUint64List{user.GetUin(), user2.GetUin()})
	assert.True(t, ok)
	// after delete
	playerMapRet, ok = store.MultiGet(key.KeyUint64List{user.GetUin(), user2.GetUin()}, now)
	assert.True(t, ok)
	assert.Equal(t, 0, len(playerMapRet))

}
