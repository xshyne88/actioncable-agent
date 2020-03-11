package actioncable

import (
	"os"
	"reflect"
	"testing"
	"time"
)

func skipCircleCI(t *testing.T) {
	t.Helper()
	if os.Getenv("CIRCLECI") != "true" {
		t.Skip()
	}
}

func assert(act, exp interface{}, t *testing.T) {
	t.Helper()
	if act != exp {
		t.Errorf("Act: %v: Exp: %v", act, exp)
	}
}

func assertDeep(act, exp interface{}, t *testing.T) {
	t.Helper()
	if !reflect.DeepEqual(exp, act) {
		t.Errorf("Act: %v: Exp: %v", act, exp)
	}
}

func checkError(err error, t *testing.T) {
	t.Helper()
	if err != nil {
		t.Error(err)
	}
}

func receiveSleepMs(ms time.Duration, done chan struct{}, t *testing.T) {
	t.Helper()
	select {
	case <-time.After(ms * time.Millisecond):
		t.Fatal("Test timed out")
	case <-done:
	}
}

func sleepMs(t time.Duration) {
	time.Sleep(time.Millisecond * t)
}
