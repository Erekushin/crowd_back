package user

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"crowdfund/pkg/database"
	"crowdfund/pkg/helpers"
	"crowdfund/pkg/helpers/data"
	"crowdfund/pkg/oauth"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type User struct {
	Id                 uint           `json:"id,string" gorm:"primaryKey"`
	CivilId            uint           `json:"civil_id,string"`
	RegNo              string         `json:"reg_no" gorm:"not null"`
	FamilyName         string         `json:"family_name" gorm:"type:varchar(80)"`
	LastName           string         `json:"last_name" gorm:"type:varchar(150)"`
	FirstName          string         `json:"first_name" gorm:"type:varchar(150)"`
	Username           string         `json:"username" gorm:"type:varchar(200)"`
	Password           string         `json:"-" gorm:"type:varchar(100)"`
	RootAccount        uint           `json:"root_account,string"`
	Email              string         `json:"email" gorm:"type:varchar(50)"`
	PhoneNo            string         `json:"phone_no" gorm:"type:varchar(20)"`
	Gender             int            `json:"gender"`
	BirthDate          string         `json:"birth_date" gorm:"type:varchar(10)"`
	IsForeign          uint           `json:"is_foreign" gorm:"default:0"`
	Hash               string         `json:"hash,omitempty" gorm:"type:varchar(200)"`
	AimagCode          string         `json:"aimag_code" gorm:"type:varchar(5)"`
	AimagName          string         `json:"aimag_name" gorm:"type:varchar(255)"`
	SumCode            string         `json:"sum_code" gorm:"type:varchar(5)"`
	SumName            string         `json:"sum_name" gorm:"type:varchar(255)"`
	BagCode            string         `json:"bag_code" gorm:"type:varchar(5)"`
	BagName            string         `json:"bag_name" gorm:"type:varchar(255)"`
	Address            string         `json:"address" gorm:"type:varchar(600)"`
	AddressType        uint           `json:"address_type,omitempty"`
	AddressTypeName    string         `json:"address_type_name,omitempty" gorm:"type:varchar(255)"`
	ProfileImage       string         `json:"profile_image"`
	CountryCode        string         `json:"country_code,omitempty" gorm:"type:varchar(10);default:MNG"`
	CountryName        string         `json:"country_name,omitempty" gorm:"type:varchar(500);default:'Монгол улс'"`
	Nationality        string         `json:"nationality,omitempty" gorm:"varchar(500);default:Mongolia"`
	FirstNameEn        string         `json:"first_name_en,omitempty" gorm:"type:varchar(100)"`
	LastNameEn         string         `json:"last_name_en,omitempty" gorm:"type:varchar(100)"`
	FamilyNameEn       string         `json:"family_name_en,omitempty" gorm:"type:varchar(100)"`
	CountryNameEn      string         `json:"country_name_en,omitempty" gorm:"type:varchar(255);default:Mongolia"`
	CLevel             uint           `json:"c_level"`
	CreatedAt          data.LocalTime `json:"created_date" gorm:"autoCreateTime"`
	CreatedBy          uint           `json:"created_by,omitempty"`
	UpdatedAt          data.LocalTime `json:"updated_date" gorm:"autoUpdateTime"`
	UpdatedBy          uint           `json:"updated_by,omitempty"`
	DeletedAt          gorm.DeletedAt `json:"-" gorm:"index"`
	DeletedBy          uint           `json:"-"`
	IsConfirmedPhoneNo uint           `json:"is_confirmed_phone_no" gorm:"type:int4;default:0"`
	IsConfirmedEmail   uint           `json:"is_confirmed_email" gorm:"type:int4;default:0"`
}

type ApiUser struct {
	Id           uint   `json:"id,string"`
	LastName     string `json:"last_name"`
	FirstName    string `json:"first_name"`
	PhoneNo      string `json:"phone_no"`
	ProfileImage string `json:"profile_image"`
}

func (*User) TableName() string {
	return os.Getenv("DB_SCHEMA") + ".tbd_users"
}

