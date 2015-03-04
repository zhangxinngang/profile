package profile

import (
	"bytes"
	"database/sql"
	"fmt"
	"github.com/0studio/databasetemplate"
	"github.com/0studio/idgen"
	key "github.com/0studio/storage_key"
	log "github.com/cihub/seelog"
	"github.com/dropbox/godropbox/errors"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
	"sync"
	"time"
)

type DBUserStorage struct {
	databasetemplate.GenericDaoImpl
	Platform uint64
	Server   uint64
	Process  uint64
	IdGen    *idgen.IdGen
}

var initUserDaoOnce sync.Once
var userDao DBUserStorage

func GetUserDao() *DBUserStorage {
	return &userDao
}
func InitDBUserStorage(db *sql.DB, platform uint64, server uint64, process uint64) *DBUserStorage {
	done := make(chan bool)
	go func() {
		initUserDaoOnce.Do(func() {
			dbTemplate := &databasetemplate.DatabaseTemplateImpl{db}
			userDao = DBUserStorage{GenericDaoImpl: databasetemplate.GenericDaoImpl{dbTemplate}, Platform: platform, Server: server, Process: process}
			userDao.CreateTable()
			userDao.IdGen = idgen.NewIdgen(PLATFORM_BIT, platform, SERVER_BIT, server, SYSTYPE_BIT, process, 0)
			maxUin, err := userDao.GetMaxUin(
				userDao.IdGen.GetSequenceMask(),
				platform,
				userDao.IdGen.GetServerMask(), userDao.IdGen.GetServerShift(), server,
				userDao.IdGen.GetSysTypeMask(), userDao.IdGen.GetSysTypeShift(), process)
			log.Info("max_user_uin", maxUin)
			if err == nil {
				userDao.IdGen.SetSequence(userDao.IdGen.GetIdSequence(maxUin))
				go userDao.IdGen.Recv()
			}
		})

		done <- true
	}()
	<-done
	return &userDao
}

func (this *DBUserStorage) GetMaxUin(seqMask uint64, platform uint64,
	serverMask uint64, serverShift uint64, serverid uint64,
	processIdxMask uint64, processIdxShift uint64, processIdx uint64) (MaxUserId uint64, err error) {
	sql := "select ifnull(Uin,0) from player where platform=? and ( (Uin&?) >>?)=? and ( (Uin&?) >>?)=? order by (Uin&?) desc limit 0,1"
	object, err := this.DatabaseTemplate.QueryObject(sql, this.mapRowUserId,
		platform, serverMask, serverShift, serverid, processIdxMask, processIdxShift, processIdx, seqMask,
	)

	if object != nil {
		MaxUserId = object.(uint64)
	}
	return
}

func (this *DBUserStorage) GetNewUin() (id uint64) {
	id = this.IdGen.GetNewId()
	log.Info("newid", id)
	return
}

func (this *DBUserStorage) CreateTable() {
	query := `CREATE TABLE if not exists player (
	  Uin bigint(20) NOT NULL COMMENT '唯一ID',
	  AccountId varchar(64) NOT NULL DEFAULT '' COMMENT '设备ID',
	  AccountName varchar(64) NOT NULL DEFAULT '' COMMENT '设备ID',
	  gender tinyint NOT NULL DEFAULT 0 COMMENT '性别',
	  platform smallint NOT NULL DEFAULT 0 COMMENT '平台',
	  serverid smallint NOT NULL DEFAULT 0 COMMENT 'serverid',
	  channel int NOT NULL DEFAULT 0 COMMENT '渠道',
	  uuid varchar(64) NOT NULL DEFAULT '' COMMENT '设备ID',
	  os tinyint NOT NULL DEFAULT 0 COMMENT '系统',
	  osVersion int NULL DEFAULT 0 COMMENT '系统',
      deviceModel varchar(64) NOT NULL DEFAULT '' COMMENT '设备型号  iPhone 4S、iPhone 5S', 
      Vip tinyINT NOT NULL DEFAULT 0 COMMENT 'VIP',
  	  Level smallint NOT NULL DEFAULT 0 COMMENT '当前等级',
	  createTime timestamp NOT NULL DEFAULT '0000-00-00 00:00:00' COMMENT '注册时间时间',
	  PRIMARY KEY (Uin),
	  KEY account_idx (AccountId,AccountName)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8`
	err := this.DatabaseTemplate.Exec(query)
	if err != nil {
		log.Warn(errors.Wrap(err, "create table player error!!!"))
	}
}

