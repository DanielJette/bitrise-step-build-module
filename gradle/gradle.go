package gradle

import (
    "fmt"
    "github.com/bitrise-io/go-steputils/stepconf"
    "github.com/bitrise-io/go-utils/log"
    "github.com/bitrise-steplib/bitrise-step-build-router-start/execcmd"
    "github.com/bitrise-steplib/bitrise-step-build-router-start/util"
)

type BuildConfig struct {
    Module            string          `env:"module,required"`
    Variant           string          `env:"variant,required"`
}

var gradlew = "./gradlew"

func BuildAPK() {

    var cfg BuildConfig
    if err := stepconf.Parse(&cfg); err != nil {
        util.Failf("Issue with an input: %s", err)
    }

    log.Infof("Building %s %s", cfg.Module, cfg.Variant)

    cmd := fmt.Sprintf("%s:%s", cfg.Module, cfg.Variant)

    execcmd.ExecuteRelativeCommand(gradlew, cmd)
}