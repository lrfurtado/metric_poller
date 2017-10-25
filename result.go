package main

import (
	"encoding/json"
	"io"
)

//Sample json payload
//{"Application":"Cache0","Version":"1.1.2","Uptime":8069351542,"Request_Count":6664489160,"Error_Count":2221104722,"Success_Count":4443384438},

type Result struct {
	Application  string
	Version      string
	Uptime       int64
	RequestCount int64   `json:"Request_Count"`
	ErrorCount   int64   `json:"Error_Count"`
	SuccessCount int64   `json:"Success_Count"`
	SuccessRate  float64 `json:Success_Rate",omitempty"`
}

func parseResult(body io.ReadCloser) (*Result, error) {
	defer body.Close()
	decoder := json.NewDecoder(body)

	var payload Result
	err := decoder.Decode(&payload)
	if err != nil {
		return nil, err
	}

	return &payload, nil

}
