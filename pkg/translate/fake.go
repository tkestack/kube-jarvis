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
package translate

// Fake is a Translator that just return ID as message directly
type Fake struct {
}

// NewFake return a Fake Translator
func NewFake() Translator {
	return &Fake{}
}

// Message get translated message from Translator
// t.module will be add before ID
// example:
//         ID = "message"  and module = "diagnostics.example"
//         then real ID will be "diagnostics.example.message"
func (f *Fake) Message(ID string, templateData map[string]interface{}) Message {
	return Message(ID)
}

// WithModule attach a module label to a Translator
// module will be add before ID when you call Translator.Message
func (f *Fake) WithModule(module string) Translator {
	return f
}
