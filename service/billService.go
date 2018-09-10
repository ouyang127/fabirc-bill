package service

import (
	"github.com/hyperledger/fabric-sdk-go/api/apitxn/chclient"
	"encoding/json"
	"fmt"
)

// 发布票据
func (t *FabricSetupService) SaveBill(bill Bill) (string, error)  {
	// 调用链码层
	var args []string
	args = append(args, "issue")

	b, err := json.Marshal(bill)
	if err != nil {
		return "", fmt.Errorf("指定票据对象序列化错误.")
	}
	// 指定调用链码的请求参数
	req := chclient.Request{ChaincodeID:t.Setup.ChaincodeID, Fcn:args[0], Args:[][]byte{b}}

	// 根据指定的请求参数调用链码
	response, err := t.Setup.Client.Execute(req)
	if err != nil {
		return "", fmt.Errorf("保存票据信息失败: %v", err)
	}

	fmt.Println("链码层返回内容: " + string(response.Payload))

	return response.TransactionID.ID, nil
}

// 根据持票人的证件号码查询票据
func (t *FabricSetupService) QueryBills(holderCmId string) ([]byte, error)  {
	var args []string
	args = append(args, "queryMyBills")
	args = append(args, holderCmId)

	req := chclient.Request{ChaincodeID:t.Setup.ChaincodeID, Fcn:args[0], Args:[][]byte{[]byte(args[1])}}
	response, err := t.Setup.Client.Query(req)

	if err != nil{
		return nil, fmt.Errorf(err.Error())
	}

	b := response.Payload

	return b[:], nil

}

// 根据票据号码查询票据详情
func (t *FabricSetupService) QueryBillByNo(billNo string) ([]byte, error){
	var args []string
	args = append(args, "queryBillByNo")
	args = append(args, billNo)

	req := chclient.Request{ChaincodeID:t.Setup.ChaincodeID, Fcn:args[0], Args:[][]byte{[]byte(args[1])}}
	response, err := t.Setup.Client.Query(req)
	if err != nil{
		return nil, fmt.Errorf(err.Error())
	}

	return response.Payload, nil
}



