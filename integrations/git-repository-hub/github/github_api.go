package github

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/go-github/v56/github"
	"golang.org/x/oauth2"
)

type githubAPI struct {
	Context context.Context
	config  *oauth2.Config
	api     *github.Client
}

func InitGithubAPI() githubAPI {

	ctx := context.Background()
	config, err := readConfigFromJSON("web_server/cmd/web_app/credentials_github.json")
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	return githubAPI{Context: ctx, config: config}
}

func (g *githubAPI) GetAuthLink(redirectURI string, state string) (string, error) {
	// work with oauth state!!!
	g.config.RedirectURL = redirectURI
	authURL := g.config.AuthCodeURL(state, oauth2.AccessTypeOffline)
	return authURL, nil
}

func (g *githubAPI) GetToken(authCode string) (*oauth2.Token, error) {
	tok, err := g.config.Exchange(context.TODO(), authCode)
	if err != nil {
		return nil, err
	}
	return tok, nil
}

func (g *githubAPI) SetupClient(token *oauth2.Token) error {
	tokenSource := g.config.TokenSource(context.Background(), token)

	newToken, err := tokenSource.Token()
	if err != nil {
		return err
	}

	client := oauth2.NewClient(context.Background(), tokenSource)
	g.api = github.NewClient(client)

	*token = *newToken
	return nil
}

func (g *githubAPI) GetRepoBranchCommitsFromTime(owner, repoName string, fromTime time.Time, branch string) ([]*github.RepositoryCommit, error) {
	opt := &github.CommitsListOptions{SHA: branch, Since: fromTime}
	commits, _, err := g.api.Repositories.ListCommits(g.Context, owner, repoName, opt)
	if err == nil {
		return commits, nil
	}
	return nil, err
}

func (g *githubAPI) GetRepoBranches(owner, repoName string) ([]*github.Branch, error) {
	opt := &github.BranchListOptions{}
	branches, _, err := g.api.Repositories.ListBranches(g.Context, owner, repoName, opt)
	if err == nil {
		return branches, nil
	}
	return nil, err
}

func readConfigFromJSON(filename string) (*oauth2.Config, error) {
	b, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	type cred struct {
		ClientID     string   `json:"client_id"`
		ClientSecret string   `json:"client_secret"`
		RedirectURIs []string `json:"redirect_uris"`
		AuthURI      string   `json:"auth_uri"`
		TokenURI     string   `json:"token_uri"`
	}
	var j struct {
		Web       *cred `json:"web"`
		Installed *cred `json:"installed"`
	}

	if err := json.Unmarshal(b, &j); err != nil {
		return nil, err
	}

	var c *cred
	switch {
	case j.Web != nil:
		c = j.Web
	case j.Installed != nil:
		c = j.Installed
	default:
		return nil, fmt.Errorf("no credentials found")
	}

	if len(c.RedirectURIs) < 1 {
		return nil, errors.New(" missing redirect URL in the client_credentials.json")
	}

	return &oauth2.Config{
		ClientID:     c.ClientID,
		ClientSecret: c.ClientSecret,
		RedirectURL:  c.RedirectURIs[0],
		Endpoint: oauth2.Endpoint{
			AuthURL:  c.AuthURI,
			TokenURL: c.TokenURI,
		},
	}, nil
}
