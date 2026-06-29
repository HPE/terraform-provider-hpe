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

	// List wikis
	req := &morpheus.Request{}
	wikiResp, err := client.ListWikis(req)
	if err != nil {
		log.Fatal(err)
	}
	result := wikiResp.Result.(*morpheus.ListWikisResult)
	log.Println(result.Wikis)

	// List wiki categories
	wikiCategoriesReq := &morpheus.Request{}
	wikiCategoryResp, err := client.ListWikiCategories(wikiCategoriesReq)
	if err != nil {
		log.Fatal(err)
	}
	wikiCategoryResult := wikiCategoryResp.Result.(*morpheus.ListWikiCategoriesResult)
	log.Println(wikiCategoryResult.WikiCategories)

	// Get instance wiki returns the wiki for the first instance found
	instancesResp, err := client.ListInstances(&morpheus.Request{})
	if err != nil {
		log.Fatal(err)
	}
	instancesResult := instancesResp.Result.(*morpheus.ListInstancesResult)
	instanceCount := len(*instancesResult.Instances)
	if instanceCount > 0 {
		instance := (*instancesResult.Instances)[0]

		instanceWikiReq := &morpheus.Request{}
		instanceWikiResp, err := client.GetInstanceWiki(instance.ID, instanceWikiReq)
		if err != nil {
			log.Fatal(err)
		}
		instanceWikiResult := instanceWikiResp.Result.(*morpheus.GetWikiResult)
		log.Println(instanceWikiResult.Wiki)
	} else {
		log.Println("No instances found")
	}
}
