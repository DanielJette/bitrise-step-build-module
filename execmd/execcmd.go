package execmd

import (
    "fmt"
    // "path/filepath"
    "os"
    "os/exec"
    "github.com/bitrise-io/go-utils/log"
    "github.com/bitrise-steplib/bitrise-step-build-router-start/util"
)

func ExecuteRelativeCommand(executablePath string, a ...string) {
    args := append([]string {executablePath}, a...)
    log.Infof("Executing %s", args)
    cmd := &exec.Cmd {
        Path: executablePath,
        Args: args,
        Stdout: os.Stdout,
        Stderr: os.Stdout,
    }

    if err := cmd.Run(); err != nil {
        util.Failf("Error %s", err)
    }
    log.Infof("OK")
}

func ExecuteCommand(executable string, a ...string) {
    executablePath, _ := exec.LookPath( executable )
    ExecuteRelativeCommand(executablePath, a...)
}

func ExecuteShellScript(script string) string {

    // dir, err := os.Getwd()
    // if err != nil {
    //     fmt.Println(err)
    // }

    dir := os.Getenv("BITRISE_SOURCE_DIR")

    log.Infof(dir)

    // dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
    // if err != nil {
    //     util.Failf("Error %s", err)
    // }

    // dir := os.Args[0]
    // fmt.Println(dir)

    scriptPath := fmt.Sprintf("%s/%s", dir, script)
    log.Infof("Executing script %s", scriptPath)
    cmd, err := exec.Command("/bin/sh", scriptPath).Output()
    if err != nil {
        util.Failf("Error %s", err)
    }

    output := string(cmd)
    log.Infof("OK")
    return output
}
