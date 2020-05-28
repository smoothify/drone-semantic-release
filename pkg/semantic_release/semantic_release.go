package semantic_release

import (
  "fmt"
  "github.com/smoothify/drone-semantic-release/pkg/tags"
  "github.com/sirupsen/logrus"
  "io"
  "io/ioutil"
  "os"
  "os/exec"
  "strings"
)

const semanticReleaseExe = "/usr/local/bin/semantic-release"
//const semanticReleaseHome = "/root/.docker/"

type (
	// Repo defines git repo parameters.
	Repo struct {
	  GitUserName    string
	  GitUserEmail   string
		GitHubToken    string
		BitBucketToken string
		GitCredentials string
		Branch         string
	}

	// Release defines release parameters.
	Release struct {
    DryRun              bool
		GenerateChangelog   bool
		GeneratePackageJson bool
		CommitMsgTemplate   string
    CustomVersion       string
    VersionFilename     string
		Assets              []string
		TagsBuild           bool
    TagsAddBranch       bool
    TagsRequireVersion  bool
		TagsFilename        string
    TagsMaster          string
	}

  // Plugin defines the plugin parameters.
  Plugin struct {
    Repo           Repo
    Release        Release
    ConfigPath     string
    ErrorNoRelease bool
  }
)

// Exec executes the plugin step
func (p Plugin) Exec() error {
  var version string

  if p.Release.CustomVersion == "" {
    cmd := commandBuild(p)

    cmd.Stderr = os.Stderr
    cmd.Stdout = os.Stdout
    trace(cmd)

    err := removeVersionFile(p)
    if err != nil {
      return err
    }

    err = cmd.Run()
    if err != nil {
      return err
    }

    version = retrieveVersion(p)

  } else {
    version = p.Release.CustomVersion
  }

  if version == "" && p.ErrorNoRelease {
    logrus.Fatal("No release required, no version created")
  }

  if p.Release.TagsBuild {
    versionTags := tags.VersionTags(version)

    if len(versionTags) > 0 || ! p.Release.TagsRequireVersion {

      if p.Release.TagsAddBranch {
        branch := p.Repo.Branch
        if branch == "master" {
          branch = p.Release.TagsMaster
        }
        if branch != "" {
          versionTags = append(versionTags, branch)
        }
      }

      err := writeTags(p, versionTags)
      if err != nil {
        return err
      }
      logrus.Infof("Tags file %s written to disk: %s", p.Release.TagsFilename, strings.Join(versionTags, ","))
    } else {
      logrus.Info("No tags found, skipping creation of tags file")
    }
  }

  return nil
}

// helper function to create the semantic-release command.
func commandBuild(plugin Plugin) *exec.Cmd {
  args := []string{
    "--extends", plugin.ConfigPath,
  }

	if plugin.Release.DryRun {
		args = append(args, "--dry-run")
	}

  cmd := exec.Command(semanticReleaseExe, args...)

  env := append(os.Environ(),
    "CI=true",
    fmt.Sprintf("DRY_RUN=%t", plugin.Release.DryRun),
    fmt.Sprintf("GIT_AUTHOR_NAME=%s", plugin.Repo.GitUserName),
    fmt.Sprintf("GIT_COMMITTER_NAME=%s", plugin.Repo.GitUserName),
    fmt.Sprintf("GIT_AUTHOR_EMAIL=%s", plugin.Repo.GitUserEmail),
    fmt.Sprintf("GIT_COMMITTER_EMAIL=%s", plugin.Repo.GitUserEmail),
    fmt.Sprintf("GENERATE_CHANGELOG=%t", plugin.Release.GenerateChangelog),
    fmt.Sprintf("GENERATE_PACKAGEJSON=%t", plugin.Release.GeneratePackageJson ),
    fmt.Sprintf("GIT_ASSETS=%s", getAssetList(plugin)),
    fmt.Sprintf("GIT_COMMIT_TEMPLATE=%s", plugin.Release.CommitMsgTemplate ),
    fmt.Sprintf("VERSION_FILENAME=%s", plugin.Release.VersionFilename),
    fmt.Sprintf("DRONE_BRANCH=%s", plugin.Repo.Branch),
    fmt.Sprintf("DRONE_REPO_BRANCH=%s", plugin.Repo.Branch),
  )

  if plugin.Repo.GitHubToken != "" {
    env = append(env, fmt.Sprintf("GH_TOKEN=%s", plugin.Repo.GitHubToken))
  }

  if plugin.Repo.BitBucketToken != "" {
    env = append(env, fmt.Sprintf("BB_TOKEN=%s", plugin.Repo.BitBucketToken))
  }

  if plugin.Repo.GitCredentials != "" {
    env = append(env, fmt.Sprintf("GIT_CREDENTIALS=%s", plugin.Repo.GitCredentials))
  }

  cmd.Env = env

	return cmd
}

func getAssetList(plugin Plugin) string {
  if plugin.Release.GenerateChangelog && ! contains(plugin.Release.Assets, "CHANGELOG.md") {
    plugin.Release.Assets = append(plugin.Release.Assets, "CHANGELOG.md")
  }
  if plugin.Release.GeneratePackageJson && ! contains(plugin.Release.Assets, "package.json") {
    plugin.Release.Assets = append(plugin.Release.Assets, "package.json")
  }

  return strings.Join(plugin.Release.Assets, ",")
}

func removeVersionFile(plugin Plugin) error  {
  filename := plugin.Release.VersionFilename
  _, err := os.Stat(filename)
  if ! os.IsNotExist(err) {
    e := os.Remove(filename)
    if e != nil {
      return e
    }
  }
  return nil
}

func contains(s []string, e string) bool {
  for _, a := range s {
    if a == e {
      return true
    }
  }
  return false
}

func retrieveVersion(plugin Plugin) string {
  file, err := os.Open(plugin.Release.VersionFilename)

  if err != nil {
    return ""
  }
  defer file.Close()

  contents, err := ioutil.ReadAll(file)
  return strings.TrimSpace(string(contents))
}

func writeTags(plugin Plugin, tags []string) error {
  file, err := os.Create(plugin.Release.TagsFilename)
  if err != nil {
    return err
  }
  defer file.Close()

  _, err = io.WriteString(file, strings.Join(tags, ","))
  if err != nil {
    return err
  }
  return file.Sync()
}

// trace writes each command to stdout with the command wrapped in an xml
// tag so that it can be extracted and displayed in the logs.
func trace(cmd *exec.Cmd) {
	fmt.Fprintf(os.Stdout, "+ %s\n", strings.Join(cmd.Args, " "))
}
