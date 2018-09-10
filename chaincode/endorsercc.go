package main

import (
	"github.com/hyperledger/fabric/protos/peer"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"encoding/json"
)

// 处理背书请求
// args: billNo, WaitEndorseCmID, WaitEndorseAcct
func (t *BillChainCode) endorse(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 3 {
		return shim.Error("必须指定票据号码, 待背书人证件号码及待背书人名称")
	}

	// 根据指定的票据号码查询对应的信息
	bill, bl := t.getBill(stub, args[0])
	if !bl {
		return shim.Error("根据指定的票据号码查询信息时发生错误")
	}

	// 检查当前待背书的票据是否为与持有人是同一个人
	if bill.HoldrCmID == args[1] {
		return shim.Error("被背书人不能是当前持票人")
	}

	// 当前待背书的票据不能是票据流转历史中的持有人
	iterator, err := stub.GetHistoryForKey(bill.BillInfoID)
	if err != nil {
		return shim.Error("获取票据流转历史信息时发生错误")
	}
	defer iterator.Close()

	var hisBill Bill
	for iterator.HasNext() {
		hisData, err := iterator.Next()
		if err != nil {
			return shim.Error("获取历史数据时发生错误")
		}
		json.Unmarshal(hisData.Value, &hisBill)
		if hisData.Value == nil {
			var empty Bill
			hisBill = empty
		}

		if bill.HoldrCmID == args[1] {
			return shim.Error("被背书人不能是该票据的历史持有人")
		}

	}

	// 更改票据状态, 待背书人信息及拒绝背书人信息
	bill.State = BillInfo_State_EndorseWaitSign
	bill.WaitEndorseCmID = args[1]
	bill.WaitEndorseAcct = args[2]
	bill.RejectEndorseAcct = ""
	bill.RejectEndorseCmID = ""

	// 保存票据信息
	_, bl = t.pubBill(stub, bill)
	if !bl {
		return shim.Error("票据背书请求失败, 保存票据信息时发生错误")
	}

	// 根据待背书人的证件号码及票据号码创建复合键, 以方便批量查询(查询待背书票据)
	waitEndorserCmIdBillInfoIdIndexkey, err := stub.CreateCompositeKey(IndexName, []string{bill.WaitEndorseCmID, bill.BillInfoID})
	if err != nil {
		return shim.Error("根据待背书人的证件号码及票据号码创建复合键失败")
	}
	stub.PutState(waitEndorserCmIdBillInfoIdIndexkey, []byte{0x00})

	return shim.Success([]byte("发起背书请求成功, 此票据待被背书人处理"))

}

// 签收票据
// args : billNo, WaitEndorseCmID, WaitEndorseAcct
func (t *BillChainCode) accept(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 3 {
		return shim.Error("必须指定票据号码, 待背书人证件号码及待背书人名称")
	}

	// 查询票据信息
	bill, bl := t.getBill(stub, args[0])
	if !bl {
		return shim.Error("根据指定的票据号码查询信息时发生错误")
	}

	//删除复合键, 以便于前手持票人不能再查询到信息
	holderCmIdBillInfoIdIndexKey, err := stub.CreateCompositeKey(IndexName, []string{bill.HoldrCmID, bill.BillInfoID})
	if err != nil{
		return shim.Error("复合键创建失败")
	}
	stub.DelState(holderCmIdBillInfoIdIndexKey)

	// 更改票据状态
	bill.State = BillInfo_State_EndorseSigned
	bill.HoldrCmID = args[1]
	bill.HoldrAcct = args[2]
	bill.RejectEndorseCmID = ""
	bill.RejectEndorseAcct = ""

	_, bl = t.pubBill(stub, bill)
	if !bl{
		return shim.Error("票据背书签收失败, 保存票据时发生错误")
	}

	// 构建复合键

	return shim.Success([]byte("票据背书签收成功"))

}

// 票据背书拒签
// args: billNo, WaitEndorseCmID, WaitEndorseAcct
func (t *BillChainCode) reject(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 3{
		return shim.Error("必须指定票据号码, 拒签人证件号码及拒签人名称")
	}

	bill, bl := t.getBill(stub, args[0])
	if !bl {
		return shim.Error("根据指定的票据号码查询信息时发生错误")
	}

	waitEndorseCmIDBillInfoIndexKey, err := stub.CreateCompositeKey(IndexName, []string{bill.WaitEndorseCmID, bill.BillInfoID})
	if err != nil {
		return shim.Error("复合键创建失败")
	}
	stub.DelState(waitEndorseCmIDBillInfoIndexKey)

	// 更改当前票据的状态
	bill.State = BillInfo_State_EndorseReject
	bill.RejectEndorseAcct = args[2]
	bill.RejectEndorseCmID = args[1]
	bill.WaitEndorseCmID = ""
	bill.WaitEndorseAcct = ""

	_, bl = t.pubBill(stub, bill)
	if !bl {
		return shim.Error("票据拒签失败, 保存票据信息时发生错误")
	}

	return shim.Success([]byte("票据拒签成功"))
}
