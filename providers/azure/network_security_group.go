// Copyright 2019 The Terraformer Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package azure

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2019-08-01/network"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/GoogleCloudPlatform/terraformer/terraform_utils"
)

type NetworkSecurityGroupGenerator struct {
	AzureService
}

func (g NetworkSecurityGroupGenerator) createResources(securityGroupListResultPage network.SecurityGroupListResultPage) []terraform_utils.Resource {
	var resources []terraform_utils.Resource
	for securityGroupListResultPage.NotDone() {
		nsgs := securityGroupListResultPage.Values()
		for _, nsg := range nsgs {
			resources = append(resources, terraform_utils.NewSimpleResource(
				*nsg.ID,
				*nsg.Name,
				"azurerm_network_security_group",
				"azurerm",
				[]string{}))
		}
		if err := securityGroupListResultPage.Next(); err != nil {
			log.Println(err)
			break
		}
	}
	return resources
}

func (g *NetworkSecurityGroupGenerator) InitResources() error {
	ctx := context.Background()
	securityGroupsClient := network.NewSecurityGroupsClient(g.Args["subscription"].(string))
	authorizer, err := auth.NewAuthorizerFromEnvironment()
	if err != nil {
		return err
	}
	securityGroupsClient.Authorizer = authorizer
	output, err := securityGroupsClient.ListAll(ctx)
	if err != nil {
		return err
	}
	g.Resources = g.createResources(output)
	return nil
}
