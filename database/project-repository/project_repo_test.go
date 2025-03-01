package projectrepository

import (
	"fmt"
	"mvp-2-spms/database"
	accountrepository "mvp-2-spms/database/account-repository"
	meetingrepository "mvp-2-spms/database/meeting-repository"
	"mvp-2-spms/database/models"
	studentrepository "mvp-2-spms/database/student-repository"
	taskrepository "mvp-2-spms/database/task-repository"
	domainaggregate "mvp-2-spms/domain-aggregate"
	usecasemodels "mvp-2-spms/services/models"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var dsn = "root:root@tcp(127.0.0.1:3308)/student_project_management_test?parseTime=true"

func connectDB() *database.Database {
	gdb, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // use singular table name, table for `User` would be `user` with this option enabled
		},
	})
	db := database.InitDatabade(gdb)
	return db
}

func TestProjectRepo_CreateProject(t *testing.T) {
	db := connectDB()

	t.Run("fail, bad prof id", func(t *testing.T) {
		// arrange
		pr := InitProjectRepository(*db)
		sr := studentrepository.InitStudentRepository(*db)

		stud, err := sr.CreateStudent(domainaggregate.Student{
			Person: domainaggregate.Person{
				Name:       "dsf",
				Surname:    "sdf",
				Middlename: "sdf",
			},
			Cource: 1,
		})

		assert.NoError(t, err)

		// act
		_, err = pr.CreateProject(domainaggregate.Project{
			Theme:        time.Now().Format(time.RFC3339),
			SupervisorId: "123",
			StudentId:    stud.Id,
			Year:         2023,
			Stage:        domainaggregate.Analysis,
			Status:       domainaggregate.ProjectInProgress,
		})

		// assert
		assert.Error(t, err)

		// cleanup

		studId, err := strconv.Atoi(stud.Id)
		assert.NoError(t, err)
		err = sr.DeleteStudent(studId)
		assert.NoError(t, err)
	})

	t.Run("fail, bad student id", func(t *testing.T) {
		// arrange
		pr := InitProjectRepository(*db)
		ar := accountrepository.InitAccountRepository(*db)

		resProfChan := ar.AddProfessor(domainaggregate.Professor{
			Person: domainaggregate.Person{
				Name:       "dsf",
				Surname:    "sdf",
				Middlename: "sdf",
			},
			ScienceDegree: time.Now().Format(time.RFC3339),
		})
		resProf := <-resProfChan
		assert.NoError(t, resProf.Err)

		// act
		_, err := pr.CreateProject(domainaggregate.Project{
			Theme:        time.Now().Format(time.RFC3339),
			SupervisorId: resProf.Professor.Id,
			StudentId:    "122",
			Year:         2023,
			Stage:        domainaggregate.Analysis,
			Status:       domainaggregate.ProjectInProgress,
		})

		// assert
		assert.Error(t, err)

		// cleanup
		profId, err := strconv.Atoi(resProf.Professor.Id)
		assert.NoError(t, err)
		resErrChan := ar.DeleteProfessor(profId)
		resErr := <-resErrChan
		assert.NoError(t, resErr.Err)
	})

	t.Run("ok", func(t *testing.T) {
		// arrange
		pr := InitProjectRepository(*db)
		ar := accountrepository.InitAccountRepository(*db)
		sr := studentrepository.InitStudentRepository(*db)

		stud, err := sr.CreateStudent(domainaggregate.Student{
			Person: domainaggregate.Person{
				Name:       "dsf",
				Surname:    "sdf",
				Middlename: "sdf",
			},
			Cource: 1,
		})

		assert.NoError(t, err)

		resProfChan := ar.AddProfessor(domainaggregate.Professor{
			Person: domainaggregate.Person{
				Name:       "dsf",
				Surname:    "sdf",
				Middlename: "sdf",
			},
			ScienceDegree: time.Now().Format(time.RFC3339),
		})
		resProf := <-resProfChan
		assert.NoError(t, resProf.Err)

		// act
		proj, err := pr.CreateProject(domainaggregate.Project{
			Theme:        time.Now().Format(time.RFC3339),
			SupervisorId: resProf.Professor.Id,
			StudentId:    stud.Id,
			Year:         2023,
			Stage:        domainaggregate.Analysis,
			Status:       domainaggregate.ProjectInProgress,
		})

		// assert
		assert.NoError(t, err)
		foundProj, err := pr.GetProjectById(proj.Id)
		assert.NoError(t, err)
		assert.Equal(t, proj.Id, foundProj.Id)
		assert.Equal(t, proj.Stage, foundProj.Stage)
		assert.Equal(t, proj.Status, foundProj.Status)
		assert.Equal(t, proj.StudentId, foundProj.StudentId)
		assert.Equal(t, proj.SupervisorId, foundProj.SupervisorId)
		assert.Equal(t, proj.Theme, foundProj.Theme)
		assert.Equal(t, proj.Year, foundProj.Year)

		// cleanup
		profId, err := strconv.Atoi(resProf.Professor.Id)
		assert.NoError(t, err)
		resErrChan := ar.DeleteProfessor(profId)
		resErr := <-resErrChan
		assert.NoError(t, resErr.Err)

		studId, err := strconv.Atoi(stud.Id)
		assert.NoError(t, err)
		err = sr.DeleteStudent(studId)
		assert.NoError(t, err)
	})

}

