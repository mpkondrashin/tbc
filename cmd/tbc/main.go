package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/viper"

	"github.com/mpkondrashin/tbcheck/pkg/sms"
)

const TBCheckMarker = "#TBC#"

const (
	FirstFilterNumber = 1
	LastFilterNumber  = 600000
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
	smsClient      *sms.SMS
	profile        string
	actionset      string
	actionsetRefID string
}

func NewApplication(smsClient *sms.SMS, profile, actionset string) *Application {
	return &Application{
		smsClient:      smsClient,
		profile:        profile,
		actionset:      actionset,
		actionsetRefID: "unknown",
	}
}

func (a *Application) Run() error {
	a.GetActionSetRefIDs()
	for n := FirstFilterNumber; n <= LastFilterNumber; n++ {
		err := a.processFilter(n)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *Application) GetActionSetRefIDs() error {
	var err error
	a.actionsetRefID, err = a.smsClient.GetActionSetRefID(a.actionset)
	return err
}

func (a *Application) getFilterComment(number int) (string, error) {
	//fmt.Println("actionset", a.actionset, "ref ID: ", a.actionsetRefID)
	body := sms.GetFilters{
		Profile: sms.Profile{Name: a.profile},
		Filter:  []sms.Filter{{Number: strconv.Itoa(number)}},
	}
	f, err := a.smsClient.GetFilters(&body)
	if err != nil {
		return "", err
	}
	//fmt.Println("Result for comment: ", f)
	return f.Filter[0].Comment, nil
}

func (a *Application) updateFilter(number int, comment string) error {
	if comment != "" {
		comment += "\n\n"
	}
	body := sms.SetFilters{
		Profile: sms.Profile{Name: a.profile},
		Filter: []sms.Filter{{
			Number:  strconv.Itoa(number),
			Comment: comment + TBCheckMarker,
			Actionset: &sms.Actionset{
				Refid: a.actionsetRefID,
			},
		}},
	}
	return a.smsClient.SetFilters(&body)
}

func (a *Application) processFilter(number int) error {
	comment, err := a.getFilterComment(number)
	if err != nil {
		return err
	}
	//fmt.Println(comment)
	if strings.Contains(comment, TBCheckMarker) {
		fmt.Printf("Filter #%d: \"%s\" marker found - skip\n", number, TBCheckMarker)
		return nil
	}
	err = a.updateFilter(number, comment)
	if err != nil {
		return fmt.Errorf("processing filter #%d: %w", number, err)
	}
	fmt.Printf("Filter #%d: done\n", number)
	return nil
}

func main() {
	config()
	url := viper.GetString("URL")
	apiKey := viper.GetString("APIKey")
	insecureSkipVerify := viper.GetBool("SkipTLSVerify")
	profile := viper.GetString("Profile")
	action := viper.GetString("Actionset")

	auth := sms.NewAPIKeyAuthorization(apiKey)
	smsClient := sms.New(url, auth).SetInsecureSkipVerify(insecureSkipVerify)
	app := NewApplication(smsClient, profile, action)
	err := app.Run()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Done")
}