func (this *DBUserStorage) Get(key key.KeyUint64, now time.Time) (user User, ok bool) {
	query := `select Uin,AccountId,AccountName,gender,platform,channel,serverid,uuid,os,osVersion,deviceModel,Vip,Level,createTime from player where Uin=?`
	var obj interface{}
	var err error
	obj, err = this.DatabaseTemplate.QueryObject(query, this.mapRow, key)
	if err != nil || obj == nil {
		return
	}

	user = obj.(User)
	ok = true
	return
}
func (this *DBUserStorage) MultiGet(keys key.KeyUint64List, now time.Time) (userMap UserMap, ok bool) {
	userMap = NewUserMap()
	if len(keys) == 0 {
		ok = true
		return
	}

	query := fmt.Sprintf(`select Uin,AccountId,AccountName,gender,platform,channel,serverid,uuid,os,osVersion,deviceModel,Vip,Level,createTime from player where Uin in (%s)`, keys.Join(","))
	var obj interface{}
	var err error
	var arr []interface{}
	var user User
	arr, err = this.DatabaseTemplate.QueryArray(query, this.mapRow)
	if err != nil {
		return
	}
	for _, obj = range arr {
		user = obj.(User)
		userMap[user.GetUin()] = user

	}
	ok = true

	return
}

func (this *DBUserStorage) Auth(accountId, accountName string, platform, channel, server uint64) (user User, ok bool) {
	query := `select Uin,AccountId,AccountName,gender,platform,channel,serverid,uuid,os,osVersion,deviceModel,Vip,Level,createTime from player where AccountId=? and AccountName=? and platform=? and channel=? and serverid=?`
	var obj interface{}
	var err error
	obj, err = this.DatabaseTemplate.QueryObject(query, this.mapRow, accountId, accountName, platform, channel, server)
	if err != nil || obj == nil {
		return
	}
	ok = true
	user = obj.(User)
	return

}

func (this *DBUserStorage) Add(user *User, now time.Time) bool {
	Sql := `insert into player(Uin,AccountId,AccountName,gender,platform,channel,serverid,uuid,os,osVersion,deviceModel,Vip,Level,createTime )values(?,?,?,?,?,?,?,?,?,?,?,?,?,?)`
	log.Info("add_user_sql", user.GetUin())
	err := this.DatabaseTemplate.Exec(Sql,
		user.GetUin(),
		user.GetAccountId(),
		user.GetAccountName(),
		user.GetGender(),
		user.GetPlatform(),
		user.GetChannel(),
		user.GetServer(),
		user.GetUuid(),
		user.GetOs(),
		user.GetOsVersion(),
		user.GetDeviceModel(),
		user.GetVip(),
		user.GetLevel(),
		user.GetCreateTime(),
	)
	if err != nil {
		log.Info("useradd", err)
		fmt.Println("useradd", err)
	}

	return err == nil
}

func (this *DBUserStorage) Set(user *User, now time.Time) bool {
	if !user.IsDirty() {
		return true
	}
	var isFirstField bool = true
	var updateBuffer bytes.Buffer
	if user.Flag.IsPosTrue(USER_FLAG_POS_ACCOUNTID) {
		if !isFirstField {
			updateBuffer.WriteString(",")
		}
		isFirstField = false

		updateBuffer.WriteString("AccountId='")
		updateBuffer.WriteString(user.GetAccountId())
		updateBuffer.WriteString("'")
	}
	if user.Flag.IsPosTrue(USER_FLAG_POS_ACCOUNTNAME) {
		if !isFirstField {
			updateBuffer.WriteString(",")
		}
		isFirstField = false
		updateBuffer.WriteString("AccountName='")
		updateBuffer.WriteString(user.GetAccountName())
		updateBuffer.WriteString("'")
	}
	if user.Flag.IsPosTrue(USER_FLAG_POS_GENDER) {
		if !isFirstField {
			updateBuffer.WriteString(",")
		}
		isFirstField = false
		updateBuffer.WriteString("gender=")
		updateBuffer.WriteString(strconv.Itoa(int(user.GetGender())))
	}
	if user.Flag.IsPosTrue(USER_FLAG_POS_VIP) {
		if !isFirstField {
			updateBuffer.WriteString(",")
		}
		isFirstField = false
		updateBuffer.WriteString("vip=")
		updateBuffer.WriteString(strconv.Itoa(int(user.GetVip())))
	}
	if user.Flag.IsPosTrue(USER_FLAG_POS_LEVEL) {
		if !isFirstField {
			updateBuffer.WriteString(",")
		}
		isFirstField = false
		updateBuffer.WriteString("level=")
		updateBuffer.WriteString(strconv.Itoa(int(user.GetLevel())))
	}

	if user.Flag.IsPosTrue(USER_FLAG_POS_PLATFORM) {
		if !isFirstField {
			updateBuffer.WriteString(",")
		}
		isFirstField = false
		updateBuffer.WriteString("platform=")
		updateBuffer.WriteString(strconv.FormatUint(user.GetPlatform(), 10))
	}

	if user.Flag.IsPosTrue(USER_FLAG_POS_CHANNEL) {
		if !isFirstField {
			updateBuffer.WriteString(",")
		}
		isFirstField = false
		updateBuffer.WriteString("channel=")
		updateBuffer.WriteString(strconv.FormatUint(user.GetChannel(), 10))
	}

	if user.Flag.IsPosTrue(USER_FLAG_POS_SERVER) {
		if !isFirstField {
			updateBuffer.WriteString(",")
		}
		isFirstField = false
		updateBuffer.WriteString("serverid=")
		updateBuffer.WriteString(strconv.FormatUint(user.GetServer(), 10))
	}

	if user.Flag.IsPosTrue(USER_FLAG_POS_UUID) {
		if !isFirstField {
			updateBuffer.WriteString(",")
		}
		isFirstField = false
		updateBuffer.WriteString("uuid='")
		updateBuffer.WriteString(user.GetUuid())
		updateBuffer.WriteString("'")
	}

	if user.Flag.IsPosTrue(USER_FLAG_POS_OS) {
		if !isFirstField {
			updateBuffer.WriteString(",")
		}
		isFirstField = false
		updateBuffer.WriteString("os=")
		updateBuffer.WriteString(strconv.Itoa(int(user.GetOs())))
	}
	if user.Flag.IsPosTrue(USER_FLAG_POS_OS_VERSION) {
		if !isFirstField {
			updateBuffer.WriteString(",")
		}
		isFirstField = false
		updateBuffer.WriteString("osVersion=")
		updateBuffer.WriteString(strconv.Itoa(int(user.GetOsVersion())))
	}
	if user.Flag.IsPosTrue(USER_FLAG_POS_DEVICE_MODEL) {
		if !isFirstField {
			updateBuffer.WriteString(",")
		}
		isFirstField = false
		updateBuffer.WriteString("deviceModel='")
		updateBuffer.WriteString(user.GetDeviceModel())
		updateBuffer.WriteString("'")
	}
	if user.Flag.IsPosTrue(USER_FLAG_POS_CREATE_TIME) {
		if !isFirstField {
			updateBuffer.WriteString(",")
		}
		isFirstField = false
		updateBuffer.WriteString("createTime=")
		updateBuffer.WriteString(user.GetCreateTime().Format("20060102150405"))
	}

	// var timeStr = time.Now().Format("20060102150405")
	var sql string = fmt.Sprintf("update  player set %s where Uin=?", updateBuffer.String())
	fmt.Println(sql)

	err := this.DatabaseTemplate.Exec(sql, user.GetUin())
	if err != nil {
		return false
	}

	user.ClearFlag()
	return true

}

