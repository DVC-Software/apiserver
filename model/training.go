package model

import (
	"github.com/jinzhu/gorm"
	"time"
)

type TrainingSession struct {
	gorm.Model
	Trainers  []Member  `gorm:"many2many:session_trainers"`
	Trainees  []Member  `gorm:"many2many:session_trainees"`
	StartTime time.Time `gorm:"not null"`
	EndTime   time.Time `gorm:"not null"`
	Duration  int       `gorm:"not null; default: 0"`
	Topic     string    `gorm:"not null"`
	Type      int       `gorm:"not null"`
}

type TrainingSessionInfo struct {
	Trainers  []string
	Trainees  []string
	StartTime time.Time
	EndTime   time.Time
	Topic     string
	Type      string
}

func (info TrainingSessionInfo) getIntType() int {
	if info.Type == "Group Session" {
		return 0
	} else if info.Type == "Personal Session" {
		return 1
	} else {
		return -1
	}
}

func (info TrainingSessionInfo) ValidateInfo() bool {
	if info.Topic == "" || len(info.Trainers) == 0 || len(info.Trainees) == 0 {
		return false
	}
	return true
}

func (session TrainingSession) calcDuration() int {
	day := session.StartTime.Day() - session.EndTime.Day()
	hour := session.StartTime.Hour() - session.EndTime.Hour()
	minute := session.StartTime.Minute() - session.EndTime.Minute()
	return day*24*60 + hour*60 + minute
}

func (session TrainingSession) toSessionInfo() TrainingSessionInfo {
	var info TrainingSessionInfo
	info.StartTime = session.StartTime
	info.EndTime = session.EndTime
	info.Topic = session.Topic
	if session.Type == 0 {
		info.Type = "Group Session"
	} else {
		info.Type = "Personal Session"
	}
	for _, trainer := range session.Trainers {
		info.Trainers = append(info.Trainers, trainer.Profile.Name)
	}
	for _, trainee := range session.Trainees {
		info.Trainees = append(info.Trainees, trainee.Profile.Name)
	}
	return info
}

func CreateTrainingSession(db *gorm.DB, info TrainingSessionInfo) (bool, TrainingSessionInfo, string) {
	var session TrainingSession
	// look for trainer profiles
	for _, trainer := range info.Trainers {
		var trainerProfile Profile
		var trainerMember Member
		err := db.Model(&session).Where("name = ?", trainer).First(&trainerProfile).Error
		if err != nil {
			return false, TrainingSessionInfo{}, "Trainer with name " + trainer + " does not exist"
		}
		err = db.Model(&trainerMember).Where("id = ?", trainerProfile.MemberID).First(&trainerMember).Error
		if err != nil {
			return false, TrainingSessionInfo{}, "Can't match member " + trainer + " with profile"
		}
		trainerMember.Profile = trainerProfile
		session.Trainers = append(session.Trainers, trainerMember)
	}
	// look for trainee profiles
	for _, trainee := range info.Trainees {
		var traineeProfile Profile
		var traineeMember Member
		err := db.Model(&session).Where("name = ?", trainee).First(&traineeProfile).Error
		if err != nil {
			return false, TrainingSessionInfo{}, "Trainee with name " + trainee + " does not exist"
		}
		err = db.Model(&traineeMember).Where("id = ?", traineeProfile.MemberID).First(&traineeMember).Error
		if err != nil {
			return false, TrainingSessionInfo{}, "Can't match member " + trainee + " with profile"
		}
		traineeMember.Profile = traineeProfile
		session.Trainees = append(session.Trainees, traineeMember)
	}
	session.StartTime = info.StartTime
	session.EndTime = info.EndTime
	session.Duration = session.calcDuration()
	session.Topic = info.Topic
	session.Type = info.getIntType()
	if session.Type == -1 {
		return false, TrainingSessionInfo{}, "Invalid session type " + info.Type
	}
	err := db.Create(&session).Error
	if err != nil {
		return false, TrainingSessionInfo{}, err.Error()
	}
	return true, session.toSessionInfo(), ""
}
