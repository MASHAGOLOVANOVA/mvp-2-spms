package manageprojects

import (
	"errors"
	"fmt"
	"html/template"
	domainaggregate "mvp-2-spms/domain-aggregate"
	"mvp-2-spms/services/interfaces"
	"mvp-2-spms/services/manage-projects/inputdata"
	"mvp-2-spms/services/manage-projects/outputdata"
	"mvp-2-spms/services/models"
	"os"
	"strconv"

	"golang.org/x/oauth2"
)

type ProjectInteractor struct {
	projectRepo interfaces.IProjetRepository
	studentRepo interfaces.IStudentRepository
	uniRepo     interfaces.IUniversityRepository
	accountRepo interfaces.IAccountRepository
}

func InitProjectInteractor(projRepo interfaces.IProjetRepository, stRepo interfaces.IStudentRepository,
	uniRepo interfaces.IUniversityRepository, accRepo interfaces.IAccountRepository) *ProjectInteractor {
	return &ProjectInteractor{
		projectRepo: projRepo,
		studentRepo: stRepo,
		uniRepo:     uniRepo,
		accountRepo: accRepo,
	}
}

// returns all professor projects (basic information)
func (p *ProjectInteractor) GetProfessorProjects(input inputdata.GetProfessorProjects) (outputdata.GetProfessorProjects, error) {
	// get from database
	projEntities := []outputdata.GetProfessorProjectsEntities{}

	var projects []domainaggregate.Project
	var err error
	if input.FilterStatus != nil {
		projects, err = p.projectRepo.GetProfessorProjectsWithFilters(fmt.Sprint(input.ProfessorId), *input.FilterStatus)

	} else {
		projects, err = p.projectRepo.GetProfessorProjects(fmt.Sprint(input.ProfessorId))
	}

	if err != nil {
		return outputdata.GetProfessorProjects{}, err
	}

	for _, project := range projects {
		student, err := p.studentRepo.GetStudentById(project.StudentId)
		if err != nil {
			return outputdata.GetProfessorProjects{}, err
		}

		projEntities = append(projEntities, outputdata.GetProfessorProjectsEntities{
			Project: project,
			Student: student,
		})
	}
	output := outputdata.MapToGetProfessorProjects(projEntities)
	return output, nil
}

// returns all commits from all branches from specific time
func (p *ProjectInteractor) GetProjectCommits(input inputdata.GetProjectCommits, gitRepositoryHub interfaces.IGitRepositoryHub) (outputdata.GetProjectCommits, error) {
	// get project repo id
	repo, err := p.projectRepo.GetProjectRepository(fmt.Sprint(input.ProjectId))
	if err != nil {
		return outputdata.GetProjectCommits{}, err
	}

	// getting professor repo hub info, should be checked for existance later
	resChan := p.accountRepo.GetAccountRepoHubData(fmt.Sprint(input.ProfessorId))
	resRepo := <-resChan
	if resRepo.Err != nil {
		return outputdata.GetProjectCommits{}, resRepo.Err
	}

	//////////////////////////////////////////////////////////////////////////////////////////////////////
	// check for access token first????????????????????????????????????????????
	token := &oauth2.Token{
		RefreshToken: resRepo.BaseIntegration.ApiKey,
	}

	err = gitRepositoryHub.Authentificate(token)
	if err != nil {
		return outputdata.GetProjectCommits{}, err
	}

	// update token
	resRepo.BaseIntegration.ApiKey = token.RefreshToken
	// save new token
	resChan1 := p.accountRepo.UpdateAccountRepoHubIntegration(resRepo.BaseIntegration)
	resErr := <-resChan1
	if resErr.Err != nil {
		return outputdata.GetProjectCommits{}, resErr.Err
	}

	// get from repo
	commits, err := gitRepositoryHub.GetRepositoryCommitsFromTime(repo, input.From)
	if err != nil {
		return outputdata.GetProjectCommits{}, err
	}

	output := outputdata.MapToGetProjectCommits(commits)
	return output, nil
}

