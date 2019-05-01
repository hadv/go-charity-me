package service

import "github.com/spf13/viper"

var signingKey = []byte(viper.GetString("SIGNING_KEY"))
