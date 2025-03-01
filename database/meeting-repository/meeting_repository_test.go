package meetingrepository

import (
	"log"
	"mvp-2-spms/database"
	accountrepository "mvp-2-spms/database/account-repository"
	studentrepository "mvp-2-spms/database/student-repository"
	domainaggregate "mvp-2-spms/domain-aggregate"
	"mvp-2-spms/services/models"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var dsn = "root:root@tcp(127.0.0.1:3306)/student_project_management_testing?parseTime=true"

func connectDB() *database.Database {
	gdb, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // use singular table name, table for `User` would be `user` with this option enabled
		},
	})
	db := database.InitDatabade(gdb)
	return db
}

func TestMeetingRepo_CreateMeeting(t *testing.T) {
	db := connectDB()

	t.Run("fail, bad organizer id", func(t *testing.T) {
		// arrange
		mr := InitMeetingRepository(*db)
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
		_, err = mr.CreateMeeting(domainaggregate.Meeting{
			Name:          time.Now().Format(time.RFC3339),
			Description:   "lwefl",
			OrganizerId:   "123",
			ParticipantId: stud.Id,
			Time:          time.Now().UTC(),
			IsOnline:      true,
			Status:        domainaggregate.MeetingPlanned,
		})

		// assert
		assert.Error(t, err)

		// cleanup

		studId, err := strconv.Atoi(stud.Id)
		assert.NoError(t, err)
		err = sr.DeleteStudent(studId)
		assert.NoError(t, err)
	})

	t.Run("fail, bad participant id", func(t *testing.T) {
		// arrange
		mr := InitMeetingRepository(*db)
		ar := accountrepository.InitAccountRepository(*db)

		prof, err := ar.AddProfessor(domainaggregate.Professor{
			Person: domainaggregate.Person{
				Name:       "dsf",
				Surname:    "sdf",
				Middlename: "sdf",
			},
			ScienceDegree: time.Now().Format(time.RFC3339),
		})
		assert.NoError(t, err)

		// act
		_, err = mr.CreateMeeting(domainaggregate.Meeting{
			Name:          time.Now().Format(time.RFC3339),
			Description:   "lwefl",
			OrganizerId:   prof.Id,
			ParticipantId: "123",
			Time:          time.Now().UTC(),
			IsOnline:      true,
			Status:        domainaggregate.MeetingPlanned,
		})

		// assert
		assert.Error(t, err)

		// cleanup
		profId, err := strconv.Atoi(prof.Id)
		assert.NoError(t, err)
		err = ar.DeleteProfessor(profId)
		assert.NoError(t, err)
	})

	t.Run("fail, bad project id", func(t *testing.T) {
		// arrange
		mr := InitMeetingRepository(*db)
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

		prof, err := ar.AddProfessor(domainaggregate.Professor{
			Person: domainaggregate.Person{
				Name:       "dsf",
				Surname:    "sdf",
				Middlename: "sdf",
			},
			ScienceDegree: time.Now().Format(time.RFC3339),
		})
		assert.NoError(t, err)

		// act
		_, err = mr.CreateMeeting(domainaggregate.Meeting{
			Name:          time.Now().Format(time.RFC3339),
			Description:   "lwefl",
			OrganizerId:   prof.Id,
			ParticipantId: stud.Id,
			ProjectId:     "1243",
			Time:          time.Now().UTC(),
			IsOnline:      true,
			Status:        domainaggregate.MeetingPlanned,
		})

		// assert
		assert.Error(t, err)

		// cleanup
		profId, err := strconv.Atoi(prof.Id)
		assert.NoError(t, err)
		err = ar.DeleteProfessor(profId)
		assert.NoError(t, err)

		studId, err := strconv.Atoi(stud.Id)
		assert.NoError(t, err)
		err = sr.DeleteStudent(studId)
		assert.NoError(t, err)
	})

	t.Run("ok", func(t *testing.T) {
		// arrange
		mr := InitMeetingRepository(*db)
		ar := accountrepository.InitAccountRepository(*db)
		sr := studentrepository.InitStudentRepository(*db)

		prof, err := ar.AddProfessor(domainaggregate.Professor{
			Person: domainaggregate.Person{
				Name:       "dsf",
				Surname:    "sdf",
				Middlename: "sdf",
			},
			ScienceDegree: time.Now().Format(time.RFC3339),
		})
		assert.NoError(t, err)

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
		meet, err := mr.CreateMeeting(domainaggregate.Meeting{
			Name:          time.Now().Format(time.RFC3339),
			Description:   "lwefl",
			OrganizerId:   prof.Id,
			ParticipantId: stud.Id,
			Time:          time.Now().UTC().Round(time.Second),
			IsOnline:      true,
			Status:        domainaggregate.MeetingPlanned,
		})

		// assert
		assert.NoError(t, err)
		foundMeet, err := mr.GetMeetingById(meet.Id)
		assert.NoError(t, err)
		assert.Equal(t, meet.Id, foundMeet.Id)
		assert.Equal(t, meet.Description, foundMeet.Description)
		assert.Equal(t, meet.Name, foundMeet.Name)
		assert.Equal(t, meet.IsOnline, foundMeet.IsOnline)
		assert.Equal(t, meet.OrganizerId, foundMeet.OrganizerId)
		assert.Equal(t, meet.ParticipantId, foundMeet.ParticipantId)
		assert.Equal(t, meet.Status, foundMeet.Status)
		assert.Equal(t, meet.Time, foundMeet.Time)

		// cleanup

		profId, err := strconv.Atoi(prof.Id)
		assert.NoError(t, err)
		err = ar.DeleteProfessor(profId)
		assert.NoError(t, err)

		studId, err := strconv.Atoi(stud.Id)
		assert.NoError(t, err)
		err = sr.DeleteStudent(studId)
		assert.NoError(t, err)
	})
}

