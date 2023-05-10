package main

import (
	"fmt"
	"loadbalancer/APIRequest"
	"loadbalancer/manager"
)

func main() {
	m := manager.NewManager("ROUND_ROBIN", "TOKEN_BUCKET")
	m.ApplyConfigs("some endpoint", []string{"3.3.3.3"})
	req := APIRequest.NewAPIRequest("1.1.1.1", "2.2.2.2", "some endpoint")
	err := m.HandleRequest(req)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(req)
}
