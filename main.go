package main

import (
	"fmt"

	"github.com/bw-global/openapi-sign/utils"
)

func main() {
	method := "baiwang.oversea.invoice.acquireInvoiceNumber"
	version := "v1"
	appkey := "appkey"
	appsecret := "appsecret"
	scope := "TAXCODE"
	scopeValue := "292221003212"

	type SignatureBody struct {
		KeyNumber string `json:"keyNumber"`
		Count     string `json:"count"`
		InState   string `json:"inState"`
		Uuid      string `json:"uuid"`
		Tin       string `json:"tin"`
	}
	body := SignatureBody{
		KeyNumber: "",
		Count:     "",
		InState:   "",
		Uuid:      "",
		Tin:       "",
	}
	signature, err := utils.Signature(method, version, appkey, appsecret, scope, scopeValue, body)
	if err != nil {
		fmt.Println("Error generating signature:", err)
		return
	}
	fmt.Println("Generated signature:", signature)
}
