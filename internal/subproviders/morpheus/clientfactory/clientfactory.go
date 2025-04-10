// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package clientfactory

import (
	"context"
	"errors"
	"net/http"
	"sync"

	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/client"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/model"
)

var (
	cf          ClientFactory
	once        sync.Once
	initialized bool
)

type ClientFactory struct {
	model.SubModel
}

func GetClientFactory() (ClientFactory, error) {
	if !initialized {
		msg := "morpheus client factory is not initialized"

		return ClientFactory{}, errors.New(msg)
	}

	return cf, nil
}

func SetClientFactory(config model.SubModel) {
	once.Do(func() {
		cf = ClientFactory{config}
		initialized = true
	})
}

func (cf ClientFactory) New(_ context.Context) client.Client {
	return client.Client{Client: http.Client{}}
}
