package main

import (
    "fmt"
    "os"
    "strings"
    "time"
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
    var root = fmt.Sprintf("./%s", testPath)
    if _, err := os.Stat(root); !os.IsNotExist(err) {
        log.Infof("OK. Found %s!\n", testPath)
        return true
    }
    return false
}

func isSkippable(module string) bool {

    testModuleDir := strings.TrimPrefix(module, "feature-")
    testPath := fmt.Sprintf("features/%s/src/androidTest", testModuleDir)
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
    timestamp()
    gradle.Assemble()
    timestamp()
    gradle.PrepareForDeploy()
    timestamp()
    env.SetTargetEnv()
    timestamp()
    deploy.Deploy()
    timestamp()
    trigger.TriggerWorkflow()
    timestamp()
}

var startTime int64 = 0

func timestamp() {
    log.Infof("[Time] %d", time.Now().UnixNano() / int64(time.Millisecond) - startTime)
}

func main() {
    startTime = time.Now().UnixNano() / int64(time.Millisecond)
    timestamp()
    var cfg PathConfig
    if err := stepconf.Parse(&cfg); err != nil {
        util.Failf("Issue with an input: %s", err)
    }
    timestamp()
    // DisplayInfo()

    if isSkippable(cfg.Module) {
        os.Exit(0)
    }

    log.Infof("Building %s", cfg.Module)
    timestamp()
    buildAndTrigger()

    os.Exit(0)
}

