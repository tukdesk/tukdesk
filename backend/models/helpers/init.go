package helpers

import (
	"github.com/tukdesk/tukdesk/backend/models"

	"github.com/tukdesk/mgoutils"
)

var currentBrand *models.Brand

func InitWithStorage(stg *mgoutils.MgoPool) error {
	models.SetStorage(stg)

	if err := models.EnsureIndexes(); err != nil {
		return err
	}

	if err := LoadCurrentBrand(); err != nil {
		return err
	}

	return nil
}

func setCurrentBrand(brand *models.Brand) {
	if brand == nil {
		return
	}

	currentBrand = brand
	return
}

func LoadCurrentBrand() error {
	query := M{"on": true}
	brand, err := BrandFindOne(query)
	if IsNotFound(err) {
		return nil
	}

	if err != nil {
		return err
	}

	setCurrentBrand(brand)
	return nil
}

func CurrentBrand() *models.Brand {
	return currentBrand
}
