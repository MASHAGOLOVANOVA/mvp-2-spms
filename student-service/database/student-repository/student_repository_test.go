package studentrepository

import (
	"fmt"
	"mvp-2-spms/database"
	domainaggregate "mvp-2-spms/domain-aggregate"
	"strconv"
	"testing"

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

func TestStudentRepo_GetStudentById(t *testing.T) {
	db := connectDB()

	t.Run("fail, student not found", func(t *testing.T) {
		// arrange
		sr := InitStudentRepository(*db)

		// act
		_, err := sr.GetStudentById("-1")

		// assert
		assert.Error(t, err)
	})

	t.Run("ok", func(t *testing.T) {
		// arrange
		sr := InitStudentRepository(*db)

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
		foundStud, err := sr.GetStudentById(stud.Id)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, stud.Id, foundStud.Id)
		assert.Equal(t, stud.Cource, foundStud.Cource)
		assert.Equal(t, stud.Middlename, foundStud.Middlename)
		assert.Equal(t, stud.Name, foundStud.Name)
		assert.Equal(t, stud.Surname, foundStud.Surname)

		// cleanup

		studId, err := strconv.Atoi(stud.Id)
		assert.NoError(t, err)
		err = sr.DeleteStudent(studId)
		assert.NoError(t, err)
	})
}

func TestStudentRepo_CreateStudent(t *testing.T) {
	db := connectDB()

	t.Run("fail, bad ed programme id", func(t *testing.T) {
		// arrange
		sr := InitStudentRepository(*db)

		// act
		_, err := sr.CreateStudent(domainaggregate.Student{
			Person: domainaggregate.Person{
				Name:       "dsf",
				Surname:    "sdf",
				Middlename: "sdf",
			},
			Cource:                 1,
			EducationalProgrammeId: "12",
		})

		// assert
		assert.Error(t, err)
	})

	t.Run("ok", func(t *testing.T) {
		// arrange
		sr := InitStudentRepository(*db)

		// act
		stud, err := sr.CreateStudent(domainaggregate.Student{
			Person: domainaggregate.Person{
				Name:       "dsf",
				Surname:    "sdf",
				Middlename: "sdf",
			},
			Cource: 1,
		})

		// assert
		assert.NoError(t, err)
		foundStud, err := sr.GetStudentById(stud.Id)
		assert.NoError(t, err)
		assert.Equal(t, stud.Id, foundStud.Id)
		assert.Equal(t, stud.Cource, foundStud.Cource)
		assert.Equal(t, stud.Middlename, foundStud.Middlename)
		assert.Equal(t, stud.Name, foundStud.Name)
		assert.Equal(t, stud.Surname, foundStud.Surname)

		// cleanup

		studId, err := strconv.Atoi(stud.Id)
		assert.NoError(t, err)
		err = sr.DeleteStudent(studId)
		assert.NoError(t, err)
	})
}

func TestStudentRepo_GetStudents(t *testing.T) {
	db := connectDB()

	t.Run("ok, list", func(t *testing.T) {
		// arrange
		sr := InitStudentRepository(*db)

		studs := []domainaggregate.Student{}
		for i := 0; i < 10; i++ {
			stud, err := sr.CreateStudent(domainaggregate.Student{
				Person: domainaggregate.Person{
					Name:       fmt.Sprint(i),
					Surname:    "sdf",
					Middlename: "sdf",
				},
				Cource: 1,
			})
			assert.NoError(t, err)
			studs = append(studs, stud)
		}

		// act
		foundStuds, err := sr.GetStudents()

		// assert
		assert.NoError(t, err)
		assert.NotEqual(t, 0, len(foundStuds))

		// cleanup
		for _, stud := range studs {
			studId, err := strconv.Atoi(stud.Id)
			assert.NoError(t, err)
			err = sr.DeleteStudent(studId)
			assert.NoError(t, err)
		}
	})
}
