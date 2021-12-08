package http_models

import (
	"github.com/microcosm-cc/bluemonday"
	"patreon/internal/app/models"
)

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

func (req *RequestComment) Sanitize(sanitizer bluemonday.Policy) {
	req.Body = sanitizer.Sanitize(req.Body)
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
func (req *RequestChangeNickname) Sanitize(sanitizer bluemonday.Policy) {
	req.OldNickname = sanitizer.Sanitize(req.OldNickname)
	req.NewNickname = sanitizer.Sanitize(req.NewNickname)
}

func (req *RequestAttach) Sanitize(sanitizer bluemonday.Policy) {
	req.Value = sanitizer.Sanitize(req.Value)
	req.Status = sanitizer.Sanitize(req.Status)
	req.Type = models.DataType(sanitizer.Sanitize(string(req.Type)))
}

func (req *RequestAttaches) Sanitize(sanitizer bluemonday.Policy) {
	for _, attach := range req.Attaches {
		attach.Sanitize(sanitizer)
	}
}