func TestMeetingRepo_AssignPlannerMeeting(t *testing.T) {
	db := connectDB()

	t.Run("fail, meeting not found", func(t *testing.T) {
		// arrange
		mr := InitMeetingRepository(*db)

		// act
		err := mr.AssignPlannerMeeting(models.PlannerMeeting{
			Meeting:          domainaggregate.Meeting{Id: "123"},
			MeetingPlannerId: time.Now().Format(time.RFC3339),
		})

		// assert
		assert.ErrorIs(t, err, models.ErrMeetingNotFound)
	})

	t.Run("ok", func(t *testing.T) {
		// arrange
		mr := InitMeetingRepository(*db)
		ar := accountrepository.InitAccountRepository(*db)
		sr := studentrepository.InitStudentRepository(*db)

		prof, err := ar.AddProfessor(domainaggregate.Professor{
			Person: domainaggregate.Person{
				Name:       "dsf",
				Surname:    "sdf",
				Middlename: "sdf",
			},
			ScienceDegree: time.Now().Format(time.RFC3339),
		})
		assert.NoError(t, err)

		stud, err := sr.CreateStudent(domainaggregate.Student{
			Person: domainaggregate.Person{
				Name:       "dsf",
				Surname:    "sdf",
				Middlename: "sdf",
			},
			Cource: 1,
		})

		assert.NoError(t, err)

		meet, err := mr.CreateMeeting(domainaggregate.Meeting{
			Name:          time.Now().Format(time.RFC3339),
			Description:   "lwefl",
			OrganizerId:   prof.Id,
			ParticipantId: stud.Id,
			Time:          time.Now().UTC(),
			IsOnline:      true,
			Status:        domainaggregate.MeetingPlanned,
		})
		assert.NoError(t, err)

		plannerId := time.Now().Format(time.RFC3339)

		// act
		err = mr.AssignPlannerMeeting(models.PlannerMeeting{
			Meeting:          meet,
			MeetingPlannerId: plannerId,
		})

		// assert
		assert.NoError(t, err)
		foundPlId, err := mr.GetMeetingPlannerId(meet.Id)
		assert.NoError(t, err)
		assert.Equal(t, foundPlId, plannerId)

		// cleanup

		profId, err := strconv.Atoi(prof.Id)
		assert.NoError(t, err)
		err = ar.DeleteProfessor(profId)
		assert.NoError(t, err)

		studId, err := strconv.Atoi(stud.Id)
		assert.NoError(t, err)
		err = sr.DeleteStudent(studId)
		assert.NoError(t, err)
	})
}

