package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
	"text/template"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	"golang.org/x/oauth2"
)

func HandlePlanningData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"test": "ok}`))
}

func HandleBuildData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"test": "ok}`))
}

func Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"test": "ok}`))
}

func route() (n *negroni.Negroni, rt *mux.Router) {
	router := mux.NewRouter()

	router.HandleFunc("/", Handler).Methods("GET")
	router.HandleFunc("/", HandlePlanningData).Methods("POST")
	router.HandleFunc("/", HandleBuildData).Methods("POST")
	n = negroni.New(negroni.NewRecovery(), negroni.NewLogger())
	return n, router
}

func main() {
	CallGithubGraphqlApi()
	n, router := route()
	n.UseHandler(router)
	n.Run(":8080")
}

type graphQLRequest struct {
	Query     string `json:"query"`
	Variables string `json:"Variables"`
}

type Projectid struct {
	Val string
	Id  int
}

func CallGithubGraphqlApi() {
	client := oauth2.NewClient(
		context.TODO(),
		oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: os.Getenv("TOKEN")},
		))

	query := `{
		node(id: "{{.Val}}") {
		  ... on ProjectV2 {
			fields(first: {{.Id}}) {
			  nodes {
				... on ProjectV2Field {
				  id
				  name
				}
				... on ProjectV2IterationField {
				  id
				  name
				  configuration {
					iterations {
					  startDate
					  id
					}
				  }
				}
				... on ProjectV2SingleSelectField {
				  id
				  name
				  options {
					id
					name
				  }
				}
			  }
			}
		  }
		}
	  }`

	// parametes to the query
	ids := Projectid{"PVT_kwDOCR9iiM4AbwsU", 20}
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
	fmt.Println(string(b))

}
