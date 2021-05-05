package gigi

type Config struct {
	Owner, Repository    string
	PullRequestNumber    int
	GitHubToken          string
	MaxAllowAdditionSize int
}