func TestMeetingRepo_GetMeetingById(t *testing.T) {
	db := connectDB()

	t.Run("fail, meeting not found", func(t *testing.T) {
		// arrange
		mr := InitMeetingRepository(*db)

		// act
		_, err := mr.GetMeetingById("123")

		// assert
		assert.ErrorIs(t, err, models.ErrMeetingNotFound)
	})

	t.Run("ok", func(t *testing.T) {
		// arrange
		mr := InitMeetingRepository(*db)
		ar := accountrepository.InitAccountRepository(*db)
		sr := studentrepository.InitStudentRepository(*db)

		prof, err := ar.AddProfessor(domainaggregate.Professor{
			Person: domainaggregate.Person{
				Name:       "dsf",
				Surname:    "sdf",
				Middlename: "sdf",
			},
			ScienceDegree: time.Now().Format(time.RFC3339),
		})
		assert.NoError(t, err)

		stud, err := sr.CreateStudent(domainaggregate.Student{
			Person: domainaggregate.Person{
				Name:       "dsf",
				Surname:    "sdf",
				Middlename: "sdf",
			},
			Cource: 1,
		})

		assert.NoError(t, err)

		meet, err := mr.CreateMeeting(domainaggregate.Meeting{
			Name:          time.Now().Format(time.RFC3339),
			Description:   "lwefl",
			OrganizerId:   prof.Id,
			ParticipantId: stud.Id,
			Time:          time.Now().UTC().Round(time.Second),
			IsOnline:      true,
			Status:        domainaggregate.MeetingPlanned,
		})
		assert.NoError(t, err)

		// act
		foundMeet, err := mr.GetMeetingById(meet.Id)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, meet.Id, foundMeet.Id)
		assert.Equal(t, meet.Description, foundMeet.Description)
		assert.Equal(t, meet.Name, foundMeet.Name)
		assert.Equal(t, meet.IsOnline, foundMeet.IsOnline)
		assert.Equal(t, meet.OrganizerId, foundMeet.OrganizerId)
		assert.Equal(t, meet.ParticipantId, foundMeet.ParticipantId)
		assert.Equal(t, meet.Status, foundMeet.Status)
		assert.Equal(t, meet.Time, foundMeet.Time)

		// cleanup

		profId, err := strconv.Atoi(prof.Id)
		assert.NoError(t, err)
		err = ar.DeleteProfessor(profId)
		assert.NoError(t, err)

		studId, err := strconv.Atoi(stud.Id)
		assert.NoError(t, err)
		err = sr.DeleteStudent(studId)
		assert.NoError(t, err)
	})
}

func TestMeetingRepo_GetMeetingPlannerId(t *testing.T) {
	db := connectDB()

	t.Run("fail, get non existent meeting", func(t *testing.T) {
		// arrange
		mr := InitMeetingRepository(*db)

		// act
		_, err := mr.GetMeetingPlannerId("3")

		// assert
		assert.ErrorIs(t, err, models.ErrMeetingNotFound)
	})

	t.Run("ok", func(t *testing.T) {
		// arrange
		mr := InitMeetingRepository(*db)
		ar := accountrepository.InitAccountRepository(*db)
		sr := studentrepository.InitStudentRepository(*db)

		prof, err := ar.AddProfessor(domainaggregate.Professor{
			Person: domainaggregate.Person{
				Name:       "dsf",
				Surname:    "sdf",
				Middlename: "sdf",
			},
			ScienceDegree: time.Now().Format(time.RFC3339),
		})
		assert.NoError(t, err)

		stud, err := sr.CreateStudent(domainaggregate.Student{
			Person: domainaggregate.Person{
				Name:       "dsf",
				Surname:    "sdf",
				Middlename: "sdf",
			},
			Cource: 1,
		})

		assert.NoError(t, err)

		meet, err := mr.CreateMeeting(domainaggregate.Meeting{
			Name:          time.Now().Format(time.RFC3339),
			Description:   "lwefl",
			OrganizerId:   prof.Id,
			ParticipantId: stud.Id,
			Time:          time.Now(),
			IsOnline:      true,
			Status:        domainaggregate.MeetingPlanned,
		})
		assert.NoError(t, err)

		plannerId := time.Now().Format(time.RFC3339)
		err = mr.AssignPlannerMeeting(models.PlannerMeeting{
			Meeting:          meet,
			MeetingPlannerId: plannerId,
		})
		assert.NoError(t, err)

		// act
		foundPlId, err := mr.GetMeetingPlannerId(meet.Id)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, foundPlId, plannerId)

		// cleanup

		profId, err := strconv.Atoi(prof.Id)
		assert.NoError(t, err)
		err = ar.DeleteProfessor(profId)
		assert.NoError(t, err)

		studId, err := strconv.Atoi(stud.Id)
		assert.NoError(t, err)
		err = sr.DeleteStudent(studId)
		assert.NoError(t, err)
	})

}

