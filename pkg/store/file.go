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

import (
	"fmt"
	"io/ioutil"
	"os"
)

// File save data in files, the data of different space saved in different file
type File struct {
	// Dir is the directory all data files in
	Dir string
}

func init() {
	registerStore("file", func(string) Store {
		return &File{}
	})
}

func (f *File) spacePath(name string) string {
	return fmt.Sprintf("%s/%s", f.Dir, name)
}

func (f *File) dataPath(space, key string) string {
	return fmt.Sprintf("%s/%s/%s", f.Dir, space, key)
}

// Complete do Initialize
func (f *File) Complete() error {
	if f.Dir == "" {
		f.Dir = "data"
	}
	_ = os.MkdirAll(f.Dir, 0755)
	return nil
}

// CreateSpace create a new namespace for specific data set
func (f *File) CreateSpace(name string) (created bool, err error) {
	if exists(f.spacePath(name)) {
		return false, nil
	}

	if err := os.MkdirAll(f.spacePath(name), 0755); err != nil {
		return false, err
	}

	return true, nil
}

// Set update a value of key
func (f *File) Set(space string, key, value string) error {
	if !exists(f.spacePath(space)) {
		return SpaceNotFound
	}

	if err := ioutil.WriteFile(f.dataPath(space, key), []byte(value), 0644); err != nil {
		return err
	}

	return nil
}

// Get return target value of key
func (f *File) Get(space string, key string) (value string, exist bool, err error) {
	if !exists(f.spacePath(space)) {
		return "", false, SpaceNotFound
	}

	if !exists(f.dataPath(space, key)) {
		return "", false, nil
	}

	data, err := ioutil.ReadFile(f.dataPath(space, key))
	if err != nil {
		return "", false, err
	}

	return string(data), true, nil
}

// Delete delete target key
func (f *File) Delete(space string, key string) error {
	if !exists(f.spacePath(space)) {
		return SpaceNotFound
	}

	if !exists(f.dataPath(space, key)) {
		return nil
	}

	return os.Remove(f.dataPath(space, key))
}

// DeleteSpace Delete whole namespace
func (f *File) DeleteSpace(name string) error {
	if !exists(f.spacePath(name)) {
		return SpaceNotFound
	}
	return os.Remove(f.spacePath(name))
}

func exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}
