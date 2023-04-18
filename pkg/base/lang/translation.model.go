package lang

import (
	"os"
	"sort"

	"crowdfund/pkg/database"
	"crowdfund/pkg/helpers/convertor"
	"crowdfund/pkg/helpers/data"
	"crowdfund/pkg/oauth"

	"github.com/gofiber/fiber/v2"
)

type Translation struct {
	Id        uint           `json:"id,string" gorm:"primaryKey"`
	KeyId     uint           `json:"key_id,string" validate:"required"`
	LangId    uint           `json:"lang_id,string" validate:"required"`
	Text      string         `json:"text" gorm:"type:varchar(2000)" validate:"required"`
	CreatedAt data.LocalTime `json:"created_date" gorm:"autoCreateTime"`
	CreatedBy uint           `json:"created_by,omitempty"`
}

func (*Translation) TableName() string {
	return os.Getenv("DB_SCHEMA") + ".tbl_translations"
}

type ResTranslations struct {
	LangId         uint   `json:"lang_id,string"`
	LangCode       string `json:"lang_code"`
	LangName       string `json:"lang_name"`
	TranslatedText string `json:"translated_text"`
}

type ResLang struct {
	Id   uint   `json:"id,string"`
	Code string `json:"code"`
	Name string `json:"name"`
}

type ResKeys struct {
	KeyId   uint   `json:"key_id,string"`
	KeyCode string `json:"key_code"`
}

type RequestItem struct {
	LangId         uint   `json:"lang_id,string"`
	TranslatedText string `json:"translated_text"`
}
type RequestTranslations struct {
	KeyId uint          `json:"key_id,string"`
	Items []RequestItem `json:"items"`
}

func (rt *RequestTranslations) TranslationSet(c *fiber.Ctx) error {
	db := database.DBconn

	db.Unscoped().Delete(&Translation{}, "key_id = ?", rt.KeyId)
	createdBy := oauth.GetSessionUserId(c)

	for _, item := range rt.Items {
		db.Save(&Translation{KeyId: rt.KeyId, LangId: item.LangId, Text: item.TranslatedText, CreatedBy: createdBy})
	}
	return nil
}

type ResTranslationItem struct {
	KeyId        uint              `json:"key_id,string"`
	KeyCode      string            `json:"key_code"`
	Translations []ResTranslations `json:"translations"`
}

func TranslationList(c *fiber.Ctx) (*data.Pagination, error) {
	var totalRow int64
	db := database.DBconn
	keys := make([]ResKeys, 0)
	key_id := convertor.StringToInt(c.Query("key_id"))

	langs := make([]ResLang, 0)
	db.Model(Lang{}).Find(&langs)

	tx := db.Table(os.Getenv("DB_SCHEMA") + ".tbl_keys as k").Select("k.id key_id, k.code key_code")

	if key_id != 0 {
		tx.Where("id=?", key_id)
	}

	tx.Count(&totalRow)

	p := data.Paginate(c, totalRow)

	err := tx.Offset(p.Offset).Limit(p.PageSize).Find(&keys).Error
	if err != nil {
		return nil, err
	}

	var res []ResTranslationItem
	for _, k := range keys {
		keyItem := new(ResTranslationItem)
		keyItem.KeyId = k.KeyId
		keyItem.KeyCode = k.KeyCode
		transRes := make([]ResTranslations, 0)
		for _, lang := range langs {
			translationItem := ResTranslations{
				LangId:         lang.Id,
				LangCode:       lang.Code,
				LangName:       lang.Name,
				TranslatedText: "",
			}
			var t Translation
			db.Where("lang_id=? AND key_id=?", lang.Id, k.KeyId).First(&t)
			translationItem.TranslatedText = t.Text
			transRes = append(transRes, translationItem)
		}

		keyItem.Translations = transRes
		res = append(res, *keyItem)
	}

	sort.Slice(res, func(i, j int) bool {
		return res[i].KeyCode < res[j].KeyCode
	})

	p.Items = res
	return p, nil
}
