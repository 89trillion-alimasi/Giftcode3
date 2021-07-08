package service

import (
	"GiftCode2/model"
	"reflect"
	"testing"
)

//存储单元测试
func TestSaveGiftCode(t *testing.T) {

	type GiftPackage struct {
		Name string
		Num  int
	}
	type tests struct {
		gift model.GiftCode
		want error
	}

	test := []tests{{gift: model.GiftCode{Type: 1, ReceivingUser: "alms", ValidPeriod: 1625628980, CreateUser: "admin", AvailableTimes: 2,
		GiftPackages: []model.GiftPackage{{Name: "金币", Num: 1000}, {Name: "钻石", Num: 30}}}, want: nil},
		{gift: model.GiftCode{Type: 2, ReceivingUser: "alms", ValidPeriod: 1625628980, CreateUser: "admin1", AvailableTimes: 5,
			GiftPackages: []model.GiftPackage{{Name: "金币", Num: 1000}, {Name: "钻石", Num: 50}}}, want: nil},
		{gift: model.GiftCode{Type: 2, ReceivingUser: "alms", ValidPeriod: 1625628980, CreateUser: "admin1", AvailableTimes: 5,
			GiftPackages: []model.GiftPackage{{Name: "金币", Num: 1000}, {Name: "钻石", Num: 50}}}, want: nil},
	}
	for _, v := range test {
		got := SaveGiftCode(&v.gift)
		if !reflect.DeepEqual(got, v.want) {
			t.Error("expect: %v,got: %v", v.want, got)
		}
	}

}

//查询单元测试
func TestQueryGiftCode(t *testing.T) {
	type tests struct {
		giftcode string
		want     *model.GiftCode
	}

	test := []tests{
		{giftcode: "0TO8B8MO", want: &model.GiftCode{Description: "测试type1", Type: 1,
			ReceivingUser: "alms", AvailableTimes: 1, ValidPeriod: 1625628980,
			CreatTime: 1625628980, CreateUser: "admin",
			GiftPackages: []model.GiftPackage{{Name: "金币", Num: 10},
				{Name: "钻石", Num: 20}}, ReceivedUsers: nil, ReceivedCount: 0, Code: "0TO8B8MO"}},
	}
	for _, v := range test {
		got := QueryGiftCode(v.giftcode)
		if !reflect.DeepEqual(got, v.want) {
			t.Error("expect: %v,got: %v", v.want, got)
		}
	}

}

//验证单元测试
func TestVerifyGiftCode(t *testing.T) {
	type tests struct {
		verrify model.VerifyRequest
		want    *model.GiftCode
	}

	test := []tests{{verrify: model.VerifyRequest{Code: "0TO8B8MO", User: "alms"}, want: &model.GiftCode{Description: "测试type1", Type: 1,
		ReceivingUser: "alms", AvailableTimes: 1, ValidPeriod: 1625628980,
		CreatTime: 1625628780, CreateUser: "admin",
		GiftPackages: []model.GiftPackage{{Name: "金币", Num: 10},
			{Name: "钻石", Num: 20}}, ReceivedUsers: nil, ReceivedCount: 0, Code: "0TO8B8MO"}},
		{verrify: model.VerifyRequest{Code: "thuVRONS", User: "alms"}, want: &model.GiftCode{Description: "测试type1", Type: 1,
			ReceivingUser: "alms", AvailableTimes: 1, ValidPeriod: 1625628980,
			CreatTime: 1625628780, CreateUser: "admin",
			GiftPackages: []model.GiftPackage{{Name: "金币", Num: 10},
				{Name: "钻石", Num: 20}}, ReceivedUsers: nil, ReceivedCount: 0, Code: "0TO8B8MO"}}}

	for _, v := range test {
		got, _ := VerifyGiftCode(v.verrify)
		if !reflect.DeepEqual(got, v.want) {
			t.Error("expect: %v,got: %v", v.want, got)
		}
	}
}
