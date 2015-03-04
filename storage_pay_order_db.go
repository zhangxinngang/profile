package profile

import (
	"database/sql"
	"fmt"
	"github.com/0studio/databasetemplate"
	key "github.com/0studio/storage_key"
	log "github.com/cihub/seelog"
	"github.com/dropbox/godropbox/errors"
	_ "github.com/go-sql-driver/mysql"
	"strings"
	"sync"
	"time"
)

type DBPayOrderStorage struct {
	databasetemplate.GenericDaoImpl
}

var initPayOrderDaoOnce sync.Once
var payOrderDao DBPayOrderStorage

func InitDBPayOrderStorage(db *sql.DB) *DBPayOrderStorage {
	done := make(chan bool)
	go func() {
		initPayOrderDaoOnce.Do(func() {
			dbTemplate := &databasetemplate.DatabaseTemplateImpl{db}
			payOrderDao = DBPayOrderStorage{GenericDaoImpl: databasetemplate.GenericDaoImpl{dbTemplate}}
			payOrderDao.CreateTable()
		})

		done <- true
	}()
	<-done
	return &payOrderDao
}

func (this *DBPayOrderStorage) CreateTable() {
	sql := `create table if not exists pay_order (
                uin BIGINT NOT NULL,
                account_id varchar(128) NOT NULL ,
                user_level smallINT DEFAULT 0,
                vip_level tinyINT DEFAULT 0,
                order_id varchar(128) NOT NULL,
                product_id VARCHAR(32) NOT NULL,
                product_name VARCHAR(20) NOT NULL,
                product_desc VARCHAR(50) NOT NULL,
                money INT DEFAULT 0,
                channel int DEFAULT 0,
                serverId smallint DEFAULT 0,
                create_time timestamp DEFAULT 0,
                status tinyint DEFAULT 0 comment '定单状态 0未处理，1游戏服务器已经处理,即钻石已加到玩家身上',
                order_type  tinyint not null default 0 comment '定单类型，0普通定单，1沙盒测试定单' ,
                receipt_data  varchar(10240) not null default '' comment 'app store receipt-data' ,
                PRIMARY KEY (uin,order_id))
             ENGINE = innodb DEFAULT CHARACTER SET utf8;`

	err := this.DatabaseTemplate.Exec(sql)
	if err != nil {
		log.Error(errors.Wrap(err, "create table pay_order error!!!"))
	}
}

func (this *DBPayOrderStorage) mapRow(resultSet *sql.Rows) (interface{}, error) {
	var (
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
	)
	payOrder := PayOrder{}
	err := resultSet.Scan(
		&uin,
		&accountId,
		&orderId,
		&orderType,
		&productId,
		&productName,
		&productDesc,
		&money,
		&channel,
		&serverId,
		&recvData,
		&status,
		&createTime)
	if err != nil {
		return nil, err
	}
	payOrder.SetUin(uin)
	payOrder.SetAccountId(accountId)
	payOrder.SetOrderId(orderId)
	payOrder.SetOrderType(orderType)
	payOrder.SetProductName(productName)
	payOrder.SetProductId(productId)
	payOrder.SetProductDesc(productDesc)
	payOrder.SetRecvData(recvData)
	payOrder.SetMoney(money)
	payOrder.SetChannel(channel)
	payOrder.SetServerId(serverId)
	payOrder.SetStatus(status)
	payOrder.SetCreateTime(createTime)
	payOrder.ClearFlag()

	return payOrder, nil
}

func (this *DBPayOrderStorage) Get(uin key.KeyUint64, orderId string) (payOrder PayOrder, ok bool) {
	sql := "select  uin,account_id,order_id,order_type,product_id,product_name,product_desc,money,channel,serverId,receipt_data,status,create_time from pay_order where uin=? and order_id=? "
	var obj interface{}
	var err error
	obj, err = this.DatabaseTemplate.QueryObject(sql, this.mapRow, uin, orderId)
	if err != nil {
		return
	}
	if obj == nil {
		return
	}

	payOrder = obj.(PayOrder)
	ok = true
	return
}

func (this *DBPayOrderStorage) MultiGet(uin key.KeyUint64, orderIdList []string) (payOrderMap PayOrderMap, ok bool) {
	payOrderMap = make(PayOrderMap)
	if len(orderIdList) == 0 {
		return
	}

	inStmt := strings.Join(orderIdList, ",")
	sql := fmt.Sprintf("select uin,account_id,order_id,order_type,product_id,product_name,product_desc,money,channel,serverId,receipt_data,status,create_time from pay_order where uin=? and order_id in (%s) ", inStmt)
	var obj interface{}
	var arr []interface{}
	var err error
	var payOrder PayOrder
	arr, err = this.DatabaseTemplate.QueryArray(sql, this.mapRow, uin)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, obj = range arr {
		if obj == nil {
			continue
		}
		payOrder = obj.(PayOrder)
		payOrderMap[payOrder.GetOrderId()] = payOrder
	}
	ok = true
	return
}

func (this *DBPayOrderStorage) GetAllUnhandledOrder(uin key.KeyUint64) (list []PayOrder) {
	list = make([]PayOrder, 0)
	sql := "select  uin,account_id,order_id,order_type,product_id,product_name,product_desc,money,channel,serverId,receipt_data,status,create_time from pay_order where uin=? and status=? order by create_time asc"
	var obj interface{}
	var err error
	var arr []interface{}
	var payOrder PayOrder
	arr, err = this.DatabaseTemplate.QueryArray(sql, this.mapRow, uin, PAY_ORDER_STATUS_UNHANDLED)
	if err != nil {
		return
	}
	for _, obj = range arr {
		payOrder = obj.(PayOrder)
		list = append(list, payOrder)
	}

	return
}
func (this *DBPayOrderStorage) UpdateStatus(order PayOrder) bool {
	sql := "update pay_order set status =? where Uin=?  and order_id=?;"
	err := this.DatabaseTemplate.Exec(sql, order.GetStatus(), order.GetUin(), order.GetOrderId())
	return err == nil
}
func (this *DBPayOrderStorage) Set(order PayOrder) bool {
	sql := "update pay_order set status =? where Uin=?  and order_id=?;"
	err := this.DatabaseTemplate.Exec(sql, order.GetStatus(), order.GetUin(), order.GetOrderId())
	return err == nil
}

func (this *DBPayOrderStorage) Add(order PayOrder) bool {
	sql := "insert ignore into  pay_order (uin,account_id,order_id,product_id,product_name,product_desc,money,channel,serverId,receipt_data,create_time,status,order_type) values (?,?,?,?,?,?,?,?,?,?,?,?,?)"
	err := this.DatabaseTemplate.Exec(sql,
		order.GetUin(),
		order.GetAccountId(),
		order.GetOrderId(),
		order.GetProductId(),
		order.GetProductName(),
		order.GetProductDesc(),
		order.GetMoney(),
		order.GetChannel(),
		order.GetServerId(),
		order.GetRecvData(),
		order.GetCreateTime(),
		order.GetStatus(),
		order.GetOrderType())
	return err == nil
}