func TestProjectRepo_CreateProjectWithRepository(t *testing.T) {
	db := connectDB()

	t.Run("fail, bad prof id", func(t *testing.T) {
		// arrange
		pr := InitProjectRepository(*db)
		sr := studentrepository.InitStudentRepository(*db)

		stud, err := sr.CreateStudent(domainaggregate.Student{
			Person: domainaggregate.Person{
				Name:       "dsf",
				Surname:    "sdf",
				Middlename: "sdf",
			},
			Cource: 1,
		})
		assert.NoError(t, err)

		repo := usecasemodels.Repository{
			RepoId:    "23",
			RepoType:  int(usecasemodels.GitHub),
			OwnerName: "klf",
		}

		// act
		_, err = pr.CreateProjectWithRepository(domainaggregate.Project{
			Theme:        time.Now().Format(time.RFC3339),
			SupervisorId: "123",
			StudentId:    stud.Id,
			Year:         2023,
			Stage:        domainaggregate.Analysis,
			Status:       domainaggregate.ProjectInProgress,
		}, repo)

		// assert
		assert.Error(t, err)

		// cleanup

		studId, err := strconv.Atoi(stud.Id)
		assert.NoError(t, err)
		err = sr.DeleteStudent(studId)
		assert.NoError(t, err)
	})

	t.Run("fail, bad student id", func(t *testing.T) {
		// arrange
		pr := InitProjectRepository(*db)
		ar := accountrepository.InitAccountRepository(*db)

		resProfChan := ar.AddProfessor(domainaggregate.Professor{
			Person: domainaggregate.Person{
				Name:       "dsf",
				Surname:    "sdf",
				Middlename: "sdf",
			},
			ScienceDegree: time.Now().Format(time.RFC3339),
		})
		resProf := <-resProfChan
		assert.NoError(t, resProf.Err)

		repo := usecasemodels.Repository{
			RepoId:    "23",
			RepoType:  int(usecasemodels.GitHub),
			OwnerName: "klf",
		}

		// act
		_, err := pr.CreateProjectWithRepository(domainaggregate.Project{
			Theme:        time.Now().Format(time.RFC3339),
			SupervisorId: resProf.Professor.Id,
			StudentId:    "122",
			Year:         2023,
			Stage:        domainaggregate.Analysis,
			Status:       domainaggregate.ProjectInProgress,
		}, repo)

		// assert
		assert.Error(t, err)

		// cleanup
		profId, err := strconv.Atoi(resProf.Professor.Id)
		assert.NoError(t, err)
		resErrChan := ar.DeleteProfessor(profId)
		resErr := <-resErrChan
		assert.NoError(t, resErr.Err)
	})

	t.Run("ok", func(t *testing.T) {
		// arrange
		pr := InitProjectRepository(*db)
		ar := accountrepository.InitAccountRepository(*db)
		sr := studentrepository.InitStudentRepository(*db)

		stud, err := sr.CreateStudent(domainaggregate.Student{
			Person: domainaggregate.Person{
				Name:       "dsf",
				Surname:    "sdf",
				Middlename: "sdf",
			},
			Cource: 1,
		})
		assert.NoError(t, err)

		resProfChan := ar.AddProfessor(domainaggregate.Professor{
			Person: domainaggregate.Person{
				Name:       "dsf",
				Surname:    "sdf",
				Middlename: "sdf",
			},
			ScienceDegree: time.Now().Format(time.RFC3339),
		})
		resProf := <-resProfChan
		assert.NoError(t, resProf.Err)

		repo := usecasemodels.Repository{
			RepoId:    "23",
			RepoType:  int(usecasemodels.GitHub),
			OwnerName: "klf",
		}

		// act
		proj, err := pr.CreateProjectWithRepository(domainaggregate.Project{
			Theme:        time.Now().Format(time.RFC3339),
			SupervisorId: resProf.Professor.Id,
			StudentId:    stud.Id,
			Year:         2023,
			Stage:        domainaggregate.Analysis,
			Status:       domainaggregate.ProjectInProgress,
		}, repo)

		// assert
		assert.NoError(t, err)
		foundProj, err := pr.GetProjectById(proj.Id)
		assert.NoError(t, err)
		assert.Equal(t, proj.Id, foundProj.Id)
		assert.Equal(t, proj.Stage, foundProj.Stage)
		assert.Equal(t, proj.Status, foundProj.Status)
		assert.Equal(t, proj.StudentId, foundProj.StudentId)
		assert.Equal(t, proj.SupervisorId, foundProj.SupervisorId)
		assert.Equal(t, proj.Theme, foundProj.Theme)
		assert.Equal(t, proj.Year, foundProj.Year)

		foundRepo, err := pr.GetProjectRepository(proj.Id)
		assert.NoError(t, err)
		assert.Equal(t, repo.OwnerName, foundRepo.OwnerName)
		assert.Equal(t, repo.RepoId, foundRepo.RepoId)
		assert.Equal(t, repo.RepoType, foundRepo.RepoType)

		// cleanup
		profId, err := strconv.Atoi(resProf.Professor.Id)
		assert.NoError(t, err)
		resErrChan := ar.DeleteProfessor(profId)
		resErr := <-resErrChan
		assert.NoError(t, resErr.Err)

		studId, err := strconv.Atoi(stud.Id)
		assert.NoError(t, err)
		err = sr.DeleteStudent(studId)
		assert.NoError(t, err)

		db.DB.Where("id > 0").Delete(&models.Repository{})
	})

}

