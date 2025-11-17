// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

// Package job provides a Morpheus job data source for Terraform.
package job

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	morpheus "github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"

	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/convert"
	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/helpers"
)

// DataSourceJob returns the schema for the Morpheus job data source.
func DataSourceJob() *schema.Resource {
	return &schema.Resource{
		Description: "Provides a Morpheus job data source.",
		ReadContext: dataSourceJobRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:          schema.TypeInt,
				Description:   "The ID of the job",
				Optional:      true,
				ConflictsWith: []string{"name"},
				Computed:      true,
			},
			"name": {
				Type:          schema.TypeString,
				Description:   "The name of the job",
				Optional:      true,
				ConflictsWith: []string{"id"},
			},
		},
	}
}

//nolint:revive // ctx is required by the ReadContext signature
func dataSourceJobRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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

	var id int
	if v, ok := d.Get("id").(int); ok {
		id = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("id", d.Get("id")))
	}

	var resp *morpheus.Response
	var err error
	if id == 0 && name != "" {
		resp, err = client.FindJobByName(name)
	} else if id != 0 {
		resp, err = client.GetJob(int64(id), &morpheus.Request{})
	} else {
		return diag.Errorf("Job cannot be read without name or id")
	}
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("API 404: %s - %v", resp, err)

			return nil
		}

		log.Printf("API FAILURE: %s - %v", resp, err)

		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	if resp.Result == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Result"))
	}

	var result *morpheus.GetJobResult
	if v, ok := resp.Result.(*morpheus.GetJobResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("Result", resp.Result))
	}

	if result.Job == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Job"))
	}

	job := result.Job
	d.SetId(convert.Int64ToString(job.ID))
	if err := d.Set("name", job.Name); err != nil {
		return diag.FromErr(err)
	}

	return diags
}
