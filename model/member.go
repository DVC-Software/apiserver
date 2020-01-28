package model

import (
	"github.com/jinzhu/gorm"
)

type Member struct {
	gorm.Model
	SlackUserID   string  `gorm:"unique;"`
	DiscordUserID string  `gorm:"unique;"`
	IsStaff       bool    `gorm:"not null;default:false"`
	Profile       Profile `gorm:"unique;not null;foreignkey:MemberID"`
}

// Training status should be converted by the front end since
// gorm doesn't have good enum support

type Profile struct {
	gorm.Model
	MemberID       uint
	Name           string      `gorm:"unique;not null"`
	IsVerified     bool        `gorm:"not null;default:false"`
	TrainingStatus int         `gorm:"not null;default:0"`
	Positions      []*Position `gorm:"many2many:profile_positions;"`
}

type Position struct {
	gorm.Model
	Role        string     `gorm:"unique;not null;default:'General Member'"`
	Description string     `gorm:"not null;default:'No Description'"`
	Profile     []*Profile `gorm:"many2many:profile_positions;"`
}

type MemberInit struct {
	Name          string
	SlackUserID   string
	DiscordUserID string
	CreatedFrom   string
}

func CreateMember(db *gorm.DB, initMember MemberInit) (bool, Member) {
	var slackRecord, discordRecord Member
	slackErr := db.Model(&slackRecord).Where("slack_user_id = ?", initMember.SlackUserID).First(&slackRecord).Error
	discordErr := db.Model(&discordRecord).Where("discord_user_id = ?", initMember.DiscordUserID).First(&discordRecord).Error
	// Check if already exist
	if slackErr == nil || discordErr == nil {
		return false, slackRecord
	}
	// initialize empty profile and initial position
	name := initMember.Name
	var position Position
	var profile Profile
	var member Member
	profile.Name = name
	// override default for staff member
	if initMember.CreatedFrom == "slack" {
		member.IsStaff = true
		position.Role = "General Staff"
		profile.TrainingStatus = 2
		profile.IsVerified = true
	}
	// initialize the member struct
	member.SlackUserID = initMember.SlackUserID
	member.DiscordUserID = initMember.DiscordUserID
	member.Profile = profile
	// save record to db
	db.Create(&member)
	return true, member
}
