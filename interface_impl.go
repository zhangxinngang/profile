package profile

import (
	"database/sql"
	"github.com/0studio/storage_key"
	"github.com/dropbox/godropbox/memcache"
	"time"
)

func GetUserService(db *sql.DB, memcacheClient memcache.Client, platform uint64, server uint64, process uint64) UserService {
	return newUserServiceImpl(db, memcacheClient, platform, server, process)
}

type UserServiceImpl struct {
	local LocalUserStorage
	mc    MCUserStorage
	db    *DBUserStorage
}

func (this UserServiceImpl) GetNewUin() uint64 {
	return this.db.GetNewUin()
}
func (this UserServiceImpl) Offline(user *User, now time.Time) (ok bool) {
	this.Set(user, now)
	ok = this.local.Delete(user.GetUin()) // clear process cache
	return
}

func (this UserServiceImpl) Auth(accountId, accountName string, platform, channel, server uint64) (user User, ok bool) {
	return this.db.Auth(accountId, accountName, platform, channel, server)
}
func (this UserServiceImpl) Get(uin key.KeyUint64, now time.Time) (user User, ok bool) {
	return this.GetProxy().Get(uin, now)
}

func (this UserServiceImpl) Set(user *User, now time.Time) (ok bool) {
	return this.GetProxy().Set(user, now)
}

func (this UserServiceImpl) Add(user *User, now time.Time) bool {
	return this.GetProxy().Add(user, now)
}
func (this UserServiceImpl) MultiGet(keys key.KeyUint64List, now time.Time) (userMap UserMap, ok bool) {
	return this.GetProxy().MultiGet(keys, now)
}

func (this UserServiceImpl) MultiUpdate(userMap UserMap, now time.Time) (ok bool) {
	return this.GetProxy().MultiUpdate(userMap, now)
}
func (this UserServiceImpl) MultiDelete(keys key.KeyUint64List) (ok bool) {
	return this.GetProxy().MultiDelete(keys)
}
func (this UserServiceImpl) Delete(uin key.KeyUint64) (ok bool) {
	return this.GetProxy().Delete(uin)
}

//
//
//

func (this UserServiceImpl) getLocalMCProxy() UserStorage {
	return NewStorageProxy(this.local, this.mc)
}
func (this UserServiceImpl) GetProxy() UserStorage {
	return NewStorageProxy(this.local, NewStorageProxy(this.mc, this.db))
}
func (this UserServiceImpl) GetMCDBProxy() UserStorage {
	return NewStorageProxy(this.mc, this.db)
}

func newUserServiceImpl(dbInstance *sql.DB, memcacheClient memcache.Client, platform uint64, server uint64, process uint64) (this UserServiceImpl) {
	this = UserServiceImpl{
		local: NewLocalUserStorage(),
		mc:    NewMCUserStorage(memcacheClient),
		db:    InitDBUserStorage(dbInstance, platform, server, process),
	}
	return

}
