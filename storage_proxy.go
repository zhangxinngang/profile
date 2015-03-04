package profile

import (
	key "github.com/0studio/storage_key"
	"time"
)

type UserStorageProxy struct {
	preferedStorage UserStorage
	backupStorage   UserStorage
}

func NewStorageProxy(prefered, backup UserStorage) UserStorageProxy {
	return UserStorageProxy{
		preferedStorage: prefered,
		backupStorage:   backup,
	}
}

func (this UserStorageProxy) Get(uin key.KeyUint64, now time.Time) (user User, ok bool) {
	user, ok = this.preferedStorage.Get(uin, now)
	if ok {
		return
	}
	user, ok = this.backupStorage.Get(uin, now)
	if !ok {
		return
	}
	this.preferedStorage.Set(&user, now)
	return
}

func (this UserStorageProxy) Set(user *User, now time.Time) (ok bool) {
	ok = this.backupStorage.Set(user, now)
	if !ok {
		return ok
	}
	ok = this.preferedStorage.Set(user, now)
	return
}

func (this UserStorageProxy) Add(user *User, now time.Time) (ok bool) {
	ok = this.backupStorage.Add(user, now)
	if !ok {
		return ok
	}
	ok = this.preferedStorage.Add(user, now)
	return
}

func (this UserStorageProxy) MultiGet(keys key.KeyUint64List, now time.Time) (userMap UserMap, ok bool) {
	userMap, ok = this.preferedStorage.MultiGet(keys, now)
	missedKeyCount := 0
	for _, uin := range keys {
		if _, find := userMap[uin]; !find {
			missedKeyCount++
		}
	}
	if missedKeyCount == 0 {
		return
	}

	missedKeys := make(key.KeyUint64List, missedKeyCount)
	i := 0
	for _, uin := range keys {
		if _, find := userMap[uin]; !find {
			missedKeys[i] = uin
			i++
		}
	}

	missedMap, ok := this.backupStorage.MultiGet(missedKeys, now)
	if !ok {
		return
	}
	this.preferedStorage.MultiUpdate(missedMap, now)
	for k, v := range missedMap {
		userMap[k] = v
	}
	return
}

func (this UserStorageProxy) MultiUpdate(userMap UserMap, now time.Time) (ok bool) {
	ok = this.backupStorage.MultiUpdate(userMap, now)
	if !ok {
		return
	}
	ok = this.preferedStorage.MultiUpdate(userMap, now)
	return
}

func (this UserStorageProxy) Delete(uin key.KeyUint64) (ok bool) {
	ok = this.backupStorage.Delete(uin)
	if !ok {
		return
	}
	ok = this.preferedStorage.Delete(uin)
	return
}

func (this UserStorageProxy) MultiDelete(keys key.KeyUint64List) (ok bool) {
	ok = this.backupStorage.MultiDelete(keys)
	if !ok {
		return
	}
	ok = this.preferedStorage.MultiDelete(keys)
	return

}
