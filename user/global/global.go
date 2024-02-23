package global

import (
	ut "github.com/go-playground/universal-translator"
	"go-apis/user/config"
)

var (
	Trans        ut.Translator
	ServerConfig *config.ServerConfig = &config.ServerConfig{}
)
