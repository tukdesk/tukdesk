package helpers

import (
	"github.com/tukdesk/httputils/tools"
	"github.com/tukdesk/tukdesk/backend/models"
)

const (
	BrandNameMaxLength = 20
	BrandAPIKeyLength  = 16
	BrandSaltLength    = 10
)

func BrandFindOne(query map[string]interface{}) (*models.Brand, error) {
	brand := &models.Brand{}
	return brand, brand.FindOne(query)
}

func BrandInit(name string) (*models.Brand, error) {
	brand := models.NewBrand(name)
	brand.Authorization.APIKey = BrandNewAPIKey()
	brand.Authorization.Salt = BrandNewSalt()
	if err := brand.Insert(); err != nil {
		return nil, err
	}

	setCurrentBrand(brand)
	return brand, nil
}

func BrandNewAPIKey() string {
	return tools.RandString(BrandAPIKeyLength)
}

func BrandNewSalt() string {
	return tools.RandString(BrandSaltLength)
}

func BrandUpdateCurrent(change map[string]interface{}) error {
	query := M{"_id": currentBrand.Id}
	return currentBrand.FindAndModify(query, change)
}