func TestMeetingRepo_GetProfessorMeetings(t *testing.T) {
	db := connectDB()

	t.Run("ok, empty", func(t *testing.T) {
		// arrange
		mr := InitMeetingRepository(*db)

		// act
		meets, err := mr.GetProfessorMeetings("3", time.Now(), time.Time{})

		// assert
		assert.NoError(t, err)
		assert.Equal(t, 0, len(meets))
	})

	t.Run("ok, list", func(t *testing.T) {
		// arrange
		mr := InitMeetingRepository(*db)
		ar := accountrepository.InitAccountRepository(*db)
		sr := studentrepository.InitStudentRepository(*db)

		prof, err := ar.AddProfessor(domainaggregate.Professor{
			Person: domainaggregate.Person{
				Name:       "dsf",
				Surname:    "sdf",
				Middlename: "sdf",
			},
			ScienceDegree: time.Now().Format(time.RFC3339),
		})
		assert.NoError(t, err)

		stud, err := sr.CreateStudent(domainaggregate.Student{
			Person: domainaggregate.Person{
				Name:       "dsf",
				Surname:    "sdf",
				Middlename: "sdf",
			},
			Cource: 1,
		})

		assert.NoError(t, err)

		meets := []domainaggregate.Meeting{}
		fromTime := time.Now().UTC().Round(time.Second)
		meet, err := mr.CreateMeeting(domainaggregate.Meeting{
			Name:          time.Now().Format(time.RFC3339),
			Description:   "lwefl",
			OrganizerId:   prof.Id,
			ParticipantId: stud.Id,
			Time:          fromTime.Add(time.Hour * 4),
			IsOnline:      true,
			Status:        domainaggregate.MeetingPlanned,
		})
		assert.NoError(t, err)
		meets = append(meets, meet)

		meet, err = mr.CreateMeeting(domainaggregate.Meeting{
			Name:          time.Now().Format(time.RFC3339),
			Description:   "lwefl",
			OrganizerId:   prof.Id,
			ParticipantId: stud.Id,
			Time:          fromTime.Add(time.Hour * 5),
			IsOnline:      true,
			Status:        domainaggregate.MeetingPlanned,
		})
		assert.NoError(t, err)
		meets = append(meets, meet)

		// act
		foundMeets, err := mr.GetProfessorMeetings(prof.Id, fromTime, time.Time{})

		// assert
		assert.NoError(t, err)
		assert.Equal(t, len(foundMeets), len(meets))
		for i := range meets {
			assert.Equal(t, meets[i].Id, foundMeets[i].Id)
			assert.Equal(t, meets[i].Description, foundMeets[i].Description)
			assert.Equal(t, meets[i].Name, foundMeets[i].Name)
			assert.Equal(t, meets[i].IsOnline, foundMeets[i].IsOnline)
			assert.Equal(t, meets[i].OrganizerId, foundMeets[i].OrganizerId)
			assert.Equal(t, meets[i].ParticipantId, foundMeets[i].ParticipantId)
			assert.Equal(t, meets[i].Status, foundMeets[i].Status)
			assert.Equal(t, meets[i].Time, foundMeets[i].Time)
		}

		// cleanup
		profId, err := strconv.Atoi(prof.Id)
		assert.NoError(t, err)
		err = ar.DeleteProfessor(profId)
		assert.NoError(t, err)

		studId, err := strconv.Atoi(stud.Id)
		assert.NoError(t, err)
		err = sr.DeleteStudent(studId)
		assert.NoError(t, err)
	})

	t.Run("ok, list part", func(t *testing.T) {
		// arrange
		mr := InitMeetingRepository(*db)
		ar := accountrepository.InitAccountRepository(*db)
		sr := studentrepository.InitStudentRepository(*db)

		prof, err := ar.AddProfessor(domainaggregate.Professor{
			Person: domainaggregate.Person{
				Name:       "dsf",
				Surname:    "sdf",
				Middlename: "sdf",
			},
			ScienceDegree: time.Now().Format(time.RFC3339),
		})
		assert.NoError(t, err)

		stud, err := sr.CreateStudent(domainaggregate.Student{
			Person: domainaggregate.Person{
				Name:       "dsf",
				Surname:    "sdf",
				Middlename: "sdf",
			},
			Cource: 1,
		})

		assert.NoError(t, err)

		meetsBeforeTo := []domainaggregate.Meeting{}
		fromTime := time.Now().UTC().Round(time.Second)
		meet, err := mr.CreateMeeting(domainaggregate.Meeting{
			Name:          time.Now().Format(time.RFC3339),
			Description:   "lwefl",
			OrganizerId:   prof.Id,
			ParticipantId: stud.Id,
			Time:          fromTime.Add(time.Hour * 2),
			IsOnline:      true,
			Status:        domainaggregate.MeetingPlanned,
		})
		assert.NoError(t, err)
		meetsBeforeTo = append(meetsBeforeTo, meet)

		meet, err = mr.CreateMeeting(domainaggregate.Meeting{
			Name:          time.Now().Format(time.RFC3339),
			Description:   "lwefl",
			OrganizerId:   prof.Id,
			ParticipantId: stud.Id,
			Time:          fromTime.Add(time.Hour * 3),
			IsOnline:      true,
			Status:        domainaggregate.MeetingPlanned,
		})
		assert.NoError(t, err)
		meetsBeforeTo = append(meetsBeforeTo, meet)

		meetsAfterTo := []domainaggregate.Meeting{}
		meet, err = mr.CreateMeeting(domainaggregate.Meeting{
			Name:          time.Now().Format(time.RFC3339),
			Description:   "lwefl",
			OrganizerId:   prof.Id,
			ParticipantId: stud.Id,
			Time:          fromTime.Add(time.Hour * 8),
			IsOnline:      true,
			Status:        domainaggregate.MeetingPlanned,
		})
		assert.NoError(t, err)
		meetsAfterTo = append(meetsAfterTo, meet)

		log.Printf("%v", meetsAfterTo)
		// act
		foundMeets, err := mr.GetProfessorMeetings(prof.Id, fromTime, fromTime.Add(time.Hour*5))

		// assert
		assert.NoError(t, err)
		assert.Equal(t, len(foundMeets), len(meetsBeforeTo))
		for i := range meetsBeforeTo {
			assert.Equal(t, meetsBeforeTo[i].Id, foundMeets[i].Id)
			assert.Equal(t, meetsBeforeTo[i].Description, foundMeets[i].Description)
			assert.Equal(t, meetsBeforeTo[i].Name, foundMeets[i].Name)
			assert.Equal(t, meetsBeforeTo[i].IsOnline, foundMeets[i].IsOnline)
			assert.Equal(t, meetsBeforeTo[i].OrganizerId, foundMeets[i].OrganizerId)
			assert.Equal(t, meetsBeforeTo[i].ParticipantId, foundMeets[i].ParticipantId)
			assert.Equal(t, meetsBeforeTo[i].Status, foundMeets[i].Status)
			assert.Equal(t, meetsBeforeTo[i].Time, foundMeets[i].Time)
		}

		// cleanup
		profId, err := strconv.Atoi(prof.Id)
		assert.NoError(t, err)
		err = ar.DeleteProfessor(profId)
		assert.NoError(t, err)

		studId, err := strconv.Atoi(stud.Id)
		assert.NoError(t, err)
		err = sr.DeleteStudent(studId)
		assert.NoError(t, err)
	})

}
