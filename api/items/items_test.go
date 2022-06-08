package items

import (
	appi18n "github.com/mecamon/shoppingify-server/i18n"
	"github.com/mecamon/shoppingify-server/models"
	"testing"
)

func wordsGenerator(wordLength int) string {
	var word = ""
	for i := 0; i < wordLength; i++ {
		word += "a"
	}
	return word
}

var itemsTests = []struct {
	testName       string
	item           models.Item
	fileInfo       FileInfo
	expectedErrors int
	expectedResult bool
}{
	{testName: "name-longer-than-30",
		item: models.Item{Name: wordsGenerator(32), Note: wordsGenerator(40), CategoryID: 2}, expectedErrors: 1, expectedResult: false,
	},
	{testName: "note-and-name-longer-than-100-and-30",
		item: models.Item{Name: wordsGenerator(32), Note: wordsGenerator(120), CategoryID: 2}, expectedErrors: 2, expectedResult: false,
	},
	{testName: "empty-category",
		item: models.Item{Name: wordsGenerator(23), Note: wordsGenerator(45)}, expectedErrors: 1, expectedResult: false,
	},
	{testName: "valid-item",
		item: models.Item{Name: wordsGenerator(25), Note: wordsGenerator(70), CategoryID: 2}, expectedErrors: 0, expectedResult: true,
	},
}

func TestItemDom_ValidItem(t *testing.T) {
	for _, tt := range itemsTests {
		t.Log(tt.testName)
		{
			appLocales := appi18n.GetLocales("en-EN")
			itemDom := ItemDom{
				appLocales: appLocales,
				item:       tt.item,
				fileInfo:   tt.fileInfo,
			}
			isValid, errMap := itemDom.validItem()
			if isValid != tt.expectedResult {
				t.Errorf("got %v when expected was %v", isValid, tt.expectedResult)
			}
			if len(errMap) != tt.expectedErrors {
				t.Errorf("got %d when expected errors were %d", len(errMap), tt.expectedErrors)
			}
		}
	}
}

var fileTests = []struct {
	testName    string
	contentType string
	fileSize    int64
	hasError    bool
}{
	{testName: "invalid-file-size", contentType: "image/jpeg", fileSize: 30000000, hasError: true},
	{testName: "invalid-file-type", contentType: "text/pdf", fileSize: 4200000, hasError: true},
	{testName: "valid-file", contentType: "image/jpeg", fileSize: 3250000, hasError: false},
}

func TestItemDom_EvaluateFile(t *testing.T) {
	for _, tt := range fileTests {
		t.Log(tt.testName)
		{
			appLocales := appi18n.GetLocales("en-EN")
			itemDom := ItemDom{
				appLocales: appLocales,
				item:       models.Item{},
				fileInfo:   FileInfo{Size: tt.fileSize, ContentType: tt.contentType},
			}
			err := itemDom.evaluateFile()
			hasError := err != nil

			if hasError != tt.hasError {
				t.Errorf("got %v when expected was %v", hasError, tt.hasError)
			}
		}
	}
}

func TestItemDom_CompletedItemInfo(t *testing.T) {
	appLocales := appi18n.GetLocales("en-EN")
	itemDom := ItemDom{
		appLocales: appLocales,
		item: models.Item{
			Name:       "Red Apple",
			Note:       "Imported from Chile",
			CategoryID: 3,
		},
		fileInfo: FileInfo{},
	}
	ci := itemDom.completedItemInfo()
	if ci.CreatedAt == 0 || ci.UpdatedAt == 0 {
		t.Error("item is not completed")
	}
}