func TestProjectRepo_GetProfessorProjects(t *testing.T) {
	db := connectDB()

	t.Run("ok, empty", func(t *testing.T) {
		// arrange
		pr := InitProjectRepository(*db)

		// act
		foundProj, err := pr.GetProfessorProjects("123")

		// assert
		assert.NoError(t, err)
		assert.Equal(t, 0, len(foundProj))
	})

	t.Run("ok, list", func(t *testing.T) {
		// arrange
		pr := InitProjectRepository(*db)
		ar := accountrepository.InitAccountRepository(*db)
		sr := studentrepository.InitStudentRepository(*db)

		stud, err := sr.CreateStudent(domainaggregate.Student{
			Person: domainaggregate.Person{
				Name:       "dsf",
				Surname:    "sdf",
				Middlename: "sdf",
			},
			Cource: 1,
		})
		assert.NoError(t, err)

		resProfChan := ar.AddProfessor(domainaggregate.Professor{
			Person: domainaggregate.Person{
				Name:       "dsf",
				Surname:    "sdf",
				Middlename: "sdf",
			},
			ScienceDegree: time.Now().Format(time.RFC3339),
		})
		resProf := <-resProfChan
		assert.NoError(t, resProf.Err)

		projs := []domainaggregate.Project{}
		for i := 0; i < 10; i++ {
			proj, err := pr.CreateProject(domainaggregate.Project{
				Theme:        fmt.Sprint(i),
				SupervisorId: resProf.Professor.Id,
				StudentId:    stud.Id,
				Year:         2023,
				Stage:        domainaggregate.Analysis,
				Status:       domainaggregate.ProjectInProgress,
			})
			assert.NoError(t, err)
			projs = append(projs, proj)
		}

		// act
		foundProjs, err := pr.GetProfessorProjects(resProf.Professor.Id)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, len(projs), len(foundProjs))
		for i := range projs {
			assert.Equal(t, projs[i].Id, foundProjs[i].Id)
			assert.Equal(t, projs[i].Stage, foundProjs[i].Stage)
			assert.Equal(t, projs[i].Status, foundProjs[i].Status)
			assert.Equal(t, projs[i].StudentId, foundProjs[i].StudentId)
			assert.Equal(t, projs[i].SupervisorId, foundProjs[i].SupervisorId)
			assert.Equal(t, projs[i].Theme, foundProjs[i].Theme)
			assert.Equal(t, projs[i].Year, foundProjs[i].Year)
		}

		// cleanup

		profId, err := strconv.Atoi(resProf.Professor.Id)
		assert.NoError(t, err)
		resErrChan := ar.DeleteProfessor(profId)
		resErr := <-resErrChan
		assert.NoError(t, resErr.Err)

		studId, err := strconv.Atoi(stud.Id)
		assert.NoError(t, err)
		err = sr.DeleteStudent(studId)
		assert.NoError(t, err)
	})

}

func TestProjectRepo_GetProjectById(t *testing.T) {
	db := connectDB()

	t.Run("fail, not found", func(t *testing.T) {
		// arrange
		pr := InitProjectRepository(*db)

		// act
		_, err := pr.GetProjectById("123")

		// assert
		assert.ErrorIs(t, err, usecasemodels.ErrProjectNotFound)
	})

	t.Run("ok", func(t *testing.T) {
		// arrange
		pr := InitProjectRepository(*db)
		ar := accountrepository.InitAccountRepository(*db)
		sr := studentrepository.InitStudentRepository(*db)

		stud, err := sr.CreateStudent(domainaggregate.Student{
			Person: domainaggregate.Person{
				Name:       "dsf",
				Surname:    "sdf",
				Middlename: "sdf",
			},
			Cource: 1,
		})
		assert.NoError(t, err)

		resProfChan := ar.AddProfessor(domainaggregate.Professor{
			Person: domainaggregate.Person{
				Name:       "dsf",
				Surname:    "sdf",
				Middlename: "sdf",
			},
			ScienceDegree: time.Now().Format(time.RFC3339),
		})
		resProf := <-resProfChan
		assert.NoError(t, resProf.Err)
		proj, err := pr.CreateProject(domainaggregate.Project{
			Theme:        time.Now().Format(time.RFC3339),
			SupervisorId: resProf.Professor.Id,
			StudentId:    stud.Id,
			Year:         2023,
			Stage:        domainaggregate.Analysis,
			Status:       domainaggregate.ProjectInProgress,
		})
		assert.NoError(t, err)

		// act
		foundProj, err := pr.GetProjectById(proj.Id)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, proj.Id, foundProj.Id)
		assert.Equal(t, proj.Stage, foundProj.Stage)
		assert.Equal(t, proj.Status, foundProj.Status)
		assert.Equal(t, proj.StudentId, foundProj.StudentId)
		assert.Equal(t, proj.SupervisorId, foundProj.SupervisorId)
		assert.Equal(t, proj.Theme, foundProj.Theme)
		assert.Equal(t, proj.Year, foundProj.Year)

		// cleanup

		profId, err := strconv.Atoi(resProf.Professor.Id)
		assert.NoError(t, err)
		resErrChan := ar.DeleteProfessor(profId)
		resErr := <-resErrChan
		assert.NoError(t, resErr.Err)

		studId, err := strconv.Atoi(stud.Id)
		assert.NoError(t, err)
		err = sr.DeleteStudent(studId)
		assert.NoError(t, err)
	})

}

