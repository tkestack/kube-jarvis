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
package diagnose

import (
	"fmt"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"tkestack.io/kube-jarvis/pkg/translate"
)

func CommonDeafer(c chan *Result) {
	close(c)
	if err := recover(); err != nil {
		c <- &Result{
			Level:   HealthyLevelFailed,
			ObjName: "*",
			Title:   "Failed",
			Desc:    translate.Message(fmt.Sprintf("%v", err)),
		}
	}
}

type MetaObject interface {
	schema.ObjectKind
	v1.Object
}

func GetRootOwner(obj MetaObject, uid2obj map[types.UID]MetaObject) MetaObject {
	ownerReferences := obj.GetOwnerReferences()
	if len(ownerReferences) > 0 {
		for _, owner := range ownerReferences {
			if owner.Controller != nil && *owner.Controller == true {
				if parent, ok := uid2obj[owner.UID]; ok {
					return GetRootOwner(parent, uid2obj)
				}
			}
		}
	}
	return obj
}
