package main

import (
	"fmt"
	ibapi "github.com/Purelightme/goibapi"
	"time"
)

var fclient *FClient

type FClient struct {
	*ibapi.EClient
}

type FWrapper struct {
	ibapi.Wrapper
}

func (w FWrapper) ConnectAck() {
	fmt.Println("IB Gateway ConnectAck")
}

func (w FWrapper) ConnectionClosed() {
	fmt.Println("IB Gateway ConnectionClosed")
	go NewFClientUntilSuccess()
}

func NewFClientUntilSuccess() {
	t := time.Now()
	for {
		_, err := NewFClient()
		if err == nil {
			break
		}
		fmt.Println(err.Error())
		time.Sleep(time.Second * 10)
		if time.Since(t) > time.Hour {
			fmt.Println("NewFClientUntilSuccessNotOK")
		}
		if time.Since(t) > time.Hour*72 {
			fmt.Println("NewFClientUntilSuccessTimeout")
			break
		}
	}
}

func NewFClient() (*FClient, error) {
	fclient = &FClient{ibapi.NewEClient(FWrapper{})}
	err := fclient.Connect("127.0.0.1", 4002, 1)
	if err != nil {
		return nil, err
	}
	fclient.ReqMarketDataType(3)
	contract := &ibapi.Contract{
		Symbol:   "USD",
		Currency: "JPY",
		SecType:  "CASH",
		Exchange: "IDEALPRO",
	}
	fclient.ReqMktData(1, contract, "", false, false, nil)

	return fclient, nil
}

func main() {
	NewFClientUntilSuccess()
	select {}
}
