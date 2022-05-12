package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/mpkondrashin/tbcheck/pkg/sms"
)

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
	a.actionsetsWithBlock, err = a.smsClient.GetActionSetRefIDsForAction("DENY")
	if err != nil {
		return
	}
	if len(a.actionsetsWithBlock) == 0 {
		log.Println("No filters with \"DENY\" action found")
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
	if filter.Actionset == nil {
		log.Printf("Filter #%d: Has no ActionSet", number)
	} else {
		log.Printf("ActionSet: %v", filter.Actionset)
		if a.HasBlockingAction(filter.Actionset.Refid) {
			log.Printf("Filter #%d: Change ActionSet", number)
			f.Actionset = &sms.Actionset{
				Refid: a.actionsetRefID,
			}
		} else {
			log.Printf("Filter #%d: Does not have DENY action in its %s ActionSet",
				number, filter.Actionset.Name)
		}
	}
	body := sms.SetFilters{
		Profile: sms.Profile{Name: a.profile},
		Filter:  []sms.Filter{f},
	}
	return a.smsClient.SetFilters(&body)
}

func (a *Application) HasBlockingAction(id string) bool {
	for _, each := range a.actionsetsWithBlock {
		log.Println("HasBlockingAction", each, "==", id)
		if each == id {
			return true
		}
	}
	return false
}

func (a *Application) processFilter(number int) error {
	filters, err := a.getFilter(number)
	if err != nil {
		log.Printf("Filter #%d: %v\n", number, err)
		return nil
	}
	fmt.Println("FILTER", filters)
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
