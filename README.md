### profile 基本信息###
```
userService:= GetUserService(db *sql.DB, memcacheClient memcache.Client, platform uint64, server uint64, process uint64)
newUin=userService.GetNewUin()
user=userService.Get(uin,now)
user.SetGender(1)
userService.Set(&user,now)
userService.Offline(&user,now)
```


### pay order 定单###
```
InitDBPayOrderStorage(db *sql.DB) *DBPayOrderStorage 
payOrderService:=InitDBPayOrderStorage(db *sql.DB) *DBPayOrderStorage 
payOrderService.Get(uin,orderId)

```