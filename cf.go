package main

import "net/http"

type cf struct {
	config map[string]string
	client *http.Client
}

func newCF() (*cf, error) {
	m, err := load()
	if err != nil {
		return nil, err
	}
	return &cf{config: m}, err
}
