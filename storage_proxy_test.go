package profile

import (
	"database/sql"
	"github.com/0studio/databasetemplate"
	"github.com/0studio/storage_key"
	"github.com/dropbox/godropbox/memcache"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"time"
)

func setupTest() {

}

func getMockDB() (db *sql.DB) {
	db, _ = databasetemplate.NewDBInstance(databasetemplate.DBConfig{
		Host: "127.0.0.1",
		User: "th_dev",
		Pass: "th_devpass",
		Name: "test",
	}, false)

	return
}
func TestUserStorageProxyGet(t *testing.T) {
	setupTest()
	Uin := key.KeyUint64(rand.Int63())
	now := time.Now()
	user := User{}
	user.SetUin(Uin)

	localStore := NewLocalUserStorage()
	mcStore := NewMCUserStorage(memcache.NewMockClient())
	dbStore := InitDBUserStorage(getMockDB(), 1, 1, 1)
	proxy1 := NewStorageProxy(localStore, mcStore)
	proxy := NewStorageProxy(proxy1, dbStore)

	ok := proxy.Add(&user, now)
	assert.True(t, ok)

	userRet, ok := proxy.Get(user.GetUin(), now)
	assert.True(t, ok)
	assert.Equal(t, userRet.GetUin(), user.GetUin())

	userRet, ok = proxy1.Get(user.GetUin(), now)
	assert.True(t, ok)
	assert.Equal(t, userRet.GetUin(), user.GetUin())

	userRet, ok = localStore.Get(user.GetUin(), now)
	assert.True(t, ok)
	assert.Equal(t, userRet.GetUin(), user.GetUin())

	userRet, ok = mcStore.Get(user.GetUin(), now)
	assert.True(t, ok)
	assert.Equal(t, userRet.GetUin(), user.GetUin())

	userRet, ok = dbStore.Get(user.GetUin(), now)
	assert.True(t, ok)
	assert.Equal(t, userRet.GetUin(), user.GetUin())

	ok = proxy.Delete(user.GetUin())
	assert.True(t, ok)

}

func TestUserStorageProxyDelete(t *testing.T) {
	setupTest()
	Uin := key.KeyUint64(rand.Int63())
	now := time.Now()
	user := User{}
	user.SetUin(Uin)

	localStore := NewLocalUserStorage()
	mcStore := NewMCUserStorage(memcache.NewMockClient())
	dbStore := InitDBUserStorage(getMockDB(), 1, 1, 1)
	proxy1 := NewStorageProxy(localStore, mcStore)
	proxy := NewStorageProxy(proxy1, dbStore)

	ok := proxy.Add(&user, now)
	assert.True(t, ok)

	userRet, ok := proxy.Get(user.GetUin(), now)
	assert.True(t, ok)
	assert.Equal(t, userRet.GetUin(), user.GetUin())

	userRet, ok = proxy1.Get(user.GetUin(), now)
	assert.True(t, ok)
	assert.Equal(t, userRet.GetUin(), user.GetUin())

	ok = proxy.Delete(user.GetUin())
	assert.True(t, ok)

	userRet, ok = proxy.Get(user.GetUin(), now)
	assert.False(t, ok)
	userRet, ok = proxy1.Get(user.GetUin(), now)
	assert.False(t, ok)

	userRet, ok = localStore.Get(user.GetUin(), now)
	assert.False(t, ok)

	userRet, ok = mcStore.Get(user.GetUin(), now)
	assert.False(t, ok)

	userRet, ok = dbStore.Get(user.GetUin(), now)
	assert.False(t, ok)
}

func TestUserStorageProxyMultiGet(t *testing.T) {
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

	localStore := NewLocalUserStorage()
	mcStore := NewMCUserStorage(memcache.NewMockClient())
	dbStore := InitDBUserStorage(getMockDB(), 1, 1, 1)
	proxy1 := NewStorageProxy(localStore, mcStore)
	proxy := NewStorageProxy(proxy1, dbStore)

	ok := proxy.Add(&user, now)
	assert.True(t, ok)
	ok = proxy.Add(&user2, now)
	assert.True(t, ok)

	userMapRet, ok := proxy.MultiGet(key.KeyUint64List{user.GetUin(), user2.GetUin()}, now)
	assert.True(t, ok)
	assert.Equal(t, len(userMapRet), 2)

	userMapRet, ok = proxy1.MultiGet(key.KeyUint64List{user.GetUin(), user2.GetUin()}, now)
	assert.True(t, ok)
	assert.Equal(t, len(userMapRet), 2)

	userMapRet, ok = localStore.MultiGet(key.KeyUint64List{user.GetUin(), user2.GetUin()}, now)
	assert.True(t, ok)
	assert.Equal(t, len(userMapRet), 2)

	userMapRet, ok = mcStore.MultiGet(key.KeyUint64List{user.GetUin(), user2.GetUin()}, now)
	assert.True(t, ok)
	assert.Equal(t, len(userMapRet), 2)

	userMapRet, ok = dbStore.MultiGet(key.KeyUint64List{user.GetUin(), user2.GetUin()}, now)
	assert.True(t, ok)
	assert.Equal(t, len(userMapRet), 2)

	ok = proxy.MultiDelete(key.KeyUint64List{user.GetUin(), user2.GetUin()})
	assert.True(t, ok)
}
