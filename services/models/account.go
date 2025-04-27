package models

import "fmt"

type GetRepoHubName int

const (
	GitHub GetRepoHubName = iota
	GitLab
)

func (gr GetRepoHubName) String() string {
	switch gr {
	case GetRepoHubName(GitHub):
		return "GitHub"
	case GetRepoHubName(GitLab):
		return "GitLab"
	default:
		return fmt.Sprintf("%d", int(gr))
	}
}

type CloudDriveName int

const (
	GoogleDrive CloudDriveName = iota
	YandexDisk
	OneDrive
)

func (cd CloudDriveName) String() string {
	switch cd {
	case CloudDriveName(GoogleDrive):
		return "GoogleDrive"
	case CloudDriveName(YandexDisk):
		return "YandexDisk"
	case CloudDriveName(OneDrive):
		return "OneDrive"
	default:
		return fmt.Sprintf("%d", int(cd))
	}
}

type PlannerName int

const (
	GoogleCalendar PlannerName = iota
	YandexCalendar
)

func (pn PlannerName) String() string {
	switch pn {
	case PlannerName(GoogleCalendar):
		return "GoogleCalendar"
	case PlannerName(YandexCalendar):
		return "YandexCalendar"
	default:
		return fmt.Sprintf("%d", int(pn))
	}
}

type Account struct {
	Id    string
	Login string
	Hash  []byte
	Salt  string
}

type StudentAccount struct {
	Id        string
	Login     string
	StudentId string
}

type BaseIntegration struct {
	AccountId string
	ApiKey    string
	Type      int
}

func (bi BaseIntegration) GetRepoHubTypeAsString() string {
	return GetRepoHubName.String(GetRepoHubName(bi.Type))
}

type PlannerIntegration struct {
	BaseIntegration
	PlannerData
}

func (pi PlannerIntegration) GetTypeAsString() string {
	return PlannerName.String(PlannerName(pi.Type))
}

type CloudDriveIntegration struct {
	BaseIntegration
	DriveData
}

func (cdi CloudDriveIntegration) GetTypeAsString() string {
	return CloudDriveName.String(CloudDriveName(cdi.Type))
}
