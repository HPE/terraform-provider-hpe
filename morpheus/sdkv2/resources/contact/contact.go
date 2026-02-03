// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package contact

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	morpheus "github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"

	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/convert"
	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/helpers"
)

func ResourceContact() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus contact resource.",
		CreateContext: resourceContactCreate,
		ReadContext:   resourceContactRead,
		UpdateContext: resourceContactUpdate,
		DeleteContext: resourceContactDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the contact",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the contact",
				Required:    true,
			},
			"email_address": {
				Type:        schema.TypeString,
				Description: "The email address associated with the contact",
				Optional:    true,
				Computed:    true,
			},
			"mobile_number": {
				Type:        schema.TypeString,
				Description: "The mobile phone number associated with the contact",
				Optional:    true,
				Computed:    true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceContactCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics

	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

	var name string
	if v, ok := d.Get("name").(string); ok {
		name = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("name", d.Get("name")))
	}

	var emailAddress string
	if v, ok := d.Get("email_address").(string); ok {
		emailAddress = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("email_address", d.Get("email_address")))
	}

	var mobileNumber string
	if v, ok := d.Get("mobile_number").(string); ok {
		mobileNumber = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("mobile_number", d.Get("mobile_number")))
	}

	req := &morpheus.Request{
		Body: map[string]any{
			"contact": map[string]any{
				"name":         name,
				"emailAddress": emailAddress,
				"smsAddress":   mobileNumber,
			},
		},
	}

	resp, err := client.CreateContact(req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)

		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	if resp.Result == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Result"))
	}

	var result *morpheus.CreateContactResult
	if v, ok := resp.Result.(*morpheus.CreateContactResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("Result", resp.Result))
	}

	if result.Contact == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Contact"))
	}

	contact := result.Contact
	d.SetId(convert.Int64ToString(contact.ID))

	diags = append(diags, resourceContactRead(ctx, d, meta)...)

	return diags
}

func resourceContactRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics

	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

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
		resp, err = client.FindContactByName(name)
	} else if id != "" {
		resp, err = client.GetContact(convert.StringToInt64(id), &morpheus.Request{})
	} else {
		return diag.Errorf("Contact cannot be read without name or id")
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

	var result *morpheus.GetContactResult
	if v, ok := resp.Result.(*morpheus.GetContactResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("Result", resp.Result))
	}

	if result.Contact == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Contact"))
	}

	contact := result.Contact
	d.SetId(convert.Int64ToString(contact.ID))
	d.Set("name", contact.Name)
	d.Set("email_address", contact.EmailAddress)
	d.Set("mobile_number", contact.SmsAddress)

	return diags
}

func resourceContactUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

	id := d.Id()

	var name string
	if v, ok := d.Get("name").(string); ok {
		name = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("name", d.Get("name")))
	}

	var emailAddress string
	if v, ok := d.Get("email_address").(string); ok {
		emailAddress = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("email_address", d.Get("email_address")))
	}

	var mobileNumber string
	if v, ok := d.Get("mobile_number").(string); ok {
		mobileNumber = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("mobile_number", d.Get("mobile_number")))
	}

	req := &morpheus.Request{
		Body: map[string]any{
			"contact": map[string]any{
				"name":         name,
				"emailAddress": emailAddress,
				"smsAddress":   mobileNumber,
			},
		},
	}

	resp, err := client.UpdateContact(convert.StringToInt64(id), req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)

		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	if resp.Result == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Result"))
	}

	var result *morpheus.UpdateContactResult
	if v, ok := resp.Result.(*morpheus.UpdateContactResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("Result", resp.Result))
	}

	if result.Contact == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Contact"))
	}

	contact := result.Contact
	d.SetId(convert.Int64ToString(contact.ID))

	return resourceContactRead(ctx, d, meta)
}

func resourceContactDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics

	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

	id := d.Id()
	req := &morpheus.Request{}
	resp, err := client.DeleteContact(convert.StringToInt64(id), req)
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
