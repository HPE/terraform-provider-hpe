package main

import (
	"fmt"
	"log"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"
)

func main() {
	client := morpheus.NewClient("https://yourmorpheus.com")
	client.SetUsernameAndPassword("username", "password")
	resp, err := client.Login()
	if err != nil {
		fmt.Println("LOGIN ERROR: ", err)
	}
	fmt.Println("LOGIN RESPONSE:", resp)

	// Create tenant
	createReq := &morpheus.Request{
		Body: map[string]interface{}{
			"account": map[string]interface{}{
				"name":        "gotest",
				"description": "go test teannt",
				"active":      true,
				"subdomain":   "gotest",
				"role": map[string]interface{}{
					"id": 2,
				},
				"currency":       "USD",
				"accountNumber":  "1234",
				"accountName":    "gotest",
				"customerNumber": "customer-1234",
			},
		},
	}
	createResp, err := client.CreateTenant(createReq)
	if err != nil {
		fmt.Println(err)
	}
	createResult := createResp.Result.(*morpheus.CreateTenantResult)
	fmt.Println(createResult)

	// Create subtenant group
	createSubtenantGroupReq := &morpheus.Request{
		Body: map[string]interface{}{
			"group": map[string]interface{}{
				"name":        "gotest",
				"description": "go test teannt",
				"code":        "testing",
				"location":    "chicago",
			},
		},
	}
	createSubtenantGroupResp, err := client.CreateSubtenantGroup(createResult.Tenant.ID, createSubtenantGroupReq)
	if err != nil {
		fmt.Println(err)
	}
	createSubtenantGroupResult := createSubtenantGroupResp.Result.(*morpheus.CreateSubtenantGroupResult)
	fmt.Println(createSubtenantGroupResult)

	// Get tenants
	req := &morpheus.Request{}
	tenantsResponse, err := client.ListTenants(req)
	if err != nil {
		log.Fatal(err)
	}
	result := tenantsResponse.Result.(*morpheus.ListTenantsResult)
	log.Println(result)
}
