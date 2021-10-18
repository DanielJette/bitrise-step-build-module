package gh

import (
    "context"
    "fmt"
    "github.com/bitrise-io/go-steputils/stepconf"
    "github.com/bitrise-steplib/bitrise-step-build-router-start/util"
    "github.com/google/go-github/github"
    "golang.org/x/oauth2"
    "math"
    "os"
    "strconv"
    "strings"
)

const owner = "neofinancial"
const repo = "neo-android"

type GitHubConfig struct {
    Token    string    `env:"github_access_token,required"`
}

func getNumPages(numChangedFiles int) int {
    return int(math.Ceil(float64(numChangedFiles) / 30.0))
}

func getModuleName(filename string) string {
    const feature_dir_prefix = "features/"
    var path string
    if strings.HasPrefix(filename, feature_dir_prefix) {
        path = "feature-" + strings.TrimPrefix(filename, feature_dir_prefix)
    } else {
        path = filename
    }
    return strings.Split(path, "/")[0]
}

func GetChangedModules() map[string]bool {
    var cfg GitHubConfig
    if err := stepconf.Parse(&cfg); err != nil {
        util.Failf("Issue with an input: %s", err)
    }

    modulesChanged := map[string]bool{}
    if cfg.Token == "testing" {
        return modulesChanged
    }

    ctx := context.Background()
    ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: cfg.Token})
    tc := oauth2.NewClient(ctx, ts)
    client := github.NewClient(tc)

    prNumber, _ := strconv.Atoi(os.Getenv("PULL_REQUEST_ID"))
    pr, _, _ := client.PullRequests.Get(ctx, owner, repo, prNumber)

    numPages := getNumPages(*pr.ChangedFiles)

    for i := 1; i <= numPages; i++ {
        fmt.Println("Fetching page", i, "...")
        opts := &github.ListOptions{Page: i}
        files, _, _ := client.PullRequests.ListFiles(ctx, owner, repo, prNumber, opts)

        for _, s := range files {
           filename := getModuleName(*s.Filename)
            modulesChanged[filename] = true
        }
    }

    fmt.Println("Changes detected in:")
    for key, _ := range modulesChanged {
        fmt.Println(" - [", key, "]")
    }

    return modulesChanged
}
