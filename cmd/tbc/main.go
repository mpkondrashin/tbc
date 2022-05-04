package main

import (
	"fmt"
	"os"

	"github.com/spf13/viper"

	"github.com/mpkondrashin/tbcheck/pkg/sms"
)

func config() {
	viper.SetConfigName("tbcheck")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/tbcheck/")
	viper.AddConfigPath("$HOME/.tbcheck")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Print(fmt.Errorf("Fatal error config file: %w \n", err))
		os.Exit(1)
	}
}

func main() {
	config()
	url := viper.GetString("url") // "http://127.0.0.1"
	auth := sms.NewAPIKeyAuthorization("ASDFA")
	sms := sms.New(url, auth)
	f, err := sms.GetFilters()
	if err != nil {
		panic(err)
	}
	fmt.Println(f)
}
