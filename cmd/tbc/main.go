package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	//"github.com/spf13/viper"

	//"github.com/mpkondrashin/tbc/pkg/sms"
	"github.com/mpkondrashin/tbcheck/pkg/sms"
	"github.com/spf13/viper"
)

const TBCheckMarker = "#TBC#"

const (
	FirstFilterNumber = 1
	LastFilterNumber  = 60000
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
	smsClient                *sms.SMS
	profile                  string
	actionset                string
	actionsetRefID           string
	distributionPriority     sms.DistributionPiority
	distributionSegmentGroup string
}

func NewApplication(smsClient *sms.SMS, profile, actionset string) *Application {
	return &Application{
		smsClient:                smsClient,
		profile:                  profile,
		actionset:                actionset,
		actionsetRefID:           "unknown",
		distributionPriority:     sms.PriorityLow,
		distributionSegmentGroup: "unknown",
	}
}

func (a *Application) ConfigDistribution(
	distributionSegmentGroup string,
	distributionPriority sms.DistributionPiority) *Application {
	a.distributionSegmentGroup = distributionSegmentGroup
	a.distributionPriority = distributionPriority
	return a
}

func (a *Application) Run() error {
	a.GetActionSetRefIDs()
	for n := FirstFilterNumber; n <= LastFilterNumber; n++ {
		err := a.processFilter(n)
		if err != nil {
			return err
		}
	}
	if a.distributionSegmentGroup == "" {
		log.Printf("SegmentGroup is missing - skip profile distribution")
		return nil
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
		fmt.Printf("Filter #%d: %v\n", number, err)
		return nil
	}
	//fmt.Println(comment)
	if strings.Contains(comment, TBCheckMarker) {
		log.Printf("Filter #%d: \"%s\" marker found - skip\n", number, TBCheckMarker)
		return nil
	}
	err = a.updateFilter(number, comment)
	if err != nil {
		log.Printf("Filter #%d: %v\n", number, err)
		return nil
	}
	log.Printf("Filter #%d: done\n", number)
	return nil
}

func (a *Application) distributeProfile(segmentGroup string, priority sms.DistributionPiority) error {
	segmengGroupId, err := a.smsClient.GetSegmentGroupId(segmentGroup)
	if err != nil {
		return err
	}

	body := sms.Distribution{
		Profile: sms.Profile{
			Name: a.profile,
		},
		Priority: priority.String(),
		SegmentGroup: &sms.SegmentGroup{
			ID: segmengGroupId,
		},
	}
	err := a.smsClient.DistributeProfile(&body)
	fmt.Println("error:", err)
	return fmt.Errorf("distributeProfile(%s, %v): %w", segmentGroup, priority, err)
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
	auth := sms.NewAPIKeyAuthorization(apiKey)
	smsClient := sms.New(url, auth).SetInsecureSkipVerify(insecureSkipVerify)
	app := NewApplication(smsClient, profile, action)
	err := app.distributeProfile(distributionSegmentGroup, distributionPriority)
	if err != nil {
		log.Print(err)
	}
	return
	/*err := app.Run()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Done")*/
}
