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

	// Get backup jobs
	req := &morpheus.Request{}
	backupJobsResponse, err := client.ListBackupJobs(req)
	if err != nil {
		log.Fatal(err)
	}
	result := backupJobsResponse.Result.(*morpheus.ListBackupJobsResult)
	for _, backupJob := range *result.BackupJobs {
		fmt.Println(backupJob)
	}
}
