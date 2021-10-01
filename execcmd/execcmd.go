package execcmd

import (
    "os"
    "os/exec"
    "github.com/bitrise-steplib/bitrise-step-build-router-start/util"
)

func ExecuteRelativeCommand(executablePath string, a ...string) {
    args := append([]string {executablePath}, a...)
    cmd := &exec.Cmd {
        Path: executablePath,
        Args: args,
        Stdout: os.Stdout,
        Stderr: os.Stdout,
    }

    if err := cmd.Run(); err != nil {
        util.Failf("Error", err)
    }
}

func ExecuteCommand(executable string, a ...string) {
    executablePath, _ := exec.LookPath( executable )
    ExecuteRelativeCommand(executablePath, a...)
}
