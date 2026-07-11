// SPDX-FileCopyrightText: Copyright Hewlett Packard Enterprise Development LP
// SPDX-License-Identifier: Apache-2.0

package client

// Site represents a physical or logical site for grouping resources
type Site struct {
	Id             int64          `json:"id,omitempty"`
	Uuid           string         `json:"uuid,omitempty"`
	Name           string         `json:"name"`
	Path           string         `json:"path"`
	Parent         *SiteParentRef `json:"parent,omitempty"`
	Description    string         `json:"description"`
	Address        string         `json:"address"`
	State          string         `json:"state"`
	City           string         `json:"city"`
	Country        string         `json:"country"`
	Zip            string         `json:"zip"`
	PrimaryContact *SiteContact   `json:"primaryContact,omitempty"`
	PhoneNumber    string         `json:"phoneNumber"`
	PhoneExtension string         `json:"phoneExtension"`
	FilterCriteria *SiteFilter    `json:"filterCriteria,omitempty"`
	Resources      []SiteResource `json:"resources,omitempty"`
}

// SiteChild represents a site entry in a site's child list.
type SiteChild struct {
	Uuid string `json:"uuid"`
	Name string `json:"name,omitempty"`
}

// SiteContact represents the primary contact for a site
type SiteContact struct {
	Id        string `json:"id"`
	LoginName string `json:"loginName,omitempty"`
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
	Email     string `json:"email,omitempty"`
}

// SiteFilter represents filter criteria for site resource selection
type SiteFilter struct {
	SearchQuery string `json:"searchQuery"`
}

// SiteResource represents a resource reference attached to a site.
type SiteResource struct {
	Id string `json:"id,omitempty"`
}

// SiteMinimal represents minimal site info for listing
type SiteMinimal struct {
	Id          int64  `json:"id"`
	Uuid        string `json:"uuid"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Address     string `json:"address,omitempty"`
}

type SiteParentRef struct {
	Id int64 `json:"id"`
}
