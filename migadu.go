package main

import (
	"github.com/MrLemur/migadu-go"
)

func NewMigaduClient(domain string, adminEmail string, APIKey string) *migadu.Client {

	client, err := migadu.New(adminEmail, APIKey, domain)
	if err != nil {
		panic(err)
	}

	return client
}