func ById(id uint) *User {
	db := database.DBconn
	var v *User
	if err := db.First(&v, id).Error; err != nil {
		return nil
	}
	return v
}

func (p *User) Create() error {
	db := database.DBconn
	return db.Create(p).Error
}

func (p *User) Update() error {
	db := database.DBconn
	return db.Omit("email", "phone_no", "password", "c_level", "country_code", "reg_no", "username", "birth_date", "gender", "is_foreign", "civil_id", "hash", "root_account").Updates(p).Error
}

func (p *User) Delete() error {
	db := database.DBconn
	return db.Delete(p).Error
}

func (p *User) DeletePermanently() error {
	db := database.DBconn
	return db.Unscoped().Delete(p).Error
}

func List(c *fiber.Ctx) (*data.Pagination, error) {
	var totalRow int64
	db := database.DBconn
	users := make([]User, 0)
	search_text := strings.ToLower(c.Query("search_text"))

	tx := db.Model(User{})
	if search_text != "" {
		tx.Where("reg_no LIKE ? OR lower(first_name) LIKE ?", "%"+search_text+"%", "%"+search_text+"%")
	}
	tx.Order("created_at desc")

	tx.Count(&totalRow)

	p := data.Paginate(c, totalRow)

	err := tx.Offset(p.Offset).Limit(p.PageSize).Find(&users).Error
	if err != nil {
		return nil, err
	}
	p.Items = users
	return p, nil
}

func (p *User) ChangePassword(old, new string) error {
	db := database.DBconn
	err := db.Where("id = ?", p.Id).First(&p).Error
	if err != nil {
		return err
	}
	password := helpers.GeneratePassword(old)
	if p.Password != password {
		return errors.New("wrong password")
	}

	return db.Model(&p).Update("password", helpers.GeneratePassword(new)).Error
}

func (p *User) UpdatePassword(new string) error {
	db := database.DBconn
	err := db.Where("id = ?", p.Id).First(&p).Error
	if err != nil {
		return err
	}
	return db.Model(&p).Update("password", helpers.GeneratePassword(new)).Error
}

func (p *User) ChangeUsername(username string) error {
	db := database.DBconn
	err := db.Where("id = ?", p.Id).First(&p).Error
	if err != nil {
		return err
	}
	return db.Model(&p).Update("username", username).Error
}

func (p *User) Save() error {
	db := database.DBconn
	return db.Updates(p).Error
}

func (u *User) FindById() error {
	db := database.DBconn
	return db.Where("id = ?", u.Id).First(&u).Error
}

func (u *User) FindByUsername() error {
	db := database.DBconn
	return db.Where("username=?", u.Username).First(&u).Error
}

func (u *User) ExistWithEmail() bool {
	db := database.DBconn
	if err := db.First(&u, "email=?", u.Email).Error; err != nil {
		return false
	}
	if u.Id == 0 {
		return false
	}
	return true
}

func (u *User) ExistWithPhone() bool {
	db := database.DBconn
	if err := db.First(&u, "phone_no=?", u.PhoneNo).Error; err != nil {
		return false
	}
	if u.Id == 0 {
		return false
	}
	return true
}

func (u *User) MergeSessionUser(s *oauth.Session) error {
	fmt.Println("MergeSessionUser user id:", s.UserId)
	if s.UserId >= 10000000 {
		return nil
	}

	oldUser := User{Id: s.UserId}
	oldUser.FindById()

	u.Email = oldUser.Email
	u.Username = oldUser.Username
	u.PhoneNo = oldUser.PhoneNo
	u.Password = oldUser.Password
	u.IsConfirmedEmail = oldUser.IsConfirmedEmail
	u.IsConfirmedPhoneNo = oldUser.IsConfirmedPhoneNo
	oldUser.DeletePermanently()
	s.UserId = u.Id
	if err := s.Save(); err != nil {
		return err
	}
	if err := u.Save(); err != nil {
		return err
	}

	return nil
}

type ReqEmailOrChange struct {
	Otp      string `json:"otp" validate:"required"`
	Identity string `json:"identity" validate:"required"`
}
