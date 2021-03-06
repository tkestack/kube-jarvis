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
package util

import (
	"testing"
	"time"
)

func TestRetryUntilTimeout(t *testing.T) {
	// direct return
	done := false
	go func() {
		<-time.After(time.Second * 3)
		if !done {
			t.Fatalf("not done")
		}
	}()
	if err := RetryUntilTimeout(time.Hour, time.Hour, func() error {
		done = true
		return nil
	}); err != nil {
		t.Fatalf(err.Error())
	}

	// check retry
	count := 0
	if err := RetryUntilTimeout(0, 0, func() error {
		count++
		if count == 3 {
			return nil
		}
		return RetryAbleErr
	}); err != nil {
		t.Fatalf(err.Error())
	}

	// check timeout
	if err := RetryUntilTimeout(time.Second, time.Second*2, func() error {
		return RetryAbleErr
	}); err == nil {
		t.Fatalf("should return an error")
	}
}
