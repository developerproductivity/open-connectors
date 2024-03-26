package main

import "time"

type Project struct {
	URL          string `json:"url"`
	Key          string `json:"key"`
	Organization string `json:"organization"`
}

type SprintInformation struct {
	goal  string
	state string
	id    string
	name  string
}

type ProjectApiResponse struct {
	Id string `json:"id"`
}

type Config struct {
	Organizations []Project `json:"organizations"`
}

type Sprint struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Goal  string `json:"goal"`
	State string `json:"state"`
}

type Planning struct {
	ID                 string   `json:"id"`
	Origin             string   `json:"origin"`
	CreatedAt          int      `json:"createdAt"`
	Creator            Creator  `json:"creator"`
	Reporter           Reporter `json:"reporter"`
	AssignedAt         int      `json:"assignedAt"`
	InProgressAt       int      `json:"inProgressAt"`
	ResolvedAt         int      `json:"resolvedAt"`
	Assignee           Assignee `json:"assignee"`
	Resolver           Resolver `json:"resolver"`
	Resolution         string   `json:"resolution"`
	Type               string   `json:"type"`
	Status             string   `json:"status"`
	StatusCategory     string   `json:"statusCategory"`
	Summary            string   `json:"summary"`
	Description        string   `json:"description"`
	URL                string   `json:"url"`
	Labels             []string `json:"labels"`
	SprintKeys         []string `json:"sprintKeys"`
	Priority           string   `json:"priority"`
	StoryPointEstimate int      `json:"storyPointEstimate"`
	ParentIssue        string   `json:"parentIssue"`
	Events             []Event  `json:"events"`
}

type Creator struct {
	Name         string `json:"name"`
	Email        string `json:"email"`
	AccountID    string `json:"accountId"`
	LastActivity int    `json:"lastActivity"`
}

type Reporter struct {
	Name         string `json:"name"`
	Email        string `json:"email"`
	AccountID    string `json:"accountId"`
	LastActivity int    `json:"lastActivity"`
}

type Assignee struct {
	Name         string `json:"name"`
	Email        string `json:"email"`
	AccountID    string `json:"accountId"`
	LastActivity int    `json:"lastActivity"`
}

type Resolver struct {
	Name         string `json:"name"`
	Email        string `json:"email"`
	AccountID    string `json:"accountId"`
	LastActivity int    `json:"lastActivity"`
}

type Event struct {
	Type   string `json:"type"`
	Author struct {
		Name         string `json:"name"`
		Email        string `json:"email"`
		AccountID    string `json:"accountId"`
		LastActivity int    `json:"lastActivity"`
	} `json:"author"`
	CreatedAt int `json:"createdAt"`
	To        struct {
		Name         string `json:"name"`
		Email        string `json:"email"`
		AccountID    string `json:"accountId"`
		LastActivity int    `json:"lastActivity"`
	} `json:"to,omitempty"`
	From string `json:"from,omitempty"`
}

type ProjectQueryResponse struct {
	Data struct {
		Node struct {
			ID    string `json:"id"`
			Title string `json:"title"`
			Items struct {
				TotalCount int `json:"totalCount"`
				Edges      []struct {
					Cursor string `json:"cursor"`
					Node   struct {
						ItemID string `json:"itemID"`
						Type   string `json:"type"`
						Status struct {
							Name      string    `json:"name"`
							UpdatedAt time.Time `json:"updatedAt"`
							Creator   struct {
								Typename string `json:"__typename"`
								Login    string `json:"login"`
								Name     string `json:"name"`
							} `json:"creator"`
						} `json:"status"`
						Content struct {
							Typename string `json:"__typename"`
							Title    string `json:"title"`
							Author   struct {
								Typename string `json:"__typename"`
								Login    string `json:"login"`
								Name     string `json:"name"`
							} `json:"author"`
							Body         string      `json:"body"`
							BodyHTML     string      `json:"bodyHTML"`
							Closed       bool        `json:"closed"`
							ClosedAt     time.Time   `json:"closedAt"`
							CreatedAt    time.Time   `json:"createdAt"`
							ID           string      `json:"id"`
							Number       int         `json:"number"`
							LastEditedAt interface{} `json:"lastEditedAt"`
							Milestone    interface{} `json:"milestone"`
							State        string      `json:"state"`
							StateReason  string      `json:"stateReason"`
							UpdatedAt    time.Time   `json:"updatedAt"`
							URL          string      `json:"url"`
						} `json:"content"`
					} `json:"node"`
				} `json:"edges"`
			} `json:"items"`
		} `json:"node"`
	} `json:"data"`
}

type graphQLRequest struct {
	Query     string `json:"query"`
	Variables string `json:"Variables"`
}

type Parameters struct {
	Val string
}
