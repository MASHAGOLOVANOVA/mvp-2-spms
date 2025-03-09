package universityrepository

import (
	"fmt"
	"mvp-2-spms/database"
	"mvp-2-spms/database/models"
	domainaggregate "mvp-2-spms/domain-aggregate"
	usecasemodels "mvp-2-spms/services/models"
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

func TestUniversityRepo_GetEducationalProgrammeById(t *testing.T) {
	db := connectDB()

	t.Run("fail, not found", func(t *testing.T) {
		// arrange
		ur := InitUniversityRepository(*db)

		// act
		_, err := ur.GetEducationalProgrammeById("124")

		// assert
		assert.ErrorIs(t, err, usecasemodels.ErrEdProgrammmeNotFound)
	})

	t.Run("ok", func(t *testing.T) {
		// arrange
		ur := InitUniversityRepository(*db)
		edProgDb := models.EducationalProgramme{
			Name:      time.Now().Format(time.RFC3339),
			FacultyId: 1,
			EdLevel:   0,
		}
		db.DB.Create(&edProgDb)
		edProg := edProgDb.MapToEntity()

		// act
		foundEdProg, err := ur.GetEducationalProgrammeById(fmt.Sprint(edProgDb.Id))

		// assert
		assert.NoError(t, err)
		assert.Equal(t, edProg.Id, foundEdProg.Id)
		assert.Equal(t, edProg.FacultyId, foundEdProg.FacultyId)
		assert.Equal(t, edProg.EducationalLevel, foundEdProg.EducationalLevel)
		assert.Equal(t, edProg.Name, foundEdProg.Name)

		// cleanup
		db.DB.Delete(&edProgDb)
	})
}

func TestUniversityRepo_GetUniversityById(t *testing.T) {
	db := connectDB()

	t.Run("fail, not found", func(t *testing.T) {
		// arrange
		ur := InitUniversityRepository(*db)

		// act
		_, err := ur.GetUniversityById("124")

		// assert
		assert.ErrorIs(t, err, usecasemodels.ErrUniNoFound)
	})

	t.Run("ok", func(t *testing.T) {
		// arrange
		ur := InitUniversityRepository(*db)

		// act
		_, err := ur.GetUniversityById("1")

		// assert
		assert.NoError(t, err)
	})
}

func TestUniversityRepo_GetUniversityEducationalProgrammes(t *testing.T) {
	db := connectDB()

	t.Run("ok, empty", func(t *testing.T) {
		// arrange
		ur := InitUniversityRepository(*db)

		// act
		_, err := ur.GetUniversityEducationalProgrammes("1")

		// assert
		assert.NoError(t, err)
	})

	t.Run("ok, list", func(t *testing.T) {
		// arrange
		ur := InitUniversityRepository(*db)
		edprogs := []domainaggregate.EducationalProgramme{}
		edprogsdb := []models.EducationalProgramme{}
		for i := 0; i < 10; i++ {
			edProgDb := models.EducationalProgramme{
				Name:      time.Now().Format(time.RFC3339),
				FacultyId: 1,
				EdLevel:   0,
			}
			db.DB.Create(&edProgDb)
			edProg := edProgDb.MapToEntity()
			edprogs = append(edprogs, edProg)
			edprogsdb = append(edprogsdb, edProgDb)
		}
		// act
		foundedProgs, err := ur.GetUniversityEducationalProgrammes("1")

		// assert
		assert.NoError(t, err)
		assert.Equal(t, len(edprogs), len(foundedProgs))
		for i := range edprogs {
			assert.Equal(t, edprogs[i].EducationalLevel, foundedProgs[i].EducationalLevel)
			assert.Equal(t, edprogs[i].FacultyId, foundedProgs[i].FacultyId)
			assert.Equal(t, edprogs[i].Id, foundedProgs[i].Id)
			assert.Equal(t, edprogs[i].Name, foundedProgs[i].Name)
		}

		// cleanup

		for _, p := range edprogsdb {
			db.DB.Delete(&p)
		}
	})
}
