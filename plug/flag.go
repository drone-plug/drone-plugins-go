package plug

import (
	"encoding/json"
	"flag"
	"fmt"
	"strings"

	"github.com/go-pa/fenv"
)

// FlagSet adds drone plugin specific functionality to a wrapped flag.FlagSet
type FlagSet struct {
	*flag.FlagSet
	es             *fenv.EnvSet
	envFiles       []string
	envFilesActive bool
}

// FlagEnv replaces automatically generated environment variable names
// (PLUGIN_FLAG_NAME) for a specific flag and supports multiple names for a
// flag. Environment variable prefixes are ignored for manually specified
// names. If an emtpy string is given as an envName argument it will be
// replaced with the default generated environment variable name.
//
func (fs *FlagSet) Env(flagVar interface{}, envName ...string) {
	fs.es.Var(flagVar, envName...)
}

func (fs *FlagSet) StringSliceVar(value *[]string, name, usage string) {
	fs.Var((*stringSliceFlag)(value), name, usage)
}

func (fs *FlagSet) StringMapVar(value *map[string]string, name, usage string) {
	fs.Var((*stringMapFlag)(value), name, usage)
}

const envfileFlagName = "env_file"

func (fs *FlagSet) EnvFiles(names ...string) {
	if len(names) > 0 {
		fs.envFiles = names
	}
	if !fs.envFilesActive {
		fs.StringSliceVar(&fs.envFiles, envfileFlagName, "source env file")
	}
	fs.envFilesActive = true
}

func (fs *FlagSet) RepoVar(r *Repo) {
	fs.RepoOwnerVar(&r.Owner)
	fs.RepoNameVar(&r.Name)
	fs.RepoLinkVar(&r.Link)
	fs.RepoAvatarVar(&r.Avatar)
	fs.RepoBranchVar(&r.Branch)
	fs.RepoPrivateVar(&r.Private)
	fs.RepoTrustedVar(&r.Trusted)
}

func (fs *FlagSet) CommitVar(c *Commit) {
	fs.CommitMessageVar(&c.Message)
	fs.CommitShaVar(&c.Sha)
	fs.CommitRefVar(&c.Ref)
	fs.CommitLinkVar(&c.Link)
	fs.CommitBranchVar(&c.Branch)
	fs.CommitAuthorEmailVar(&c.Author.Email)
	fs.CommitAuthorNameVar(&c.Author.Name)
	fs.CommitAuthorAvatarVar(&c.Author.Avatar)
}

func (fs *FlagSet) BuildVar(b *Build) {
	fs.BuildNumberVar(&b.Number)
	fs.BuildEventVar(&b.Event)
	fs.BuildStatusVar(&b.Status)
	fs.DroneDeployToVar(&b.Deploy)
	fs.BuildCreatedVar(&b.Created)
	fs.BuildStartedVar(&b.Started)
	fs.BuildFinishedVar(&b.Finished)
	fs.BuildLinkVar(&b.Link)
}

// RepoFullNameVar defines a string flag for DRONE_REPO.
func (fs *FlagSet) RepoFullNameVar(v *string) {
	fs.droneFlag("repo", v, "repository full name")
}

// RepoOwnerVar defines a string flag for DRONE_REPO_OWNER
func (fs *FlagSet) RepoOwnerVar(v *string) {
	fs.droneFlag("repo.owner", v, "repository owner")
}

// RepoNameVar defines a string flag for DRONE_REPO_NAME
func (fs *FlagSet) RepoNameVar(v *string) {
	fs.droneFlag("repo.name", v, "repository name")
}

// RepoLinkVar defines a string flag for DRONE_REPO_LINK
func (fs *FlagSet) RepoLinkVar(v *string) {
	fs.droneFlag("repo.link", v, "repository link")
}

// RepoAvatarVar defines a string flag for DRONE_REPO_AVATAR
func (fs *FlagSet) RepoAvatarVar(v *string) {
	fs.droneFlag("repo.avatar", v, "repository avatar")
}

// RepoSCMVar defines a string flag for DRONE_REPO_SCM
func (fs *FlagSet) RepoSCMVar(v *string) {
	fs.droneFlag("repo.scm", v, "repository scm (git)")
}

// RepoBranchVar defines a string flag for DRONE_REPO_BRANCH
func (fs *FlagSet) RepoBranchVar(v *string) {
	fs.droneFlag("repo.branch", v, "repository default branch (master)")
}

// RepoPrivateVar defines a bool flag for DRONE_REPO_PRIVATE
func (fs *FlagSet) RepoPrivateVar(v *bool) {
	fs.droneFlag("repo.private", v, "repository is private")
}

// RepoTrustedVar defines a bool flag for DRONE_REPO_TRUSTED
func (fs *FlagSet) RepoTrustedVar(v *bool) {
	fs.droneFlag("repo.trusted", v, "repository is trusted")
}

// BuildNumberVar defines a int flag for DRONE_BUILD_NUMBER
func (fs *FlagSet) BuildNumberVar(v *int64) {
	fs.droneFlag("build.number", v, "build number")
}

// BuildEventVar defines a string flag for DRONE_BUILD_EVENT
func (fs *FlagSet) BuildEventVar(v *string) {
	fs.droneFlag("build.event", v, "build event (push, pull_request, tag)")
}

