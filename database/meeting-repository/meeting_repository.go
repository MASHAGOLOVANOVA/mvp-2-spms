package meetingrepository

import (
	"errors"
	"gorm.io/gorm"
	"mvp-2-spms/database"
	"mvp-2-spms/database/models"
	entities "mvp-2-spms/domain-aggregate"
	usecasemodels "mvp-2-spms/services/models"
)

type MeetingRepository struct {
	dbContext database.Database
}

func InitMeetingRepository(dbcxt database.Database) *MeetingRepository {
	return &MeetingRepository{
		dbContext: dbcxt,
	}
}

func (r *MeetingRepository) GetProjectMeetingByStudMeetingId(studMeetingId string) (entities.ProjectMeeting, error) {

	var meetingDb models.ProjectMeeting
	tx := r.dbContext.DB.Select("*").Where(" stud_meeting_id = ?", studMeetingId).Take(&meetingDb)

	if tx.Error != nil {
		return entities.ProjectMeeting{}, nil
	}
	return meetingDb.MapToEntity(), nil
}

func (r *MeetingRepository) GetProfessorStudentMeetings(profId string) ([]entities.StudMeeting, error) {
	slots, _ := r.GetProfessorSlots(profId, "")
	meetings := []entities.StudMeeting{}
	for _, slot := range slots {

		var meetingsDb []models.StudMeeting

		res := r.dbContext.DB.Select("*").Where(" slot_id = ?", slot.Id).Find(&meetingsDb)
		if res.Error != nil {
			continue
		}

		for _, m := range meetingsDb {
			meetings = append(meetings, m.MapToEntity())
			break
		}
	}
	return meetings, nil
}

func (r *MeetingRepository) GetStudentMeetings(studId string) ([]entities.StudMeeting, error) {
	var meetingsDb []models.StudMeeting

	result := r.dbContext.DB.Select("*").Where(" student_id = ?", studId).Find(&meetingsDb)

	if result.Error != nil {
		return []entities.StudMeeting{}, result.Error
	}
	meetings := []entities.StudMeeting{}
	for _, m := range meetingsDb {
		meetings = append(meetings, m.MapToEntity())
	}
	return meetings, nil
}

func (r *MeetingRepository) GetProfessorSlots(profId string, filter string) ([]entities.Slot, error) {
	var meetingsDb []models.Slot

	result := r.dbContext.DB.Select("*").Where(" professor_id = ?", profId).Find(&meetingsDb)

	if result.Error != nil {
		return []entities.Slot{}, result.Error
	}

	meetings := []entities.Slot{}
	for _, m := range meetingsDb {
		if filter == "free" {
			var studmeetDb []models.StudMeeting
			res := r.dbContext.DB.Select("*").Where(" slot_id = ?", m.Id).Find(&studmeetDb)
			if res.Error == nil {
				if len(studmeetDb) > 0 {
					continue
				}
			}
		}
		meetings = append(meetings, m.MapToEntity())
	}
	return meetings, nil
}

func (r *MeetingRepository) DeleteSlot(slotId int) error {
	dbmeeting := models.Slot{Id: uint(slotId)}

	result := r.dbContext.DB.Delete(&dbmeeting)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *MeetingRepository) ChooseSlot(slotId int, studId int) error {
	dbStudMeeting := models.StudMeeting{}
	dbStudMeeting.MapToThis(slotId, studId)
	result := r.dbContext.DB.Create(&dbStudMeeting)

	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *MeetingRepository) BindSlotToProject(slotId int, projectId int) error {
	dbProjectMeeting := models.ProjectMeeting{}
	dbProjectMeeting.MapToThis(slotId, projectId)
	result := r.dbContext.DB.Create(&dbProjectMeeting)

	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *MeetingRepository) AddSlot(meeting entities.Slot, planerId string) (entities.Slot, error) {
	dbmeeting := models.Slot{}
	dbmeeting.MapEntityToThis(meeting, planerId)
	result := r.dbContext.DB.Create(&dbmeeting)
	if result.Error != nil {
		return entities.Slot{}, result.Error
	}
	return dbmeeting.MapToEntity(), nil
}

func (r *MeetingRepository) GetSlotById(slotId int) (entities.Slot, error) {
	meetingsDb := models.Slot{}

	result := r.dbContext.DB.Select("*").Where("id = ?", slotId).Find(&meetingsDb)
	if result.Error != nil {
		return entities.Slot{}, result.Error
	}

	return meetingsDb.MapToEntity(), nil
}

func (r *MeetingRepository) CreateMeeting(meeting entities.Meeting) (entities.Meeting, error) {
	dbmeeting := models.Meeting{}
	dbmeeting.MapEntityToThis(meeting)
	result := r.dbContext.DB.Create(&dbmeeting)
	if result.Error != nil {
		return entities.Meeting{}, result.Error
	}
	return dbmeeting.MapToEntity(), nil
}

func (r *MeetingRepository) DeleteMeeting(id int) error {
	dbmeeting := models.Meeting{Id: uint(id)}

	result := r.dbContext.DB.Delete(&dbmeeting)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *MeetingRepository) AssignPlannerMeeting(meeting usecasemodels.PlannerMeeting) error {
	err := r.dbContext.DB.Transaction(func(tx *gorm.DB) error {
		result := r.dbContext.DB.Model(&models.Meeting{}).Where("id = ?", meeting.Meeting.Id).Update("planner_id", meeting.MeetingPlannerId)
		if result.Error != nil {
			return result.Error
		}

		if result.RowsAffected == 0 {
			return usecasemodels.ErrMeetingNotFound
		}

		return nil
	})

	return err
}

func (r *MeetingRepository) GetProfessorMeetings(profId string) ([]entities.Slot, error) {
	var meetingsDb []models.Slot

	query := r.dbContext.DB.Select("*")
	query = query.Where("professor_id = ? ", profId)

	meetings := []entities.Slot{}
	for _, m := range meetingsDb {
		meetings = append(meetings, m.MapToEntity())
	}
	return meetings, nil
}

func (r *MeetingRepository) GetMeetingById(meetId string) (entities.Meeting, error) {
	var dbMeeting models.Meeting

	result := r.dbContext.DB.Select("*").Where("id = ?", meetId).Take(&dbMeeting)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return entities.Meeting{}, usecasemodels.ErrMeetingNotFound
		}
		return entities.Meeting{}, result.Error
	}

	return dbMeeting.MapToEntity(), nil
}

func (r *MeetingRepository) GetMeetingPlannerId(meetId string) (string, error) {
	meeting := models.Meeting{}

	result := r.dbContext.DB.Select("planner_id").Where("id = ?", meetId).Take(&meeting)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return "", usecasemodels.ErrMeetingNotFound
		}
		return "", result.Error
	}

	return meeting.PlannerId.String, nil
}
