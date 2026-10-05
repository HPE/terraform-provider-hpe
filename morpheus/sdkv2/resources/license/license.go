// (C) Copyright 2025-2026 Hewlett Packard Enterprise Development LP

package license

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	morpheus "github.com/HPE/terraform-provider-hpe/internal/sdk/legacy"
	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/helpers"
)

// Morpheus never returns a licence key, but /api/license reports each
// installed licence with its keyId (the key's first eight characters) and hash
// (the SHA-256 of the key), so the resource can tell whether the key it is
// configured with is the one installed. The key is sent trimmed, as the API
// trims it before hashing (MORPH-18310).

func ResourceLicense() *schema.Resource {
	return &schema.Resource{
		Description: "Provides a Morpheus license resource: installs the licence key on the appliance and " +
			"keeps it installed.\n\n" +
			"The resource reads the installed licences back and compares them with `key`: when the " +
			"configured key is no longer installed -- the appliance was rebuilt, or someone installed a " +
			"different key -- the next plan re-installs it. The licence's details are exposed as read-only " +
			"attributes.\n\n" +
			"~> **Destroying this resource uninstalls the licence** (the one matching `key`; other " +
			"installed licences are kept), leaving the appliance unlicensed if it was the only one. To stop " +
			"managing a licence without uninstalling it, remove it from state instead: " +
			"`terraform state rm hpe_morpheus_license.<name>`.",
		CreateContext: resourceLicenseCreate,
		ReadContext:   resourceLicenseRead,
		UpdateContext: resourceLicenseUpdate,
		DeleteContext: resourceLicenseDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type: schema.TypeString,
				Description: "The licence's key id: the first eight characters of the key, which is how Morpheus " +
					"identifies an installed licence (`keyId` in `GET /api/license`, `DELETE /api/license?keyId=`). " +
					"Stable for as long as the same key is used -- unlike Morpheus's internal record id, which " +
					"changes each time the key is re-installed.",
				Computed: true,
			},
			"key": {
				Type:             schema.TypeString,
				Description:      "The Morpheus license key",
				Required:         true,
				Sensitive:        true,
				DiffSuppressFunc: keyDiffSuppress,
			},
			"product_tier": {
				Type:        schema.TypeString,
				Description: "The product tier the licence grants (for example `enterprise`, `capacity`)",
				Computed:    true,
			},
			"account_name": {
				Type:        schema.TypeString,
				Description: "The account the licence was issued to",
				Computed:    true,
			},
			"start_date": {
				Type:        schema.TypeString,
				Description: "When the licence becomes valid (RFC 3339)",
				Computed:    true,
			},
			"end_date": {
				Type:        schema.TypeString,
				Description: "When the licence expires (RFC 3339); empty for a perpetual licence",
				Computed:    true,
			},
			"max_instances": {
				Type:        schema.TypeInt,
				Description: "The maximum number of instances the licence allows; 0 when unlimited",
				Computed:    true,
			},
			"max_memory": {
				Type:        schema.TypeInt,
				Description: "The maximum managed memory, in bytes, the licence allows; 0 when unlimited",
				Computed:    true,
			},
			"max_storage": {
				Type:        schema.TypeInt,
				Description: "The maximum managed storage, in bytes, the licence allows; 0 when unlimited",
				Computed:    true,
			},
			"hard_limit": {
				Type:        schema.TypeBool,
				Description: "Whether Morpheus refuses to provision beyond the licence's limits rather than warning",
				Computed:    true,
			},
			"free_trial": {
				Type:        schema.TypeBool,
				Description: "Whether the licence is a free trial",
				Computed:    true,
			},
			"multi_tenant": {
				Type:        schema.TypeBool,
				Description: "Whether the licence allows sub-tenants",
				Computed:    true,
			},
			"support_level": {
				Type:        schema.TypeString,
				Description: "The support level the licence carries",
				Computed:    true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// licenseKeyHash is the identifier Morpheus reports for an installed licence:
// the hex SHA-256 of the trimmed key (ApplianceLicenseService.getLicenseKeyHash).
func licenseKeyHash(key string) string {
	sum := sha256.Sum256([]byte(strings.TrimSpace(key)))

	return hex.EncodeToString(sum[:])
}

// licenseKeyID is the short identifier Morpheus uses in the API for an
// installed licence: the first eight characters of the trimmed key.
func licenseKeyID(key string) string {
	k := strings.TrimSpace(key)
	if len(k) > 8 {
		return k[:8]
	}

	return k
}

// keyDiffSuppress is the DiffSuppressFunc of key. Morpheus trims the key
// before storing it, so a difference of surrounding whitespace is no change.
// Everything else is a change -- in particular an empty prior value, which is
// what Read leaves behind when the configured key is not installed, so that
// the next plan re-installs it; and a different key, which installs the new
// one.
func keyDiffSuppress(_, oldValue, newValue string, _ *schema.ResourceData) bool {
	return oldValue != "" && strings.TrimSpace(oldValue) == strings.TrimSpace(newValue)
}

// findInstalledLicense returns the installed licence matching the hash (when
// hash is set) or else the key id, looking at the active licence and then at
// every installed one, or nil when none matches.
func findInstalledLicense(result *morpheus.GetLicenseResult, hash, keyID string) *morpheus.ApplianceLicense {
	if result == nil {
		return nil
	}

	matches := func(l *morpheus.ApplianceLicense) bool {
		if l == nil {
			return false
		}

		if hash != "" {
			return strings.EqualFold(l.Hash, hash)
		}

		return keyID != "" && l.KeyID == keyID
	}

	if matches(result.License) {
		return result.License
	}

	for i := range result.InstalledLicenses {
		if matches(&result.InstalledLicenses[i]) {
			return &result.InstalledLicenses[i]
		}
	}

	return nil
}

// setLicenseAttributes writes the installed licence's details into state.
func setLicenseAttributes(d *schema.ResourceData, l *morpheus.ApplianceLicense) error {
	formatTime := func(t time.Time) string {
		if t.IsZero() {
			return ""
		}

		return t.UTC().Format(time.RFC3339)
	}

	values := map[string]any{
		"product_tier":  l.Producttier,
		"account_name":  l.Accountname,
		"start_date":    formatTime(l.Startdate),
		"end_date":      formatTime(l.Enddate),
		"max_instances": l.Maxinstances,
		"max_memory":    l.Maxmemory,
		"max_storage":   l.Maxstorage,
		"hard_limit":    l.Hardlimit,
		"free_trial":    l.Freetrial,
		"multi_tenant":  l.Multitenant,
		"support_level": l.Supportlevel,
	}
	for k, v := range values {
		if err := d.Set(k, v); err != nil {
			return err
		}
	}

	return nil
}

// clearLicenseAttributes empties the installed licence's details when the
// configured key is not installed.
func clearLicenseAttributes(d *schema.ResourceData) error {
	for _, k := range []string{"product_tier", "account_name", "start_date", "end_date", "support_level"} {
		if err := d.Set(k, ""); err != nil {
			return err
		}
	}

	for _, k := range []string{"max_instances", "max_memory", "max_storage"} {
		if err := d.Set(k, 0); err != nil {
			return err
		}
	}

	for _, k := range []string{"hard_limit", "free_trial", "multi_tenant"} {
		if err := d.Set(k, false); err != nil {
			return err
		}
	}

	return nil
}

func clientFromMeta(meta any) (*morpheus.Client, error) {
	if v, ok := meta.(*morpheus.Client); ok {
		return v, nil
	}

	return nil, helpers.TypeAssertFailError("client", meta)
}

func keyFromData(d *schema.ResourceData) (string, error) {
	if v, ok := d.Get("key").(string); ok {
		return strings.TrimSpace(v), nil
	}

	return "", helpers.TypeAssertFailError("key", d.Get("key"))
}

// installLicense posts the key. Posting a key that is already installed is a
// no-op on the appliance, so Create and Update share it.
func installLicense(client *morpheus.Client, key string) error {
	req := &morpheus.Request{
		Body: map[string]any{
			"license": key,
		},
	}

	resp, err := client.InstallLicense(req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)

		return err
	}

	log.Printf("API RESPONSE: %s", resp)

	if resp.Result == nil {
		return helpers.NotFoundInResponseError("Result")
	}

	if _, ok := resp.Result.(*morpheus.GetLicenseResult); !ok {
		return helpers.TypeAssertFailError("Result", resp.Result)
	}

	return nil
}

func resourceLicenseCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client, err := clientFromMeta(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	key, err := keyFromData(d)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := installLicense(client, key); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(licenseKeyID(key))

	return resourceLicenseRead(ctx, d, meta)
}

func resourceLicenseRead(_ context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client, err := clientFromMeta(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	var diags diag.Diagnostics

	resp, err := client.GetLicense(&morpheus.Request{})
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

	result, ok := resp.Result.(*morpheus.GetLicenseResult)
	if !ok {
		return diag.FromErr(helpers.TypeAssertFailError("Result", resp.Result))
	}

	// Match the configured key by its hash. On import there is no key yet, so
	// the licence is matched by the key id the resource was imported with; the
	// key is then set by the first apply (posting an installed key is a no-op).
	key, err := keyFromData(d)
	if err != nil {
		return diag.FromErr(err)
	}

	hash := ""
	if key != "" {
		hash = licenseKeyHash(key)
	}

	installed := findInstalledLicense(result, hash, d.Id())
	if installed == nil {
		// The configured key is not installed (the appliance was rebuilt and
		// runs its setup licence, someone installed a different key, or the
		// licence was uninstalled). Clearing key in state makes the next plan
		// an in-place update that re-installs the configured key.
		log.Printf("license %s: not among the installed licences; the next apply re-installs it", d.Id())

		if err := d.Set("key", ""); err != nil {
			return diag.FromErr(err)
		}

		if err := clearLicenseAttributes(d); err != nil {
			return diag.FromErr(err)
		}

		return diags
	}

	d.SetId(installed.KeyID)

	if err := setLicenseAttributes(d, installed); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceLicenseUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client, err := clientFromMeta(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	key, err := keyFromData(d)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := installLicense(client, key); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(licenseKeyID(key))

	return resourceLicenseRead(ctx, d, meta)
}

func resourceLicenseDelete(_ context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client, err := clientFromMeta(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	var diags diag.Diagnostics

	// Uninstall this licence only: without a key id the API deactivates every
	// installed licence, including ones this resource never managed.
	resp, err := client.UninstallLicense(d.Id(), &morpheus.Request{})
	if err != nil {
		if resp != nil && (resp.StatusCode == 404 || strings.Contains(string(resp.Body), "No active license found")) {
			// Already gone.
			log.Printf("API %d: %s - %s", resp.StatusCode, resp, err)
			d.SetId("")

			return diags
		}

		log.Printf("API FAILURE: %s - %s", resp, err)

		return diag.FromErr(err)
	}

	log.Printf("API RESPONSE: %s", resp)
	d.SetId("")

	return diags
}
