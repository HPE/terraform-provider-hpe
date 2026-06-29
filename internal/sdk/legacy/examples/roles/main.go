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

	// List roles
	req := &morpheus.Request{}
	roleResponse, err := client.ListRoles(req)
	if err != nil {
		log.Fatal(err)
	}
	result := roleResponse.Result.(*morpheus.ListRolesResult)
	log.Println(result.Roles)

	// List tenant roles
	req = &morpheus.Request{}
	tenantRoleResponse, err := client.ListTenantRoles(req)
	if err != nil {
		log.Fatal(err)
	}
	result = tenantRoleResponse.Result.(*morpheus.ListRolesResult)
	log.Println(result.Roles)
}