// returns detailed project data (with student data and ed programme)
func (p *ProjectInteractor) GetProjectById(input inputdata.GetProjectById) (outputdata.GetProjectById, error) {
	// get project by id
	project, err := p.projectRepo.GetProjectById(fmt.Sprint(input.ProjectId))
	if err != nil {
		return outputdata.GetProjectById{}, err
	}

	cloudFolder, err := p.projectRepo.GetProjectFolderLink(fmt.Sprint(input.ProjectId))
	if err != nil {
		if !errors.Is(err, models.ErrProjectCloudFolderNotFound) {
			return outputdata.GetProjectById{}, err
		}
		cloudFolder = "" // поменять на нил
	}

	// getting student info
	student, err := p.studentRepo.GetStudentById(project.StudentId)
	if err != nil {
		return outputdata.GetProjectById{}, err
	}

	edProg, err := p.uniRepo.GetEducationalProgrammeById(student.EducationalProgrammeId)
	if err != nil {
		return outputdata.GetProjectById{}, err
	}

	output := outputdata.MapToGetProjectsById(project, student, edProg, cloudFolder)
	return output, nil
}

// returns project statistics
func (p *ProjectInteractor) GetProjectSupReport(input inputdata.GetProjectSupReport) (string, error) {
	report := models.SupReport{}
	report.Comment = input.Comment
	projId := fmt.Sprint(input.ProjectId)

	project, err := p.projectRepo.GetProjectById(projId)
	if err != nil {
		return "", err
	}
	report.Theme = project.Theme

	projectGrading, err := p.projectRepo.GetProjectGradingById(projId)
	if err != nil {
		return "", err
	}
	report.Date = projectGrading.SupervisorReview.CreationDate.Format("02.01.2006")
	for i, c := range projectGrading.SupervisorReview.Criterias {
		report.Items = append(report.Items, models.Ctiteria{
			Num:   fmt.Sprint(i + 1),
			Name:  c.Description,
			Grade: fmt.Sprint(c.Grade),
		})
	}
	report.SupRewGrade = fmt.Sprint(projectGrading.SupervisorReview.GetGrade())

	student, err := p.studentRepo.GetStudentById(project.StudentId)
	if err != nil {
		return "", err
	}
	report.StudentName = student.FullNameToString()
	report.Course = fmt.Sprint(student.Cource)

	resChan := p.accountRepo.GetProfessorById(fmt.Sprint(input.ProfessorId))
	resProf := <-resChan
	if resProf.Err != nil {
		return "", resProf.Err
	}
	report.ProfName = fmt.Sprint([]rune(resProf.Professor.Name)[0], ".", []rune(resProf.Professor.Middlename)[0], ". ", resProf.Professor.Surname)
	report.ScienceDegree = resProf.Professor.ScienceDegree

	ep, err := p.uniRepo.GetEducationalProgrammeFullById(fmt.Sprint(student.EducationalProgrammeId))
	if err != nil {
		return "", err
	}
	report.EdProgramme = ep.Name
	report.Dept = ep.Dept
	report.Faculty = ep.Faculty

	template, err := template.New("./report.docx").ParseFiles("./report.docx")
	if err != nil {
		return "", err
	}

	// filename := fmt.Sprint("", time.Now().UTC().Format(time.RFC3339Nano), "-", input.ProfessorId, "-", input.ProjectId, ".docx")
	filename := fmt.Sprint("abc", ".docx")
	var f *os.File
	f, err = os.Create(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()

	if err := template.Execute(f, report); err != nil {
		return "", err
	}

	return filename, nil
}

// returns project statistics
func (p *ProjectInteractor) GetProjectStatsById(input inputdata.GetProjectStatsById) (outputdata.GetProjectStatsById, error) {
	stats := models.ProjectStats{}
	projId := fmt.Sprint(input.ProjectId)

	projectGrading, err := p.projectRepo.GetProjectGradingById(projId)
	if err != nil {
		return outputdata.GetProjectStatsById{}, err
	}
	stats.ProjectGrading = projectGrading

	stats.MeetingInfo, err = p.projectRepo.GetProjectMeetingInfoById(projId)
	if err != nil {
		return outputdata.GetProjectStatsById{}, err
	}

	stats.TasksInfo, err = p.projectRepo.GetProjectTaskInfoById(projId)
	if err != nil {
		return outputdata.GetProjectStatsById{}, err
	}

	output := outputdata.MapToGetProjectStatsById(stats)
	return output, nil
}

func (p *ProjectInteractor) UpdateProject(input inputdata.UpdateProject, cloudDrive interfaces.ICloudDrive) error {
	project, err := p.projectRepo.GetProjectById(fmt.Sprint(input.Id))
	if err != nil {
		return err
	}

	supId, err := strconv.Atoi(project.SupervisorId)
	if err != nil {
		return err
	}
	if supId != *input.ProfessorId {
		return models.ErrProjectNotProfessors
	}

	projPointer := &project
	err = input.UpdateProjectEntity(projPointer)
	if err != nil {
		return err
	}

	err = p.projectRepo.UpdateProject(*projPointer)
	if err != nil {
		return err
	}
	return nil
}

func (p *ProjectInteractor) UpdateProjectGrading(input inputdata.UpdateProjectGrading) error {
	project, err := p.projectRepo.GetProjectById(fmt.Sprint(input.ProjId))
	if err != nil {
		return err
	}

	supId, err := strconv.Atoi(project.SupervisorId)
	if err != nil {
		return err
	}
	if supId != *input.ProfessorId {
		return models.ErrProjectNotProfessors
	}

	if input.DefenctGrade != nil {
		err = p.projectRepo.UpdateProjectDefenceGrade(fmt.Sprint(input.ProjId), *input.DefenctGrade)
		if err != nil {
			return err
		}
	} else {
		err = p.projectRepo.UpdateProjectDefenceGrade(fmt.Sprint(input.ProjId), 0)
		if err != nil {
			return err
		}
	}

	if input.SupervisorReview != nil {
		sr := domainaggregate.SupervisorReview{}
		if input.SupervisorReview.Id != nil {
			sr.Id = uint(*input.SupervisorReview.Id)
		}
		if input.SupervisorReview.Criterias != nil {
			sr.Criterias = []domainaggregate.Criteria{}
			for _, c := range *input.SupervisorReview.Criterias {
				cr := domainaggregate.Criteria{
					Description: c.Criteria,
					Weight:      c.Weight,
				}
				if c.Grade != nil {
					cr.Grade = *c.Grade
				}
				sr.Criterias = append(sr.Criterias, cr)
			}
		}
		sr.CreationDate = input.SupervisorReview.CreationDate

		err = p.projectRepo.UpdateProjectSupRew(fmt.Sprint(input.ProjId), sr)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *ProjectInteractor) AddProject(input inputdata.AddProject, cloudDrive interfaces.ICloudDrive) (outputdata.AddProject, error) {
	// add to db with repository
	proj, err := p.projectRepo.CreateProjectWithRepository(input.MapToProjectEntity(), input.MapToRepositoryEntity())
	if err != nil {
		return outputdata.AddProject{}, err
	}

	// getting professor drive info, should be checked for existance later
	found := true
	resChan := p.accountRepo.GetAccountDriveData(fmt.Sprint(input.ProfessorId))
	resDrive := <-resChan
	if resDrive.Err != nil {
		if !errors.Is(resDrive.Err, models.ErrAccountDriveDataNotFound) {
			return outputdata.AddProject{}, resDrive.Err
		}
		found = false
	}

	if found {
		//////////////////////////////////////////////////////////////////////////////////////////////////////
		// check for access token first????????????????????????????????????????????
		token := &oauth2.Token{
			RefreshToken: resDrive.CloudDriveIntegration.ApiKey,
		}
		err := cloudDrive.Authentificate(token)
		if err != nil {
			return outputdata.AddProject{}, err
		}

		// add folder to cloud
		driveProject, err := cloudDrive.AddProjectFolder(proj.Project, resDrive.CloudDriveIntegration)
		if err != nil {
			return outputdata.AddProject{}, err
		}

		// add folder id from drive
		err = p.projectRepo.AssignDriveFolder(driveProject)
		if err != nil {
			return outputdata.AddProject{}, err
		}
	}

	// returning id
	output := outputdata.MapToAddProject(proj.Project)
	return output, nil
}
