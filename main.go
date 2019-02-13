package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2018-05-01/resources"
	"github.com/Azure/go-autorest/autorest/azure/auth"
)

type credentials struct {
	ClientID                       string `json:"clientId"`
	ClientSecret                   string `json:"clientSecret"`
	TenantID                       string `json:"tenantId"`
	SubscriptionID                 string `json:"subscriptionId"`
	ActiveDirectoryEndpointURL     string `json:"activeDirectoryEndpointUrl"`
	ResourceManagerEndpointURL     string `json:"resourceManagerEndpointUrl"`
	ActiveDirectoryGraphResourceID string `json:"activeDirectoryGraphResourceId"`
}

func main() {
	azureSecretData, _ := ioutil.ReadFile("azure-provider-key.json")
	creds := credentials{}
	err := json.Unmarshal(azureSecretData, &creds)
	if err != nil {
		msg := fmt.Errorf("failed to unmarshal azure client secret data: %+v", err)
		log.Fatal(msg)
	}

	config := auth.NewClientCredentialsConfig(creds.ClientID, creds.ClientSecret, creds.TenantID)
	config.AADEndpoint = creds.ActiveDirectoryEndpointURL
	config.Resource = creds.ResourceManagerEndpointURL

	authorizer, err := config.Authorizer()
	if err != nil {
		msg := fmt.Errorf("failed to get authorizer from config: %+v", err)
		log.Fatal(msg)
	}
	gc := resources.NewGroupsClient(creds.SubscriptionID)
	gc.Authorizer = authorizer
	group := resources.Group{}
	location := "Central US"
	group.Location = &location
	_, err = gc.CreateOrUpdate(context.TODO(), "my-centralus-rg", group)

	if err != nil {
		log.Fatal(err)
	}
}