func TestProjectRepo_GetProjectMeetingInfoById(t *testing.T) {
	db := connectDB()
	t.Run("ok, zero", func(t *testing.T) {
		// arrange
		pr := InitProjectRepository(*db)
		ar := accountrepository.InitAccountRepository(*db)
		sr := studentrepository.InitStudentRepository(*db)

		stud, err := sr.CreateStudent(domainaggregate.Student{
			Person: domainaggregate.Person{
				Name:       "dsf",
				Surname:    "sdf",
				Middlename: "sdf",
			},
			Cource: 1,
		})
		assert.NoError(t, err)

		resProfChan := ar.AddProfessor(domainaggregate.Professor{
			Person: domainaggregate.Person{
				Name:       "dsf",
				Surname:    "sdf",
				Middlename: "sdf",
			},
			ScienceDegree: time.Now().Format(time.RFC3339),
		})

		resProf := <-resProfChan
		assert.NoError(t, resProf.Err)
		proj, err := pr.CreateProject(domainaggregate.Project{
			Theme:        time.Now().Format(time.RFC3339),
			SupervisorId: resProf.Professor.Id,
			StudentId:    stud.Id,
			Year:         2023,
			Stage:        domainaggregate.Analysis,
			Status:       domainaggregate.ProjectInProgress,
		})
		assert.NoError(t, err)

		// act
		info, err := pr.GetProjectMeetingInfoById(proj.Id)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, 0, info.PassedCount)

		// cleanup

		profId, err := strconv.Atoi(resProf.Professor.Id)
		assert.NoError(t, resProf.Err)
		resErrChan := ar.DeleteProfessor(profId)
		resErr := <-resErrChan
		assert.NoError(t, resErr.Err)

		studId, err := strconv.Atoi(stud.Id)
		assert.NoError(t, err)
		err = sr.DeleteStudent(studId)
		assert.NoError(t, err)
	})

	t.Run("ok, non zero", func(t *testing.T) {
		// arrange
		pr := InitProjectRepository(*db)
		ar := accountrepository.InitAccountRepository(*db)
		sr := studentrepository.InitStudentRepository(*db)
		mr := meetingrepository.InitMeetingRepository(*db)

		stud, err := sr.CreateStudent(domainaggregate.Student{
			Person: domainaggregate.Person{
				Name:       "dsf",
				Surname:    "sdf",
				Middlename: "sdf",
			},
			Cource: 1,
		})
		assert.NoError(t, err)

		resProfChan := ar.AddProfessor(domainaggregate.Professor{
			Person: domainaggregate.Person{
				Name:       "dsf",
				Surname:    "sdf",
				Middlename: "sdf",
			},
			ScienceDegree: time.Now().Format(time.RFC3339),
		})
		resProf := <-resProfChan
		assert.NoError(t, resProf.Err)
		proj, err := pr.CreateProject(domainaggregate.Project{
			Theme:        time.Now().Format(time.RFC3339),
			SupervisorId: resProf.Professor.Id,
			StudentId:    stud.Id,
			Year:         2023,
			Stage:        domainaggregate.Analysis,
			Status:       domainaggregate.ProjectInProgress,
		})
		assert.NoError(t, err)

		meet, err := mr.CreateMeeting(domainaggregate.Meeting{
			Name:          time.Now().Format(time.RFC3339),
			Description:   "lwefl",
			OrganizerId:   resProf.Professor.Id,
			ParticipantId: stud.Id,
			Time:          time.Now().UTC().Add(time.Hour * 8),
			IsOnline:      true,
			Status:        domainaggregate.MeetingPassed,
			ProjectId:     proj.Id,
		})
		assert.NoError(t, err)

		// act
		info, err := pr.GetProjectMeetingInfoById(proj.Id)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, 1, info.PassedCount)

		// cleanup

		profId, err := strconv.Atoi(resProf.Professor.Id)
		assert.NoError(t, err)
		resErrChan := ar.DeleteProfessor(profId)
		resErr := <-resErrChan
		assert.NoError(t, resErr.Err)

		studId, err := strconv.Atoi(stud.Id)
		assert.NoError(t, err)
		err = sr.DeleteStudent(studId)
		assert.NoError(t, err)

		meetId, err := strconv.Atoi(meet.Id)
		assert.NoError(t, err)
		err = mr.DeleteMeeting(meetId)
		assert.NoError(t, err)
	})
}

