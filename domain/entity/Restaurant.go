package entity

import (
	"strings"

	"github.com/smarest/smarest-common/domain/value"
)

type Restaurant struct {
	ID              int64           `db:"id" json:"id"`
	Name            string          `db:"name" json:"name"`
	AccessKey       string          `db:"access_key" json:"-"`
	Description     *string         `db:"description" json:"description"`
	Available       bool            `db:"available" json:"-"`
	Image           *string         `db:"image" json:"image"`
	Address         *string         `db:"address" json:"address"`
	Phone           *string         `db:"phone" json:"phone"`
	Creator         string          `db:"creator" json:"creator"`
	CreatedDate     value.DateTime  `db:"created_date" json:"createdDate"`
	Updater         *string         `db:"updater" json:"updater"`
	LastUpdatedDate *value.DateTime `db:"last_updated_date" json:"lastUpdatedDate"`
}

func (item *Restaurant) ToSlide(fields string) map[string]interface{} {
	result := make(map[string]interface{})
	// Loop over the parts from the string.
	for _, field := range strings.Split(fields, ",") {
		switch field {
		case "id":
			result[field] = item.ID
		case "name":
			result[field] = item.Name
		case "description":
			result[field] = item.Description
		case "image":
			result[field] = item.Image
		case "address":
			result[field] = item.Address
		case "phone":
			result[field] = item.Phone
		case "creator":
			result[field] = item.Creator
		case "createdDate":
			result[field] = item.CreatedDate
		case "updater":
			result[field] = item.Updater
		case "lastUpdatedDate":
			result[field] = item.LastUpdatedDate
		default:
		}
	}
	return result
}
