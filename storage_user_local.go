package profile

import (
	"github.com/0studio/cachemap"
	"github.com/0studio/storage_key"
	"time"
)

const (
	// 	SESSION_PLAYER_EXPIRE_SECONDS   = 60 * 60 // 60min, 从db加载user算起， 如果玩家不在线， 60分钟就清除cache
	NOT_LOGIN_PLAYER_EXPIRE_SECONDS = 5 * 60 // 未登录玩家缓存的User 有效期
)

type LocalUserStorage struct {
	cache cachemap.Uint64CacheMap
}

func NewLocalUserStorage() (localUserStorage LocalUserStorage) {
	localUserStorage = LocalUserStorage{cache: make(cachemap.Uint64CacheMap)}
	startCleaner(&localUserStorage)
	return
}

// user.UserStorage 接口实现
func (m LocalUserStorage) Get(uin key.KeyUint64, now time.Time) (user User, ok bool) {
	cacheObj, ok := m.cache.Get(uint64(uin), now)
	if !ok {
		return
	}
	user = cacheObj.(User)
	// log.Debug("get_user_from_process_cache", uin)
	return
}
func (m LocalUserStorage) Set(user *User, now time.Time) bool {
	// if isLoginUser(user) {
	// m.cache.Put(uint64(user.GetUin()), cachemap.NewCacheObject(*user, now, SESSION_PLAYER_EXPIRE_SECONDS))
	// } else {					//
	m.cache.Put(uint64(user.GetUin()), cachemap.NewCacheObject(*user, now, NOT_LOGIN_PLAYER_EXPIRE_SECONDS))
	// }
	return true
}
func (m LocalUserStorage) Add(user *User, now time.Time) bool {
	return m.Set(user, now)
}

func (m LocalUserStorage) MultiGet(keys key.KeyUint64List, now time.Time) (userMap UserMap, ok bool) {
	userMap = make(UserMap)
	var user User
	for _, uin := range keys {
		user, ok = m.Get(uin, now)
		if ok {
			userMap[uin] = user
		}
	}
	ok = true
	return
}
func (m LocalUserStorage) MultiUpdate(userMap UserMap, now time.Time) (ok bool) {
	for _, user := range userMap {
		m.Set(&user, now)
	}
	return true
}

func (m LocalUserStorage) Delete(uin key.KeyUint64) (ok bool) {
	// log.Debug("try_delete_user_from_process_cache", uin)
	return m.cache.Delete(uint64(uin))
}
func (m LocalUserStorage) MultiDelete(keys key.KeyUint64List) (ok bool) {
	for _, uin := range keys {
		m.Delete(uin)
	}
	return true
}
func (m LocalUserStorage) Len() int {
	return len(m.cache)
}
func (m LocalUserStorage) GetAll() (userMap UserMap) {
	// get all not outdate user
	userMap = make(UserMap)
	now := time.Now()
	var user User
	for uin, _ := range m.cache {
		userObj, ok := m.cache.Get(uin, now)
		if ok {
			user = userObj.(User)
			userMap[user.GetUin()] = user
		}
	}

	return
}
