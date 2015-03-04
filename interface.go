package profile

import (
	"github.com/0studio/storage_key"
	"time"
)

type UserStorage interface {
	Get(uin key.KeyUint64, now time.Time) (user User, ok bool)
	Set(user *User, now time.Time) (ok bool)
	Add(user *User, now time.Time) bool
	MultiGet(keys key.KeyUint64List, now time.Time) (userMap UserMap, ok bool)
	MultiUpdate(userMap UserMap, now time.Time) (ok bool)
	Delete(uin key.KeyUint64) (ok bool)
	MultiDelete(keys key.KeyUint64List) (ok bool)
}
type UserService interface {
	UserStorage
	GetNewUin() uint64
	Offline(user *User, now time.Time) bool
	Auth(accountId, accountName string, platform, channel, server uint64) (user User, ok bool)
}
type PayOrderService interface {
	Get(uin key.KeyUint64, orderId string) (payOrder PayOrder, ok bool)
	MultiGet(uin key.KeyUint64, orderIdList []string) (payOrderMap PayOrderMap, ok bool)
	GetAllUnhandledOrder(uin key.KeyUint64) (list []PayOrder)
	UpdateStatus(order PayOrder) bool
	Set(order PayOrder) bool
	Add(order PayOrder) bool
}
