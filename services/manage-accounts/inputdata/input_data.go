package inputdata

type GetAccountInfo struct {
	AccountId uint
}

type GetAccountIntegrations struct {
	AccountId uint
}

type GetPlannerIntegration struct {
	AccountId uint
}

type GetDriveIntegration struct {
	AccountId uint
}

type GetRepoHubIntegration struct {
	AccountId uint
}

type SetPlannerIntegration struct {
	AccountId uint
	AuthCode  string
	Type      int
}

type SetDriveIntegration struct {
	AccountId uint
	AuthCode  string
	Type      int
}

type SetRepoHubIntegration struct {
	AccountId uint
	AuthCode  string
	Type      int
}
