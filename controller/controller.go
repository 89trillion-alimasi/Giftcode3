package controller

import (
	"GiftCode2/model"
	"GiftCode2/redis"
	"GiftCode2/service"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

const (
	MESSAGE = "message"
)

//创建礼品码
func createGiftCode(c *gin.Context) {
	var gift model.GiftCode
	// 绑定post参数
	if err := c.Bind(&gift); err != nil {
		c.JSON(ParameterBindingIsUnsuccessful, gin.H{MESSAGE: StatusText(ParameterBindingIsUnsuccessful)})
		return
	}

	// 检查礼品码类型
	switch gift.Type {
	case 1: // 1 - 指定用户一次性消耗
		// 一次性消耗强制设置领取次数为 1
		gift.AvailableTimes = 1
		if gift.ReceivingUser == "" {
			c.JSON(ReceivingUserIsEmpty, gin.H{MESSAGE: StatusText(ReceivingUserIsEmpty)})
			return
		}
	case 2: // 2 - 不指定用户限制兑换次数
		if gift.AvailableTimes <= 0 {
			c.JSON(SpecifyTheNumberOfRedemptions, gin.H{MESSAGE: StatusText(SpecifyTheNumberOfRedemptions)})
			return
		}
	case 3: // 3 - 不限用户不限次数兑换 无用户限制 无兑换次数限制这里不做处理
	default: // 非法填写的礼品码类型
		c.JSON(Invalidgiftcodetype, gin.H{MESSAGE: StatusText(Invalidgiftcodetype)})
		return
	}

	// 检查礼品码描述
	if gift.Description == "" {
		c.JSON(GiftCodeDescription, gin.H{MESSAGE: StatusText(GiftCodeDescription)})
		return
	}

	// 检查礼品码有效期是否为空
	if gift.ValidPeriod == "" {
		c.JSON(PleaseEnterAValidTime, gin.H{MESSAGE: StatusText(PleaseEnterAValidTime)})
		return
	}

	// 检查礼品码有效期是否正确
	vp, err := time.ParseInLocation("2006-01-02 15:04:05", gift.ValidPeriod, time.Local)
	if err != nil {
		c.JSON(IncorrectTimeFormat, gin.H{MESSAGE: StatusText(IncorrectTimeFormat)})
		return
	}
	exp := vp.Sub(time.Now())
	if exp <= 0 {
		c.JSON(InvalidTime, gin.H{MESSAGE: StatusText(InvalidTime)})
		return
	}

	// 检查礼品码包含的礼品包是否为空
	if gift.GiftPackages == nil || len(gift.GiftPackages) == 0 {
		c.JSON(PackageContent, gin.H{MESSAGE: StatusText(PackageContent)})
		return
	}
	var n = 5
	for {
		// 生成 8 位礼品码
		gift.Code = service.RandStringBytesMask(8)
		cmd := redis.RDB.Get(gift.Code)

		if cmd.Err() != nil {
			break
		} else {
			if n > 0 {
				n--
				continue
			} else {
				c.String(CreateGiftCodeFaied, StatusText(CreateGiftCodeFaied))
				return
			}

		}

	}
	logrus.Infof("生成新的礼品码: %s", gift.Code)

	// 将礼品码存储到 redis
	if err := service.SaveGiftCode(&gift); err != nil {
		c.JSON(InsertionFailed, gin.H{MESSAGE: StatusText(InsertionFailed)})
		return
	}

	// 返回礼品码
	c.JSON(CreatedSuccessfully, gin.H{"code": StatusText(CreatedSuccessfully)})
}

//查询礼品码
func queryGiftCode(c *gin.Context) {
	// 检测是否输入了礼品码
	code := c.Query("code")
	if code == "" {
		c.JSON(GiftCodeErr, gin.H{MESSAGE: StatusText(GiftCodeErr)})
		return
	}

	// 尝试从 redis 中查找该礼品码
	giftCode := service.QueryGiftCode(code)
	if giftCode == nil {
		c.JSON(GiftCodeHasExpired, gin.H{MESSAGE: StatusText(GiftCodeHasExpired)})
		return
	}

	// 返回礼品码
	c.JSON(200, giftCode)
}

//验证礼品码
func verifyGiftCode(c *gin.Context) {
	var req model.VerifyRequest
	if err := c.Bind(&req); err != nil {
		c.JSON(ParameterBindingIsUnsuccessful, gin.H{MESSAGE: StatusText(PleaseEnterAValidTime)})
		return
	}

	// 检查礼品码
	if req.Code == "" {
		c.JSON(GiftCodeErr, gin.H{MESSAGE: StatusText(GiftCodeErr)})
		return
	}

	// 检查领取用户
	if req.User == "" {
		c.JSON(ReceivingUserIsEmpty, gin.H{MESSAGE: StatusText(ReceivingUserIsEmpty)})
		return
	}

	gifCode, err := service.VerifyGiftCode(req)
	if err != nil {
		c.JSON(FailedToClaim, gin.H{MESSAGE: StatusText(FailedToClaim)})
		return
	}
	c.JSON(http.StatusInsufficientStorage, gin.H{"GiftPackages": gifCode.GiftPackages})
}
