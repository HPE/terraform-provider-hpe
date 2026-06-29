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

	// Get backups
	req := &morpheus.Request{}
	backupsResponse, err := client.ListBackups(req)
	if err != nil {
		log.Fatal(err)
	}
	result := backupsResponse.Result.(*morpheus.ListBackupsResult)
	for _, backup := range *result.Backups {
		fmt.Println(backup.BackupProvider)
	}
}
