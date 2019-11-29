package log

import (
	"github.com/jinzhu/gorm"
)

// BoxLog Model
type BoxLog struct {
	gorm.Model
	UID  uint
	Name string
	Desc string
}
