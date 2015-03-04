package profile

import (
	"github.com/0studio/storage_key"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"time"
)

func TestLocalUserStorage(t *testing.T) {
	Uin := key.KeyUint64(rand.Int63())
	now := time.Now()
	user := User{}
	user.SetUin(key.KeyUint64(Uin))

	store := NewLocalUserStorage()

	ok := store.Set(&user, now)
	assert.True(t, ok)

	userRet, ok := store.Get(Uin, now)
	assert.True(t, ok)
	assert.Equal(t, userRet.GetUin(), user.GetUin())

	ok = store.Delete(Uin)
	assert.True(t, ok)

	userRet, ok = store.Get(Uin, now)
	assert.False(t, ok)
}
func TestLocalUserStorageMulti(t *testing.T) {
	Uin := key.KeyUint64(rand.Int63())
	Uin2 := Uin + 1
	now := time.Now()
	user := User{}
	user.SetUin(Uin)
	user2 := User{}
	user2.SetUin(Uin2)

	userMap := NewUserMap()
	userMap[user.GetUin()] = user
	userMap[user2.GetUin()] = user2

	store := NewLocalUserStorage()

	ok := store.MultiUpdate(userMap, now)
	assert.True(t, ok)

	userMapRet, ok := store.MultiGet(key.KeyUint64List{user.GetUin(), user2.GetUin()}, now)
	assert.True(t, ok)
	assert.Equal(t, 2, len(userMapRet))

	ok = store.MultiDelete(key.KeyUint64List{user.GetUin(), user2.GetUin()})
	assert.True(t, ok)
	// after delete
	userMapRet, ok = store.MultiGet(key.KeyUint64List{user.GetUin(), user2.GetUin()}, now)
	assert.True(t, ok)
	assert.Equal(t, 0, len(userMapRet))

}
