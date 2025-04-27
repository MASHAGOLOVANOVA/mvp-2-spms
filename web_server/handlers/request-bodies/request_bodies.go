package requestbodies

import "time"

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type CredentialsBot struct {
	Phone string `json:"phone_number"`
}

type SignUp struct {
	Credentials
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Middlename string `json:"middlename"`
}

type StudentSignUp struct {
	Login      string `json:"login"`
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Middlename string `json:"middlename"`
	EdProgName string `json:"ed_prog_name"`
	University string `json:"university"`
	Course     uint   `json:"course"`
}

type SetProfessorPlanner struct {
	Id string `json:"planner_id"`
}

type AddStudent struct {
	Name                   string `json:"name"`
	Surname                string `json:"surname"`
	Middlename             string `json:"middlename"`
	Cource                 int    `json:"cource"`
	EducationalProgrammeId int    `json:"education_programme_id"`
}

type AddMeeting struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	MeetingTime time.Time `json:"meeting_time"`
	StudentId   int       `json:"student_participant_id"`
	ProjectId   int       `json:"project_id,omitempty"`
	IsOnline    bool      `json:"is_online"`
}

type AddSlot struct {
	Description string    `json:"description"`
	MeetingTime time.Time `json:"meeting_time"`
	Duration    int       `json:"duration"`
	IsOnline    bool      `json:"is_online"`
}

type AddProject struct {
	Theme          string   `json:"theme"`
	StudentIds     []string `json:"student_ids"`
	Year           int      `json:"year"`
	RepoOwner      string   `json:"repository_owner_login"`
	RepositoryName string   `json:"repository_name"`
}

type UpdateProject struct {
	Theme          *string `json:"theme,omitempty"`
	StudentId      *int    `json:"student_id,omitempty"`
	Year           *int    `json:"year,omitempty"`
	RepoOwner      *string `json:"repository_owner_login,omitempty"`
	RepositoryName *string `json:"repository_name,omitempty"`
	Status         *int    `json:"status,omitempty"`
	Stage          *int    `json:"stage,omitempty"`
}

type AddTask struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Deadline    time.Time `json:"deadline"`
}

type Apply struct {
	StudentId   *int `json:"student_id,omitempty"`
	ProfessorId *int `json:"professor_id,omitempty"`
}
