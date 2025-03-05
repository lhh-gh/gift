package test

import (
	"gift/database"
	"gift/util"
	"strconv"
	"testing"
	"time"
)

func TestProvisionalOrder(t *testing.T) {
	util.InitLog("log")
	database.InitGiftInventory()

	uid := 1
	giftId := 10
	orderLife := 3 //订单多少秒后自动失效
	cache := database.NewTimeoutCache()
	database.CreateProvisionalOrder(uid, giftId, orderLife, cache) //创建订单
	value := database.GetProvisionalOrder(uid)                     //查询订单
	if value != strconv.Itoa(giftId) {
		t.Errorf("没有得到预期有商品ID: %s", value)
	}
	time.Sleep(time.Duration(orderLife)*time.Second + time.Second)
	value = database.GetProvisionalOrder(uid) //查询订单
	if value != "" {
		t.Errorf("订单应该失效，但是却拿到了商品ID: %s", value)
	}
}

// go test -v .\database\test\ -run=^TestProvisionalOrder$ -count=1
