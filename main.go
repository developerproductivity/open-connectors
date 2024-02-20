package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

func CallGithubApi() {

	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("TOKEN")},
	)
	httpClient := oauth2.NewClient(context.Background(), src)

	client := githubv4.NewClient(httpClient)

	var q struct {
		Repository struct {
			Description string
			Name        string
		} `graphql:"repository(owner: \"open-connectors\", name: \"open-connectors\")"`
	}

	err := client.Query(context.Background(), &q, nil)
	if err != nil {
		// Handle error.
	}
	fmt.Println(q.Repository.Description)
	printJSON(q)

	var query struct {
		Repository struct {
			Issue struct {
				ID        githubv4.ID
				Reactions struct {
					ViewerHasReacted githubv4.Boolean
				} `graphql:"reactions(content:$reactionContent)"`
			} `graphql:"issue(number:$issueNumber)"`
		} `graphql:"repository(owner:$repositoryOwner,name:$repositoryName)"`
	}
	variables := map[string]interface{}{
		"repositoryOwner": githubv4.String("shurcooL-test"),
		"repositoryName":  githubv4.String("test-repo"),
		"issueNumber":     githubv4.Int(2),
		"reactionContent": githubv4.ReactionContentThumbsUp,
	}
	err = client.Query(context.Background(), &q, variables)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("already reacted:", query.Repository.Issue.Reactions.ViewerHasReacted)

	var listquery struct {
		Organization struct {
			Projectv2 struct {
				Nodes []struct {
					Id    string
					Title string
				}
			} `graphql:"projectsV2(first: 20)"`
		} `graphql:"organization(login: $orgname)"`
	}
	variables = map[string]interface{}{
		"orgname": githubv4.String("containers"),
	}
	err = client.Query(context.Background(), &listquery, variables)
	if err != nil {
		fmt.Println(err)
	}
	nodes := listquery.Organization.Projectv2.Nodes
	for _, node := range nodes {
		fmt.Println(node.Id)
		fmt.Println(node.Title)
	}

	type (
		ProjectV2IterationFieldFragment struct {
			id            string
			name          string
			configuration struct {
				iterations struct {
					startDate string
					id        string
				}
			}
		}
		ProjectV2SingleSelectFieldFragment struct {
			id      string
			name    string
			options struct {
				id   string
				name string
			}
		}
	)

	var projquery struct {
		Node struct {
			fields struct {
				Nodes []struct {
					Id    string
					Title string
				} `graphql:"... on ProjectV2Field "`
				ProjectV2IterationFieldFragment    `graphql:"... on ProjectV2IterationField "`
				ProjectV2SingleSelectFieldFragment `graphql:"... on ProjectV2SingleSelectField "`
			} `graphql:"... on ProjectV2"`
		} `graphql:"node(id:\"PVT_kwDOAFmk9s4AB47o\")"`
	}
	// 	variables = map[string]interface{}{
	// 	"id": githubv4.String("PVT_kwDOAFmk9s4AB47o"),
	// }
	err = client.Query(context.Background(), &projquery, nil)
	if err != nil {
		fmt.Println(err)
	}
}

// printJSON prints v as JSON encoded with indent to stdout. It panics on any error.
func printJSON(v interface{}) {
	w := json.NewEncoder(os.Stdout)
	w.SetIndent("", "\t")
	err := w.Encode(v)
	if err != nil {
		panic(err)
	}
}

func main() {
	CallGithubApi()
}
