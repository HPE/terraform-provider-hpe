package main

import (
	"fmt"
	"log"
	"os"

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

	// Upload plugin
	filePath := "/plugin/file/path/"
	fileName := "morpheus-netbox-plugin-0.0.1.jar"
	fullPath := filePath + fileName
	data, err := os.ReadFile(fullPath)
	if err != nil {
		fmt.Println(err)
	}

	var filePayloads []*morpheus.FilePayload
	filePayload := &morpheus.FilePayload{
		ParameterName: "file",
		FileName:      fileName,
		FileContent:   data,
	}
	filePayloads = append(filePayloads, filePayload)
	response, err := client.UploadPlugin(filePayloads, &morpheus.Request{})
	if err != nil {
		fmt.Printf("API FAILURE: %s - %s", response, err)
	}
	fmt.Println(response)

	// List plugins
	req := &morpheus.Request{}
	pluginResponse, err := client.ListPlugins(req)
	if err != nil {
		log.Fatal(err)
	}
	result := pluginResponse.Result.(*morpheus.ListPluginsResult)
	log.Println(result.Plugins)
}
