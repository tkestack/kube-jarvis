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
package logger

import (
	"fmt"
	"log"
	"sort"
)

type loggerInfo struct {
	labels map[string]string
}

// NewLogger create a logger that just print logs using golang fmt logger
func NewLogger() Logger {
	return &loggerInfo{
		labels: map[string]string{},
	}
}

func (l *loggerInfo) With(labels map[string]string) Logger {
	nLogger := &loggerInfo{
		labels: map[string]string{},
	}
	for k, v := range l.labels {
		nLogger.labels[k] = v
	}

	for k, v := range labels {
		nLogger.labels[k] = v
	}

	return nLogger
}

func (l *loggerInfo) Message(prefix string, format string, args ...interface{}) string {
	message := prefix + " "
	message += fmt.Sprintf(format, args...)
	message += "  "

	keys := make([]string, 0)
	for k := range l.labels {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		message += fmt.Sprintf("%s = %s | ", k, l.labels[k])
	}

	return message
}

func (l *loggerInfo) Infof(format string, args ...interface{}) {
	log.Println(l.Message("[INFO]", format, args...))
}
func (l *loggerInfo) Debugf(format string, args ...interface{}) {
	log.Println(l.Message("[DEBUG]", format, args...))

}
func (l *loggerInfo) Errorf(format string, args ...interface{}) {
	log.Println(l.Message("[ERROR]", format, args...))
}
