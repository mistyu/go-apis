package global

import (
	ut "github.com/go-playground/universal-translator"
	"go-apis/user/config"
	"go-apis/user/proto"
)

var (
	Trans            ut.Translator
	ServerConfig     *config.ServerConfig = &config.ServerConfig{}
	UserServerClient proto.UserClient
)
