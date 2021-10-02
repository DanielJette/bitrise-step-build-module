package gradle

import (
    "fmt"
    "github.com/bitrise-io/go-steputils/stepconf"
    "github.com/bitrise-io/go-utils/log"
    "github.com/bitrise-steplib/bitrise-step-build-router-start/execmd"
    "github.com/bitrise-steplib/bitrise-step-build-router-start/util"
)

type BuildConfig struct {
    Module            string          `env:"module,required"`
    Variant           string          `env:"variant,required"`
}

var gradlew = "./gradlew"

func Assemble() {

    var cfg BuildConfig
    if err := stepconf.Parse(&cfg); err != nil {
        util.Failf("Issue with an input: %s", err)
    }

    log.Infof("Building %s %s", cfg.Module, cfg.Variant)

    cmd := fmt.Sprintf("%s:assemble%s", cfg.Module, cfg.Variant)

    execmd.ExecuteRelativeCommand(gradlew, cmd)
}

func Deploy() {

    // features/login/build/outputs/apk/androidTest/internal/debug/feature-login-internal-debug-androidTest.apk

}
