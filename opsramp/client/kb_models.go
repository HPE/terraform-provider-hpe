// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package client

// KB (Knowledge Base) models

// KBCategory represents a knowledge base category.
type KBCategory struct {
	Id             string         `json:"id,omitempty"`
	Name           string         `json:"name"`
	Scope          string         `json:"scope,omitempty"`
	Description    string         `json:"description,omitempty"`
	ParentCategory *KBCategoryRef `json:"parentCategory,omitempty"`
	State          string         `json:"state,omitempty"`
}

// KBCategoryRef is a minimal reference to a category (used in nested fields).
type KBCategoryRef struct {
	Id string `json:"id"`
}

type KBArticleAttachment struct {
	Id string `json:"id"`
}

type KBArticleRef struct {
	Id string `json:"id"`
}

// KBArticle represents a knowledge base article.
type KBArticle struct {
	Id             string                `json:"id,omitempty"`
	Subject        string                `json:"subject"`
	Content        string                `json:"content"`
	Category       *KBCategoryRef        `json:"category,omitempty"`
	State          string                `json:"state,omitempty"`
	Attachments    []KBArticleAttachment `json:"attachments,omitempty"`
	ExpiryDate     string                `json:"expiryDate,omitempty"`
	LinkedArticles []KBArticleRef        `json:"linkedArticles,omitempty"`
}
