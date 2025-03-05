package handler

import (
	"github/lhh-gh/gift/database"
	"github/lhh-gh/gift/util"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 获取所有奖品信息，用于初始化轮盘
func GetAllGifts(ctx *gin.Context) {
	gifts := database.GetAllGiftsV1()
	if len(gifts) == 0 {
		ctx.JSON(http.StatusInternalServerError, nil)
	} else {
		//抹掉敏感信息
		for _, gift := range gifts {
			gift.Count = 0
		}
		ctx.JSON(http.StatusOK, gifts)
	}
}

// 抽奖
func Lottery(ctx *gin.Context) {
	for try := 0; try < 10; try++ { //最多重试10次
		gifts := database.GetAllGiftInventory() //获取所有奖品剩余的库存量
		ids := make([]int, 0, len(gifts))
		probs := make([]float64, 0, len(gifts))
		for _, gift := range gifts {
			if gift.Count > 0 { //先确保redis返回的库存量大小0，因为抽奖算法Lottery不支持抽中概率为0的奖品
				ids = append(ids, gift.Id)
				probs = append(probs, float64(gift.Count))
			}
		}
		if len(ids) == 0 {
			CloseChannel() //关闭channel
			// go CloseMQ()                               //关闭写mq的连接
			ctx.String(http.StatusOK, strconv.Itoa(0)) //0表示所有奖品已抽完
			return
		}
		index := util.Lottery(probs) //抽中第index个奖品
		giftId := ids[index]
		err := database.ReduceInventory(giftId) // 先从redis上减库存
		if err != nil {
			util.LogRus.Warnf("奖品%d减库存失败", giftId)
			continue //减库存失败，则重试
		} else {
			// 用户ID写死为1，
			PutOrder(1, giftId) //把订单信息写入channel
			// ProduceOrder(1, giftId)                         //把订单信息写入mq
			ctx.String(http.StatusOK, strconv.Itoa(giftId)) //减库存成功后才给前端返回奖品ID
			return
		}
	}
	ctx.String(http.StatusOK, strconv.Itoa(database.EMPTY_GIFT)) //如果10次之后还失败，则返回“谢谢参与”
}
