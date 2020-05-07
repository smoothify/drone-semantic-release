package main

import (
  "errors"
  "fmt"
  "github.com/smoothify/drone-semantic-release/pkg/semantic_release"
  "github.com/joho/godotenv"
  "github.com/sirupsen/logrus"
  "github.com/urfave/cli/v2"
  "net/url"
  "os"
  "path"
)

var (
	version = "unknown"
)

func main() {
	// Load env-file if it exists first
	if env := os.Getenv("PLUGIN_ENV_FILE"); env != "" {
		godotenv.Load(env)
	}

	app := cli.NewApp()
	app.Name = "semantic-release plugin"
	app.Usage = "semantic-release plugin"
	app.Action = run
	app.Version = version
	app.Flags = []cli.Flag{
    &cli.BoolFlag{
			Name:   "dry-run",
			Usage:  "dry run disables git push/release",
			EnvVars: []string{"PLUGIN_DRY_RUN"},
		},
    &cli.StringFlag{
			Name:   "git.user.name",
			Usage:  "git user name",
			Value:  "drone-ci",
			EnvVars: []string{"PLUGIN_GIT_USER_NAME"},
		},
    &cli.StringFlag{
      Name:   "git.user.email",
      Usage:  "git user email",
      EnvVars: []string{"PLUGIN_GIT_USER_EMAIL"},
    },
    &cli.StringFlag{
      Name:   "git.login",
      Usage:  "git login",
      EnvVars: []string{"PLUGIN_GIT_LOGIN"},
    },
    &cli.StringFlag{
      Name:   "git.password",
      Usage:  "git password",
      EnvVars: []string{"PLUGIN_GIT_PASSWORD"},
    },
    &cli.StringFlag{
      Name:   "git.credentials",
      Usage:  "git credentials",
      EnvVars: []string{"PLUGIN_GIT_CREDENTIALS"},
    },
    &cli.StringFlag{
      Name:   "github.token",
      Usage:  "github token",
      EnvVars: []string{"PLUGIN_GITHUB_TOKEN","GITHUB_TOKEN"},
    },
    &cli.StringFlag{
      Name:   "bitbucket.token",
      Usage:  "bitbucket token",
      EnvVars: []string{"PLUGIN_BITBUCKET_TOKEN","BITBUCKET_TOKEN"},
    },
    &cli.StringFlag{
      Name:   "repo.branch",
      Usage:  "repo branch",
      EnvVars: []string{"PLUGIN_BRANCH","DRONE_REPO_BRANCH", "DRONE_BRANCH"},
    },
    &cli.BoolFlag{
			Name:   "changelog",
			Usage:  "generate change log",
      Value: true,
			EnvVars: []string{"PLUGIN_CHANGELOG"},
		},
    &cli.BoolFlag{
      Name:   "package-json",
      Usage:  "generate package.json",
      Value: true,
      EnvVars: []string{"PLUGIN_PACKAGE_JSON"},
    },
    &cli.StringFlag{
      Name:   "version.filename",
      Usage:  "version file filename",
      Value:  ".release-version",
      EnvVars: []string{"PLUGIN_VERSION_FILENAME"},
    },
    &cli.StringFlag{
			Name:   "commit.template",
			Usage:  "commit msg template",
			Value:  "chore(release): ${nextRelease.version} [skip ci]",
			EnvVars: []string{"PLUGIN_DOCKERFILE"},
		},
    &cli.StringSliceFlag{
      Name:     "assets",
      Usage:    "git assets",
      Value:    &cli.StringSlice{},
      EnvVars:  []string{"PLUGIN_ASSETS"},
      FilePath: ".assets",
    },
    &cli.BoolFlag{
      Name:   "tags.build",
      Usage:  "build tag list",
      Value:  true,
      EnvVars: []string{"PLUGIN_TAGS_BUILD"},
    },
    &cli.BoolFlag{
      Name:   "tags.add-branch",
      Usage:  "add branch to tag list",
      Value:  true,
      EnvVars: []string{"PLUGIN_TAGS_ADD_BRANCH"},
    },
    &cli.BoolFlag{
      Name:    "tags.require-version",
      Usage:   "require version to create tags file",
      EnvVars: []string{"PLUGIN_TAGS_REQUIRE_VERSION"},
    },
    &cli.StringFlag{
      Name:   "tags.filename",
      Usage:  "output tags filename",
      Value:  ".tags",
      EnvVars: []string{"PLUGIN_TAGS_FILENAME"},
    },
    &cli.StringFlag{
      Name:   "tags.master",
      Usage:  "tag to use for master branch",
      Value:  "latest",
      EnvVars: []string{"PLUGIN_TAGS_MASTER"},
    },
    &cli.BoolFlag{
      Name:   "error-no-release",
      Usage:  "error on no new release",
      EnvVars: []string{"PLUGIN_ERROR_NO_RELEASE"},
    },
    &cli.StringFlag{
      Name:   "config-path",
      Usage:  "default config path",
      Value:  "/semantic-release/default.config.js",
      EnvVars: []string{"PLUGIN_CONFIG_PATH"},
    },
	}

	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}

func getExecDir() string {
  ex, err := os.Executable()
  if err != nil {
    logrus.Fatal(err)
  }
  dir := path.Dir(ex)
  return dir
}

func run(c *cli.Context) error {
  if c.String("git.login") != "" && c.String("git.password") != "" {
    _ = c.Set("git.credentials", fmt.Sprintf("%s:%s", url.PathEscape(c.String("git.login")), url.PathEscape(c.String("git.password"))))
  }
  if c.String("github.token") == "" && c.String("bitbucket.token") == "" && c.String("git.credentials") == "" {
    return errors.New("one of github.token, bitbucket.token or git credentials must be supplied")
  }

  plugin := semantic_release.Plugin{
    Repo: semantic_release.Repo{
      GitUserName:    c.String("git.user.name"),
      GitUserEmail:   c.String("git.user.email"),
      GitHubToken:    c.String("github.token"),
      BitBucketToken: c.String("bitbucket.token"),
      GitCredentials: c.String("git.credentials"),
      Branch:         c.String("repo.branch"),
    },
    Release: semantic_release.Release{
      DryRun:              c.Bool("dry-run"),
      GenerateChangelog:   c.Bool("changelog"),
      GeneratePackageJson: c.Bool("package-json"),
      VersionFilename:     c.String("version.filename"),
      CommitMsgTemplate:   c.String("commit.template"),
      Assets:              c.StringSlice("assets"),
      TagsBuild:           c.Bool("tags.build"),
      TagsAddBranch:       c.Bool("tags.add-branch"),
      TagsRequireVersion:  c.Bool("tags.require-version"),
      TagsFilename:        c.String("tags.filename"),
      TagsMaster:          c.String("tags.master"),
    },
    ConfigPath:     c.String("config-path"),
    ErrorNoRelease: c.Bool("error-no-release"),
  }

  return plugin.Exec()
}
