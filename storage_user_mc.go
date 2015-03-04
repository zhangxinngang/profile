package profile

import (
	key "github.com/0studio/storage_key"
	"github.com/dropbox/godropbox/memcache"
	"strings"
	"time"
)

// 默认缓存过期时间

type MCUserStorage struct {
	expireSeconds uint32
	prefix        string
	client        memcache.Client
}

func NewMCUserStorage(client memcache.Client) MCUserStorage {
	return MCUserStorage{expireSeconds: 60 * 30, prefix: "u", client: client}
}
func (m MCUserStorage) getMCKey(Key key.KeyUint64) string {
	return strings.Join([]string{m.prefix, Key.ToString()}, "_")
}

func (m MCUserStorage) getRawKey(PrefixKey string) (keyUint64 key.KeyUint64) {
	keys := strings.Split(PrefixKey, "_")
	keyUint64.FromString(keys[len(keys)-1])
	return
}
func (m MCUserStorage) Get(uin key.KeyUint64, now time.Time) (user User, ok bool) {
	item := m.client.Get(m.getMCKey(key.KeyUint64(uin)))
	if item.Error() != nil || item.Status() != memcache.StatusNoError {
		return
	}
	byteData := item.Value()
	if byteData == nil {
		return
	}

	user.UnSerial(byteData)
	if uin != user.GetUin() {
		ok = false
		return
	}

	ok = true
	return
}
func (m MCUserStorage) Set(user *User, now time.Time) (ok bool) {
	item := memcache.Item{Key: m.getMCKey(user.GetUin()), Value: user.Serial(), Expiration: m.expireSeconds}
	response := m.client.Set(&item)
	return response.Error() == nil

}
func (m MCUserStorage) Add(user *User, now time.Time) bool {
	return m.Set(user, now)
}

func (m MCUserStorage) MultiGet(keys key.KeyUint64List, now time.Time) (userMap UserMap, ok bool) {
	prefixKeys := make([]string, len(keys))
	for idx, Key := range keys {
		prefixKeys[idx] = m.getMCKey(key.KeyUint64(Key))
	}
	itemMap := m.client.GetMulti(prefixKeys)

	userMap = make(UserMap)
	var user User
	for k, item := range itemMap {
		if len(item.Value()) == 0 {
			continue
		}
		if user.UnSerial(item.Value()) {
			userMap[m.getRawKey(k)] = user
		}
	}
	ok = true
	return
}
func (m MCUserStorage) MultiUpdate(userMap UserMap, now time.Time) (ok bool) {
	items := make([]memcache.Item, 0, len(userMap))
	itemsPointer := make([]*memcache.Item, 0, len(userMap))
	var idx int = 0
	for Key, user := range userMap {
		items = append(items, memcache.Item{Key: m.getMCKey(key.KeyUint64(Key)), Value: user.Serial(), Expiration: m.expireSeconds})
		itemsPointer = append(itemsPointer, &(items[idx]))
		idx++
	}
	responses := m.client.SetMulti(itemsPointer)
	for _, response := range responses {
		if response.Error() != nil {
			return false
		}
	}
	return true
}
func (m MCUserStorage) Delete(Key key.KeyUint64) (ok bool) {
	response := m.client.Delete(m.getMCKey(key.KeyUint64(Key)))
	return response.Error() == nil
}
func (m MCUserStorage) MultiDelete(keys key.KeyUint64List) (ok bool) {
	Prefixkeys := make([]string, len(keys), len(keys))
	for idx, Key := range keys {
		Prefixkeys[idx] = m.getMCKey(key.KeyUint64(Key))
	}
	responses := m.client.DeleteMulti(Prefixkeys)
	for _, res := range responses {
		if res.Error() != nil {
			return false
		}
	}
	return true
}
