package profile

import (
	key "github.com/0studio/storage_key"
)

type UserMap map[key.KeyUint64]User

func (this UserMap) Len() int {
	return len(this)
}

func (this UserMap) Delete(id key.KeyUint64) bool {
	delete(this, id)
	return true
}
func (this UserMap) MultiGet(keys key.KeyUint64List) (userMap UserMap, ok bool) {
	userMap = make(UserMap)
	var user User
	for _, uin := range keys {
		user, ok = this[uin]
		if ok {
			userMap[uin] = user
		}
	}
	ok = true
	return
}

func (this UserMap) MultiDelete(keys key.KeyUint64List) (ok bool) {
	for _, uin := range keys {
		this.Delete(uin)
	}
	return true
}

func NewUserMap() (this UserMap) {
	this = make(UserMap)
	return
}
