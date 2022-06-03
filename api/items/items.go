package items

import (
	"errors"
	appi18n "github.com/mecamon/shoppingify-server/i18n"
	"github.com/mecamon/shoppingify-server/models"
	"time"
)

type ItemDom struct {
	appLocales appi18n.AppLocales
	item       models.Item
	fileInfo   FileInfo
}

type FileInfo struct {
	Size        int64
	ContentType string
}

const (
	nameMaxLength       = 30
	noteMaxLength       = 100
	maxFileSize   int64 = 5242880
)

var allowedContentTypes = []string{"image/jpeg", "image/png"}

func (m *ItemDom) validItem() (bool, models.ErrorMap) {
	errMap := models.ErrorMap{}

	if len(m.item.Name) > nameMaxLength {
		td := map[string]interface{}{"Limit": nameMaxLength}
		errMap["name"] = m.appLocales.GetMsg("InvalidItemName", td)
	}
	if len(m.item.Note) > noteMaxLength {
		td := map[string]interface{}{"Limit": noteMaxLength}
		errMap["note"] = m.appLocales.GetMsg("InvalidItemNote", td)
	}
	if m.item.CategoryID == 0 {
		errMap["category"] = m.appLocales.GetMsg("CategoryMandatory", nil)
	}
	if err := m.evaluateFile(); err != nil {
		errMap["image"] = err.Error()
	}
	return len(errMap) == 0, errMap
}

func (m *ItemDom) evaluateFile() error {
	var isTypeAllowed bool

	if m.fileInfo.Size == 0 && m.fileInfo.ContentType == "" {
		return nil
	}
	if m.fileInfo.Size > maxFileSize {
		td := map[string]interface{}{"MaxSize": maxFileSize}
		return errors.New(m.appLocales.GetMsg("FileTooBig", td))
	}
	for _, fileType := range allowedContentTypes {
		if isTypeAllowed = fileType == m.fileInfo.ContentType; isTypeAllowed {
			break
		}
	}
	if !isTypeAllowed {
		td := map[string]interface{}{"Types": allowedContentTypes[0] + "or" + allowedContentTypes[1]}
		return errors.New(m.appLocales.GetMsg("WrongFileType", td))
	}
	return nil
}

func (m *ItemDom) completedItemInfo() models.Item {
	return models.Item{
		Name:       m.item.Name,
		Note:       m.item.Note,
		CategoryID: m.item.CategoryID,
		CreatedAt:  time.Now().Unix(),
		UpdatedAt:  time.Now().Unix(),
	}
}
