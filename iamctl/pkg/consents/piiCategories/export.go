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

package piiCategories

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/wso2-extensions/identity-tools-cli/iamctl/pkg/utils"
)

func ExportAll(exportFilePath string, format string) error {

	log.Println("Exporting PII categories...")
	exportFilePath = filepath.Join(exportFilePath, utils.PII_CATEGORIES.String())

	if _, err := os.Stat(exportFilePath); os.IsNotExist(err) {
		os.MkdirAll(exportFilePath, 0700)
	} else {
		if utils.TOOL_CONFIGS.AllowDelete {
			deployedCatNames := getDeployedPIICategories()
			utils.RemoveDeletedLocalResources(exportFilePath, deployedCatNames)
		}
	}

	piiCategories, err := getPIICategorieList()
	if err != nil {
		log.Println("Error: when exporting PII categories.", err)
	} else {
		for _, cat := range piiCategories {
			err := exportPIICategory(cat.Id, cat.Category, exportFilePath, format)
			if err != nil {
				log.Printf("Error while exporting PII category: %s. %s", cat.Category, err)
			}
		}
	}
}

func exportPIICategory(categoryId, category, outputDirPath, formatString string) error {

	scope, err := utils.GetResourceData(utils.PII_CATEGORIES, categoryId)
	if err != nil {
		return fmt.Errorf("error while getting PII category: %w", err)
	}

	format := utils.FormatFromString(formatString)
	exportedFileName := utils.GetExportedFilePath(outputDirPath, category, format)

	scopeKeywordMapping := getOidcScopeKeywordMapping(scopeName)
	modifiedScope, err := utils.ProcessExportedData(scope, exportedFileName, format, scopeKeywordMapping, utils.OIDC_SCOPES)
	if err != nil {
		return fmt.Errorf("error while processing exported content: %w", err)
	}

	modifiedFile, err := utils.Serialize(modifiedScope, format, utils.OIDC_SCOPES)
	if err != nil {
		return fmt.Errorf("error while serializing scope: %w", err)
	}

	err = os.WriteFile(exportedFileName, modifiedFile, 0644)
	if err != nil {
		return fmt.Errorf("error when writing exported content to file: %w", err)
	}

	return nil
}
