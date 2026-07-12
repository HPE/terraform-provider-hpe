// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package client

// ScriptCategory models

// ScriptCategory represents an RBA category for creation, update, and retrieval.
// The API returns a tree structure where each category may have child categories.
type ScriptCategory struct {
	Uuid   string                   `json:"uuid,omitempty"`
	Name   string                   `json:"name"`
	Parent *ScriptCategoryParentRef `json:"parent,omitempty"`
	Childs []ScriptCategory         `json:"childs,omitempty"`
}

// ScriptCategoryParentRef is the nested parent object used in create/update requests.
type ScriptCategoryParentRef struct {
	Uuid string `json:"uuid"`
}