func (this *DBUserStorage) MultiUpdate(userMap UserMap, now time.Time) (ok bool) {
	for _, user := range userMap {
		this.Set(&user, now)
	}
	ok = true
	return
}

func (this *DBUserStorage) MultiDelete(keys key.KeyUint64List) (ok bool) {
	if len(keys) == 0 {
		return true
	}
	err := this.DatabaseTemplate.Exec(fmt.Sprintf(`delete from player where Uin in (%s)`, keys.Join(",")))
	return err == nil
}
func (this *DBUserStorage) Delete(key key.KeyUint64) (ok bool) {
	query := `delete from player where Uin=?`
	var err error
	err = this.DatabaseTemplate.Exec(query, key)

	return err == nil
}

func (this *DBUserStorage) Truncate() {
	this.DatabaseTemplate.Exec("truncate table player")
}

func (this *DBUserStorage) mapRow(resultSet *sql.Rows) (interface{}, error) {
	var (
		uin         key.KeyUint64
		accountId   string
		accountName string
		gender      int32
		platform    uint64
		channel     uint64 // 1:91 2:360  拇指玩 豌豆荚
		server      uint64
		uuid        string // 设备id
		os          int32  // 1:IOS 2:ANDROID 13:WP8
		osVersion   int32  //  系统版本  610 712
		deviceModel string //  设备型号  iPhone 4S、iPhone 5S
		createTime  time.Time
		vip         int32
		level       int32
	)
	user := User{}
	err := resultSet.Scan(
		&uin,
		&accountId,
		&accountName,
		&gender,
		&platform,
		&channel,
		&server,
		&uuid,
		&os,
		&osVersion,
		&deviceModel,
		&vip,
		&level,
		&createTime,
	)
	if err != nil {
		return nil, err
	}
	user.SetUin(uin)
	user.SetAccountId(accountId)
	user.SetAccountName(accountName)
	user.SetGender(gender)
	user.SetPlatform(platform)
	user.SetChannel(channel)
	user.SetServer(server)
	user.SetUuid(uuid)
	user.SetOs(os)
	user.SetOsVersion(osVersion)
	user.SetDeviceModel(deviceModel)
	user.SetCreateTime(createTime)
	user.SetVip(vip)
	user.SetLevel(level)
	user.ClearFlag()

	return user, nil
}

func (this *DBUserStorage) mapRowUserId(resultSet *sql.Rows) (interface{}, error) {
	var uin uint64
	err := resultSet.Scan(&uin)
	if err != nil {
		return nil, err
	}

	return uin, nil
}
