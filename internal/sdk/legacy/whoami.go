package morpheus

import (
	_ "fmt"
	_ "log"
)

var (
	// WhoamiPath is the API endpoint for whoami
	WhoamiPath = "/api/whoami"
)

type WhoamiResult struct {
	User            *WhoamiUserResult         `json:"user"`
	IsMasterAccount bool                      `json:"isMasterAccount"`
	Permissions     *[]WhoamiPermissionObject `json:"permissions"`
	Appliance       WhoamiApplianceResult     `json:"appliance"`
}

type WhoamiPermissionObject struct {
	Name   string `json:"name"`
	Code   string `json:"code"`
	Access string `json:"access"`
}

type WhoamiUserResult struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}
type WhoamiApplianceResult struct {
	BuildVersion string `json:"buildVersion"`
}

func (client *Client) Whoami() (*Response, error) {
	whoamiRequest := &Request{
		Method:  "GET",
		Path:    WhoamiPath,
		Timeout: 10,
		Result:  &WhoamiResult{},
	}

	resp, err := client.Execute(whoamiRequest)
	// whoamiResult := resp.Result.(*WhoamiResult)
	// if (whoamiResult.User != nil) {
	// 	if  whoamiResult.User.Username != client.Username {
	// 		client.SetUsername(whoamiResult.User.Username)
	// 	}
	// 	var perms []string
	// 	for k,v := range whoamiResult.Permissions {
	// 		perms = append(perms, fmt.Sprintf("%v (%v)", k, v))
	// 	}
	// 	log.Printf("Permissions: %v", strings.Join(perms, ", "))
	// 	log.Printf("Access Token: %v", client.AccessToken)
	// } else {
	// 	err = errors.New("Unable to parse whoami result from api response")
	// 	log.Fatalf(err)
	// }
	// client.setLastWhoamiResult(whoamiResult)
	return resp, err
}
