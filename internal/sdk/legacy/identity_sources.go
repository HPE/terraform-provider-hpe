package morpheus

import (
	"fmt"
	"time"
)

var (
	// IdentitySourcesPath is the API endpoint for identity sources
	IdentitySourcesPath       = "/api/user-sources"
	CreateIdentitySourcesPath = "/api/accounts"
)

// IdentitySource structures for use in request and response payloads
type IdentitySource struct {
	ID                  int64  `json:"id"`
	Name                string `json:"name"`
	Description         string `json:"description"`
	Code                string `json:"code"`
	Type                string `json:"type"`
	Active              bool   `json:"active"`
	Deleted             bool   `json:"deleted"`
	AutoSyncOnLogin     bool   `json:"autoSyncOnLogin"`
	ExternalLogin       bool   `json:"externalLogin"`
	AllowCustomMappings bool   `json:"allowCustomMappings"`
	Account             struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"account"`
	DefaultAccountRole struct {
		ID        int64  `json:"id"`
		Name      string `json:"name"`
		Authority string `json:"authority"`
	} `json:"defaultAccountRole"`
	Config struct {
		URL                            string `json:"url"`
		Domain                         string `json:"domain"`
		UseSSL                         string `json:"useSSL"`
		BindingUsername                string `json:"bindingUsername"`
		BindingPassword                string `json:"bindingPassword"`
		BindingPasswordHash            string `json:"bindingPasswordHash"`
		RequiredGroup                  string `json:"requiredGroup"`
		SearchMemberGroups             bool   `json:"searchMemberGroups"`
		RequiredGroupDN                string `json:"requiredGroupDN"`
		UserFqnExpression              string `json:"userFqnExpression"`
		RequiredRoleFqn                string `json:"requiredRoleFqn"`
		UsernameAttribute              string `json:"usernameAttribute"`
		CommonNameAttribute            string `json:"commonNameAttribute"`
		FirstNameAttribute             string `json:"firstNameAttribute"`
		LastNameAttribute              string `json:"lastNameAttribute"`
		EmailAttribute                 string `json:"emailAttribute"`
		UniqueMemberAttribute          string `json:"uniqueMemberAttribute"`
		MemberOfAttribute              string `json:"memberOfAttribute"`
		OrganizationID                 string `json:"organizationId"`
		RequiredRole                   string `json:"requiredRole"`
		AdministratorAPIToken          string `json:"administratorAPIToken"`
		RequiredGroupID                string `json:"requiredGroupId"`
		Subdomain                      string `json:"subdomain"`
		Region                         string `json:"region"`
		ClientSecret                   string `json:"clientSecret"`
		ClientID                       string `json:"clientId"`
		RequiredRoleID                 string `json:"requiredRoleId"`
		RoleAttributeName              string `json:"roleAttributeName"`
		RequiredAttributeValue         string `json:"requiredAttributeValue"`
		GivenNameAttribute             string `json:"givenNameAttribute"`
		SurnameAttribute               string `json:"surnameAttribute"`
		LogoutURL                      string `json:"logoutUrl"`
		LoginURL                       string `json:"loginUrl"`
		EncryptionAlgorithim           string `json:"encryptionAlgo"`
		EncryptionKey                  string `json:"encryptionKey"`
		APIStyle                       string `json:"apiStyle"`
		DoNotIncludeSAMLRequest        bool   `json:"doNotIncludeSAMLRequest"`
		DoNotValidateSignature         bool   `json:"doNotValidateSignature"`
		DoNotValidateStatusCode        bool   `json:"doNotValidateStatusCode"`
		DoNotValidateDestination       bool   `json:"doNotValidateDestination"`
		DoNotValidateIssueInstants     bool   `json:"doNotValidateIssueInstants"`
		DoNotValidateAssertions        bool   `json:"doNotValidateAssertions"`
		DoNotValidateAuthStatements    bool   `json:"doNotValidateAuthStatements"`
		DoNotValidateSubject           bool   `json:"doNotValidateSubject"`
		DoNotValidateConditions        bool   `json:"doNotValidateConditions"`
		DoNotValidateAudiences         bool   `json:"doNotValidateAudiences"`
		DoNotValidateSubjectRecipients bool   `json:"doNotValidateSubjectRecipients"`
		SAMLSignatureMode              string `json:"SAMLSignatureMode"`
		Endpoint                       string `json:"endpoint"`
		Logout                         string `json:"logout"`
		Request509Certificate          string `json:"request509Certificate"`
		RequestPrivateKey              string `json:"requestPrivateKey"`
		PublicKey                      string `json:"publicKey"`
		PrivateKey                     string `json:"privateKey"`
	} `json:"config"`
	RoleMappings []struct {
		SourceRoleName string `json:"sourceRoleName"`
		SourceRoleFqn  string `json:"sourceRoleFqn"`
		MappedRole     struct {
			ID        int64  `json:"id"`
			Name      string `json:"string"`
			Authority string `json:"authority"`
		} `json:"mappedRole"`
	} `json:"roleMappings"`
	Subdomain        string `json:"subdomain"`
	LoginURL         string `json:"loginURL"`
	ProviderSettings struct {
	} `json:"providerSettings"`
	DateCreated time.Time `json:"dateCreated"`
	LastUpdated time.Time `json:"lastUpdated"`
}

// ListIdentitySourcesResult structure parses the list identity source response payload
type ListIdentitySourcesResult struct {
	IdentitySources *[]IdentitySource `json:"userSources"`
	Meta            *MetaResult       `json:"meta"`
}

type GetIdentitySourceResult struct {
	IdentitySource *IdentitySource `json:"userSource"`
}

type CreateIdentitySourceResult struct {
	Success        bool              `json:"success"`
	Message        string            `json:"msg"`
	Errors         map[string]string `json:"errors"`
	IdentitySource *IdentitySource   `json:"userSource"`
}

type UpdateIdentitySourceResult struct {
	CreateIdentitySourceResult
}

type DeleteIdentitySourceResult struct {
	DeleteResult
}

// Client request methods

func (client *Client) ListIdentitySources(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        IdentitySourcesPath,
		QueryParams: req.QueryParams,
		Result:      &ListIdentitySourcesResult{},
	})
}

func (client *Client) GetIdentitySource(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/%d", IdentitySourcesPath, id),
		QueryParams: req.QueryParams,
		Result:      &GetIdentitySourceResult{},
	})
}

func (client *Client) CreateIdentitySource(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "POST",
		Path:        fmt.Sprintf("%s/%d/user-sources", CreateIdentitySourcesPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &CreateIdentitySourceResult{},
	})
}

func (client *Client) UpdateIdentitySource(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d", IdentitySourcesPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdateIdentitySourceResult{},
	})
}

func (client *Client) UpdateIdentitySourceSubdomain(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d/subdomain", IdentitySourcesPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdateIdentitySourceResult{},
	})
}

func (client *Client) DeleteIdentitySource(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "DELETE",
		Path:        fmt.Sprintf("%s/%d", IdentitySourcesPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &DeleteIdentitySourceResult{},
	})
}

// helper functions
func (client *Client) FindIdentitySourceByName(name string) (*Response, error) {
	// Find by name, then get by ID
	resp, err := client.ListIdentitySources(&Request{
		QueryParams: map[string]string{
			"name": name,
		},
	})
	if err != nil {
		return resp, err
	}
	listResult := resp.Result.(*ListIdentitySourcesResult)
	userSourceCount := len(*listResult.IdentitySources)
	if userSourceCount != 1 {
		return resp, fmt.Errorf("found %d Identity Sources for %v", userSourceCount, name)
	}
	firstRecord := (*listResult.IdentitySources)[0]
	userSourceID := firstRecord.ID
	return client.GetIdentitySource(userSourceID, &Request{})
}
