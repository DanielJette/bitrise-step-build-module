package main

import (
    "os"
    "github.com/bitrise-io/go-utils/log"
    // "github.com/bitrise-steplib/bitrise-step-build-router-start/env"
    "github.com/bitrise-steplib/bitrise-step-build-router-start/execmd"
    "github.com/bitrise-steplib/bitrise-step-build-router-start/gradle"
    // "github.com/bitrise-steplib/bitrise-step-build-router-start/deploy"
    // "github.com/bitrise-steplib/bitrise-step-build-router-start/trigger"
)

func DisplayInfo() {
    log.Infof("=== Display environment info ===")
    execmd.ExecuteCommand("go", "version")
    execmd.ExecuteCommand("git", "--version")
    execmd.ExecuteCommand("adb", "--version")
    execmd.ExecuteRelativeCommand("./gradlew", "--version")
}

func main() {
    DisplayInfo()

    gradle.PrintModifiedModules()
    // gradle.Assemble()
    // gradle.PrepareForDeploy()
    // env.SetTargetEnv()
    // deploy.Deploy()
    // trigger.TriggerWorkflow()

    os.Exit(0)
}

