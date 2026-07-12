// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package client

// Script models

// Script represents an RBA script for creation, update, and retrieval
type Script struct {
	Uuid           string                   `json:"uuid,omitempty"`
	Name           string                   `json:"name"`
	Description    string                   `json:"description"`
	Category       *ScriptCategoryParentRef `json:"category,omitempty"`
	Platforms      []string                 `json:"platforms"`
	InstallTimeout int                      `json:"installTimeout"`
	Parameters     []ScriptParameter        `json:"parameters,omitempty"`
	ExecutionType  string                   `json:"executionType,omitempty"`
	Attachment     *ScriptAttachment        `json:"attachment,omitempty"`
	ScriptVersion  string                   `json:"scriptVersion,omitempty"`
}
type ScriptCreationResponse struct {
	StatusType string `json:"statusType"` // CREATED, OK
	Entity     string `json:"entity"`
	EntityType string `json:"entityType"`
	Status     int    `json:"status"`
}

// ScriptParameter represents a single parameter of a script
type ScriptParameter struct {
	Id           int    `json:"id,omitempty"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	DefaultValue string `json:"defaultValue"`
	Type         string `json:"type"`     // OPTIONAL, REQUIRED
	DataType     string `json:"dataType"` // INTEGER, STRING, PASSWORD
}

// ScriptAttachment represents the file attachment of a script
type ScriptAttachment struct {
	Id   int    `json:"id,omitempty"`
	Name string `json:"name"`
	File string `json:"file"`
}
