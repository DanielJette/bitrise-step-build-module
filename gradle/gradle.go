package gradle

import (
    "fmt"
    "github.com/bitrise-io/go-steputils/stepconf"
    "github.com/bitrise-io/go-utils/log"
    "github.com/bitrise-steplib/bitrise-step-build-router-start/execmd"
    "github.com/bitrise-steplib/bitrise-step-build-router-start/util"
)

type BuildConfig struct {
    Module      string     `env:"module,required"`
    Variant     string     `env:"variant,required"`
    DeployDir   string     `env:"deploy_path,required"`
    APK         string     `env:"target_apk,required"`
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

func PrepareForDeploy() {
    var cfg BuildConfig
    if err := stepconf.Parse(&cfg); err != nil {
        util.Failf("Issue with an input: %s", err)
    }
    execmd.ExecuteCommand("find", ".", "-name", cfg.APK, "-exec", "cp", "{}", cfg.DeployDir, ";")
}
