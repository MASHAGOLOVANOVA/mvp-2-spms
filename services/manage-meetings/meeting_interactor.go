package managemeetings

import (
	"errors"
	"fmt"
	domainaggregate "mvp-2-spms/domain-aggregate"
	"mvp-2-spms/services/interfaces"
	"mvp-2-spms/services/manage-meetings/inputdata"
	"mvp-2-spms/services/manage-meetings/outputdata"
	"mvp-2-spms/services/models"
	"slices"

	"golang.org/x/oauth2"
)

type MeetingInteractor struct {
	meetingRepo interfaces.IMeetingRepository
	accountRepo interfaces.IAccountRepository
	projectRepo interfaces.IProjetRepository
	studentRepo interfaces.IStudentRepository
}

func InitMeetingInteractor(mtRepo interfaces.IMeetingRepository, accRepo interfaces.IAccountRepository,
	sRepo interfaces.IStudentRepository, projRepo interfaces.IProjetRepository) *MeetingInteractor {
	return &MeetingInteractor{
		meetingRepo: mtRepo,
		accountRepo: accRepo,
		studentRepo: sRepo,
		projectRepo: projRepo,
	}
}

func (m *MeetingInteractor) AddMeeting(input inputdata.AddMeeting, planner interfaces.IPlannerService) (outputdata.AddMeeting, error) {
	// adding meeting to db, returns created meeting (with id)
	meeting, err := m.meetingRepo.CreateMeeting(input.MapToMeetingEntity())
	if err != nil {
		return outputdata.AddMeeting{}, err
	}

	// getting calendar info, should be checked for existance later
	plannerFound := true
	resChan := m.accountRepo.GetAccountPlannerData(fmt.Sprint(input.ProfessorId))
	resPlanner := <-resChan
	if resPlanner.Err != nil {
		if !errors.Is(resPlanner.Err, models.ErrAccountPlannerDataNotFound) {
			return outputdata.AddMeeting{}, resPlanner.Err
		}
		plannerFound = false
	}

	meeitngPlanner := models.PlannerMeeting{Meeting: meeting}
	if plannerFound {
		//////////////////////////////////////////////////////////////////////////////////////////////////////
		// check for access token first????????????????????????????????????????????
		token := &oauth2.Token{
			RefreshToken: resPlanner.PlannerIntegration.ApiKey,
		}

		err = planner.Authentificate(token)
		if err != nil {
			return outputdata.AddMeeting{}, err
		}

		// add meeting to calendar
		meeitngPlanner, err = planner.AddMeeting(meeting, resPlanner.PlannerIntegration)
		if err != nil {
			return outputdata.AddMeeting{}, err
		}
	}

	// add meeting id from planner
	err = m.meetingRepo.AssignPlannerMeeting(meeitngPlanner)
	if err != nil {
		return outputdata.AddMeeting{}, err
	}

	// returning id
	output := outputdata.MapToAddMeeting(meeting)
	return output, nil
}

func (m *MeetingInteractor) GetProfessorMeetings(input inputdata.GetProfessorMeetings, planner interfaces.IPlannerService) (outputdata.GetProfesorMeetings, error) {
	// get from db
	meetings, err := m.meetingRepo.GetProfessorMeetings(fmt.Sprint(input.ProfessorId), input.From, input.To)
	if err != nil {
		return outputdata.GetProfesorMeetings{}, err
	}

	meetEntities := []outputdata.GetProfesorMeetingsEntities{}
	// getting calendar info, should be checked for existance later
	plannerFound := true
	resChan := m.accountRepo.GetAccountPlannerData(fmt.Sprint(input.ProfessorId))
	resPlanner := <-resChan
	if resPlanner.Err != nil {
		if !errors.Is(resPlanner.Err, models.ErrAccountPlannerDataNotFound) {
			return outputdata.GetProfesorMeetings{}, resPlanner.Err
		}
		plannerFound = false
	}

	plannerMetingsIds := []string{}
	if plannerFound {
		//////////////////////////////////////////////////////////////////////////////////////////////////////
		// check for access token first????????????????????????????????????????????
		token := &oauth2.Token{
			RefreshToken: resPlanner.PlannerIntegration.ApiKey,
		}

		err = planner.Authentificate(token)
		if err != nil {
			return outputdata.GetProfesorMeetings{}, err
		}

		plannerMetingsIds, err = planner.GetScheduleMeetingIds(input.From, resPlanner.PlannerIntegration)
		if err != nil {
			return outputdata.GetProfesorMeetings{}, err
		}
	}

	for _, meet := range meetings {
		student, err := m.studentRepo.GetStudentById(meet.ParticipantId)
		if err != nil {
			return outputdata.GetProfesorMeetings{}, err
		}

		proj, err := m.projectRepo.GetStudentCurrentProject(meet.ParticipantId)
		if err != nil {
			if !errors.Is(err, models.ErrStudentHasNoCurrentProject) {
				return outputdata.GetProfesorMeetings{}, err
			}
			proj = domainaggregate.Project{} // change to nil
		}

		// getting planner meeting id
		plannerId, err := m.meetingRepo.GetMeetingPlannerId(meet.Id)
		if err != nil {
			return outputdata.GetProfesorMeetings{}, err
		}

		// check if meeting exists in planner
		hasPlanner := slices.Contains(plannerMetingsIds, plannerId)
		meetEntities = append(meetEntities, outputdata.GetProfesorMeetingsEntities{
			Meeting:           meet,
			Student:           student,
			Project:           proj,
			HasPlannerMeeting: hasPlanner,
		})
	}

	output := outputdata.MapToGetProfesorMeetings(meetEntities)
	return output, nil
}
