package main

import (
	"fmt"
	"os"
)

type options struct {
	configurationFilename string
	subscriptionID        string
	timezone              string
}

func (o *options) validate() error {
	if len(o.configurationFilename) == 0 {
		return fmt.Errorf("configuration file name not specified")
	}

	if len(o.subscriptionID) == 0 {
		return fmt.Errorf("azure subscription id name not specified")
	}

	if len(o.timezone) == 0 {
		return fmt.Errorf("timezone not specified")
	}

	return nil
}

func newOptions() *options {
	o := options{}
	o.configurationFilename = os.Getenv("CONFIGURATION_FILE")
	o.subscriptionID = os.Getenv("AZURE_SUBSCRIPTION_ID")
	o.timezone = os.Getenv("TZ")
	return &o
}