func TestProjectRepo_GetProjectTaskInfoById(t *testing.T) {
	db := connectDB()
	t.Run("ok, zero", func(t *testing.T) {
		// arrange
		pr := InitProjectRepository(*db)
		ar := accountrepository.InitAccountRepository(*db)
		sr := studentrepository.InitStudentRepository(*db)

		stud, err := sr.CreateStudent(domainaggregate.Student{
			Person: domainaggregate.Person{
				Name:       "dsf",
				Surname:    "sdf",
				Middlename: "sdf",
			},
			Cource: 1,
		})
		assert.NoError(t, err)

		resProfChan := ar.AddProfessor(domainaggregate.Professor{
			Person: domainaggregate.Person{
				Name:       "dsf",
				Surname:    "sdf",
				Middlename: "sdf",
			},
			ScienceDegree: time.Now().Format(time.RFC3339),
		})
		resProf := <-resProfChan
		assert.NoError(t, resProf.Err)
		proj, err := pr.CreateProject(domainaggregate.Project{
			Theme:        time.Now().Format(time.RFC3339),
			SupervisorId: resProf.Professor.Id,
			StudentId:    stud.Id,
			Year:         2023,
			Stage:        domainaggregate.Analysis,
			Status:       domainaggregate.ProjectInProgress,
		})
		assert.NoError(t, err)

		// act
		info, err := pr.GetProjectTaskInfoById(proj.Id)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, 0, info.NotStartedCount)
		assert.Equal(t, 0, info.InProgressCount)
		assert.Equal(t, 0, info.FinishedCount)

		// cleanup

		profId, err := strconv.Atoi(resProf.Professor.Id)
		assert.NoError(t, err)
		resErrChan := ar.DeleteProfessor(profId)
		resErr := <-resErrChan
		assert.NoError(t, resErr.Err)

		studId, err := strconv.Atoi(stud.Id)
		assert.NoError(t, err)
		err = sr.DeleteStudent(studId)
		assert.NoError(t, err)
	})

	t.Run("ok, non zero", func(t *testing.T) {
		// arrange
		pr := InitProjectRepository(*db)
		ar := accountrepository.InitAccountRepository(*db)
		sr := studentrepository.InitStudentRepository(*db)
		tr := taskrepository.InitTaskRepository(*db)

		stud, err := sr.CreateStudent(domainaggregate.Student{
			Person: domainaggregate.Person{
				Name:       "dsf",
				Surname:    "sdf",
				Middlename: "sdf",
			},
			Cource: 1,
		})
		assert.NoError(t, err)

		resProfChan := ar.AddProfessor(domainaggregate.Professor{
			Person: domainaggregate.Person{
				Name:       "dsf",
				Surname:    "sdf",
				Middlename: "sdf",
			},
			ScienceDegree: time.Now().Format(time.RFC3339),
		})
		resProf := <-resProfChan
		assert.NoError(t, resProf.Err)
		proj, err := pr.CreateProject(domainaggregate.Project{
			Theme:        time.Now().Format(time.RFC3339),
			SupervisorId: resProf.Professor.Id,
			StudentId:    stud.Id,
			Year:         2023,
			Stage:        domainaggregate.Analysis,
			Status:       domainaggregate.ProjectInProgress,
		})
		assert.NoError(t, err)

		task, err := tr.CreateTask(domainaggregate.Task{
			ProjectId:   proj.Id,
			Name:        "sdf",
			Description: "",
			Deadline:    time.Now().UTC(),
			Status:      domainaggregate.Finished,
		})
		assert.NoError(t, err)

		// act
		info, err := pr.GetProjectTaskInfoById(proj.Id)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, 1, info.FinishedCount)
		assert.Equal(t, 0, info.NotStartedCount)
		assert.Equal(t, 0, info.InProgressCount)

		// cleanup

		profId, err := strconv.Atoi(resProf.Professor.Id)
		assert.NoError(t, err)
		resErrChan := ar.DeleteProfessor(profId)
		resErr := <-resErrChan
		assert.NoError(t, resErr.Err)

		studId, err := strconv.Atoi(stud.Id)
		assert.NoError(t, err)
		err = sr.DeleteStudent(studId)
		assert.NoError(t, err)

		taskId, err := strconv.Atoi(task.Id)
		assert.NoError(t, err)
		err = tr.DeleteTask(taskId)
		assert.NoError(t, err)
	})
}
func TestProjectRepo_AssignDriveFolder(t *testing.T) {
	db := connectDB()

	t.Run("fail, not found", func(t *testing.T) {
		// arrange
		pr := InitProjectRepository(*db)

		// act
		err := pr.AssignDriveFolder(usecasemodels.DriveProject{
			Project: domainaggregate.Project{
				Id: "123",
			},
			DriveFolder: usecasemodels.DriveFolder{
				Id:   "123",
				Link: time.Now().Format(time.RFC3339),
			},
		})

		// assert
		assert.ErrorIs(t, err, usecasemodels.ErrProjectNotFound)
	})

	t.Run("ok", func(t *testing.T) {
		// arrange
		pr := InitProjectRepository(*db)
		ar := accountrepository.InitAccountRepository(*db)
		sr := studentrepository.InitStudentRepository(*db)

		stud, err := sr.CreateStudent(domainaggregate.Student{
			Person: domainaggregate.Person{
				Name:       "dsf",
				Surname:    "sdf",
				Middlename: "sdf",
			},
			Cource: 1,
		})
		assert.NoError(t, err)

		resProfChan := ar.AddProfessor(domainaggregate.Professor{
			Person: domainaggregate.Person{
				Name:       "dsf",
				Surname:    "sdf",
				Middlename: "sdf",
			},
			ScienceDegree: time.Now().Format(time.RFC3339),
		})
		resProf := <-resProfChan
		assert.NoError(t, resProf.Err)
		proj, err := pr.CreateProject(domainaggregate.Project{
			Theme:        time.Now().Format(time.RFC3339),
			SupervisorId: resProf.Professor.Id,
			StudentId:    stud.Id,
			Year:         2023,
			Stage:        domainaggregate.Analysis,
			Status:       domainaggregate.ProjectInProgress,
		})
		assert.NoError(t, err)
		folder := usecasemodels.DriveFolder{
			Id:   "1234",
			Link: time.Now().Format(time.RFC3339),
		}

		// act
		err = pr.AssignDriveFolder(usecasemodels.DriveProject{
			Project:     proj,
			DriveFolder: folder,
		})

		// assert
		assert.NoError(t, err)
		foundFolder, err := pr.GetProjectCloudFolderId(proj.Id)
		assert.NoError(t, err)
		assert.Equal(t, folder.Id, foundFolder)
		foundLink, err := pr.GetProjectFolderLink(proj.Id)
		assert.NoError(t, err)
		assert.Equal(t, folder.Link, foundLink)

		// cleanup

		profId, err := strconv.Atoi(resProf.Professor.Id)
		assert.NoError(t, err)
		resErrChan := ar.DeleteProfessor(profId)
		resErr := <-resErrChan
		assert.NoError(t, resErr.Err)

		studId, err := strconv.Atoi(stud.Id)
		assert.NoError(t, err)
		err = sr.DeleteStudent(studId)
		assert.NoError(t, err)

		db.DB.Delete(&models.CloudFolder{Id: folder.Id})
	})

}
func TestProjectRepo_GetProjectCloudFolderId(t *testing.T) {
	db := connectDB()

	t.Run("fail, proj not found", func(t *testing.T) {
		// arrange
		pr := InitProjectRepository(*db)

		// act
		_, err := pr.GetProjectCloudFolderId("123")

		// assert
		assert.ErrorIs(t, err, usecasemodels.ErrProjectNotFound)
	})

	t.Run("fail, folder id not found", func(t *testing.T) {
		// arrange
		pr := InitProjectRepository(*db)
		ar := accountrepository.InitAccountRepository(*db)
		sr := studentrepository.InitStudentRepository(*db)

		stud, err := sr.CreateStudent(domainaggregate.Student{
			Person: domainaggregate.Person{
				Name:       "dsf",
				Surname:    "sdf",
				Middlename: "sdf",
			},
			Cource: 1,
		})
		assert.NoError(t, err)

		resProfChan := ar.AddProfessor(domainaggregate.Professor{
			Person: domainaggregate.Person{
				Name:       "dsf",
				Surname:    "sdf",
				Middlename: "sdf",
			},
			ScienceDegree: time.Now().Format(time.RFC3339),
		})
		resProf := <-resProfChan
		assert.NoError(t, resProf.Err)
		proj, err := pr.CreateProject(domainaggregate.Project{
			Theme:        time.Now().Format(time.RFC3339),
			SupervisorId: resProf.Professor.Id,
			StudentId:    stud.Id,
			Year:         2023,
			Stage:        domainaggregate.Analysis,
			Status:       domainaggregate.ProjectInProgress,
		})
		assert.NoError(t, err)

		// act
		_, err = pr.GetProjectCloudFolderId(proj.Id)

		// assert
		assert.ErrorIs(t, err, usecasemodels.ErrProjectCloudFolderNotFound)

		// cleanup

		profId, err := strconv.Atoi(resProf.Professor.Id)
		assert.NoError(t, err)
		resErrChan := ar.DeleteProfessor(profId)
		resErr := <-resErrChan
		assert.NoError(t, resErr.Err)

		studId, err := strconv.Atoi(stud.Id)
		assert.NoError(t, err)
		err = sr.DeleteStudent(studId)
		assert.NoError(t, err)
	})

	t.Run("ok", func(t *testing.T) {
		// arrange
		pr := InitProjectRepository(*db)
		ar := accountrepository.InitAccountRepository(*db)
		sr := studentrepository.InitStudentRepository(*db)

		stud, err := sr.CreateStudent(domainaggregate.Student{
			Person: domainaggregate.Person{
				Name:       "dsf",
				Surname:    "sdf",
				Middlename: "sdf",
			},
			Cource: 1,
		})
		assert.NoError(t, err)

		resProfChan := ar.AddProfessor(domainaggregate.Professor{
			Person: domainaggregate.Person{
				Name:       "dsf",
				Surname:    "sdf",
				Middlename: "sdf",
			},
			ScienceDegree: time.Now().Format(time.RFC3339),
		})
		resProf := <-resProfChan
		assert.NoError(t, resProf.Err)
		proj, err := pr.CreateProject(domainaggregate.Project{
			Theme:        time.Now().Format(time.RFC3339),
			SupervisorId: resProf.Professor.Id,
			StudentId:    stud.Id,
			Year:         2023,
			Stage:        domainaggregate.Analysis,
			Status:       domainaggregate.ProjectInProgress,
		})
		assert.NoError(t, err)

		folder := usecasemodels.DriveFolder{
			Id:   "1234",
			Link: time.Now().Format(time.RFC3339),
		}
		err = pr.AssignDriveFolder(usecasemodels.DriveProject{
			Project:     proj,
			DriveFolder: folder,
		})
		assert.NoError(t, err)

		// act
		foundFolder, err := pr.GetProjectCloudFolderId(proj.Id)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, folder.Id, foundFolder)

		// cleanup

		profId, err := strconv.Atoi(resProf.Professor.Id)
		assert.NoError(t, err)
		resErrChan := ar.DeleteProfessor(profId)
		resErr := <-resErrChan
		assert.NoError(t, resErr.Err)

		studId, err := strconv.Atoi(stud.Id)
		assert.NoError(t, err)
		err = sr.DeleteStudent(studId)
		assert.NoError(t, err)

		db.DB.Delete(&models.CloudFolder{Id: folder.Id})
	})

}

