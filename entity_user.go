package profile

import (
	"github.com/0studio/bit"
	"github.com/0studio/storage_key"
	"github.com/gogo/protobuf/proto"
	"time"
)

type User struct {
	Flag        bit.BitInt
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
}

func (this *User) SetUin(value key.KeyUint64) {
	this.uin = value
}
func (this User) GetUin() key.KeyUint64 {
	return this.uin
}
func (this *User) SetAccountId(value string) {
	if value != this.GetAccountId() {
		this.Flag.SetFlag(USER_FLAG_POS_ACCOUNTID)
		this.accountId = value
	}

}
func (this User) GetAccountId() string {
	return this.accountId
}
func (this *User) SetAccountName(value string) {
	if this.GetAccountName() != value {
		this.Flag.SetFlag(USER_FLAG_POS_ACCOUNTNAME)
		this.accountName = value
	}

}
func (this User) GetAccountName() string {
	return this.accountName
}
func (this *User) SetGender(value int32) {
	if this.GetGender() != value {
		this.Flag.SetFlag(USER_FLAG_POS_GENDER)
		this.gender = value
	}

}
func (this User) GetGender() int32 {
	return this.gender
}
func (this *User) SetPlatform(value uint64) {
	if this.GetPlatform() != value {
		this.Flag.SetFlag(USER_FLAG_POS_PLATFORM)
		this.platform = value
	}
}
func (this User) GetPlatform() uint64 {
	return this.platform
}
func (this *User) SetChannel(value uint64) {
	if this.GetChannel() != value {
		this.Flag.SetFlag(USER_FLAG_POS_CHANNEL)
		this.channel = value
	}

}
func (this User) GetChannel() uint64 {
	return this.channel
}
func (this *User) SetServer(value uint64) {
	if this.GetServer() != value {
		this.Flag.SetFlag(USER_FLAG_POS_SERVER)
		this.server = value

	}

}
func (this User) GetServer() uint64 {
	return this.server
}
func (this *User) SetUuid(value string) {
	if this.GetUuid() != value {
		this.Flag.SetFlag(USER_FLAG_POS_UUID)
		this.uuid = value
	}

}
func (this User) GetUuid() string {
	return this.uuid
}
func (this *User) SetOs(value int32) {
	if this.GetOs() != value {
		this.os = value
		this.Flag.SetFlag(USER_FLAG_POS_OS)
	}

}
func (this User) GetOs() int32 {
	return this.os
}
func (this *User) SetOsVersion(value int32) {
	if this.GetOsVersion() != value {
		this.Flag.SetFlag(USER_FLAG_POS_OS_VERSION)
		this.osVersion = value
	}

}
func (this User) GetOsVersion() int32 {
	return this.osVersion
}
func (this *User) SetDeviceModel(value string) {
	if this.GetDeviceModel() != value {
		this.deviceModel = value
		this.Flag.SetFlag(USER_FLAG_POS_DEVICE_MODEL)
	}

}
func (this User) GetDeviceModel() string {
	return this.deviceModel
}
func (this *User) SetCreateTime(value time.Time) {
	this.Flag.SetFlag(USER_FLAG_POS_CREATE_TIME)
	this.createTime = value
}
func (this User) GetCreateTime() time.Time {
	return this.createTime
}
func (this *User) SetVip(value int32) {
	if this.GetVip() != value {
		this.Flag.SetFlag(USER_FLAG_POS_VIP)
		this.vip = value
	}

}
func (this User) GetVip() int32 {
	return this.vip
}
func (this *User) SetLevel(value int32) {
	if this.GetLevel() != value {
		this.level = value
		this.Flag.SetFlag(USER_FLAG_POS_LEVEL)
	}

}
func (this User) GetLevel() int32 {
	return this.level
}
func (this *User) ClearFlag() {
	this.Flag = 0
}

func (this *User) IsDirty() bool {
	return !this.Flag.IsAllZero()
}

func (user *User) Serial() (data []byte) {
	var pb ProtoUser
	pb.Uin = uint64(user.GetUin())
	pb.AccountId = user.GetAccountId()
	pb.AccountName = user.GetAccountName()
	pb.Gender = user.GetGender()
	pb.Platform = user.GetPlatform()
	pb.Channel = user.GetChannel()
	pb.Server = user.GetServer()
	pb.Uuid = user.GetUuid()
	pb.Os = user.GetOs()
	pb.OsVersion = user.GetOsVersion()
	pb.CreateTime = user.GetCreateTime().Unix()
	pb.Vip = user.GetVip()
	pb.Level = user.GetLevel()
	data, _ = proto.Marshal(&pb)
	return
}
func (user *User) UnSerial(data []byte) (ok bool) {
	var pb ProtoUser
	err := proto.Unmarshal(data, &pb)
	if err != nil {
		return false
	}
	ok = true
	user.uin = key.KeyUint64(pb.Uin)
	user.accountId = pb.AccountId
	user.accountName = pb.AccountName
	user.gender = pb.Gender
	user.platform = pb.Platform
	user.channel = pb.Channel
	user.server = pb.Server
	user.uuid = pb.Uuid
	user.os = pb.Os
	user.osVersion = pb.OsVersion
	user.createTime = time.Unix(pb.GetCreateTime(), 0)
	user.vip = pb.Vip
	user.level = pb.Level

	return
}
