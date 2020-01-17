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

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

// Mysql save data in mysql
type Mysql struct {
	// Url is the Connection URL of target mysql
	Url         string
	clusterName string
}

// Data is the table for storing data
type Data struct {
	ID      uint64 `gorm:"primary_key;AUTO_INCREMENT;column:id"`
	Cluster string `gorm:"column:cluster;index;type:varchar(30)"`
	Space   string `gorm:"column:space;index;type:varchar(30)"`
	Key     string `gorm:"column:key_name;index;type:varchar(30)"`
	Value   string `gorm:"column:value;type:longtext"`
}

// TableName is the table name in mysql
func (d *Data) TableName() string {
	return "global_store"
}

func init() {
	registerStore("mysql", func(clusterName string) Store {
		return &Mysql{
			clusterName: clusterName,
		}
	})
}

// Complete do Initialize
func (m *Mysql) Complete() error {
	if m.Url == "" {
		return fmt.Errorf("config.url must be set")
	}

	db, err := gorm.Open("mysql", m.Url)
	if err != nil {
		return err
	}
	defer func() { _ = db.Close() }()

	if db.AutoMigrate(&Data{}).Error != nil {
		return err
	}

	return nil
}

// CreateSpace create a new namespace for specific data set
func (m *Mysql) CreateSpace(name string) (created bool, err error) {
	return false, nil
}

// Set update a value of key
func (m *Mysql) Set(space string, key, value string) error {
	db, err := gorm.Open("mysql", m.Url)
	if err != nil {
		return err
	}
	defer func() { _ = db.Close() }()

	d := &Data{}
	notFound := false
	if err := db.Where("space = ? AND key_name = ? AND cluster = ?",
		space, key, m.clusterName).Find(d).Error; err != nil {
		if !gorm.IsRecordNotFoundError(err) {
			return err
		} else {
			notFound = true
		}
	}

	d.Space = space
	d.Key = key
	d.Value = value
	d.Cluster = m.clusterName

	if notFound {
		return db.Create(d).Error
	}

	return db.Save(d).Error
}

// Get return target value of key
func (m *Mysql) Get(space string, key string) (value string, exist bool, err error) {
	db, err := gorm.Open("mysql", m.Url)
	if err != nil {
		return "", false, err
	}
	defer func() { _ = db.Close() }()

	d := &Data{}
	if err := db.Where("space = ? AND key_name = ? AND cluster = ?",
		space, key, m.clusterName).Find(d).Error; err != nil {
		if !gorm.IsRecordNotFoundError(err) {
			return "", false, err
		} else {
			return "", false, nil
		}
	}

	return d.Value, true, nil
}

// Delete delete target key
func (m *Mysql) Delete(space string, key string) error {
	db, err := gorm.Open("mysql", m.Url)
	if err != nil {
		return err
	}
	defer func() { _ = db.Close() }()

	return db.Delete(Data{}, "space = ? and key_name = ? and cluster = ?",
		space, key, m.clusterName).Error
}

// DeleteSpace Delete whole namespace
func (m *Mysql) DeleteSpace(name string) error {
	db, err := gorm.Open("mysql", m.Url)
	if err != nil {
		return err
	}
	defer func() { _ = db.Close() }()
	return db.Delete(Data{}, "space = ? and cluster = ? ",
		name, m.clusterName).Error
}
