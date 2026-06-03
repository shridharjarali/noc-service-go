package notification

import (
	"log"
	"strings"

	"digit-oss/noc-services/internal/config"
	"digit-oss/noc-services/internal/domain"
	"digit-oss/noc-services/internal/repository/postgres"
	"digit-oss/noc-services/internal/service"
	"digit-oss/noc-services/internal/util"
)

// NOCNotificationService translates NOCNotificationService.java.
// Builds SMS requests and produces them to the Kafka SMS topic.
type NOCNotificationService struct {
	Cfg         *config.Config
	Producer    domain.Producer
	SvcRequest  *postgres.ServiceRequestRepository
	UserService *service.UserService
}

// Process creates and sends SMS notifications based on the NOC request.
func (n *NOCNotificationService) Process(nocReq *domain.NocRequest) {
	if !n.Cfg.IsSMSEnabled {
		return
	}

	smsRequests := n.enrichSMSRequest(nocReq)
	for _, sms := range smsRequests {
		if n.Producer != nil {
			_ = n.Producer.Push(n.Cfg.SmsNotifTopic, sms)
		}
		log.Printf("[Notification] MobileNumber: %s Message: %s", sms.MobileNumber, sms.Message)
	}
}

func (n *NOCNotificationService) enrichSMSRequest(nocReq *domain.NocRequest) []domain.SMSRequest {
	noc := nocReq.Noc

	message := n.getCustomizedMsg(noc)
	if message == "" {
		return nil
	}

	mobileToOwner := n.getUserList(nocReq)
	smsRequests := make([]domain.SMSRequest, 0, len(mobileToOwner))
	for mobile := range mobileToOwner {
		smsRequests = append(smsRequests, domain.SMSRequest{
			MobileNumber: mobile,
			Message:      message,
		})
	}
	return smsRequests
}

func (n *NOCNotificationService) getUserList(nocReq *domain.NocRequest) map[string]string {
	noc := nocReq.Noc
	criteria := domain.NocSearchCriteria{
		TenantID: noc.TenantID,
		OwnerIDs: []string{noc.AccountID},
	}
	resp, err := n.UserService.GetUser(criteria, nocReq.RequestInfo)
	if err != nil || len(resp.User) == 0 {
		return nil
	}
	return map[string]string{
		resp.User[0].MobileNumber: resp.User[0].Name,
	}
}

func (n *NOCNotificationService) getCustomizedMsg(noc *domain.Noc) string {
	var action string
	if noc.Workflow == nil {
		action = "null"
	} else {
		action = noc.Workflow.Action
	}
	messageCode := action + "_" + noc.ApplicationStatus

	switch messageCode {
	case util.ActionStatusCreated,
		util.ActionStatusInitiated,
		util.ActionStatusApproved,
		util.ActionStatusRejected:
		return n.getInitiatedMsg(noc, messageCode)
	default:
		return ""
	}
}

func (n *NOCNotificationService) getInitiatedMsg(noc *domain.Noc, code string) string {
	nocTypeLabel := "AAI"
	if strings.EqualFold(noc.NocType, util.FireNOCType) {
		nocTypeLabel = "Fire"
	}
	// Template: "{1} NOC application {2} has been <action>"
	return nocTypeLabel + " NOC application " + noc.ApplicationNo + " status: " + code
}
