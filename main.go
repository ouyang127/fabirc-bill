package main

import (
	"bill/blockchain"
	"os"
	"fmt"
	"bill/service"
	"bill/web"
	"bill/web/controller"
	"encoding/json"
)

func main() {
	setup := blockchain.FabricSetup{
		OrgAdmin:"Admin",
		OrgName:"Org1",
		ChannelID:"mychannel",

		ConfigFile:"config.yaml",

		ChannelConfig:os.Getenv("GOPATH") + "/src/bill/fixtures/artifacts/channel.tx",

		// 链码相关
		ChaincodeID:"billcc",
		ChaincodeGOPath:os.Getenv("GOPATH"),
		ChaincodePath:"/bill/chaincode",
		UserName:"User1",
	}

	err := setup.Initialized()
	if err != nil {
		fmt.Printf("初始化SDK失败: %s", err)
	}

	err = setup.InstallAndInstantiateCC()
	if err != nil {
		fmt.Println("链码安装实例化发生错误: %s", err)
	}

	serviceSetup := new(service.FabricSetupService)
	serviceSetup.Setup = &setup


	// 测试添加票据
	bill := service.Bill{
		BillInfoID: "PCO100001001",
		BillInfoAmt: "20000",
		BillInfoType: "qqq",
		BillInfoIsseDate:"2010-01-01",
		BillInfoDueDate:"2010-01-04",

		DrwrAcct:"A公司",
		DrwrCmID:"ACMID",

		AccptrAcct: "111",
		AccptrCmID:"111",
		PyeeAcct:"111",
		PyeeCmID:"111",

		HoldrAcct:"A公司",
		HoldrCmID:"ACMID",
	}


	responeStr, err := serviceSetup.SaveBill(bill)
	if err != nil{
		fmt.Println(err)
	}else {
		fmt.Println("票据发布成功: " + responeStr)
	}

	// 根据持票人的证件号码查询票据
	b, err := serviceSetup.QueryBills("ACMID")

	if err != nil {
		fmt.Errorf(err.Error())
	}else {
		fmt.Println("根据持票人的证件号码查询票据成功")
		var bills = []service.Bill{}
		json.Unmarshal([]byte(b), &bills)

		for _, temp := range bills {
			fmt.Println(temp)
		}
	}

	// 根据票据号码查询票据详情
	b, err = serviceSetup.QueryBillByNo("PCO100001001")
	if err != nil {
		fmt.Println(err.Error())
	}else{
		fmt.Println("根据票据号码查询票据详情成功")
		bill = service.Bill{}
		json.Unmarshal(b, &bill)
		fmt.Println(bill)
		for _, hisItem := range bill.History{
			fmt.Println(hisItem)
		}
	}

	// 发起背书
	result, err := serviceSetup.Endorse("PCO100001001", "BCMID", "B公司")
	if err != nil {
		fmt.Println(err.Error())
	}else{
		fmt.Println(result)
	}

	// 根据待背书人证件号码查询其对应的待背书票据
	b, err = serviceSetup.QueryMyWaitBills("BCMID")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("根据待背书人证件号码查询其对应的待背书票据成功")
		var bills = []service.Bill{}
		json.Unmarshal(b, &bills)
		for _, temp := range bills {
			fmt.Println(temp)
		}
	}

	// 签收票据
	result, err = serviceSetup.Accept("PCO100001001", "BCMID", "B公司")
	if err != nil{
		fmt.Println(err.Error())
	}else {
		fmt.Println(result)
	}

	// 拒签票据
	result, err = serviceSetup.Reject("PCO100001001", "BCMID", "B公司")
	if err != nil{
		fmt.Println(err.Error())
	}else {
		fmt.Println(result)
	}

	// 根据票据号码查询票据详情
	b, err = serviceSetup.QueryBillByNo("PCO100001001")
	if err != nil {
		fmt.Println(err.Error())
	}else{
		fmt.Println("根据票据号码查询票据详情成功")
		bill = service.Bill{}
		json.Unmarshal(b, &bill)
		fmt.Println(bill)
		for _, hisItem := range bill.History{
			fmt.Println(hisItem)
		}
	}

	app := new(controller.Application)
	app.Fabric = serviceSetup
	err = web.WebStart(app)
	if err != nil {
		fmt.Println(err.Error())
	}


}
