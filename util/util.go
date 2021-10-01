package util

import (
    "os"
    "github.com/bitrise-io/go-utils/log"
)

func Failf(s string, a ...interface{}) {
    log.Errorf(s, a...)
    os.Exit(1)
}