func TestProjectRepo_GetProjectFolderLink(t *testing.T) {
	db := connectDB()

	t.Run("fail, proj not found", func(t *testing.T) {
		// arrange
		pr := InitProjectRepository(*db)

		// act
		_, err := pr.GetProjectFolderLink("123")

		// assert
		assert.ErrorIs(t, err, usecasemodels.ErrProjectNotFound)
	})

	t.Run("fail, folder id not found", func(t *testing.T) {
		// arrange
		pr := InitProjectRepository(*db)
		ar := accountrepository.InitAccountRepository(*db)
		sr := studentrepository.InitStudentRepository(*db)

		stud, err := sr.CreateStudent(domainaggregate.Student{
			Person: domainaggregate.Person{
				Name:       "dsf",
				Surname:    "sdf",
				Middlename: "sdf",
			},
			Cource: 1,
		})
		assert.NoError(t, err)

		resProfChan := ar.AddProfessor(domainaggregate.Professor{
			Person: domainaggregate.Person{
				Name:       "dsf",
				Surname:    "sdf",
				Middlename: "sdf",
			},
			ScienceDegree: time.Now().Format(time.RFC3339),
		})
		resProf := <-resProfChan
		assert.NoError(t, err)
		proj, err := pr.CreateProject(domainaggregate.Project{
			Theme:        time.Now().Format(time.RFC3339),
			SupervisorId: resProf.Professor.Id,
			StudentId:    stud.Id,
			Year:         2023,
			Stage:        domainaggregate.Analysis,
			Status:       domainaggregate.ProjectInProgress,
		})
		assert.NoError(t, err)

		// act
		_, err = pr.GetProjectFolderLink(proj.Id)

		// assert
		assert.ErrorIs(t, err, usecasemodels.ErrProjectCloudFolderNotFound)

		// cleanup

		profId, err := strconv.Atoi(resProf.Professor.Id)
		assert.NoError(t, err)
		resErrChan := ar.DeleteProfessor(profId)
		resErr := <-resErrChan
		assert.NoError(t, resErr.Err)

		studId, err := strconv.Atoi(stud.Id)
		assert.NoError(t, err)
		err = sr.DeleteStudent(studId)
		assert.NoError(t, err)
	})

	t.Run("ok", func(t *testing.T) {
		// arrange
		pr := InitProjectRepository(*db)
		ar := accountrepository.InitAccountRepository(*db)
		sr := studentrepository.InitStudentRepository(*db)

		stud, err := sr.CreateStudent(domainaggregate.Student{
			Person: domainaggregate.Person{
				Name:       "dsf",
				Surname:    "sdf",
				Middlename: "sdf",
			},
			Cource: 1,
		})
		assert.NoError(t, err)

		resProfChan := ar.AddProfessor(domainaggregate.Professor{
			Person: domainaggregate.Person{
				Name:       "dsf",
				Surname:    "sdf",
				Middlename: "sdf",
			},
			ScienceDegree: time.Now().Format(time.RFC3339),
		})
		resProf := <-resProfChan
		assert.NoError(t, err)
		proj, err := pr.CreateProject(domainaggregate.Project{
			Theme:        time.Now().Format(time.RFC3339),
			SupervisorId: resProf.Professor.Id,
			StudentId:    stud.Id,
			Year:         2023,
			Stage:        domainaggregate.Analysis,
			Status:       domainaggregate.ProjectInProgress,
		})
		assert.NoError(t, err)

		folder := usecasemodels.DriveFolder{
			Id:   "1234",
			Link: time.Now().Format(time.RFC3339),
		}
		err = pr.AssignDriveFolder(usecasemodels.DriveProject{
			Project:     proj,
			DriveFolder: folder,
		})
		assert.NoError(t, err)

		// act
		foundLink, err := pr.GetProjectFolderLink(proj.Id)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, folder.Link, foundLink)

		// cleanup

		profId, err := strconv.Atoi(resProf.Professor.Id)
		assert.NoError(t, err)
		resErrChan := ar.DeleteProfessor(profId)
		resErr := <-resErrChan
		assert.NoError(t, resErr.Err)

		studId, err := strconv.Atoi(stud.Id)
		assert.NoError(t, err)
		err = sr.DeleteStudent(studId)
		assert.NoError(t, err)

		db.DB.Delete(&models.CloudFolder{Id: folder.Id})
	})

}

