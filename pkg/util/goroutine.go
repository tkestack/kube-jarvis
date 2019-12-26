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
	"fmt"
	"math"
	"time"
)

// RetryAbleErr should be returned if you want RetryUntilTimeout to retry
var RetryAbleErr = fmt.Errorf("retry")

// RetryUntilTimeout retry target function "do" until  timeout
func RetryUntilTimeout(interval time.Duration, timeout time.Duration, do func() error) error {
	err := do()
	if err == nil {
		return nil
	}

	if err != RetryAbleErr {
		return err
	}

	if timeout == 0 {
		timeout = time.Duration(math.MaxInt64)
	}

	t := time.NewTimer(timeout)
	for {
		select {
		case <-t.C:
			return fmt.Errorf("timeout")
		case <-time.After(interval):
			err := do()
			if err == nil {
				return nil
			}

			if err != RetryAbleErr {
				return err
			}
		}
	}
}
