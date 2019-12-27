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
