/*
* Tencent is pleased to support the open source community by making TKEStack
* available.
*
* Copyright (C) 2012-2019 Tencent. All Rights Reserved.
*
* Licensed under the Apache License, Version 2.0 (the “License”); you may not use
* this file except in compliance with the License. You may obtain a copy of the
* License at
*
* https://opensource.org/licenses/Apache-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an “AS IS” BASIS, WITHOUT
* WARRANTIES OF ANY KIND, either express or implied.  See the License for the
* specific language governing permissions and limitations under the License.
 */
package plugins

import (
	"tkestack.io/kube-jarvis/pkg/logger"
	"tkestack.io/kube-jarvis/pkg/store"
	"tkestack.io/kube-jarvis/pkg/translate"
)

// MetaData is the common attributes of a plugins
type MetaData struct {
	// Store is the global storage
	Store store.Store
	// Translator is a translator with plugins module context
	Translator translate.Translator
	// Logger is a logger with plugins module context
	Logger logger.Logger
	// Type is the type of plugins
	Type string
	// Title is the custom name of plugins
	Name string
	// Desc is the description of this plugin,
	// Desc should be set by plugin
	Desc translate.Message
}

// IsSupportedCloud return true if cloud type is supported
func IsSupportedCloud(supported []string, cloud string) bool {
	if len(supported) == 0 {
		return true
	}

	for _, c := range supported {
		if c == cloud {
			return true
		}
	}
	return false
}
