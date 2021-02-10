package users

import (
	"encoding/json"
)

type PublicUser struct {
	Id          int64  `json:"id"`
	DateCreated string `json:"date_created"`
	Status      string `json:"status"`
}

type PrivateUser struct {
	Id          int64  `json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	DateCreated string `json:"date_created"`
	Status      string `json:"status"`
}

func (u *User) Marshal(isPublic bool) interface{} {
	if isPublic {
		return PublicUser{
			Id:          u.Id,
			DateCreated: u.DateCreated,
			Status:      u.Status,
		}
	}
	// Since most of the fields in the json annotations match between the private user and user structs, this is an easier approach
	userJSON, _ := json.Marshal(u)
	var privateUser PrivateUser
	json.Unmarshal(userJSON, &privateUser)
	return privateUser
}

func (userCollection Users) Marshal(isPublic bool) []interface{} {
	res := make([]interface{}, 0)
	for _, el := range userCollection {
		res = append(res, el.Marshal(isPublic))
	}
	return res
}
