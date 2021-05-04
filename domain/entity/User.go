package entity

import (
	"strings"

	"github.com/smarest/smarest-common/domain/value"
)

type User struct {
	UserName          string          `db:"user_name" json:"userName"`
	RestaurantGroupID int64           `db:"restaurant_group_id" json:"restaurantGroupID"`
	Role              string          `db:"role" json:"role"`
	Password          string          `db:"password" json:"-"`
	Name              string          `db:"name" json:"name"`
	Available         bool            `db:"available" json:"available"`
	SalaryType        string          `db:"salary_type" json:"salaryType"`
	JoinedDate        value.DateTime  `db:"joined_date" json:"joinedDate"`
	LeftDate          *value.DateTime `db:"left_date" json:"leftDate"`
}

func (item *User) ToSlide(fields string) map[string]interface{} {
	result := make(map[string]interface{})
	for _, field := range strings.Split(fields, ",") {
		switch field {
		case "userName":
			result[field] = item.UserName
		case "restaurant_group_id":
			result[field] = item.RestaurantGroupID
		case "role":
			result[field] = item.Role
		case "name":
			result[field] = item.Name
		case "salaryType":
			result[field] = item.SalaryType
		case "joinedDate":
			result[field] = item.JoinedDate
		case "leftDate":
			result[field] = item.LeftDate
		default:
		}
	}
	return result
}