func TestProjectRepo_GetStudentCurrentProject(t *testing.T) {
	db := connectDB()

	t.Run("fail, no current proj", func(t *testing.T) {
		// arrange
		pr := InitProjectRepository(*db)
		sr := studentrepository.InitStudentRepository(*db)

		stud, err := sr.CreateStudent(domainaggregate.Student{
			Person: domainaggregate.Person{
				Name:       "dsf",
				Surname:    "sdf",
				Middlename: "sdf",
			},
			Cource: 1,
		})
		assert.NoError(t, err)

		// act
		_, err = pr.GetStudentCurrentProject(stud.Id)

		// assert
		assert.ErrorIs(t, err, usecasemodels.ErrStudentHasNoCurrentProject)

		// cleanup

		studId, err := strconv.Atoi(stud.Id)
		assert.NoError(t, err)
		err = sr.DeleteStudent(studId)
		assert.NoError(t, err)
	})

	t.Run("ok", func(t *testing.T) {
		// arrange
		pr := InitProjectRepository(*db)
		ar := accountrepository.InitAccountRepository(*db)
		sr := studentrepository.InitStudentRepository(*db)

		stud, err := sr.CreateStudent(domainaggregate.Student{
			Person: domainaggregate.Person{
				Name:       "dsf",
				Surname:    "sdf",
				Middlename: "sdf",
			},
			Cource: 1,
		})
		assert.NoError(t, err)

		resProfChan := ar.AddProfessor(domainaggregate.Professor{
			Person: domainaggregate.Person{
				Name:       "dsf",
				Surname:    "sdf",
				Middlename: "sdf",
			},
			ScienceDegree: time.Now().Format(time.RFC3339),
		})
		resProf := <-resProfChan
		assert.NoError(t, resProf.Err)

		proj, err := pr.CreateProject(domainaggregate.Project{
			Theme:        time.Now().Format(time.RFC3339),
			SupervisorId: resProf.Professor.Id,
			StudentId:    stud.Id,
			Year:         2023,
			Stage:        domainaggregate.Analysis,
			Status:       domainaggregate.ProjectInProgress,
		})
		assert.NoError(t, err)

		// act
		foundProj, err := pr.GetStudentCurrentProject(stud.Id)
		assert.Equal(t, proj.Id, foundProj.Id)
		assert.Equal(t, proj.Stage, foundProj.Stage)
		assert.Equal(t, proj.Status, foundProj.Status)
		assert.Equal(t, proj.StudentId, foundProj.StudentId)
		assert.Equal(t, proj.SupervisorId, foundProj.SupervisorId)
		assert.Equal(t, proj.Theme, foundProj.Theme)
		assert.Equal(t, proj.Year, foundProj.Year)

		// assert
		assert.NoError(t, err)

		// cleanup

		profId, err := strconv.Atoi(resProf.Professor.Id)
		assert.NoError(t, err)
		resErrChan := ar.DeleteProfessor(profId)
		resErr := <-resErrChan
		assert.NoError(t, resErr.Err)

		studId, err := strconv.Atoi(stud.Id)
		assert.NoError(t, err)
		err = sr.DeleteStudent(studId)
		assert.NoError(t, err)
	})

}

