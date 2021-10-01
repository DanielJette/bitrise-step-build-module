package execcmd

import (
    "os"
    "os/exec"

    "github.com/bitrise-io/go-utils/log"
)

func failf(s string, a ...interface{}) {
    log.Errorf(s, a...)
    os.Exit(1)
}

func ExecuteRelativeCommand(executablePath string, a ...string) {
    args := append([]string {executablePath}, a...)
    cmd := &exec.Cmd {
        Path: executablePath,
        Args: args,
        Stdout: os.Stdout,
        Stderr: os.Stdout,
    }

    if err := cmd.Run(); err != nil {
        failf("Error", err)
    }
}

func ExecuteCommand(executable string, a ...string) {
    executablePath, _ := exec.LookPath( executable )
    ExecuteRelativeCommand(executablePath, a...)
}
