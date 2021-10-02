package main

import (
    "os"
    "github.com/bitrise-io/go-utils/log"
    "github.com/bitrise-steplib/bitrise-step-build-router-start/env"
    "github.com/bitrise-steplib/bitrise-step-build-router-start/execcmd"
    "github.com/bitrise-steplib/bitrise-step-build-router-start/gradle"
    "github.com/bitrise-steplib/bitrise-step-build-router-start/deploy"
    "github.com/bitrise-steplib/bitrise-step-build-router-start/trigger"
)

func DisplayInfo() {
    log.Infof("=== Display environment info ===")
    execcmd.ExecuteCommand("go", "version")
    execcmd.ExecuteCommand("git", "--version")
    execcmd.ExecuteCommand("adb", "--version")
    execcmd.ExecuteRelativeCommand("./gradlew", "--version")
}

func main() {
    DisplayInfo()

    gradle.Assemble()
    env.SetTargetEnv()
    trigger.TriggerWorkflow()

    deploy.Deploy()

    os.Exit(0)
}

