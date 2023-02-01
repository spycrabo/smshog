package http

import "strings"

type MessageRequest struct {
	Login  string `form:"login"`
	Psw    string `form:"psw"`
	Sender string `form:"sender"`
	Phones string `form:"phones"`
	Mes    string `form:"mes"`
	Fmt    string `form:"fmt"`
}

type ListMessagesRequest struct {
	Offset int    `form:"offset,default=0"`
	Limit  int    `form:"limit,default=50"`
	Search string `form:"search"`
}

func (r *MessageRequest) GetPhones() (output []string) {
	for _, phone := range strings.Split(r.Phones, ",") {
		if len(phone) > 0 {
			output = append(output, phone)
		}
	}
	return output
}
