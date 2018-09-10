package service

import (
	"github.com/hyperledger/fabric-sdk-go/api/apitxn/chclient"
	"fmt"
)

// 发起背书
func (t *FabricSetupService) Endorse(billNo string, waitEndorseCmID string, waitEndorseAcct string) (string, error) {
	var args []string
	args = append(args, "endorse")
	args = append(args, billNo)
	args = append(args, waitEndorseCmID)
	args = append(args, waitEndorseAcct)

	req := chclient.Request{ChaincodeID:t.Setup.ChaincodeID, Fcn:args[0], Args:[][]byte{[]byte(args[1]), []byte(args[2]), []byte(args[3])}}
	response, err := t.Setup.Client.Execute(req)
	if err != nil {
		return "", fmt.Errorf(err.Error())
	}
	return string(response.Payload), nil
}

// 根据待背书人证件号码查询其对应的待背书票据
func (t *FabricSetupService) QueryMyWaitBills(waitEndorseCmId string) ([]byte, error) {
	var args []string
	args = append(args, "queryMyWaitBills")
	args = append(args, waitEndorseCmId)

	req := chclient.Request{ChaincodeID:t.Setup.ChaincodeID, Fcn:args[0], Args:[][]byte{[]byte(args[1])}}
	response, err := t.Setup.Client.Query(req)
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}
	b := response.Payload
	return b[:], nil
}

// 签收票据
func (t *FabricSetupService) Accept(billNo string, waitEndorseCmID string, waitEndorseAcct string) (string, error)  {
	var args []string
	args = append(args, "accept")
	args = append(args, billNo)
	args = append(args, waitEndorseCmID)
	args = append(args, waitEndorseAcct)

	req := chclient.Request{ChaincodeID:t.Setup.ChaincodeID, Fcn:args[0], Args:[][]byte{[]byte(args[1]), []byte(args[2]), []byte(args[3])}}
	response, err := t.Setup.Client.Execute(req)
	if err != nil {
		return "", fmt.Errorf(err.Error())
	}

	return string(response.Payload), nil
}

// 拒签票据
func (t *FabricSetupService) Reject(billNo string, waitEndorseCmID string, waitEndorseAcct string) (string, error) {
	var args []string
	args = append(args, "reject")
	args = append(args, billNo)
	args = append(args, waitEndorseCmID)
	args = append(args, waitEndorseAcct)

	req := chclient.Request{ChaincodeID:t.Setup.ChaincodeID, Fcn:args[0], Args:[][]byte{[]byte(args[1]), []byte(args[2]), []byte(args[3])}}
	response, err := t.Setup.Client.Execute(req)
	if err != nil {
		return "", fmt.Errorf(err.Error())
	}

	return string(response.Payload), nil
}
