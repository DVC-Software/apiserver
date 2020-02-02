package model

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
)

type Member struct {
	gorm.Model
	SlackUserID   sql.NullString `gorm:"unique;"`
	DiscordUserID sql.NullString `gorm:"unique;"`
	IsStaff       bool           `gorm:"not null;default:false"`
	Profile       Profile        `gorm:"unique;not null;foreignkey:MemberID"`
}

// Training status should be converted by the front end since
// gorm doesn't have good enum support

type Profile struct {
	gorm.Model
	MemberID       uint
	Name           string     `gorm:"unique;not null"`
	IsVerified     bool       `gorm:"not null;default:false"`
	TrainingStatus int        `gorm:"not null;default:0"`
	Positions      []Position `gorm:"many2many:profile_positions;"`
}

type Position struct {
	gorm.Model
	Role        string `gorm:"unique;not null;default:'General Member'"`
	Description string `gorm:"not null;default:'No Description'"`
}

type MemberInit struct {
	Name          string
	SlackUserID   string
	DiscordUserID string
	CreatedFrom   string
}

type MemberInfo struct {
	ID             uint
	Name           string
	Positions      []string
	TrainingStatus int
	IsStaff        bool
}

func CreateMember(db *gorm.DB, initMember MemberInit) (bool, MemberInfo, string) {
	var slackRecord, discordRecord Member
	var nameRecord Profile
	slackErr := db.Model(&slackRecord).Where("slack_user_id = ?", NewNullString(initMember.SlackUserID)).First(&slackRecord).Error
	discordErr := db.Model(&discordRecord).Where("discord_user_id = ?", NewNullString(initMember.DiscordUserID)).First(&discordRecord).Error
	nameErr := db.Model(&nameRecord).Where("name = ?", NewNullString(initMember.Name)).First(&nameRecord).Error
	// Check if already exist
	if slackErr == nil {
		return false, toMemberInfo(slackRecord), "A member with this slack id already exists"
	}
	if discordErr == nil {
		return false, toMemberInfo(discordRecord), "A member with this discord id already exists"
	}
	if nameErr == nil {
		return false, toMemberInfo(Member{Profile: nameRecord}), "A member with this name already exists"
	}
	// initialize empty profile and initial position
	name := initMember.Name
	var position Position
	var profile Profile
	var member Member
	profile.Name = name
	position.Role = "General Member"
	// override default for staff member
	if initMember.CreatedFrom == "slack" {
		member.IsStaff = true
		position.Role = "General Staff"
		profile.TrainingStatus = 2
		profile.IsVerified = true
	}
	// initialize the member struct
	member.SlackUserID = NewNullString(initMember.SlackUserID)
	member.DiscordUserID = NewNullString(initMember.DiscordUserID)
	tempPosition := position
	err := db.Model(&position).Where("role = ?", position.Role).Find(&position).Error
	if err != nil {
		position = tempPosition
	}
	profile.Positions = append(profile.Positions, position)
	member.Profile = profile
	// save record to db
	db.Create(&member)
	return true, toMemberInfo(member), ""
}

func IdentifyMember(db *gorm.DB, id string) (bool, MemberInfo, string) {
	var slackRecord, discordRecord, member Member
	var valid bool
	var profile Profile
	errmsg := ""
	slackErr := db.Model(&slackRecord).Where("slack_user_id = ?", id).First(&slackRecord).Error
	discordErr := db.Model(&discordRecord).Where("discord_user_id = ?", id).First(&discordRecord).Error
	if slackErr == nil {
		member = slackRecord
		valid = true
	} else if discordErr == nil {
		member = discordRecord
		fmt.Println("discord record found")
		valid = true
	} else {
		valid = false
		errmsg = "Can't find any matching record by given id"
		fmt.Println("No record!")
	}
	db.Model(&member).Preload("Positions").Related(&profile)
	member.Profile = profile
	bytes, _ := json.Marshal(member)
	fmt.Println(string(bytes))
	info := toMemberInfo(member)
	return valid, info, errmsg
}

func toMemberInfo(member Member) MemberInfo {
	var info MemberInfo
	info.ID = member.ID
	info.Name = member.Profile.Name
	for _, position := range member.Profile.Positions {
		info.Positions = append(info.Positions, position.Role)
	}
	info.IsStaff = member.IsStaff
	info.TrainingStatus = member.Profile.TrainingStatus
	return info
}
