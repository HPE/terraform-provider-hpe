//  (C) Copyright 2021-2024 Hewlett Packard Enterprise Development LP

package client

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var (
	jsonCheck  = regexp.MustCompile("(?i:[application|text]/json)")
	xmlCheck   = regexp.MustCompile("(?i:[application|text]/xml)")
	errVersion = fmt.Errorf("error, API is not supported for the current version, please contact your administrator")
)

func getURLValues(query map[string]string) url.Values {
	m := make(map[string][]string)
	for k, v := range query {
		m[k] = []string{v}
	}

	return m
}

type CustomError struct {
	Errors             string                 `json:"error,omitempty"`
	Body               map[string]interface{} `json:"body,omitempty"`
	StatusCode         int                    `json:"statuscode,omitempty"`
	RecommendedActions []string               `json:"recommendedActions,omitempty"`
}

func (c CustomError) Error() string {
	jsonObj, err := json.Marshal(c)
	if err != nil {
		return err.Error()
	}

	return string(jsonObj)
}

func ParseError(resp *http.Response) error {
	customErr := CustomError{
		StatusCode: resp.StatusCode,
	}
	err := json.NewDecoder(resp.Body).Decode(&customErr.Body)
	if err != nil {
		customErr.Errors = err.Error()
	} else if len(customErr.Body) == 0 {
		customErr.Errors = "No additional information is available"
	}

	if rAction, ok := customErr.Body["recommendedActions"]; ok {
		delete(customErr.Body, "recommendedActions")
		rActions := rAction.([]interface{})
		customErr.RecommendedActions = make([]string, 0, len(rActions))
		for _, a := range rActions {
			customErr.RecommendedActions = append(customErr.RecommendedActions, a.(string))
		}
	}

	return customErr
}

func parseVersion(version string) (int, error) {
	if version == "" {
		return 0, nil
	}
	version = strings.Split(version, "-")[0]
	versionSplit := strings.Split(version, ".")

	mul := 10000
	sum := 0
	for _, v := range versionSplit {
		vInt, err := strconv.Atoi(v)
		if err != nil {
			return 0, err
		}
		sum += mul * vInt
		mul /= 100
	}

	return sum, nil
}

// Add a file to the multipart request
func addFile(writer *multipart.Writer, fieldName, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	part, err := writer.CreateFormFile(fieldName, filepath.Base(path))
	if err != nil {
		return err
	}
	_, err = io.Copy(part, file)

	return err
}

// Set request body from an interface{}
func setBody(body interface{}, contentType string) (bodyBuf *bytes.Buffer, err error) {
	if bodyBuf == nil {
		bodyBuf = &bytes.Buffer{}
	}

	switch bodyType := body.(type) {
	case io.Reader:
		_, err = bodyBuf.ReadFrom(bodyType)
	case []byte:
		_, err = bodyBuf.Write(bodyType)
	case string:
		_, err = bodyBuf.WriteString(bodyType)
	case *string:
		_, err = bodyBuf.WriteString(*bodyType)
	default:
		if jsonCheck.MatchString(contentType) {
			err = json.NewEncoder(bodyBuf).Encode(body)
		} else if xmlCheck.MatchString(contentType) {
			err = xml.NewEncoder(bodyBuf).Encode(body)
		}
	}

	if err != nil {
		return nil, err
	}

	if bodyBuf.Len() == 0 {
		err = fmt.Errorf("invalid body type %s", contentType)

		return nil, err
	}

	return bodyBuf, nil
}

// detectContentType method is used to figure out `Request.Body` content type for request header
func detectContentType(body interface{}) string {
	contentType := "text/plain; charset=utf-8"
	kind := reflect.TypeOf(body).Kind()

	switch kind {
	case reflect.Struct, reflect.Map, reflect.Ptr:
		contentType = "application/json; charset=utf-8"
	case reflect.String:
		contentType = "text/plain; charset=utf-8"
	default:
		if b, ok := body.([]byte); ok {
			contentType = http.DetectContentType(b)
		} else if kind == reflect.Slice {
			contentType = "application/json; charset=utf-8"
		}
	}

	return contentType
}
