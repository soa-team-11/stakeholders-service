package models

import (
	"encoding/json"

	"github.com/google/uuid"
)

type Profile struct {
	UserID  uuid.UUID `json:"user_id" bson:"user_id"`
	Name    string    `json:"name" bson:"name"`
	Surname string    `json:"surname" bson:"surname"`
	Bio     string    `json:"bio" bson:"bio"`
	Motto   string    `json:"motto" bson:"motto"`
}

func (p *Profile) ToJSONString() string {
	b, err := json.Marshal(p)
	if err != nil {
		return ""
	}
	return string(b)
}

func (p *Profile) ToJSON() []byte {
	r, _ := json.Marshal(p)
	return r
}

func ProfileFromJSON(jsonString []byte) (*Profile, error) {
	var profile Profile
	err := json.Unmarshal(jsonString, &profile)
	return &profile, err
}
