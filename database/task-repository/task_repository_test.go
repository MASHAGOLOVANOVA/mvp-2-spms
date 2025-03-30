package taskrepository

import (
	"mvp-2-spms/database"
	accountrepository "mvp-2-spms/database/account-repository"
	"mvp-2-spms/database/models"
	projectrepository "mvp-2-spms/database/project-repository"
	studentrepository "mvp-2-spms/database/student-repository"
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

func createProj(db *database.Database) (domainaggregate.Project, error) {
	pr := projectrepository.InitProjectRepository(*db)
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
	if err != nil {
		return domainaggregate.Project{}, err
	}

	resprofChan := ar.AddProfessor(domainaggregate.Professor{
		Person: domainaggregate.Person{
			Name:       "dsf",
			Surname:    "sdf",
			Middlename: "sdf",
		},
		ScienceDegree: time.Now().Format(time.RFC3339),
	})
	resProf := <-resprofChan
	if resProf.Err != nil {
		return domainaggregate.Project{}, err
	}

	proj, err := pr.CreateProject(domainaggregate.Project{
		Theme:        time.Now().Format(time.RFC3339),
		SupervisorId: resProf.Professor.Id,
		StudentId:    stud.Id,
		Year:         2023,
		Stage:        domainaggregate.Analysis,
		Status:       domainaggregate.ProjectInProgress,
	})
	if err != nil {
		return domainaggregate.Project{}, err
	}

	return proj, nil
}

func deleteProj(db *database.Database, proj domainaggregate.Project) error {
	ar := accountrepository.InitAccountRepository(*db)
	sr := studentrepository.InitStudentRepository(*db)

	profId, err := strconv.Atoi(proj.SupervisorId)
	if err != nil {
		return err
	}
	resErrChan := ar.DeleteProfessor(profId)
	resErr := <-resErrChan
	if resErr.Err != nil {
		return resErr.Err
	}

	studId, err := strconv.Atoi(proj.StudentId)
	if err != nil {
		return err
	}
	err = sr.DeleteStudent(studId)
	if err != nil {
		return err
	}

	return nil
}

func TestStudentRepo_CreateTask(t *testing.T) {
	db := connectDB()

	t.Run("fail, bad project id", func(t *testing.T) {
		// arrange
		tr := InitTaskRepository(*db)

		// act
		_, err := tr.CreateTask(domainaggregate.Task{
			ProjectId:   "12",
			Name:        "af",
			Deadline:    time.Now().UTC().Round(time.Second),
			Description: "eqw",
			Status:      domainaggregate.NotStarted,
		})

		// assert
		assert.Error(t, err)
	})

	t.Run("ok", func(t *testing.T) {
		// arrange
		tr := InitTaskRepository(*db)
		proj, err := createProj(db)
		assert.NoError(t, err)

		// act
		task, err := tr.CreateTask(domainaggregate.Task{
			ProjectId:   proj.Id,
			Name:        "af",
			Deadline:    time.Now().UTC().Round(time.Second),
			Description: "eqw",
			Status:      domainaggregate.NotStarted,
		})

		// assert
		assert.NoError(t, err)
		foundTask, err := tr.GetTaskById(task.Id)
		assert.NoError(t, err)
		assert.Equal(t, task.Id, foundTask.Id)
		assert.Equal(t, task.Deadline, foundTask.Deadline)
		assert.Equal(t, task.Description, foundTask.Description)
		assert.Equal(t, task.Name, foundTask.Name)
		assert.Equal(t, task.ProjectId, foundTask.ProjectId)
		assert.Equal(t, task.Status, foundTask.Status)

		// cleanup

		taskId, err := strconv.Atoi(task.Id)
		assert.NoError(t, err)
		err = tr.DeleteTask(taskId)
		assert.NoError(t, err)
		err = deleteProj(db, proj)
		assert.NoError(t, err)
	})
}

func TestStudentRepo_GetProjectTasksWithCloud(t *testing.T) {
	db := connectDB()

	t.Run("ok, empty", func(t *testing.T) {
		// arrange
		tr := InitTaskRepository(*db)
		proj, err := createProj(db)
		assert.NoError(t, err)
		// act
		foundTasks, err := tr.GetProjectTasksWithCloud(proj.Id)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, 0, len(foundTasks))
		err = deleteProj(db, proj)
		assert.NoError(t, err)
	})

}

func TestStudentRepo_GetProjectTasks(t *testing.T) {
	db := connectDB()

	t.Run("ok, empty", func(t *testing.T) {
		// arrange
		tr := InitTaskRepository(*db)
		proj, err := createProj(db)
		assert.NoError(t, err)
		// act
		foundTasks, err := tr.GetProjectTasks(proj.Id)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, 0, len(foundTasks))

		// cleanup
		err = deleteProj(db, proj)
		assert.NoError(t, err)
	})

}

func TestStudentRepo_AssignDriveTask(t *testing.T) {
	db := connectDB()

	t.Run("fail, task not found", func(t *testing.T) {
		// arrange
		tr := InitTaskRepository(*db)

		// act
		err := tr.AssignDriveTask(usecasemodels.DriveTask{
			Task: domainaggregate.Task{
				Id: "123",
			},
			DriveFolder: usecasemodels.DriveFolder{
				Id:   "23f",
				Link: time.Now().Format(time.RFC3339),
			},
		})

		// assert
		assert.ErrorIs(t, err, usecasemodels.ErrTaskNotFound)
	})

	t.Run("ok", func(t *testing.T) {
		// arrange
		tr := InitTaskRepository(*db)
		proj, err := createProj(db)
		assert.NoError(t, err)

		task, err := tr.CreateTask(domainaggregate.Task{
			ProjectId:   proj.Id,
			Name:        "af",
			Deadline:    time.Now().UTC().Round(time.Second),
			Description: "eqw",
			Status:      domainaggregate.NotStarted,
		})
		assert.NoError(t, err)
		folder := usecasemodels.DriveFolder{
			Id:   "23f",
			Link: time.Now().Format(time.RFC3339),
		}
		// act
		err = tr.AssignDriveTask(usecasemodels.DriveTask{
			Task:        task,
			DriveFolder: folder,
		})

		// assert
		assert.NoError(t, err)

		// cleanup

		taskId, err := strconv.Atoi(task.Id)
		assert.NoError(t, err)
		err = tr.DeleteTask(taskId)
		assert.NoError(t, err)
		err = deleteProj(db, proj)
		assert.NoError(t, err)

		db.DB.Where("id = ?", folder.Id).Delete(&models.CloudFolder{})
	})
}
