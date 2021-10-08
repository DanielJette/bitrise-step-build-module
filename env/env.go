package env

import (
    "fmt"
    "os"
    "github.com/bitrise-io/go-steputils/stepconf"
    "github.com/bitrise-io/go-utils/log"
    "github.com/bitrise-steplib/bitrise-step-build-router-start/execmd"
    "github.com/bitrise-steplib/bitrise-step-build-router-start/util"
)

type TargetConfig struct {
    APK             string          `env:"target_apk,required"`
    TestPackage     string          `env:"test_package,required"`
    TestRunner      string          `env:"test_runner,required"`
    JUnit5          bool            `env:"is_junit_5,required"`
}

func SetTargetEnv() {
    log.Infof("=== Set target environment ===")

    var cfg TargetConfig
    if err := stepconf.Parse(&cfg); err != nil {
        util.Failf("Issue with an input: %s", err)
    }

    var runnerBuilder string
    if cfg.JUnit5 {
        runnerBuilder = "-e runnerBuilder de.mannodermaus.junit5.AndroidJUnit5Builder"
    } else {
        runnerBuilder = ""
    }

    adbFormatString := "\"adb shell am instrument -r -w %s %s/%s\""
    adbCommand := fmt.Sprintf(adbFormatString, runnerBuilder, cfg.TestPackage, cfg.TestRunner)
    log.Infof("Set adb command to [%s]", adbCommand)
    execmd.ExecuteCommand("envman", "add", "--key", "ADB_COMMAND", "--value", adbCommand)
    os.Setenv("ADB_COMMAND", adbCommand)

    log.Infof("Set target apk to [%s]", cfg.APK)
    execmd.ExecuteCommand("envman", "add", "--key", "TARGET_APK", "--value", cfg.APK)
    os.Setenv("TARGET_APK", cfg.APK)
}
