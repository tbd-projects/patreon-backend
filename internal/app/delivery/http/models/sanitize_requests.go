package models

import "github.com/microcosm-cc/bluemonday"

func (req *RequestCreator) Sanitize(sanitizer bluemonday.Policy) {
	req.Category = sanitizer.Sanitize(req.Category)
	req.Description = sanitizer.Sanitize(req.Description)
}

func (req *RequestLogin) Sanitize(sanitizer bluemonday.Policy) {
	req.Login = sanitizer.Sanitize(req.Login)
	req.Password = sanitizer.Sanitize(req.Password)
}

func (req *RequestChangePassword) Sanitize(sanitizer bluemonday.Policy) {
	req.OldPassword = sanitizer.Sanitize(req.OldPassword)
	req.NewPassword = sanitizer.Sanitize(req.NewPassword)
}

func (req *RequestRegistration) Sanitize(sanitizer bluemonday.Policy) {
	req.Login = sanitizer.Sanitize(req.Login)
	req.Nickname = sanitizer.Sanitize(req.Nickname)
	req.Password = sanitizer.Sanitize(req.Password)
}

func (req *RequestAwards) Sanitize(sanitizer bluemonday.Policy) {
	req.Name = sanitizer.Sanitize(req.Name)
	req.Description = sanitizer.Sanitize(req.Description)
}

func (req *RequestPosts) Sanitize(sanitizer bluemonday.Policy) {
	req.Title = sanitizer.Sanitize(req.Title)
	req.Description = sanitizer.Sanitize(req.Description)
}

func (req *RequestText) Sanitize(sanitizer bluemonday.Policy) {
	req.Text = sanitizer.Sanitize(req.Text)
}

func (req *SubscribeRequest) Sanitize(sanitizer bluemonday.Policy) {
	req.AwardName = sanitizer.Sanitize(req.AwardName)
}