func TestProjectRepo_GetProjectRepository(t *testing.T) {
	db := connectDB()

	t.Run("fail, no project", func(t *testing.T) {
		// arrange
		pr := InitProjectRepository(*db)

		// act
		_, err := pr.GetProjectRepository("123")

		// assert
		assert.ErrorIs(t, err, usecasemodels.ErrProjectNotFound)
	})

	t.Run("fail, no proj repo", func(t *testing.T) {
		// arrange
		pr := InitProjectRepository(*db)
		ar := accountrepository.InitAccountRepository(*db)
		sr := studentrepository.InitStudentRepository(*db)

		stud, err := sr.CreateStudent(domainaggregate.Student{
			Person: domainaggregate.Person{
				Name:       "dsf",
				Surname:    "sdf",
				Middlename: "sdf",
			},
			Cource: 1,
		})
		assert.NoError(t, err)

		resProfChan := ar.AddProfessor(domainaggregate.Professor{
			Person: domainaggregate.Person{
				Name:       "dsf",
				Surname:    "sdf",
				Middlename: "sdf",
			},
			ScienceDegree: time.Now().Format(time.RFC3339),
		})
		resProf := <-resProfChan
		assert.NoError(t, resProf.Err)

		proj, err := pr.CreateProject(domainaggregate.Project{
			Theme:        time.Now().Format(time.RFC3339),
			SupervisorId: resProf.Professor.Id,
			StudentId:    stud.Id,
			Year:         2023,
			Stage:        domainaggregate.Analysis,
			Status:       domainaggregate.ProjectInProgress,
		})
		assert.NoError(t, err)

		// act
		_, err = pr.GetProjectRepository(proj.Id)

		// assert
		assert.ErrorIs(t, err, usecasemodels.ErrProjectRepoNotFound)

		// cleanup

		profId, err := strconv.Atoi(resProf.Professor.Id)
		assert.NoError(t, err)
		resErrChan := ar.DeleteProfessor(profId)
		resErr := <-resErrChan
		assert.NoError(t, resErr.Err)

		studId, err := strconv.Atoi(stud.Id)
		assert.NoError(t, err)
		err = sr.DeleteStudent(studId)
		assert.NoError(t, err)
	})

	t.Run("ok", func(t *testing.T) {
		// arrange
		pr := InitProjectRepository(*db)
		ar := accountrepository.InitAccountRepository(*db)
		sr := studentrepository.InitStudentRepository(*db)

		stud, err := sr.CreateStudent(domainaggregate.Student{
			Person: domainaggregate.Person{
				Name:       "dsf",
				Surname:    "sdf",
				Middlename: "sdf",
			},
			Cource: 1,
		})
		assert.NoError(t, err)

		resProfChan := ar.AddProfessor(domainaggregate.Professor{
			Person: domainaggregate.Person{
				Name:       "dsf",
				Surname:    "sdf",
				Middlename: "sdf",
			},
			ScienceDegree: time.Now().Format(time.RFC3339),
		})
		resProf := <-resProfChan
		assert.NoError(t, resProf.Err)

		repo := usecasemodels.Repository{
			RepoId:    time.Now().Format(time.RFC3339),
			OwnerName: "123445",
			RepoType:  int(usecasemodels.GitHub),
		}
		proj, err := pr.CreateProjectWithRepository(domainaggregate.Project{
			Theme:        time.Now().Format(time.RFC3339),
			SupervisorId: resProf.Professor.Id,
			StudentId:    stud.Id,
			Year:         2023,
			Stage:        domainaggregate.Analysis,
			Status:       domainaggregate.ProjectInProgress,
		}, repo)
		assert.NoError(t, err)

		// act
		foundRepo, err := pr.GetProjectRepository(proj.Id)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, repo.RepoId, foundRepo.RepoId)
		assert.Equal(t, repo.OwnerName, foundRepo.OwnerName)
		assert.Equal(t, repo.RepoType, foundRepo.RepoType)

		// cleanup

		profId, err := strconv.Atoi(resProf.Professor.Id)
		assert.NoError(t, err)

		resErrChan := ar.DeleteProfessor(profId)
		resErr := <-resErrChan
		assert.NoError(t, resErr.Err)

		studId, err := strconv.Atoi(stud.Id)
		assert.NoError(t, err)
		err = sr.DeleteStudent(studId)
		assert.NoError(t, err)

		db.DB.Where("id > 0").Delete(&models.Repository{})
	})

}
