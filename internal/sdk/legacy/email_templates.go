package morpheus

import (
	"fmt"
)

var (
	// EmailTemplatesPath is the API endpoint for email templates
	EmailTemplatesPath = "/api/email-templates"
)

// EmailTemplate structures for use in request and response payloads
type EmailTemplate struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Owner struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"owner"`
	Accounts []struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"accounts"`
	Code     string `json:"code"`
	Template string `json:"template"`
	Enabled  bool   `json:"enabled"`
}

// ListEmailTemplatesResult structure parses the list email templates response payload
type ListEmailTemplatesResult struct {
	EmailTemplates *[]EmailTemplate `json:"emailTemplates"`
	Meta           *MetaResult      `json:"meta"`
}

// GetEmailTemplateResult structure parses the email template response payload
type GetEmailTemplateResult struct {
	EmailTemplate *EmailTemplate `json:"emailTemplate"`
}

// CreateEmailTemplateResult structure parses the create email template response payload
type CreateEmailTemplateResult struct {
	Success       bool              `json:"success"`
	Message       string            `json:"msg"`
	Errors        map[string]string `json:"errors"`
	EmailTemplate *EmailTemplate    `json:"emailTemplate"`
}

type UpdateEmailTemplateResult struct {
	CreateEmailTemplateResult
}

type DeleteEmailTemplateResult struct {
	DeleteResult
}

// Client request methods

// ListEmailTemplates get all existing email templates
func (client *Client) ListEmailTemplates(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        EmailTemplatesPath,
		QueryParams: req.QueryParams,
		Result:      &ListEmailTemplatesResult{},
	})
}

// GetEmailTemplate gets an existing email template
func (client *Client) GetEmailTemplate(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/%d", EmailTemplatesPath, id),
		QueryParams: req.QueryParams,
		Result:      &GetEmailTemplateResult{},
	})
}

// CreateEmailTemplate creates a new email template
func (client *Client) CreateEmailTemplate(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "POST",
		Path:        EmailTemplatesPath,
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &CreateEmailTemplateResult{},
	})
}

// UpdateEmailTemplate updates an existing email template
func (client *Client) UpdateEmailTemplate(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d", EmailTemplatesPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdateEmailTemplateResult{},
	})
}

// DeleteEmailTemplate deletes an existing email template
func (client *Client) DeleteEmailTemplate(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "DELETE",
		Path:        fmt.Sprintf("%s/%d", EmailTemplatesPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &DeleteEmailTemplateResult{},
	})
}

func (client *Client) FindEmailTemplateByName(name string) (*Response, error) {
	// Find by name, then get by ID
	resp, err := client.ListEmailTemplates(&Request{
		QueryParams: map[string]string{
			"name": name,
		},
	})
	if err != nil {
		return resp, err
	}
	listResult := resp.Result.(*ListEmailTemplatesResult)
	emailTemplateCount := len(*listResult.EmailTemplates)
	if emailTemplateCount != 1 {
		return resp, fmt.Errorf("found %d email templates for %v", emailTemplateCount, name)
	}
	firstRecord := (*listResult.EmailTemplates)[0]
	checkID := firstRecord.ID
	return client.GetEmailTemplate(checkID, &Request{})
}
