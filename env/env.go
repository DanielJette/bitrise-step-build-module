package env

import (
    "fmt"
    "github.com/bitrise-io/go-steputils/stepconf"
    "github.com/bitrise-steplib/bitrise-step-build-router-start/execcmd"
    "github.com/bitrise-steplib/bitrise-step-build-router-start/util"
)

type TargetConfig struct {
    APK             string          `env:"target_apk,required"`
    TestPackage     string          `env:"test_package,required"`
    TestRunner      string          `env:"test_runner,required"`
}

func SetTargetEnv() {

    var cfg TargetConfig
    if err := stepconf.Parse(&cfg); err != nil {
        util.Failf("Issue with an input: %s", err)
    }

    adbCommand := fmt.Sprintf("adb shell am instrument -w -m -e debug false %s/%s", cfg.TestPackage, cfg.TestRunner)
    execcmd.ExecuteCommand("envman", "add", "--key", "ADB_COMMAND", "--value", adbCommand)
    execcmd.ExecuteCommand("enveman", "add", "--key", "TARGET_APK", "--value", cfg.APK)
}
