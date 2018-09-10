package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"fmt"
	"github.com/hyperledger/fabric/protos/peer"
)


type BillChainCode struct {

}

func (t *BillChainCode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}

func (t *BillChainCode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	function, args := stub.GetFunctionAndParameters()

	if function == "issue" {
		return t.issue(stub, args)		// 发布票据
	}else if function == "queryMyBills" {
		return t.queryMyBills(stub, args)	// 查询当前用户的票据
	}else if function == "queryBillByNo" {
		return t.queryBillByNo(stub, args)	// 根据票据号码查询相应的票据详情
	}else  if function == "queryMyWaitBills" {
		return t.queryMyWaitBills(stub, args)	// 查询当前用户的待背书票据
	}else if function == "endorse" {
		return t.endorse(stub, args)	// 票据背书
	}else if function == "accept" {
		return t.accept(stub, args)	// 签收票据
	}else if function == "reject" {
		return t.reject(stub, args)	// 拒签票据
	}

	return shim.Error("指定的函数名称错误")
}

//b->c->s

func main()  {
	err := shim.Start(new(BillChainCode))
	if err != nil {
		fmt.Println("启动链码错误: ", err)
	}
}