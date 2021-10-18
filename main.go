package main

import (
    "fmt"
    "os"
    "path/filepath"
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

func checkIfTestsExist(testPath string) bool {
    log.Infof("Checking for tests in %s", testPath)
    var exists bool = false
    filepath.Walk("./features",
        func(path string, info os.FileInfo, err error) error {
            if err != nil {
                return err
            }
            fmt.Println("checking ", path)
            if path == testPath {
                exists = true
                fmt.Printf("Found %s!\n", path)
                return nil
            }
            return nil
        })
    return exists
}

func isSkippable(module string) bool {

    testPath := fmt.Sprintf("features/%s/src/androidTest", module)
    exists := checkIfTestsExist(testPath)
    if !exists {
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

func buildAndTrigger() {
    gradle.Assemble()
    gradle.PrepareForDeploy()
    env.SetTargetEnv()
    deploy.Deploy()
    trigger.TriggerWorkflow()
}

func main() {
    var cfg PathConfig
    if err := stepconf.Parse(&cfg); err != nil {
        util.Failf("Issue with an input: %s", err)
    }

    // DisplayInfo()

    if isSkippable(cfg.Module) {
        os.Exit(0)
    }

    log.Infof("Building %s", cfg.Module)

    // buildAndTrigger()

    os.Exit(0)
}

