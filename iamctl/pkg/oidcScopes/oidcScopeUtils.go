/**
* Copyright (c) 2026, WSO2 LLC. (https://www.wso2.com).
*
* WSO2 LLC. licenses this file to you under the Apache License,
* Version 2.0 (the "License"); you may not use this file except
* in compliance with the License.
* You may obtain a copy of the License at
*
* http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing,
* software distributed under the License is distributed on an
* "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
* KIND, either express or implied. See the License for the
* specific language governing permissions and limitations
* under the License.
 */

package oidcScopes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/wso2-extensions/identity-tools-cli/iamctl/pkg/utils"
)

type oidcScope struct {
	Name        string   `json:"name"`
	DisplayName string   `json:"displayName"`
	Description string   `json:"description,omitempty"`
	Claims      []string `json:"claims"`
}

type oidcScopeList struct {
	TotalResults int         `json:"totalResults"`
	Scopes       []oidcScope `json:"scopes"`
}

type oidcScopeConfig struct {
	Name string `yaml:"name"`
}

func getOidcScopeList() ([]oidcScope, error) {

	scopeCount, err := getTotalOidcScopeCount()
	if err != nil {
		log.Println("Error: when retrieving OIDC scope count. Retrieving only the default count.", err)
	}
	var list oidcScopeList
	resp, err := utils.SendGetListRequest(utils.OIDC_SCOPES, scopeCount)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve available OIDC scope list. %w", err)
	}
	defer resp.Body.Close()

	statusCode := resp.StatusCode
	if statusCode == 200 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("error when reading the retrieved OIDC scope list. %w", err)
		}

		err = json.Unmarshal(body, &list)
		if err != nil {
			return nil, fmt.Errorf("error when unmarshalling the retrieved OIDC scope list. %w", err)
		}
		resp.Body.Close()

		return list.Scopes, nil
	} else if error, ok := utils.ErrorCodes[statusCode]; ok {
		return nil, fmt.Errorf("error while retrieving OIDC scope list. Status code: %d, Error: %s", statusCode, error)
	}
	return nil, fmt.Errorf("error while retrieving OIDC scope list")
}

func getTotalOidcScopeCount() (count int, err error) {

	var list oidcScopeList
	resp, err := utils.SendGetListRequest(utils.OIDC_SCOPES, -1)
	if err != nil {
		return -1, fmt.Errorf("failed to retrieve available OIDC scope list. %w", err)
	}
	defer resp.Body.Close()

	statusCode := resp.StatusCode
	if statusCode == 200 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return -1, fmt.Errorf("error when reading the retrieved OIDC scope list. %w", err)
		}

		err = json.Unmarshal(body, &list)
		if err != nil {
			return -1, fmt.Errorf("error when unmarshalling the retrieved OIDC scope list. %w", err)
		}
		resp.Body.Close()

		return list.TotalResults, nil
	} else if error, ok := utils.ErrorCodes[statusCode]; ok {
		return -1, fmt.Errorf("error while retrieving OIDC scope count. Status code: %d, Error: %s", statusCode, error)
	}
	return -1, fmt.Errorf("error while retrieving OIDC scope count")
}

func getDeployedOidcScopeNames() []string {

	scopes, err := getOidcScopeList()
	if err != nil {
		return []string{}
	}

	var scopeNames []string
	for _, scope := range scopes {
		scopeNames = append(scopeNames, scope.Name)
	}
	return scopeNames
}

func getOidcScopeKeywordMapping(scopeName string) map[string]interface{} {

	if utils.KEYWORD_CONFIGS.OidcScopeConfigs != nil {
		return utils.ResolveAdvancedKeywordMapping(scopeName, utils.KEYWORD_CONFIGS.OidcScopeConfigs)
	}
	return utils.KEYWORD_CONFIGS.KeywordMappings
}
