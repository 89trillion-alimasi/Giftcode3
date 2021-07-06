package service

import (
	"GiftCode2/model"
	"GiftCode2/redis"
	"errors"
	jsoniter "github.com/json-iterator/go"
	"github.com/sirupsen/logrus"
	"time"
)

// SaveGiftCode 将礼品码存储到 redis 中
func SaveGiftCode(gift *model.GiftCode) error {
	vp, _ := time.ParseInLocation("2006-01-02 15:04:05", gift.ValidPeriod, time.Local)
	gift.CreatTime = time.Now().Format("2006-01-02 15:04:05")
	bs, err := jsoniter.Marshal(gift)
	if err != nil {
		return err
	}

	return redis.RDB.Set(gift.Code, string(bs), time.Until(vp)).Err()
}

// QueryGiftCode 查询礼品码信息
func QueryGiftCode(code string) *model.GiftCode {
	cmd := redis.RDB.Get(code)
	if cmd.Err() != nil {
		logrus.Warnf("礼品码未找到: %s", cmd.Err())
		return nil
	}
	var giftCode model.GiftCode
	err := jsoniter.Unmarshal([]byte(cmd.Val()), &giftCode)
	if cmd.Err() != nil {
		logrus.Warnf("礼品码反序列化失败: %s", err)
		return nil
	}
	return &giftCode
}

// VerifyGiftCode 验证礼品码
func VerifyGiftCode(req model.VerifyRequest) (*model.GiftCode, error) {
	// 从 redis 中获取礼品码
	gifCode := QueryGiftCode(req.Code)
	if gifCode == nil {
		return nil, errors.New("礼品码不存在或已过期")
	}

	// 检查礼品码类型
	switch gifCode.Type {
	case 1: // 1 - 指定用户一次性消耗
		if req.User != gifCode.ReceivingUser {
			return nil, errors.New("当前礼品码已经指定用户，您输入的用户无权领取")
		}
		if _, err := One_time(gifCode); err != nil {
			return nil, err
		}
		// TODO: 当前没有真实用户体系，所以这里模拟添加奖励
		logrus.Infof("用户 %s 添加奖励完成", req.User)
		return gifCode, nil

	case 2: // 2 - 不指定用户限制兑换次数

		// 如果礼品码正好还剩一次可以领取，领取后需要删除
		if gifCode.AvailableTimes == 1 {
			if _, err := One_time(gifCode); err != nil {
				return nil, err
			}
			return gifCode, nil
		} else {
			// 如果礼品码剩余领取次数大于 1，领取后对可用次数 -1
			gifCode.AvailableTimes--
			// 将领取用户加入到已 礼品码 领取列表中
			gifCode.AddReceivedUser(req.User)

			// 重新将礼品码存入 redis
			bs, err := jsoniter.Marshal(gifCode)
			if err != nil {
				return nil, err
			}

			logrus.Infof("用户 %s 添加奖励完成", req.User)
			return gifCode, redis.RDB.Set(gifCode.Code, string(bs), gifCode.Expiration).Err()
		}
	case 3: // 3 - 不限用户不限次数兑换 无用户限制 无兑换次数限制这里不做处理
		return gifCode, nil
	default: // 非法填写的礼品码类型
		logrus.Infof("用户 %s 添加奖励完成", req.User)
		return gifCode, nil
	}
}

//只能取一次礼品
func One_time(code *model.GiftCode) (*model.GiftCode, error) {
	cmd := redis.RDB.Del(code.Code)
	if cmd.Err() != nil {
		return nil, cmd.Err()
	}
	if cmd.Val() != 1 {
		return nil, errors.New("礼品码领取失败: 礼品码删除失败")
	}
	return code, nil
}
