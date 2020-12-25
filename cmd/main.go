package main

import (
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/containerservice/mgmt/2020-11-01/containerservice"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/tdc-yamada-ya/aks-scheduled-poolscaler/internal/scaler"
)

var logger *log.Logger

func sub() error {
	o := newOptions()
	err := o.validate()
	if err != nil {
		return err
	}

	logger.Printf("load location - timezone: %v", o.timezone)

	loc, err := time.LoadLocation(o.timezone)
	if err != nil {
		return err
	}

	logger.Printf("read configuration - filename: %v", o.configurationFilename)

	configuration, err := readConfiguration(o.configurationFilename)
	if err != nil {
		return err
	}

	logger.Printf("create authorizer from environment")

	authorizer, err := auth.NewAuthorizerFromEnvironment()
	if err != nil {
		return err
	}

	logger.Printf("create agent pools client - subscriptionID: %v", o.subscriptionID)

	client := containerservice.NewAgentPoolsClient(o.subscriptionID)
	client.Authorizer = authorizer
	updater := &scaler.AzPoolUpdater{Client: client}
	scaler := &scaler.Scaler{
		Logger:        logger,
		PoolUpdater:   updater,
		Configuration: configuration,
	}

	t := time.Now().In(loc)

	logger.Printf("start scaling - t: %v", t.Format(time.RFC1123))

	return scaler.Scale(t)
}

func readConfiguration(filename string) (*scaler.Configuration, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return scaler.UnmarshalConfigurationWithYaml(b)
}

func main() {
	logger = log.New(os.Stdout, "[aks-scheduled-poolscaler] ", log.LstdFlags)

	err := sub()
	if err != nil {
		logger.Printf("error: %v", err)
		os.Exit(1)
	}
}
