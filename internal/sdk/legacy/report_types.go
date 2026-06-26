package morpheus

var (
	// ReportTypesPath is the API endpoint for report types
	ReportTypesPath = "/api/report-types"
)

// ReportType structures for use in request and response payloads
type ReportType struct {
	ID                   int64  `json:"id"`
	Code                 string `json:"code"`
	Name                 string `json:"name"`
	Description          string `json:"description"`
	Category             string `json:"category"`
	Visible              bool   `json:"visible"`
	MasterOnly           bool   `json:"masterOnly"`
	OwnerOnly            bool   `json:"ownerOnly"`
	SupportsAllZoneTypes bool   `json:"supportsAllZoneTypes"`
	IsPlugin             bool   `json:"isPlugin"`
	DateCreated          string `json:"dateCreated"`
	OptionTypes          []struct {
		ID   int64  `json:"id"`
		Code string `json:"code"`
		Name string `json:"name"`
	} `json:"optionTypes"`
	SupportedZoneTypes []struct {
		ID   int64  `json:"id"`
		Code string `json:"code"`
		Name string `json:"name"`
	} `json:"supportedZoneTypes"`
}

type ListReportTypesResult struct {
	ReportTypes *[]ReportType `json:"reportTypes"`
	Meta        *MetaResult   `json:"meta"`
}

// Client request methods
func (client *Client) ListReportTypes(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        ReportTypesPath,
		QueryParams: req.QueryParams,
		Result:      &ListReportTypesResult{},
	})
}

// FindReportTypeByName gets an existing report type by name
func (client *Client) FindReportTypeByName(name string) (*Response, error) {
	// Find by name
	resp, err := client.ListReportTypes(&Request{
		QueryParams: map[string]string{
			"name": name,
		},
	})
	return resp, err
}
