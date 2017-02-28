package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strings"

	"github.com/google/go-github/github"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

var (
	ghreposCmd = &cobra.Command{
		Use:   "ghrepos",
		Short: "ghrepos prints a filtered list of GitHub repositories",
		Long:  "TODO",
		RunE:  run,
	}
)

func init() {
	ghreposCmd.PersistentFlags().String("token", "", "GitHub token to use for API authentication")
	must(viper.BindPFlag("token", ghreposCmd.PersistentFlags().Lookup("token")))
	must(viper.BindEnv("token", "GITHUB_TOKEN"))

	ghreposCmd.PersistentFlags().StringP("owner", "o", "", "User or organization filter")
	must(viper.BindPFlag("owner", ghreposCmd.PersistentFlags().Lookup("owner")))
	must(viper.BindEnv("owner", "GITHUB_USER"))

	ghreposCmd.PersistentFlags().Bool("json", false, "Print JSON array instead of human readable list")
	must(viper.BindPFlag("json", ghreposCmd.PersistentFlags().Lookup("json")))
}

func run(cmd *cobra.Command, args []string) error {
	var httpClient *http.Client
	if token := viper.GetString("token"); token != "" {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: viper.GetString("token")},
		)
		httpClient = oauth2.NewClient(oauth2.NoContext, ts)
	}
	c := github.NewClient(httpClient)

	if len(args) < 1 {
		return errors.New("You need to provide a topic")
	}

	repos, err := fetchRepositories(c, viper.GetString("owner"), args)
	if err != nil {
		return err
	}

	output := []string{}
	for _, repo := range repos {
		output = append(output, fmt.Sprintf("%s/%s", *repo.Owner.Login, *repo.Name))
	}

	switch {
	case viper.GetBool("json"):
		return json.NewEncoder(os.Stdout).Encode(output)
	default:
		fmt.Println(strings.Join(output, "\n"))
	}

	return nil
}

func fetchRepositories(c *github.Client, owner string, topics []string) ([]github.Repository, error) {
	query := []string{}

	for _, topic := range topics {
		query = append(query, fmt.Sprintf("topic:%s", topic))
	}

	if owner != "" {
		query = append(query, fmt.Sprintf("user:%s", owner))
	}

	var result []github.Repository
	opt := &github.SearchOptions{}
	for {
		repos, resp, err := c.Search.Repositories(context.Background(), strings.Join(query, " "), opt)
		if err != nil {
			return nil, fmt.Errorf("Could not perform search: %s", err)
		}
		result = append(result, repos.Repositories...)
		if resp.NextPage == 0 {
			break
		}
		opt.ListOptions.Page = resp.NextPage
	}

	sort.Sort(byName(result))

	return result, nil
}

func must(err error) {
	if err != nil {
		abort(err)
	}
}

type byName []github.Repository

func (bn byName) Len() int      { return len(bn) }
func (bn byName) Swap(i, j int) { bn[i], bn[j] = bn[j], bn[i] }
func (bn byName) Less(i, j int) bool {
	if *bn[i].Owner.Login < *bn[j].Owner.Login {
		return true
	}
	if *bn[i].Owner.Login == *bn[j].Owner.Login && *bn[i].Name < *bn[j].Name {
		return true
	}
	return false
}
