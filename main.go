package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
	"gopkg.in/yaml.v3"
)

var planningUpload = flag.Bool("planning", false, "upload planning data to Logilica")
var ciUpload = flag.Bool("ci", false, "upload CI build data to Logilica")

func main() {
	flag.Parse()
	if *planningUpload {
		// client := CreateClient()
		// LogilicaUploadPlanningData(client)
		fmt.Println("Upload Planning data")
	}
	if *ciUpload {
		fmt.Println("Upload CI Build data")
	}
}

func CreateClient() *http.Client {
	client := oauth2.NewClient(
		context.TODO(),
		oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: os.Getenv("TOKEN")},
		))
	return client
}

func LogilicaUploadPlanningData(client *http.Client) {
	config := getOrganizations()
	for _, org := range config.Organizations {
		pid := CreateProjectIdOnLogilica(org.URL, org.Organization, org.Key)
		gitbhubProjectId := GetProjectIdFromGithub(org.Organization, org.URL)
		sprintInformation := GetSprintInformation(gitbhubProjectId, client)
		CreateSprints(pid, sprintInformation)                                  // creating a sprint on logilica
		payload := PreparePayload(gitbhubProjectId, sprintInformation, client) // call graphql api's for preparing the payload
		UploadPlanningData(pid, payload)                                       // upload planning data call to logilica
	}
}

func getOrganizations() *Config {
	var c Config
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	yamlFile, err := os.ReadFile(basepath + "/config/orgs.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	fmt.Println(string(yamlFile))
	if err := yaml.Unmarshal(yamlFile, &c); err != nil {
		panic(err)
	}
	fmt.Println(c)
	return &c
}

func CreateSprints(projectId string, sprintInformation SprintInformation) {
	postBody, _ := json.Marshal(map[string]string{
		"id":    sprintInformation.id,
		"name":  sprintInformation.name,
		"state": sprintInformation.state,
		"goal":  sprintInformation.goal,
	})
	logilicaUrl := fmt.Sprintf("https://logilica.io/api/import/v1/pm/%v/sprints", projectId)
	contentType := "application/json"

	client := &http.Client{}
	req, err := http.NewRequest("POST", logilicaUrl, bytes.NewBuffer(postBody))
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("Content-Type", contentType)
	req.Header.Add("x-lgca-domain", "redhat")
	req.Header.Add("X-lgca-token", os.Getenv("LOGILICA_TOKEN"))
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(body)
}

func CreateProjectIdOnLogilica(url, organizationName, key string) string {
	postBody, _ := json.Marshal(map[string]string{
		"URL":     url,
		"orgName": organizationName,
		"key":     key,
	})

	logilicaUrl := "https://logilica.io/api/import/v1/pm/projects/create"
	contentType := "application/json"

	client := &http.Client{}
	req, err := http.NewRequest("POST", logilicaUrl, bytes.NewBuffer(postBody))
	if err != nil {
		fmt.Println(err)
		return ""
	}
	req.Header.Add("Content-Type", contentType)
	req.Header.Add("x-lgca-domain", "redhat")
	req.Header.Add("X-lgca-token", os.Getenv("LOGILICA_TOKEN"))
	resp, err := client.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	fmt.Println(resp.StatusCode)
	fmt.Println(string(body))
	var projectResponse ProjectApiResponse
	json.Unmarshal(body, &projectResponse)
	return projectResponse.Id
}

func GetProjectIdFromGithub(organizationName string, url string) string {
	var projectId string
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("TOKEN")},
	)
	httpClient := oauth2.NewClient(context.Background(), src)

	client := githubv4.NewClient(httpClient)
	var listquery struct {
		Organization struct {
			Projectv2 struct {
				Nodes []struct {
					Id    string
					Title string
					Url   string
				}
			} `graphql:"projectsV2(first: 20)"`
		} `graphql:"organization(login: $orgname)"`
	}
	variables := map[string]interface{}{
		"orgname": githubv4.String(organizationName),
	}
	err := client.Query(context.Background(), &listquery, variables)
	if err != nil {
		fmt.Println(err)
	}
	nodes := listquery.Organization.Projectv2.Nodes
	for _, node := range nodes {
		if node.Url == url {
			projectId = node.Id
			break
		}
	}
	return projectId

}

