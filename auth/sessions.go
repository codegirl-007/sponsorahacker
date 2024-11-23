package auth

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type User struct {
	NickName        string `json:"NickName"`
	Provider        string `json:"Provider"`
	Provider_userid string `json:"UserID"`
	Sid             string `json:"Sid"`
}

func HydrateUser(userJson string) *User {
	user := User{}

	if err := json.Unmarshal([]byte(userJson), &user); err != nil {
		fmt.Println("serialization in hydrating user error = ", err)
		return nil
	}

	return &user
}

func (u *User) GetSid() int {
	strSid, err := strconv.Atoi(u.Sid)

	if err != nil {
		fmt.Println("sid is invalid")
	}

	return strSid
}
