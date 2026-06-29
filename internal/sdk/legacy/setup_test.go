package morpheus_test

import (
	"github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"
	"net/http"
	"testing"
)

func TestSetupCheck(t *testing.T) {
	client := getTestClient(t)
	req := &morpheus.Request{}
	resp, err := client.SetupCheck(req)
	assertResponse(t, resp, err)
}

// func TestSetupInit(t *testing.T) {
// 	client := getTestClient(t)
// 	req := &morpheus.Request{
// 		Body: map[string]interface{}{
// 			"accountName": "root",
// 			"firstName": "Super",
// 			"lastName": "Admin",
// 			"password": "ChangeM3!",
// 		}
// 	}
// 	resp, err := client.SetupInit(req)
// 	assertResponse(t, resp, err)
// }

// It should expect to be already setup
func TestSetupInitShouldBeAlreadySetup(t *testing.T) {
	client := getTestClient(t)
	req := &morpheus.Request{
		Body: map[string]interface{}{
			"accountName": "root",
			"firstName":   "Super",
			"lastName":    "Admin",
			"password":    "ChangeM3!",
		},
	}
	resp, err := client.SetupInit(req)
	assertResponseStatusCode(t, resp, http.StatusBadRequest)
	assertErrorNotNil(t, err)
	//assertBadResponse(t, resp, err)
	// could assert resp.Msg == "Already setup" or whatever
}
