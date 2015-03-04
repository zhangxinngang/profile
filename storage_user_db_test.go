package profile

import (
	key "github.com/0studio/storage_key"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"time"
)

func TestDBUserStorage(t *testing.T) {
	setupTest()
	Uin := key.KeyUint64(rand.Int63())
	now := time.Now()
	user := User{}
	user.SetUin(Uin)

	store := InitDBUserStorage(getMockDB(), 1, 1, 1)

	ok := store.Add(&user, now)
	assert.True(t, ok)

	userRet, ok := store.Get(user.GetUin(), now)
	assert.True(t, ok)
	assert.Equal(t, userRet.GetUin(), user.GetUin())

	user.SetAccountId("222")
	user.SetAccountName("222")
	user.SetGender(1)
	user.SetChannel(1)
	user.SetDeviceModel("de")
	user.SetLevel(100)
	user.SetVip(100)
	user.SetOs(1)
	user.SetOsVersion(1)
	user.SetServer(1)
	ok = store.Set(&user, now)
	assert.True(t, ok)

	userRet, ok = store.Get(user.GetUin(), now)
	assert.True(t, ok)
	assert.Equal(t, userRet.GetUin(), user.GetUin())
	assert.Equal(t, userRet.GetAccountId(), user.GetAccountId())
	assert.Equal(t, userRet.GetAccountName(), user.GetAccountName())
	assert.Equal(t, userRet.GetGender(), user.GetGender())
	assert.Equal(t, userRet.GetChannel(), user.GetChannel())
	assert.Equal(t, userRet.GetDeviceModel(), user.GetDeviceModel())
	assert.Equal(t, userRet.GetLevel(), user.GetLevel())
	assert.Equal(t, userRet.GetOs(), user.GetOs())
	assert.Equal(t, userRet.GetOsVersion(), user.GetOsVersion())
	assert.Equal(t, userRet.GetServer(), user.GetServer())

	ok = store.Delete(user.GetUin())
	assert.True(t, ok)

	userRet, ok = store.Get(user.GetUin(), now)
	assert.False(t, ok)
}
func TestDBUserStorageMulti(t *testing.T) {
	setupTest()
	Uin := key.KeyUint64(rand.Int63())
	now := time.Now()
	user := User{}
	user2 := User{}
	user.SetUin(Uin)
	user2.SetUin(Uin + 1)

	userMap := make(UserMap)
	userMap[user.GetUin()] = user
	userMap[user2.GetUin()] = user2

	store := InitDBUserStorage(getMockDB(), 1, 1, 1)

	ok := store.Add(&user, now)
	assert.True(t, ok)
	ok = store.Add(&user2, now)
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
