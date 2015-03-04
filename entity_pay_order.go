package profile

import (
	"time"
)

type PayOrder struct {
	uin         uint64
	accountId   string
	orderId     string
	orderType   int32 // 0普通定单，1沙盒测试定单
	productName string
	productId   string
	productDesc string
	recvData    string
	money       int32
	channel     uint64
	serverId    uint64
	status      int32
	createTime  time.Time
}

func (this PayOrder) IsStatusUnHandled() bool {
	return this.status == PAY_ORDER_STATUS_UNHANDLED
}
func (this PayOrder) IsStatusHandled() bool {
	return this.status == PAY_ORDER_STATUS_HANDLED
}

func (this *PayOrder) SetStatusHandled() {
	this.status = PAY_ORDER_STATUS_HANDLED
}

type PayOrderMap map[string]PayOrder // key orderid

func (this *PayOrder) SetUin(value uint64) {
	this.uin = value
}
func (this PayOrder) GetUin() uint64 {
	return this.uin
}
func (this *PayOrder) SetAccountId(value string) {
	this.accountId = value
}
func (this PayOrder) GetAccountId() string {
	return this.accountId
}
func (this *PayOrder) SetOrderId(value string) {
	this.orderId = value
}
func (this PayOrder) GetOrderId() string {
	return this.orderId
}
func (this *PayOrder) SetOrderType(value int32) {
	this.orderType = value
}
func (this PayOrder) GetOrderType() int32 {
	return this.orderType
}
func (this *PayOrder) SetProductName(value string) {
	this.productName = value
}
func (this PayOrder) GetProductName() string {
	return this.productName
}
func (this *PayOrder) SetProductId(value string) {
	this.productId = value
}
func (this PayOrder) GetProductId() string {
	return this.productId
}
func (this *PayOrder) SetProductDesc(value string) {
	this.productDesc = value
}
func (this PayOrder) GetProductDesc() string {
	return this.productDesc
}
func (this *PayOrder) SetRecvData(value string) {
	this.recvData = value
}
func (this PayOrder) GetRecvData() string {
	return this.recvData
}
func (this *PayOrder) SetMoney(value int32) {
	this.money = value
}
func (this PayOrder) GetMoney() int32 {
	return this.money
}
func (this *PayOrder) SetChannel(value uint64) {
	this.channel = value
}
func (this PayOrder) GetChannel() uint64 {
	return this.channel
}
func (this *PayOrder) SetServerId(value uint64) {
	this.serverId = value
}
func (this PayOrder) GetServerId() uint64 {
	return this.serverId
}
func (this *PayOrder) SetStatus(value int32) {
	this.status = value
}
func (this PayOrder) GetStatus() int32 {
	return this.status
}
func (this *PayOrder) SetCreateTime(value time.Time) {
	this.createTime = value
}
func (this PayOrder) GetCreateTime() time.Time {
	return this.createTime
}

func (this *PayOrder) ClearFlag() {
}