// BuildStatusVar defines a string flag for DRONE_BUILD_STATUS
func (fs *FlagSet) BuildStatusVar(v *string) {
	fs.droneFlag("build.status", v, "build status (success, failure)")
}

// BuildCreatedVar defines a int flag for DRONE_BUILD_CREATED
func (fs *FlagSet) BuildCreatedVar(v *int64) {
	fs.droneFlag("build.created", v, "build created unix timestamp")
}

// BuildStartedVar defines a int flag for DRONE_BUILD_STARTED
func (fs *FlagSet) BuildStartedVar(v *int64) {
	fs.droneFlag("build.started", v, "build started unix timestamp")
}

// BuildFinishedVar defines a int flag for DRONE_BUILD_FINISHED
func (fs *FlagSet) BuildFinishedVar(v *int64) {
	fs.droneFlag("build.finished", v, "build finished unix timestamp")
}

// BuildLinkVar defines a string flag for DRONE_BUILD_LINK
func (fs *FlagSet) BuildLinkVar(v *string) {
	fs.droneFlag("build.link", v, "build result link")
}

// CommitShaVar defines a string flag for DRONE_COMMIT_SHA
func (fs *FlagSet) CommitShaVar(v *string) {
	fs.droneFlag("commit.sha", v, "commit sha")
}

// CommitRefVar defines a string flag for DRONE_COMMIT_REF
func (fs *FlagSet) CommitRefVar(v *string) {
	fs.droneFlag("commit.ref", v, "commit ref")
}

// CommitLinkVar defines a string flag for DRONE_COMMIT_LINK
func (fs *FlagSet) CommitLinkVar(v *string) {
	fs.droneFlag("commit.link", v, "commit link in remote")
}

// CommitBranchVar defines a string flag for DRONE_COMMIT_BRANCH
func (fs *FlagSet) CommitBranchVar(v *string) {
	fs.droneFlag("commit.branch", v, "commit branch")
}

// CommitMessageVar defines a string flag for DRONE_COMMIT_MESSAGE
func (fs *FlagSet) CommitMessageVar(v *string) {
	fs.droneFlag("commit.message", v, "commit message")
}

// CommitAuthorNameVar defines a string flag for DRONE_COMMIT_AUTHOR_NAME
func (fs *FlagSet) CommitAuthorNameVar(v *string) {
	fs.droneFlag("commit.author.name", v, "commit author username")
}

// CommitAuthorEmailVar defines a string flag for DRONE_COMMIT_AUTHOR_EMAIL
func (fs *FlagSet) CommitAuthorEmailVar(v *string) {
	fs.droneFlag("commit.author.email", v, "commit author email address")
}

// CommitAuthorAvatarVar defines a string flag for DRONE_COMMIT_AUTHOR_AVATAR
func (fs *FlagSet) CommitAuthorAvatarVar(v *string) {
	fs.droneFlag("commit.author.avatar", v, "commit author avatar")
}

// PrevBuildStatusVar defines a string flag for DRONE_PREV_BUILD_STATUS
func (fs *FlagSet) PrevBuildStatusVar(v *string) {
	fs.droneFlag("prev.build.status", v, "prior build status")
}

// PrevBuildNumberVar defines a string flag for DRONE_PREV_BUILD_NUMBER
func (fs *FlagSet) PrevBuildNumberVar(v *string) {
	fs.droneFlag("prev.build.number", v, "prior build number")
}

// DroneDeployToVar defines a string flag for DRONE_DEPLOY_TO
func (fs *FlagSet) DroneDeployToVar(v *string) {
	fs.droneFlag("deploy.to", v, "build deployment target")
}

// DroneRemoteURLVar defines a string flag for DRONE_REMOTE_URL
func (fs *FlagSet) DroneRemoteURLVar(v *string) {
	fs.droneFlag("remote.url", v, "repository clone url")
}

// DronePullRequestVar defines a string flag for DRONE_PULL_REQUEST
func (fs *FlagSet) DronePullRequestVar(v *string) {
	fs.droneFlag("pull.request", v, "pull request number")
}

var (
	flagNamePrefix = "" // for tests
)

func (fs *FlagSet) droneFlag(name string, ref interface{}, help string) {
	name = flagNamePrefix + name
	s := "drone_" + name
	s = strings.Replace(s, ".", "_", -1)
	s = strings.Replace(s, "-", "_", -1)
	s = strings.ToUpper(s)
	usage := fmt.Sprintf("%s (%s)", help, s)
	switch v := ref.(type) {
	case *string:
		fs.StringVar(v, name, "", usage)
	case *bool:
		fs.BoolVar(v, name, false, usage)
	case *int64:
		fs.Int64Var(v, name, -1, usage)
	default:
		panic(v)
	}
	fs.Env(ref, s)
}

// stringSliceFlag is a flag type which
type stringSliceFlag []string

func (s *stringSliceFlag) String() string {
	return strings.Join(*s, ",")
}

func (s *stringSliceFlag) Set(value string) error {
	*s = strings.Split(value, ",")
	return nil
}

type stringMapFlag map[string]string

func (s *stringMapFlag) String() string {
	data, err := json.Marshal(s)
	if err != nil {
		panic(err)
	}
	return string(data)
}

func (s *stringMapFlag) Set(value string) error {
	var m map[string]string
	err := json.Unmarshal([]byte(value), &m)
	if err != nil {
		return err
	}
	*s = m
	return nil
}
