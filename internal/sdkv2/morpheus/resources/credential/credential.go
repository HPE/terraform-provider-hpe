// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package credential

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/convert"
	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/helpers"

	morpheus "github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"
)

const (
	credentialTypeAccessKeySecret         = "access-key-secret" //nolint: gosec
	credentialTypeAPIKey                  = "api-key"
	credentialTypeClientIDSecret          = "client-id-secret"
	credentialTypeEmailPrivateKey         = "email-private-key"
	credentialTypeTenantUsernameKeypair   = "tenant-username-keypair"
	credentialTypeUsernameAPIKey          = "username-api-key"
	credentialTypeUsernameKeypair         = "username-keypair"
	credentialTypeUsernamePassword        = "username-password"
	credentialTypeUsernamePasswordKeypair = "username-password-keypair"
)

func ResourceContactCredential() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus credential resource.",
		CreateContext: resourceContactCredentialCreate,
		ReadContext:   resourceContactCredentialRead,
		UpdateContext: resourceContactCredentialUpdate,
		DeleteContext: resourceContactCredentialDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the credential",
				Computed:    true,
			},
			"credential_store_integration_id": {
				Type:        schema.TypeInt,
				Description: "The ID of the credential store integration",
				Optional:    true,
				ForceNew:    true,
			},
			"type": {
				Type:        schema.TypeString,
				Description: "The credential type",
				Required:    true,
				ForceNew:    true,
				ValidateFunc: validation.StringInSlice([]string{
					credentialTypeAccessKeySecret,
					credentialTypeAPIKey,
					credentialTypeClientIDSecret,
					credentialTypeEmailPrivateKey,
					credentialTypeTenantUsernameKeypair,
					credentialTypeUsernamePassword,
					credentialTypeUsernameAPIKey,
					credentialTypeUsernameKeypair,
					credentialTypeUsernamePasswordKeypair,
				}, false),
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the credential",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "The description of the credential",
				Optional:    true,
				Computed:    true,
			},
			"enabled": {
				Type:        schema.TypeBool,
				Description: "Whether the credential is enabled",
				Optional:    true,
			},
			"access_key": {
				Type:        schema.TypeString,
				Description: "The credential access key",
				Optional:    true,
			},
			"secret_key": {
				Type:        schema.TypeString,
				Description: "The credential secret key",
				Optional:    true,
				Sensitive:   true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					h := sha256.New()
					h.Write([]byte(new))
					sha256Hash := hex.EncodeToString(h.Sum(nil))

					return strings.EqualFold(old, sha256Hash)
				},
			},
			"client_id": {
				Type:        schema.TypeString,
				Description: "The credential client id",
				Optional:    true,
			},
			"client_secret": {
				Type:        schema.TypeString,
				Description: "The credential client secret",
				Optional:    true,
				Sensitive:   true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					h := sha256.New()
					h.Write([]byte(new))
					sha256Hash := hex.EncodeToString(h.Sum(nil))

					return strings.EqualFold(old, sha256Hash)
				},
			},
			"tenant": {
				Type:        schema.TypeString,
				Description: "The credential tenant",
				Optional:    true,
			},
			"email": {
				Type:        schema.TypeString,
				Description: "The credential email address",
				Optional:    true,
			},
			"username": {
				Type:        schema.TypeString,
				Description: "The credential username",
				Optional:    true,
			},
			"password": {
				Type:        schema.TypeString,
				Description: "The credential password",
				Optional:    true,
				Sensitive:   true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					h := sha256.New()
					h.Write([]byte(new))
					sha256Hash := hex.EncodeToString(h.Sum(nil))

					return strings.EqualFold(old, sha256Hash)
				},
			},
			"api_key": {
				Type:        schema.TypeString,
				Description: "The credential api key",
				Optional:    true,
				Sensitive:   true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					h := sha256.New()
					h.Write([]byte(new))
					sha256Hash := hex.EncodeToString(h.Sum(nil))

					return strings.EqualFold(old, sha256Hash)
				},
			},
			"key_pair_id": {
				Type:        schema.TypeInt,
				Description: "The ID of the credential key pair",
				Optional:    true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceContactCredentialCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

	var diags diag.Diagnostics
	credential := make(map[string]any)

	var name string
	if v, ok := d.Get("name").(string); ok {
		name = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("name", d.Get("name")))
	}
	credential["name"] = name

	var description string
	if v, ok := d.Get("description").(string); ok {
		description = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("description", d.Get("description")))
	}
	credential["description"] = description

	var enabled bool
	if v, ok := d.Get("enabled").(bool); ok {
		enabled = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("enabled", d.Get("enabled")))
	}
	credential["enabled"] = enabled

	integration := make(map[string]any)

	var credentialStoreIntegrationID int
	if v, ok := d.Get("credential_store_integration_id").(int); ok {
		credentialStoreIntegrationID = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError(
			"credential_store_integration_id",
			d.Get("credential_store_integration_id")))
	}
	if credentialStoreIntegrationID != 0 {
		integration["id"] = credentialStoreIntegrationID
	}
	credential["integration"] = integration

	var credType string
	if v, ok := d.Get("type").(string); ok {
		credType = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("type", d.Get("type")))
	}

	switch credType {
	case credentialTypeAccessKeySecret:
		credential["type"] = credentialTypeAccessKeySecret

		var accessKey string
		if v, ok := d.Get("access_key").(string); ok {
			accessKey = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("access_key", d.Get("access_key")))
		}
		credential["username"] = accessKey

		var secretKey string
		if v, ok := d.Get("secret_key").(string); ok {
			secretKey = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("secret_key", d.Get("secret_key")))
		}
		credential["password"] = secretKey
	case credentialTypeAPIKey:
		credential["type"] = credentialTypeAPIKey

		var apiKey string
		if v, ok := d.Get("api_key").(string); ok {
			apiKey = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("api_key", d.Get("api_key")))
		}
		credential["password"] = apiKey
	case credentialTypeClientIDSecret:
		credential["type"] = credentialTypeClientIDSecret

		var clientID string
		if v, ok := d.Get("client_id").(string); ok {
			clientID = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("client_id", d.Get("client_id")))
		}
		credential["username"] = clientID

		var clientSecret string
		if v, ok := d.Get("client_secret").(string); ok {
			clientSecret = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("client_secret", d.Get("client_secret")))
		}
		credential["password"] = clientSecret
	case credentialTypeEmailPrivateKey:
		credential["type"] = credentialTypeEmailPrivateKey

		var email string
		if v, ok := d.Get("email").(string); ok {
			email = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("email", d.Get("email")))
		}
		credential["username"] = email

		keypair := make(map[string]any)

		var keyPairID int
		if v, ok := d.Get("key_pair_id").(int); ok {
			keyPairID = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("key_pair_id", d.Get("key_pair_id")))
		}
		keypair["id"] = keyPairID
		credential["authKey"] = keypair
	case credentialTypeTenantUsernameKeypair:
		credential["type"] = credentialTypeTenantUsernameKeypair

		var tenant string
		if v, ok := d.Get("tenant").(string); ok {
			tenant = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("tenant", d.Get("tenant")))
		}
		credential["authPath"] = tenant

		var username string
		if v, ok := d.Get("username").(string); ok {
			username = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("username", d.Get("username")))
		}
		credential["username"] = username

		keypair := make(map[string]any)

		var keyPairID int
		if v, ok := d.Get("key_pair_id").(int); ok {
			keyPairID = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("key_pair_id", d.Get("key_pair_id")))
		}
		keypair["id"] = keyPairID
		credential["authKey"] = keypair
	case credentialTypeUsernameAPIKey:
		credential["type"] = credentialTypeUsernameAPIKey

		var username string
		if v, ok := d.Get("username").(string); ok {
			username = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("username", d.Get("username")))
		}
		credential["username"] = username

		var apiKey string
		if v, ok := d.Get("api_key").(string); ok {
			apiKey = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("api_key", d.Get("api_key")))
		}
		credential["password"] = apiKey
	case credentialTypeUsernameKeypair:
		credential["type"] = credentialTypeUsernameKeypair

		var username string
		if v, ok := d.Get("username").(string); ok {
			username = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("username", d.Get("username")))
		}
		credential["username"] = username

		keypair := make(map[string]any)

		var keyPairID int
		if v, ok := d.Get("key_pair_id").(int); ok {
			keyPairID = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("key_pair_id", d.Get("key_pair_id")))
		}
		keypair["id"] = keyPairID
		credential["authKey"] = keypair
	case credentialTypeUsernamePassword:
		credential["type"] = credentialTypeUsernamePassword

		var username string
		if v, ok := d.Get("username").(string); ok {
			username = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("username", d.Get("username")))
		}
		credential["username"] = username

		var password string
		if v, ok := d.Get("password").(string); ok {
			password = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("password", d.Get("password")))
		}
		credential["password"] = password
	case credentialTypeUsernamePasswordKeypair:
		credential["type"] = credentialTypeUsernamePasswordKeypair

		var username string
		if v, ok := d.Get("username").(string); ok {
			username = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("username", d.Get("username")))
		}
		credential["username"] = username

		var password string
		if v, ok := d.Get("password").(string); ok {
			password = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("password", d.Get("password")))
		}
		credential["password"] = password

		keypair := make(map[string]any)

		var keyPairID int
		if v, ok := d.Get("key_pair_id").(int); ok {
			keyPairID = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("key_pair_id", d.Get("key_pair_id")))
		}
		keypair["id"] = keyPairID
		credential["authKey"] = keypair
	}

	req := &morpheus.Request{
		Body: map[string]any{
			"credential": credential,
		},
	}

	resp, err := client.CreateCredential(req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)

		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	if resp.Result == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Result"))
	}

	var result *morpheus.CreateCredentialResult
	if v, ok := resp.Result.(*morpheus.CreateCredentialResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("Result", resp.Result))
	}

	if result.Credential == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Credential"))
	}

	contact := result.Credential
	d.SetId(convert.Int64ToString(contact.ID))

	diags = append(diags, resourceContactCredentialRead(ctx, d, meta)...)

	return diags
}

func resourceContactCredentialRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

	var diags diag.Diagnostics

	id := d.Id()

	var name string
	if v, ok := d.Get("name").(string); ok {
		name = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("name", d.Get("name")))
	}

	var resp *morpheus.Response
	var err error
	if id == "" && name != "" {
		resp, err = client.FindCredentialByName(name)
	} else if id != "" {
		resp, err = client.GetCredential(convert.StringToInt64(id), &morpheus.Request{})
	} else {
		return diag.Errorf("Credential cannot be read without name or id")
	}

	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("API 404: %s - %s", resp, err)
			log.Printf("Forcing recreation of resource")
			d.SetId("")

			return diags
		}

		log.Printf("API FAILURE: %s - %s", resp, err)

		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	if resp.Result == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Result"))
	}

	var result *morpheus.GetCredentialResult
	if v, ok := resp.Result.(*morpheus.GetCredentialResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("Result", resp.Result))
	}

	if result.Credential == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Credential"))
	}

	credential := result.Credential
	d.SetId(convert.Int64ToString(credential.ID))
	d.Set("name", credential.Name)
	d.Set("description", credential.Description)
	d.Set("enabled", credential.Enabled)

	if credential.Integration.ID != 0 {
		d.Set("credential_store_integration_id", credential.Integration.ID)
	}

	switch credential.Type.Code {
	case credentialTypeAccessKeySecret:
		d.Set("access_key", credential.Username)
		d.Set("secret_key", credential.PasswordHash)
	case credentialTypeAPIKey:
		d.Set("api_key", credential.PasswordHash)
	case credentialTypeClientIDSecret:
		d.Set("client_id", credential.Username)
		d.Set("client_secret", credential.PasswordHash)
	case credentialTypeEmailPrivateKey:
		d.Set("email", credential.Username)
		d.Set("key_pair_id", credential.AuthKey.ID)
	case credentialTypeTenantUsernameKeypair:
		d.Set("tenant", credential.AuthPath)
		d.Set("username", credential.Username)
		d.Set("key_pair_id", credential.AuthKey.ID)
	case credentialTypeUsernameAPIKey:
		d.Set("username", credential.Username)
		d.Set("api_key", credential.PasswordHash)
	case credentialTypeUsernameKeypair:
		d.Set("username", credential.Username)
		d.Set("key_pair_id", credential.AuthKey.ID)
	case credentialTypeUsernamePassword:
		d.Set("username", credential.Username)
		d.Set("password", credential.PasswordHash)
	case credentialTypeUsernamePasswordKeypair:
		d.Set("username", credential.Username)
		d.Set("password", credential.PasswordHash)
		d.Set("key_pair_id", credential.AuthKey.ID)
	}

	return diags
}

func resourceContactCredentialUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

	id := d.Id()
	credential := make(map[string]any)

	var name string
	if v, ok := d.Get("name").(string); ok {
		name = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("name", d.Get("name")))
	}
	credential["name"] = name

	var description string
	if v, ok := d.Get("description").(string); ok {
		description = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("description", d.Get("description")))
	}
	credential["description"] = description

	var enabled bool
	if v, ok := d.Get("enabled").(bool); ok {
		enabled = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("enabled", d.Get("enabled")))
	}
	credential["enabled"] = enabled

	var credType string
	if v, ok := d.Get("type").(string); ok {
		credType = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("type", d.Get("type")))
	}

	switch credType {
	case credentialTypeAccessKeySecret:
		credential["type"] = credentialTypeAccessKeySecret

		var accessKey string
		if v, ok := d.Get("access_key").(string); ok {
			accessKey = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("access_key", d.Get("access_key")))
		}
		credential["username"] = accessKey

		var secretKey string
		if v, ok := d.Get("secret_key").(string); ok {
			secretKey = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("secret_key", d.Get("secret_key")))
		}
		credential["password"] = secretKey
	case credentialTypeAPIKey:
		credential["type"] = credentialTypeAPIKey

		var apiKey string
		if v, ok := d.Get("api_key").(string); ok {
			apiKey = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("api_key", d.Get("api_key")))
		}
		credential["password"] = apiKey
	case credentialTypeClientIDSecret:
		credential["type"] = credentialTypeClientIDSecret

		var clientID string
		if v, ok := d.Get("client_id").(string); ok {
			clientID = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("client_id", d.Get("client_id")))
		}
		credential["username"] = clientID

		var clientSecret string
		if v, ok := d.Get("client_secret").(string); ok {
			clientSecret = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("client_secret", d.Get("client_secret")))
		}
		credential["password"] = clientSecret
	case credentialTypeEmailPrivateKey:
		credential["type"] = credentialTypeEmailPrivateKey

		var email string
		if v, ok := d.Get("email").(string); ok {
			email = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("email", d.Get("email")))
		}
		credential["username"] = email

		keypair := make(map[string]any)

		var keyPairID int
		if v, ok := d.Get("key_pair_id").(int); ok {
			keyPairID = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("key_pair_id", d.Get("key_pair_id")))
		}
		keypair["id"] = keyPairID
		credential["authKey"] = keypair
	case credentialTypeTenantUsernameKeypair:
		credential["type"] = credentialTypeTenantUsernameKeypair

		var tenant string
		if v, ok := d.Get("tenant").(string); ok {
			tenant = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("tenant", d.Get("tenant")))
		}
		credential["authPath"] = tenant

		var username string
		if v, ok := d.Get("username").(string); ok {
			username = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("username", d.Get("username")))
		}
		credential["username"] = username

		keypair := make(map[string]any)

		var keyPairID int
		if v, ok := d.Get("key_pair_id").(int); ok {
			keyPairID = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("key_pair_id", d.Get("key_pair_id")))
		}
		keypair["id"] = keyPairID
		credential["authKey"] = keypair
	case credentialTypeUsernameAPIKey:
		credential["type"] = credentialTypeUsernameAPIKey

		var username string
		if v, ok := d.Get("username").(string); ok {
			username = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("username", d.Get("username")))
		}
		credential["username"] = username

		var apiKey string
		if v, ok := d.Get("api_key").(string); ok {
			apiKey = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("api_key", d.Get("api_key")))
		}
		credential["password"] = apiKey
	case credentialTypeUsernameKeypair:
		credential["type"] = credentialTypeUsernameKeypair

		var username string
		if v, ok := d.Get("username").(string); ok {
			username = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("username", d.Get("username")))
		}
		credential["username"] = username

		keypair := make(map[string]any)

		var keyPairID int
		if v, ok := d.Get("key_pair_id").(int); ok {
			keyPairID = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("key_pair_id", d.Get("key_pair_id")))
		}
		keypair["id"] = keyPairID
		credential["authKey"] = keypair
	case credentialTypeUsernamePassword:
		credential["type"] = credentialTypeUsernamePassword

		var username string
		if v, ok := d.Get("username").(string); ok {
			username = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("username", d.Get("username")))
		}
		credential["username"] = username

		var password string
		if v, ok := d.Get("password").(string); ok {
			password = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("password", d.Get("password")))
		}
		credential["password"] = password
	case credentialTypeUsernamePasswordKeypair:
		credential["type"] = credentialTypeUsernamePasswordKeypair

		var username string
		if v, ok := d.Get("username").(string); ok {
			username = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("username", d.Get("username")))
		}
		credential["username"] = username

		var password string
		if v, ok := d.Get("password").(string); ok {
			password = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("password", d.Get("password")))
		}
		credential["password"] = password

		keypair := make(map[string]any)

		var keyPairID int
		if v, ok := d.Get("key_pair_id").(int); ok {
			keyPairID = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("key_pair_id", d.Get("key_pair_id")))
		}
		keypair["id"] = keyPairID
		credential["authKey"] = keypair
	}

	req := &morpheus.Request{
		Body: map[string]any{
			"credential": credential,
		},
	}

	resp, err := client.UpdateCredential(convert.StringToInt64(id), req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)

		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	if resp.Result == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Result"))
	}

	var result *morpheus.UpdateCredentialResult
	if v, ok := resp.Result.(*morpheus.UpdateCredentialResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("Result", resp.Result))
	}

	if result.Credential == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Credential"))
	}

	contact := result.Credential
	d.SetId(convert.Int64ToString(contact.ID))

	return resourceContactCredentialRead(ctx, d, meta)
}

func resourceContactCredentialDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

	var diags diag.Diagnostics

	id := d.Id()
	req := &morpheus.Request{}
	resp, err := client.DeleteCredential(convert.StringToInt64(id), req)
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("API 404: %s - %s", resp, err)

			return diag.FromErr(err)
		}

		log.Printf("API FAILURE: %s - %s", resp, err)

		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)
	d.SetId("")

	return diags
}
