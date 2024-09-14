package user

import "chatgpt-web-new-go/router/base"

type Request struct {
	base.Page
	Type string `form:"type"`
}
