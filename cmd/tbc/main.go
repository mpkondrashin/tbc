package main

import (
	"fmt"
	"log"
	"os"

	"github.com/mpkondrashin/tbcheck/pkg/sms"
	"github.com/spf13/viper"
)

const TBCheckMarker = "#TBC2#"

const (
	FirstFilterNumber = 51
	LastFilterNumber  = 1500
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
	log.Print("TBCheck started")
	config()
	url := viper.GetString("SMS.URL")
	log.Printf("SMS.URL: %s", url)
	apiKey := viper.GetString("SMS.APIKey")
	log.Printf("SMS.APIKey: %s", apiKey)
	insecureSkipVerify := viper.GetBool("SMS.SkipTLSVerify")
	profile := viper.GetString("Profile")
	action := viper.GetString("Actionset")
	/*
		distributionPriorityString := viper.GetString("Distribution.Priority")
		distributionSegmentGroup := viper.GetString("Distribution.SegmentGroup")
		distributionPriority := sms.PriorityLow
		if distributionPriorityString != "" {
			var err error
			distributionPriority, err = sms.DistributionPiorityFromString(distributionPriorityString)
			if err != nil {
				panic(err)
			}
		}
		log.Printf("distributionSegmentGroup = %s", distributionSegmentGroup)
	*/
	auth := sms.NewAPIKeyAuthorization(apiKey)
	smsClient := sms.New(url, auth).SetInsecureSkipVerify(insecureSkipVerify)
	err := smsClient.DownloadProfile("tbcheck", "tbcheck_d.pkg")
	fmt.Println(err)
	return
	//smsClient.DataDictionaryAll()
	//fmt.Println(smsClient.GetActionSet())
	//a, e := smsClient.DataDictionary("CATEGORY")
	//fmt.Println(e, "CATEGORY ", a)
	//return
	app := NewApplication(smsClient, profile, action)
	//err := app.distributeProfile(distributionSegmentGroup, distributionPriority)
	//if err != nil {
	//	log.Print(err)
	//}
	//return
	err := app.Run()
	if err != nil {
		log.Println(err)
	}
	log.Println("Done")
}
