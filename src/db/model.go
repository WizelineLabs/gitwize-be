package db

import (
	"time"
)

type Repository struct {
	ID              uint      `json:"id" gorm:"primary_key"`
	Name            string    `json:"name"`
	Url             string    `json:"url"`
	Status          string    `json:"status"`
	UserName        string    `json:"username"`
	PassWord        string    `json:"password"`
	CtlCreatedDate  time.Time `json:"ctl_created_date"`
	CtlCreatedBy    string    `json:"ctl_created_by"`
	CtlModifiedDate time.Time `json:"ctl_modified_date"`
	CtlModifiedBy   string    `json:"ctl_modified_by"`
}
