package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

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
	actionsetsWithBlock      []string
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

func (a *Application) Run() (err error) {
	a.actionsetRefID, err = a.smsClient.GetActionSetRefID(a.actionset)
	if err != nil {
		return
	}
	a.actionsetsWithBlock, err = a.smsClient.GetActionSetRefIDsForAction("BLOCK")
	if err != nil {
		return
	}
	if len(a.actionsetsWithBlock) == 0 {
		log.Println("No filters with \"BLOCK\" action found")
		return nil
	}
	for n := FirstFilterNumber; n <= LastFilterNumber; n++ {
		err = a.processFilter(n)
		if err != nil {
			return
		}
	}
	if a.distributionSegmentGroup == "" {
		log.Printf("SegmentGroup is missing - skip profile distribution")
		return nil
	}
	return nil
}

func (a *Application) getFilter(number int) (*sms.Filters, error) {
	body := sms.GetFilters{
		Profile: sms.Profile{Name: a.profile},
		Filter:  []sms.Filter{{Number: strconv.Itoa(number)}},
	}
	filters, err := a.smsClient.GetFilters(&body)
	if err != nil {
		return nil, err
	}
	return filters, nil
}

func (a *Application) updateFilter(number int, filter *sms.Filter) error {
	comment := filter.Comment
	if comment != "" {
		comment += "\n\n"
	}
	f := sms.Filter{
		Number:  strconv.Itoa(number),
		Comment: comment + TBCheckMarker,
	}
	if a.HasBlockingAction(filter.Actionset.Refid) {
		log.Printf("Filter #%d: Change ActionSet", number)
		f.Actionset = &sms.Actionset{
			Refid: a.actionsetRefID,
		}
	} else {
		log.Printf("Filter #%d: Does not have BLOCK action", number)
	}
	body := sms.SetFilters{
		Profile: sms.Profile{Name: a.profile},
		Filter:  []sms.Filter{f},
	}
	return a.smsClient.SetFilters(&body)
}

func (a *Application) HasBlockingAction(id string) bool {
	for _, each := range a.actionsetsWithBlock {
		if each == id {
			return true
		}
	}
	return false
}

func (a *Application) processFilter(number int) error {
	filters, err := a.getFilter(number)
	if err != nil {
		fmt.Printf("Filter #%d: %v\n", number, err)
		return nil
	}
	filter := filters.Filter[0]
	comment := filter.Comment
	if strings.Contains(comment, TBCheckMarker) {
		log.Printf("Filter #%d: \"%s\" marker found - skip\n", number, TBCheckMarker)
		return nil
	}
	err = a.updateFilter(number, &filter)
	if err != nil {
		log.Printf("Filter #%d: %v\n", number, err)
		return nil
	}
	//log.Printf("Filter #%d: done\n", number)
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
	err = a.smsClient.DistributeProfile(&body)
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
