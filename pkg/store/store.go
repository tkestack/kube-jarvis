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
package store

import "fmt"

var SpaceNotFound = fmt.Errorf("space not found")

// Store provider a default k/v storage for all plugins
type Store interface {
	// Complete do Initialize
	Complete() error
	// CreateSpace create a new namespace for specific data set
	CreateSpace(name string) (created bool, err error)
	// Set update a value of key
	Set(space string, key, value string) error
	// Get return target value of key
	Get(space string, key string) (value string, exist bool, err error)
	// Delete delete target key
	Delete(space string, key string) error
	// DeleteSpace Delete whole namespace
	DeleteSpace(name string) error
}

var factories = map[string]func() Store{}

func registerStore(typ string, creator func() Store) {
	factories[typ] = creator
}

func GetStore(typ string) Store {
	f, exsit := factories[typ]
	if !exsit {
		panic(fmt.Sprintf("cant not found store type %s", typ))
	}
	return f()
}