func PreparePayload(githubProjectId string, sprintInformation SprintInformation, client *http.Client) []Planning {
	var projectQueryResponse ProjectQueryResponse
	var payload []Planning
	query := `node(id: "{{.Val}}") {
		... on ProjectV2 {
		  id
		  title
		  items(first:100) {
			totalCount
			edges {
			  cursor
			  node {
				itemID: id
				type
				status: fieldValueByName(name: "Status") {
				  # GitHub predefined a "Status" field, type ProjectV2ItemFieldSingleSelectValue, that cannot be deleted.
				  ... on ProjectV2ItemFieldSingleSelectValue {
					name
					updatedAt
					creator {
					  __typename
					  login
					  ... on User {
						login
						name
					  }
					  ... on Organization {
						login
						name
					  }
					  ... on Bot {
						login
					  }
					  ... on EnterpriseUserAccount {
						login
						name
					  }
					  ... on Mannequin {
						login
						email
					  }
					}
				  }
				}
				content {
				  ... on Issue {
					__typename
					title
					author {
					  __typename
					  login
					  ... on User {
						login
						name
					  }
					  ... on Organization {
						login
						name
					  }
					  ... on Bot {
						login
					  }
					  ... on EnterpriseUserAccount {
						login
						name
					  }
					  ... on Mannequin {
						login
						email
					  }
					}
					body
					bodyHTML
					closed
					closedAt
					createdAt
					id
					number
					lastEditedAt

					milestone {
					  id
					  title
					  createdAt
					}
					state
					stateReason
					updatedAt
					url
				  }
				}
			  }
			}
		  }
		}
	  }
	}`
	// parametes to the query
	ids := Parameters{githubProjectId}
	tmpl, err := template.New("new").Parse(query)
	buf := &bytes.Buffer{}
	err = tmpl.Execute(buf, ids)
	if err != nil {
		panic(err)
	}
	gqlMarshalled, err := json.Marshal(graphQLRequest{Query: buf.String()})
	if err != nil {
		panic(err)
	}
	resp, err := client.Post("https://api.github.com/graphql", "application/json", strings.NewReader(string(gqlMarshalled)))
	b, _ := httputil.DumpResponse(resp, true)
	json.Unmarshal(b, &projectQueryResponse)
	// get issue information
	for _, edge := range projectQueryResponse.Data.Node.Items.Edges {
		planningObj := Planning{
			ID:        edge.Node.Content.ID,
			Origin:    "GITHUB",
			CreatedAt: int(edge.Node.Content.CreatedAt.Unix()),
			Creator: Creator{
				Name:         edge.Node.Status.Creator.Name,
				Email:        edge.Node.Status.Creator.Login,
				AccountID:    "",
				LastActivity: 0,
			},
		}
		query = `node(id: "{{.Val}}") {
		... on Issue {
		  id
		  title
		  timelineItems(first: 100) {
			edges {
			  cursor
			  node {
				... on ConvertedNoteToIssueEvent {
				  __typename
				  createdAt
				  actor {
					login
					... on User {
					  login
					  name
					}
					... on Organization {
					  login
					  name
					}
					... on Bot {
					  login
					}
					... on EnterpriseUserAccount {
					  login
					  name
					}
					... on Mannequin {
					  login
					  email
					}
				  }
				}
				... on AddedToProjectEvent {
				  __typename
				  createdAt
				  actor {
					__typename
					login
					... on User {
					  login
					  name
					}
					... on Organization {
					  login
					  name
					}
					... on Bot {
					  login
					}
					... on EnterpriseUserAccount {
					  login
					  name
					}
					... on Mannequin {
					  login
					  email
					}
				  }
				}
				... on AssignedEvent {
				  __typename
				  createdAt
				  assignee {
					__typename
					... on User {
					  login
					  name
					}
					... on Organization {
					  login
					  name
					}
					... on Bot {
					  login
					}
					... on Mannequin {
					  login
					  email
					}
				  }
				  actor {
					__typename
					... on User {
					  login
					  name
					}
				  }
				}
				... on ClosedEvent {
				  __typename
				  createdAt
				  stateReason
				  actor {
					__typename
					... on User {
					  login
					  name
					}
				  }
				}
			  }
			}
		  }
		}
	  }
	}`
		// parametes to the query
		ids = Parameters{edge.Node.Content.ID}
		tmpl, err = template.New("new").Parse(query)
		buf = &bytes.Buffer{}
		err = tmpl.Execute(buf, ids)
		if err != nil {
			panic(err)
		}
		gqlMarshalled, err = json.Marshal(graphQLRequest{Query: buf.String()})
		if err != nil {
			panic(err)
		}
		resp, err = client.Post("https://api.github.com/graphql", "application/json", strings.NewReader(string(gqlMarshalled)))
		b, _ = httputil.DumpResponse(resp, true)
		json.Unmarshal(b, &projectQueryResponse)
		payload = append(payload, planningObj)
	}
	return payload

}

func GetSprintInformation(gitbhubProjectId string, client *http.Client) SprintInformation {
	var sprintInformation SprintInformation
	query := `node(id: "{{.Val}}") {
		... on ProjectV2 {
		  items(last: 20) {
			nodes {
			  id
			  fieldValues(first: 20) {
				nodes {
				  ... on ProjectV2ItemFieldIterationValue {
					iterationId
					title
					startDate
					updatedAt
					duration
				  }
				}
			  }
			}
		  }
		}
	  }
	}`
	// parametes to the query
	ids := Parameters{gitbhubProjectId}
	tmpl, err := template.New("new").Parse(query)
	buf := &bytes.Buffer{}
	err = tmpl.Execute(buf, ids)
	if err != nil {
		panic(err)
	}
	gqlMarshalled, err := json.Marshal(graphQLRequest{Query: buf.String()})
	if err != nil {
		panic(err)
	}
	resp, err := client.Post("https://api.github.com/graphql", "application/json", strings.NewReader(string(gqlMarshalled)))
	b, _ := httputil.DumpResponse(resp, true)
	json.Unmarshal(b, &sprintInformation)
	return sprintInformation
}

func UploadPlanningData(projectId string, payload []Planning) {
	postBody, _ := json.Marshal(payload)
	logilicaUrl := fmt.Sprintf("https://logilica.io/api/import/v1/pm/%v/issues/create", projectId)
	contentType := "application/json"

	client := &http.Client{}
	req, err := http.NewRequest("POST", logilicaUrl, bytes.NewBuffer(postBody))
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("Content-Type", contentType)
	req.Header.Add("x-lgca-domain", "redhat")
	req.Header.Add("X-lgca-token", os.Getenv("LOGILICA_TOKEN"))
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(body)
}
