// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package trust

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	morpheus "github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"

	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/convert"
	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/helpers"
)

func ResourceKeyPair() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus key pair resource.",
		CreateContext: resourceKeyPairCreate,
		ReadContext:   resourceKeyPairRead,
		DeleteContext: resourceKeyPairDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the key pair",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "Name of the key pair",
				ForceNew:    true,
				Required:    true,
			},
			"public_key": {
				Type:        schema.TypeString,
				Description: "The public key of the key pair",
				ForceNew:    true,
				Required:    true,
			},
			"private_key": {
				Type:        schema.TypeString,
				Description: "The private key of the key pair",
				ForceNew:    true,
				Optional:    true,
				Sensitive:   true,
				StateFunc: func(v any) string {
					h := sha256.New()
					var str string
					if s, ok := v.(string); ok {
						str = s
					}
					h.Write([]byte(str))
					sha256Hash := hex.EncodeToString(h.Sum(nil))

					return sha256Hash
				},
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					h := sha256.New()
					h.Write([]byte(new))
					sha256Hash := hex.EncodeToString(h.Sum(nil))

					return strings.EqualFold(old, sha256Hash)
				},
			},
			"passphrase": {
				Type:        schema.TypeString,
				Description: "The passphrase for the private key of the key pair",
				ForceNew:    true,
				Optional:    true,
				Sensitive:   true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceKeyPairCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

	var diags diag.Diagnostics

	var name string
	if v, ok := d.Get("name").(string); ok {
		name = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("name", d.Get("name")))
	}

	var publicKey string
	if v, ok := d.Get("public_key").(string); ok {
		publicKey = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("public_key", d.Get("public_key")))
	}

	var privateKey string
	if v, ok := d.Get("private_key").(string); ok {
		privateKey = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("private_key", d.Get("private_key")))
	}

	var passphrase string
	if v, ok := d.Get("passphrase").(string); ok {
		passphrase = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("passphrase", d.Get("passphrase")))
	}

	keyPairPayload := make(map[string]any)
	keyPairPayload["name"] = name
	keyPairPayload["publicKey"] = publicKey
	keyPairPayload["privateKey"] = privateKey
	keyPairPayload["passphrase"] = passphrase

	req := &morpheus.Request{
		Body: map[string]any{
			"keyPair": keyPairPayload,
		},
	}

	resp, err := client.CreateKeyPair(req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)

		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	if resp.Result == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Result"))
	}

	var result *morpheus.CreateKeyPairResult
	if v, ok := resp.Result.(*morpheus.CreateKeyPairResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("Result", resp.Result))
	}

	if result.KeyPair == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("KeyPair"))
	}

	keyPair := result.KeyPair
	d.SetId(convert.Int64ToString(keyPair.ID))

	diags = append(diags, resourceKeyPairRead(ctx, d, meta)...)

	return diags
}

func resourceKeyPairRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

	var diags diag.Diagnostics

	var name string
	if v, ok := d.Get("name").(string); ok {
		name = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("name", d.Get("name")))
	}

	id := convert.StringToInt64(d.Id())

	var resp *morpheus.Response
	var err error
	if (id != 0 && name == "") || (id != 0 && name != "") {
		resp, err = client.GetKeyPair(id)
	} else if id == 0 && name != "" {
		resp, err = client.GetKeyPairByName(name)
	} else if id == 0 && name == "" {
		return diag.Errorf("Key pair cannot be read without name or id")
	}

	if err != nil {
		return diag.FromErr(err)
	}

	if resp.Result == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Result"))
	}

	var keyPair *morpheus.KeyPair
	if id != 0 {
		var result *morpheus.GetKeyPairResult
		if v, ok := resp.Result.(*morpheus.GetKeyPairResult); ok {
			result = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("Result", resp.Result))
		}

		if result.KeyPair == nil {
			return diag.FromErr(helpers.NotFoundInResponseError("KeyPair"))
		}

		keyPair = result.KeyPair
	} else if name != "" {
		var listResult *morpheus.ListKeyPairsResult
		if v, ok := resp.Result.(*morpheus.ListKeyPairsResult); ok {
			listResult = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("Result", resp.Result))
		}

		if listResult.KeyPairs == nil {
			return diag.FromErr(helpers.NotFoundInResponseError("KeyPairs"))
		}

		keyPairs := listResult.KeyPairs
		if len(*keyPairs) == 0 {
			return diag.FromErr(helpers.EmptySliceError("KeyPairs"))
		}

		keyPair = &(*keyPairs)[0]
	}

	if keyPair != nil {
		d.SetId(convert.Int64ToString(keyPair.ID))
		d.Set("name", keyPair.Name)
		d.Set("public_key", keyPair.PublicKey)
		d.Set("private_key", keyPair.PrivateKeyHash)
	} else {
		return diag.Errorf("Key pair not found in response data.")
	}

	return diags
}

func resourceKeyPairDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

	var diags diag.Diagnostics

	var idStr string
	if v, ok := d.Get("id").(string); ok {
		idStr = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("id", d.Get("id")))
	}

	id := convert.StringToInt64(idStr)

	req := &morpheus.Request{}
	resp, err := client.DeleteKeyPair(id, req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)

		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)
	d.SetId("")

	return diags
}
