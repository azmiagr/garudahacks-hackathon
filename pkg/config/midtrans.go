package config

import (
	"os"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/midtrans/midtrans-go/snap"
)

type MidtransConfig struct {
	ServerKey   string
	ClientKey   string
	Environment midtrans.EnvironmentType
	IsSandbox   bool
}

func LoadMidtransConfig() *MidtransConfig {
	serverKey := os.Getenv("MIDTRANS_SERVER_KEY")
	clientKey := os.Getenv("MIDTRANS_CLIENT_KEY")
	isSandbox := os.Getenv("MIDTRANS_ENVIRONMENT") == "sandbox"

	env := midtrans.Production
	if isSandbox {
		env = midtrans.Sandbox
	}

	return &MidtransConfig{
		ServerKey:   serverKey,
		ClientKey:   clientKey,
		Environment: env,
		IsSandbox:   isSandbox,
	}
}

func (c *MidtransConfig) NewSnapClient() snap.Client {
	var client snap.Client
	client.New(c.ServerKey, c.Environment)
	return client
}

func (c *MidtransConfig) NewCoreAPIClient() coreapi.Client {
	var client coreapi.Client
	client.New(c.ServerKey, c.Environment)
	return client
}
