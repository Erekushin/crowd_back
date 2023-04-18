package category

import (
	"crowdfund/pkg/helpers/data"

	"gorm.io/gorm"
)

type Category struct {
	Id        uint           `json:"id,string" gorm:"primaryKey"`
	Name      string         `json:"name" validate:"required"`
	CreatedAt data.LocalTime `json:"created_date" gorm:"autoCreateTime"`
	CreatedBy uint           `json:"-"`
	UpdatedAt data.LocalTime `json:"-" gorm:"autoUpdateTime"`
	UpdatedBy uint           `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	DeletedBy uint           `json:"-"`
}
