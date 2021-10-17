package main

import (
    "fmt"
    "os"
    "github.com/bitrise-io/go-utils/log"
    "github.com/bitrise-steplib/bitrise-step-build-router-start/env"
    "github.com/bitrise-steplib/bitrise-step-build-router-start/execmd"
    "github.com/bitrise-steplib/bitrise-step-build-router-start/gradle"
    "github.com/bitrise-steplib/bitrise-step-build-router-start/deploy"
    "github.com/bitrise-steplib/bitrise-step-build-router-start/trigger"
    "github.com/bitrise-steplib/bitrise-step-build-router-start/gh"
    "github.com/bitrise-io/go-steputils/stepconf"
    "github.com/bitrise-steplib/bitrise-step-build-router-start/util"
)

func DisplayInfo() {
    log.Infof("=== Display environment info ===")
    execmd.ExecuteCommand("go", "version")
    execmd.ExecuteCommand("git", "--version")
    execmd.ExecuteCommand("adb", "--version")
    execmd.ExecuteRelativeCommand("./gradlew", "--version")
}

type PathConfig struct {
    Module      string     `env:"module,required"`
}

func isSkippable(module string) bool {

    path := fmt.Sprintf("./features/%s/src/androidTest/", module)
    log.Infof("Checking for %s", path)
    if _, err := os.Stat(path); os.IsNotExist(err) {
        log.Errorf("No tests detected in %s. Skipping build", module)
        return true
    }

    modules := gh.GetChangedModules()
    if modules[module] == false {
        log.Errorf("No changes detected in %s. Skipping build", module)
        return true
    }
    log.Infof("Changes to module %s found. Running tests.", module)
    return false
}

func main() {
    var cfg PathConfig
    if err := stepconf.Parse(&cfg); err != nil {
        util.Failf("Issue with an input: %s", err)
    }

    DisplayInfo()

    if isSkippable(cfg.Module) {
        os.Exit(0)
    }

    log.Infof("Building %s", cfg.Module)

    gradle.Assemble()
    gradle.PrepareForDeploy()
    env.SetTargetEnv()
    deploy.Deploy()
    trigger.TriggerWorkflow()

    os.Exit(0)
}

