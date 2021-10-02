package env

import (
    "fmt"
    "github.com/bitrise-io/go-steputils/stepconf"
    "github.com/bitrise-io/go-utils/log"
    "github.com/bitrise-steplib/bitrise-step-build-router-start/execmd"
    "github.com/bitrise-steplib/bitrise-step-build-router-start/util"
)

type TargetConfig struct {
    APK             string          `env:"target_apk,required"`
    TestPackage     string          `env:"test_package,required"`
    TestRunner      string          `env:"test_runner,required"`
}

func SetTargetEnv() {
    log.Infof("=== Set target environment ===")

    var cfg TargetConfig
    if err := stepconf.Parse(&cfg); err != nil {
        util.Failf("Issue with an input: %s", err)
    }

    adbCommand := fmt.Sprintf("adb shell am instrument -w -m -e debug false %s/%s", cfg.TestPackage, cfg.TestRunner)
    log.Infof("Set adb command to [%s]", adbCommand)
    execmd.ExecuteCommand("envman", "add", "--key", "ADB_COMMAND", "--value", adbCommand)
    log.Infof("Set target apk to [%s]", cfg.APK)
    execmd.ExecuteCommand("envman", "add", "--key", "TARGET_APK", "--value", cfg.APK)
}
