package server

import "gorm.io/gorm"

type DbAlertData struct {
	MsgTimestamp string
	Status       string
	AlertData
	gorm.Model
}

func Create(msgTs string, a AlertData, db *gorm.DB) {
	d := DbAlertData{
		MsgTimestamp: msgTs,
		Status:       "Triggered",
		AlertData:    a,
	}
	db.Create(&d)
}

func Update(msgTs string, status string, db *gorm.DB) {
	d := DbAlertData{
		MsgTimestamp: msgTs,
	}
	db.Model(&d).Where("msg_timestamp=?", msgTs).Update("status", status)
}

func Retrieve(msgTs string, db *gorm.DB) AlertData {
	d := DbAlertData{
		MsgTimestamp: msgTs,
	}
	db.First(&d)
	return d.AlertData
}

func Delete(msgTs string, db *gorm.DB) {
	d := DbAlertData{
		MsgTimestamp: msgTs,
	}
	db.Delete(&d)
}
