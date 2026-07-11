// SPDX-FileCopyrightText: Copyright Hewlett Packard Enterprise Development LP
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"encoding/json"
	"fmt"
)

// ── KB Categories ─────────────────────────────────────────────────────────────

// ListKBCategories returns all KB categories for a tenant.
func (c *OpsRampClient) ListKBCategories(tenantId string) ([]KBCategory, error) {
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/kb/categorylist", c.BaseUrl, tenantId)

	body, err := c.NewJsonRequest("GET", apiUrl, nil)
	if err != nil {
		return nil, err
	}

	var responseBody []KBCategory
	err = json.Unmarshal([]byte(body), &responseBody)
	if err != nil {
		return nil, err
	}

	return responseBody, nil
}

// GetKBCategory retrieves a specific KB category by ID from the full list.
func (c *OpsRampClient) GetKBCategory(tenantId string, categoryId string) (*KBCategory, error) {
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/kb/category/%s", c.BaseUrl, tenantId, categoryId)
	body, err := c.NewJsonRequest("GET", apiUrl, nil)
	if err != nil {
		return nil, err
	}

	var responseBody KBCategory
	err = json.Unmarshal([]byte(body), &responseBody)
	if err != nil {
		return nil, err
	}

	return &responseBody, nil

}

// CreateKBCategory creates a new KB category.
func (c *OpsRampClient) CreateKBCategory(tenantId string, categoryData KBCategory) (*KBCategory, error) {
	rb, err := json.Marshal(categoryData)
	if err != nil {
		return nil, err
	}

	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/kb/category/create", c.BaseUrl, tenantId)
	body, err := c.NewJsonRequest("POST", apiUrl, rb)
	if err != nil {
		return nil, err
	}

	var responseBody KBCategory
	err = json.Unmarshal([]byte(body), &responseBody)
	if err != nil {
		return nil, err
	}

	return &responseBody, nil
}

// UpdateKBCategory updates an existing KB category.
func (c *OpsRampClient) UpdateKBCategory(tenantId string, categoryId string, categoryData KBCategory) (*KBCategory, error) {
	rb, err := json.Marshal(categoryData)
	if err != nil {
		return nil, err
	}

	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/kb/category/update/%s", c.BaseUrl, tenantId, categoryId)

	body, err := c.NewJsonRequest("POST", apiUrl, rb)
	if err != nil {
		return nil, err
	}

	var responseBody KBCategory
	err = json.Unmarshal([]byte(body), &responseBody)
	if err != nil {
		return nil, err
	}

	return &responseBody, nil
}

// DeleteKBCategory deletes a KB category by ID.
func (c *OpsRampClient) DeleteKBCategory(tenantId string, categoryId string) error {
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/kb/category/delete/%s", c.BaseUrl, tenantId, categoryId)

	_, err := c.NewJsonRequest("DELETE", apiUrl, nil)
	return err
}

// ── KB Articles ───────────────────────────────────────────────────────────────

// CreateKBArticle creates a new KB article.
func (c *OpsRampClient) CreateKBArticle(tenantId string, articleData KBArticle) (*KBArticle, error) {
	rb, err := json.Marshal(articleData)
	if err != nil {
		return nil, err
	}

	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/kb/article", c.BaseUrl, tenantId)

	body, err := c.NewJsonRequest("POST", apiUrl, rb)
	if err != nil {
		return nil, err
	}

	var responseBody KBArticle
	err = json.Unmarshal([]byte(body), &responseBody)
	if err != nil {
		return nil, err
	}

	return &responseBody, nil
}

// GetKBArticle retrieves a specific KB article by ID.
func (c *OpsRampClient) GetKBArticle(tenantId string, articleId string) (*KBArticle, error) {
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/kb/article/%s", c.BaseUrl, tenantId, articleId)

	body, err := c.NewJsonRequest("GET", apiUrl, nil)
	if err != nil {
		return nil, err
	}

	var responseBody KBArticle
	err = json.Unmarshal([]byte(body), &responseBody)
	if err != nil {
		return nil, err
	}

	return &responseBody, nil
}

func (c *OpsRampClient) ShareKBArticles(tenantId string, kbArticleId string) (*KBArticle, error) {

	articleData := struct {
		Id string `json:"id"`
	}{
		Id: kbArticleId,
	}

	rb, err := json.Marshal(articleData)
	if err != nil {
		return nil, err
	}

	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/kb/article/%s/share", c.BaseUrl, tenantId, kbArticleId)

	body, err := c.NewJsonRequest("POST", apiUrl, rb)
	if err != nil {
		return nil, err
	}

	var responseBody KBArticle
	err = json.Unmarshal([]byte(body), &responseBody)
	if err != nil {
		return nil, err
	}

	return &responseBody, nil
}

// UpdateKBArticle updates an existing KB article.
func (c *OpsRampClient) UpdateKBArticle(tenantId string, articleId string, articleData KBArticle) (*KBArticle, error) {
	rb, err := json.Marshal(articleData)
	if err != nil {
		return nil, err
	}

	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/kb/article/%s", c.BaseUrl, tenantId, articleId)

	body, err := c.NewJsonRequest("POST", apiUrl, rb)
	if err != nil {
		return nil, err
	}

	var responseBody KBArticle
	err = json.Unmarshal([]byte(body), &responseBody)
	if err != nil {
		return nil, err
	}

	return &responseBody, nil
}

// DeleteKBArticle deletes a KB article by ID.
func (c *OpsRampClient) DeleteKBArticle(tenantId string, articleId string) error {
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/kb/article/%s/delete", c.BaseUrl, tenantId, articleId)

	_, err := c.NewJsonRequest("DELETE", apiUrl, nil)
	return err
}
