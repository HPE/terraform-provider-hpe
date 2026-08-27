package morpheus

var (
	// ActivityPath is the API endpoint for activity
	ActivityPath = "/api/activity"
)

// Activity structures for use in request and response payloads
type Activity struct {
	ID           string `json:"_id"`
	Success      bool   `json:"success"`
	ActivityType string `json:"activityType"`
	Name         string `json:"name"`
	Message      string `json:"message"`
	ObjectType   string `json:"objectType"`
	ObjectId     int64  `json:"objectId"`
	User         struct {
		ID       int64  `json:"id"`
		UserName string `json:"username"`
	} `json:"user"`
	TS string `json:"ts"`
}

// GetActivityResult structure parses the list alerts response payload
type GetActivityResult struct {
	Activity *[]Activity `json:"activity"`
	Meta     *MetaResult `json:"meta"`
}

// GetActivity get activity
func (client *Client) GetActivity(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        ActivityPath,
		QueryParams: req.QueryParams,
		Result:      &GetActivityResult{},
	})
}
