package document

import (
	"errors"
	"fmt"
	"os"

	"crowdfund/pkg/base/user"
	"crowdfund/pkg/database"
	"crowdfund/pkg/helpers/client"
	"crowdfund/pkg/helpers/convertor"
	"crowdfund/pkg/helpers/data"
	"crowdfund/pkg/oauth"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type CoreFindResult struct {
	Document  Document    `json:"document"`
	User      user.User   `json:"user"`
	SameUsers []user.User `json:"same_users"`
}

type SetUserReq struct {
	DocumentId uint `json:"document_id"`
	UserId     uint `json:"user_id"`
}

type Document struct {
	Id             uint           `json:"id,string"`
	UserId         uint           `json:"user_id,string"`
	DocumentNumber string         `json:"document_number"`
	TypeId         uint           `json:"type_id"`
	CategoryId     uint           `json:"category_id"` //irgenii unemleh, gadaad passport, joloonii unemleh
	CountryCode    string         `json:"country_code"`
	FamilyName     string         `json:"family_name"`
	FirstName      string         `json:"first_name"`
	LastName       string         `json:"last_name"`
	BirthDate      string         `json:"birth_date"`
	Gender         int            `json:"gender"`
	Hash           string         `json:"hash"`
	DateOfIssue    string         `json:"date_of_issue"`
	DateOfExpire   string         `json:"date_of_expire"`
	RegNo          string         `json:"reg_no"`
	RegNoEn        string         `json:"reg_no_en"`
	Mrz            string         `json:"mrz"`
	CountryName    string         `json:"country_name"`
	CountryNameEn  string         `json:"country_name_en"`
	TypeName       string         `json:"type_name"`
	TypeNameEn     string         `json:"type_name_en"`
	CategoryName   string         `json:"category_name"`
	Address        string         `json:"address"`
	PeopleImage    string         `json:"people_image"`
	FullImage      string         `json:"full_image"`
	CreatedAt      data.LocalTime `json:"created_date" gorm:"autoCreateTime"`
	CreatedBy      uint           `json:"created_by,omitempty"`
	UpdatedAt      data.LocalTime `json:"updated_date" gorm:"autoUpdateTime"`
	UpdatedBy      uint           `json:"updated_by,omitempty"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`
	DeletedBy      uint           `json:"-"`
}

func (*Document) TableName() string {
	return os.Getenv("DB_SCHEMA") + ".tbd_documents"
}

func DocumentList(c *fiber.Ctx, user_id uint) (*data.Pagination, error) {
	var totalRow int64
	db := database.DBconn
	documents := make([]Document, 0)

	tx := db.Model(Document{}).Where("user_id = ?", user_id)

	tx.Count(&totalRow)

	p := data.Paginate(c, totalRow)

	err := tx.Offset(p.Offset).Limit(p.PageSize).Find(&documents).Error
	if err != nil {
		return nil, err
	}
	p.Items = documents
	return p, nil
}

func (p *Document) Create() error {
	db := database.DBconn
	return db.Create(p).Error
}

func (p *Document) Update() error {
	db := database.DBconn
	return db.Updates(p).Error
}

func (p *Document) Delete() error {
	db := database.DBconn
	return db.Delete(p).Error
}

func SetUser(user_id, document_id uint) error {
	db := database.DBconn
	return db.Model(Document{}).Where("id = ?", document_id).Update("user_id", user_id).Error
}

func (d *Document) FindByDocumentNumber() error {
	db := database.DBconn
	return db.Where("document_number = ?", d.DocumentNumber).First(&d).Error
}

func (d *Document) FindByHash() error {
	db := database.DBconn
	return db.Where("hash = ?", d.Hash).First(&d).Error
}

type ApiFindDocument struct {
	CountryCode    string `json:"country_code" query:"country_code" validate:"required"`
	RegNo          string `json:"reg_no" query:"reg_no"`
	DocumentNumber string `json:"document_number" query:"document_number"`
	TypeId         string `json:"type_id" query:"type_id"`
	CategoryId     string `json:"category_id" query:"category_id"`
	PeopleImage    string `json:"people_image" query:"people_image"`
	FullImage      string `json:"full_image" query:"full_image"`
	Mrz            string `json:"mrz" query:"mrz"`
	DateOfIssue    string `json:"date_of_issue" query:"date_of_issue"`
	DateOfExpire   string `json:"date_of_expire" query:"date_of_expire"`
	RegNoEn        string `json:"reg_no_en" query:"reg_no_en"`
	FirstName      string `json:"first_name" query:"first_name"`
	LastName       string `json:"last_name" query:"last_name"`
	Gender         string `json:"gender" query:"gender"`
	BirthDate      string `json:"birth_date" query:"birth_date"`
	UserId         string `json:"user_id" query:"user_id"`
}

func ParseSameUsers(resp interface{}) []user.User {
	result := CoreFindResult{}
	convertor.MapToStruct(resp, &result)
	return result.SameUsers
}

func ParseAndCreateDocument(resp interface{}) (*Document, error) {
	result := CoreFindResult{}
	if err := convertor.MapToStruct(resp, &result); err != nil {
		return nil, err
	}
	var document *Document
	document = &result.Document
	if document.Id == 0 {
		return document, errors.New("document id is zero")
	}
	if err := document.Create(); err != nil {
		return nil, err
	}
	return document, nil
}

func ParseUser(session *oauth.Session, resp interface{}) (*user.User, error) {
	result := CoreFindResult{}
	convertor.MapToStruct(resp, &result)
	u := result.User

	if err := u.FindById(); err != nil {
		if err := u.Create(); err != nil {
			return nil, err
		}
	}

	if err := u.MergeSessionUser(session); err != nil {
		return nil, err
	}

	return &u, nil
}

func FindDocumentFromCore(session *oauth.Session, body *ApiFindDocument) *map[string]interface{} {
	req, _ := convertor.InterfaceToMap(body)
	header := map[string]string{"message_code": os.Getenv("MESSAGE_CODE_FIND_DOCUMENT")}
	resp := client.Request(os.Getenv("URL_USER_GO"), "POST", req, header)

	fmt.Println("document find resp: ", resp)

	result := make(map[string]interface{})
	if resp.Code == 200 {
		result["user"], _ = ParseUser(session, resp.Result)
		result["document"], _ = ParseAndCreateDocument(resp.Result)
		result["same_users"] = ParseSameUsers(resp.Result)

	}
	return &result
}
