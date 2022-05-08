package main

import (
	"fmt"
	"os"
	"strconv"

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

type Application struct {
	smsClient  *sms.SMS
	policyName string
}

func NewApplication(smsClient *sms.SMS, policyName string) *Application {
	return &Application{
		smsClient:  smsClient,
		policyName: policyName,
	}
}

func (a *Application) getFilterComment(number int) (string, error) {
	body := sms.GetFilters{
		Profile: sms.Profile{Name: a.policyName},
		Filter:  []sms.Filter{{Number: strconv.Itoa(number)}},
	}
	f, err := a.smsClient.GetFilters(&body)
	if err != nil {
		return "", err
	}
	return f.Filter[0].Comment, nil
}

func (a *Application) processFilter(number int) error {
	comment, err := a.getFilterComment(number)
	if err != nil {
		return err
	}
	fmt.Println(comment)
	return nil
}

func main() {
	config()
	url := viper.GetString("URL")
	apiKey := viper.GetString("APIKey")
	insecureSkipVerify := viper.GetBool("SkipTLSVerify")
	auth := sms.NewAPIKeyAuthorization(apiKey)
	smsClient := sms.New(url, auth).SetInsecureSkipVerify(insecureSkipVerify)
	app := NewApplication(smsClient, "tmcheck")
	err := app.processFilter(51)
	fmt.Printf("Err: %s", err)
	if false {
		body := sms.GetFilters{
			Profile: sms.Profile{Name: "tbcheck"},
			Filter:  []sms.Filter{{Number: "51"}},
		}

		f, err := smsClient.GetFilters(&body)
		if err != nil {
			panic(err)
		}
		fmt.Println("Result:", f)
		fmt.Println("Result:", f.Filter[0].Name, f.Filter[0].Actionset)
	}
	if false {
		body := sms.SetFilters{
			Profile: sms.Profile{Name: "tbcheck"},
			Filter: []sms.Filter{{
				Number:  "51",
				Comment: "#TBC#",
			}},
		}
		err := smsClient.SetFilters(&body)
		if err != nil {
			panic(err)
		}
	}
	if false {
		s, err := smsClient.GetActionSetRefID("Block / Notify")
		if err != nil {
			panic(err)
		}
		fmt.Println("result", s)
	}
}
