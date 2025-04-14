// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// Client -
type BrickByBrickClient struct {
	HostURL    string
	HTTPClient *http.Client
	Token      string
}

// NewClient -
func NewClient(apiKey *string) (*BrickByBrickClient, error) {
	c := BrickByBrickClient{
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
		Token:      *apiKey,
	}
	// If API key is not provided, return empty client
	if apiKey == nil {
		return &c, nil
	}
	return &c, nil
}

func (c *BrickByBrickClient) doRequest(req *http.Request, apiKey *string) ([]byte, error) {
	token := c.Token

	if apiKey != nil {
		token = *apiKey
	}

	req.Header.Set("api_key", token)

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	return body, err
}

func (c *BrickByBrickClient) GetExercise(exerciseId string) (*Exercise, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://mlsojdnlzcsczxwkeuwy.supabase.co/functions/v1/api/exercises/%s", exerciseId), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req, nil)
	if err != nil {
		return nil, err
	}

	exercise := Exercise{}
	err = json.Unmarshal(body, &exercise)
	if err != nil {
		return nil, err
	}

	return &exercise, nil
}

func (c *BrickByBrickClient) GetExercises() ([]Exercise, error) {
	req, err := http.NewRequest("GET", "https://mlsojdnlzcsczxwkeuwy.supabase.co/functions/v1/api/exercises", nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req, nil)
	if err != nil {
		return nil, err
	}

	exercises := []Exercise{}
	err = json.Unmarshal(body, &exercises)
	if err != nil {
		return nil, err
	}

	return exercises, nil
}

func (c *BrickByBrickClient) CreateExercise(exercise Exercise) (*Exercise, error) {
	rb, err := json.Marshal(exercise)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", "https://mlsojdnlzcsczxwkeuwy.supabase.co/functions/v1/api/exercises", strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req, nil)
	if err != nil {
		return nil, err
	}

	createdExercise := Exercise{}
	err = json.Unmarshal(body, &createdExercise)
	if err != nil {
		fmt.Println("Error when unmarshaling response into an exercise struct.")
		return nil, err
	}

	return &createdExercise, nil
}

func (c *BrickByBrickClient) UpdateExercise(exerciseIdStr string, exercise Exercise) (*Exercise, error) {
	rb, err := json.Marshal(exercise)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("https://mlsojdnlzcsczxwkeuwy.supabase.co/functions/v1/api/exercises/%s", exerciseIdStr), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req, nil)
	if err != nil {
		return nil, err
	}

	updatedExercise := Exercise{}
	err = json.Unmarshal(body, &updatedExercise)
	if err != nil {
		fmt.Println("Error when unmarshaling response into an exercise struct.")
		return nil, err
	}

	return &updatedExercise, nil
}

func (c *BrickByBrickClient) DeleteExercise(exerciseIdStr string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("https://mlsojdnlzcsczxwkeuwy.supabase.co/functions/v1/api/exercises/%s", exerciseIdStr), nil)
	if err != nil {
		return err
	}
	_, err = c.doRequest(req, nil)
	return err
}
