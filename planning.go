package main

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
