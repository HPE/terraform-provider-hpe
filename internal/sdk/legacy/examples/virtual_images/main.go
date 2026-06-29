package main

import (
	b64 "encoding/base64"
	"fmt"
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
	/*
		// List virtual images
		req := &morpheus.Request{
			QueryParams: map[string]string{"filterType": "All"},
		}
		response, err := client.ListVirtualImages(req)
		if err != nil {
			log.Fatal(err)
		}
		result := response.Result.(*morpheus.ListVirtualImagesResult)
		log.Println(result.VirtualImages)
	*/

	createReq := &morpheus.Request{
		Body: map[string]interface{}{
			"virtualImage": map[string]interface{}{
				"name":            "demoimage",
				"imageType":       "iso",
				"isCloudInit":     false,
				"installAgent":    false,
				"virtioSupported": false,
			},
		},
	}

	createResp, err := client.CreateVirtualImage(createReq)

	if err != nil {
		fmt.Println(err)
	}

	createImageResult := createResp.Result.(*morpheus.CreateVirtualImageResult)
	fmt.Println(createImageResult)

	// Upload Virtual Image
	virtualImageName := "netboot.xyz.iso"
	virtualImagePath := fmt.Sprintf("/Downloads/%s", virtualImageName)
	data, err := os.ReadFile(virtualImagePath)

	if err != nil {
		fmt.Println(err)
	}

	sEnc := b64.StdEncoding.EncodeToString([]byte(data))

	uploadResp, err := client.UploadVirtualImage(createImageResult.VirtualImage.ID, &morpheus.Request{
		QueryParams: map[string]string{
			"filename": virtualImageName,
		},
		IsStream:   true,
		StreamBody: fmt.Sprintf("data:application/octet-stream;name=%s;base64,%s", virtualImageName, sEnc),
	})

	if err != nil {
		fmt.Println(err)
	}

	uploadImageResult := uploadResp.Result.(*morpheus.UploadVirtualImageResult)
	fmt.Println(uploadImageResult)

}
