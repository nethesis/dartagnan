/*
 * Copyright (C) 2018 Nethesis S.r.l.
 * http://www.nethesis.it - info@nethesis.it
 *
 * This file is part of Dartagnan project.
 *
 * Dartagnan is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License,
 * or any later version.
 *
 * Dartagnan is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with Dartagnan.  If not, see COPYING.
 *
 */

package notifications

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	gomail "gopkg.in/gomail.v2"

	"github.com/cbroglie/mustache"
	"github.com/nethesis/dartagnan/athos/configuration"
	"github.com/nethesis/dartagnan/athos/models"
	"github.com/nethesis/dartagnan/athos/utils"
	"github.com/nicksnyder/go-i18n/i18n"
)

func statusColor(status string) string {
	switch status {
	case "OK":
		return "#3f9c35"
	case "WARNING":
		return "#ec7a08"
	case "FAILURE":
		return "#cd0808"
	}
	return ""
}

func alertHtmlBody(alert models.Alert) string {
	values := make(map[string]string)
	T, _ := i18n.Tfunc("en-US", "en-US", "en-US")

	values["labelName"] = T("Name")
	values["labelAlertReport"] = T("Alert report")
	values["labelPublicIp"] = T("Public IP")
	values["labelDetails"] = T("Details")
	values["labelCustomer"] = T("Customer")
	values["labelIssue"] = T("Issue")
	values["labelStatus"] = T("Status")
	values["labelDocs"] = T("Docs")
	values["labelHelp"] = T("Help")
	values["labelAutoMail"] = T("This is an automatic mail, please do not reply")

	values["title"] = configuration.Config.Notifications.PortalTitle
	values["urlPortal"] = configuration.Config.Notifications.PortalUrl
	values["urlDocs"] = configuration.Config.Notifications.DocsUrl
	values["urlHelp"] = configuration.Config.Notifications.HelpUrl

	values["serverName"] = alert.System.UUID
	values["serverDetails"] = fmt.Sprintf(fmt.Sprintf("%s/servers/%d", configuration.Config.Notifications.PortalUrl, alert.System.ID))
	values["alertId"] = alert.NameI18n
	values["alertStatus"] = alert.Status
	values["alertStatusColor"] = statusColor(alert.Status)
	values["serverIp"] = alert.System.PublicIP

	output, err := mustache.RenderFile("templates/alert-template.mustache", values)
	if err != nil {
		fmt.Println(err)
	}

	return output
}

func alertTextBody(alert models.Alert) string {
	T, _ := i18n.Tfunc("en-US", "en-US", "en-US")

	text := fmt.Sprintf("***** %s ***** \n\n", configuration.Config.Notifications.PortalTitle)
	text += fmt.Sprintf("%s: %s\n", T("Issue"), alert.NameI18n)
	text += fmt.Sprintf("%s: %s\n", T("Status"), alert.Status)
	text += "--------------------------------------------------------------------\n"
	text += fmt.Sprintf("%s: %s\n", T("Public IP"), alert.System.PublicIP)
	text += fmt.Sprintf("%s: %s\n", T("Details"), fmt.Sprintf("%s/servers/%d", configuration.Config.Notifications.PortalUrl, alert.System.ID))
	text += "--------------------------------------------------------------------\n"

	return text
}

func expireHtmlBody(system models.System, period string) string {
	values := make(map[string]string)
	T, _ := i18n.Tfunc("en-US", "en-US", "en-US")

	values["hello"] = T("hello")
	values["year"] = strconv.Itoa(time.Now().Year())

	values["title"] = configuration.Config.Notifications.PortalTitle
	values["urlComm"] = configuration.Config.Notifications.CommunityUrl
	values["urlDocs"] = configuration.Config.Notifications.DocsUrl
	values["urlHelp"] = configuration.Config.Notifications.HelpUrl

	values["systemUUID"] = T("server_uuid")
	values["serverName"] = system.UUID
	values["serverDetails"] = fmt.Sprintf(fmt.Sprintf("%s/servers/%d", configuration.Config.Notifications.PortalUrl, system.ID))

	values["footerMessage"] = T("footer_message")

	renderFile := ""
	if system.Subscription.SubscriptionPlan.Code == "trial" {
		values["message"] = T(period + "_trial")
		values["serverDetailsMessage"] = T("server_link")

		renderFile = "templates/expiration-template-trial.mustache"
	} else {
		values["message"] = T(period+"_others", map[string]interface{}{
			"Subscription": system.Subscription.SubscriptionPlan.Name,
		})
		values["serverDetailsMessage"] = T("server_link_renew")
		values["ownerName"] = utils.GetBillingInfo(system.CreatorID).Name

		renderFile = "templates/expiration-template-others.mustache"
	}

	output, err := mustache.RenderFile(renderFile, values)
	if err != nil {
		fmt.Println(err)
	}

	return output
}

func expireTextBody(system models.System, period string) string {
	T, _ := i18n.Tfunc("en-US", "en-US", "en-US")

	text := fmt.Sprintf("***** %s ***** \n\n", configuration.Config.Notifications.PortalTitle)
	text += fmt.Sprintf("%s: %s\n", T("server_uuid"), system.UUID)
	text += fmt.Sprintf("%s: %s\n", T("subscription"), system.Subscription.SubscriptionPlan.Name)

	message := ""
	if system.Subscription.SubscriptionPlan.Code == "trial" {
		message = T(period + "_trial_plain")

	} else {
		message = T(period+"_others_plain", map[string]interface{}{
			"Subscription": system.Subscription.SubscriptionPlan.Name,
			"Name":         utils.GetBillingInfo(system.CreatorID).Name,
		})
	}

	text += "\n" + message
	text += fmt.Sprintf("\n%s\n", fmt.Sprintf(fmt.Sprintf("%s/servers/%d", configuration.Config.Notifications.PortalUrl, system.ID)))

	text += fmt.Sprintf("\n%s: %s\n", "Blog", configuration.Config.Notifications.CommunityUrl)
	text += fmt.Sprintf("%s: %s\n", "Support", strings.Split(configuration.Config.Notifications.HelpUrl, ":")[1])
	text += fmt.Sprintf("%s: %s\n", "Docs", configuration.Config.Notifications.DocsUrl)
	return text
}

func MailNotificationAlert(address string, alert models.Alert, isNew bool) bool {
	i18n.MustLoadTranslationFile("templates/en-US-alert.json")
	T, _ := i18n.Tfunc("en-US", "en-US", "en-US")
	status := true
	m := gomail.NewMessage()
	m.SetHeader("From", configuration.Config.Notifications.Email.From)
	m.SetHeader("To", address)
	subject := fmt.Sprintf("[%s][%s] %s", T("ALERT"), T(alert.Status), alert.NameI18n)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", alertTextBody(alert))
	if isNew {
		m.SetHeader("Message-Id", fmt.Sprintf("%d@dartagnan.project", alert.ID))
	} else {
		m.SetHeader("Message-Id", fmt.Sprintf("%d@dartagnan.project", time.Now().UnixNano()))
		m.SetHeader("References", fmt.Sprintf("%d@dartagnan.project", alert.ID))
	}

	m.AddAlternative("text/html", alertHtmlBody(alert))

	d := gomail.NewDialer(
		configuration.Config.Notifications.Email.SMTPHost,
		configuration.Config.Notifications.Email.SMTPPort,
		configuration.Config.Notifications.Email.SMTPUser,
		configuration.Config.Notifications.Email.SMTPPassword,
	)

	// send the email
	if err := d.DialAndSend(m); err != nil {
		status = false
	}

	return status
}

func MailNotificationExpire(address string, system models.System, period string) bool {
	i18n.MustLoadTranslationFile("templates/en-US-alert.json")
	T, _ := i18n.Tfunc("en-US", "en-US", "en-US")
	status := true
	m := gomail.NewMessage()
	m.SetHeader("From", configuration.Config.Notifications.Email.From)
	m.SetHeader("To", address)
	subject := fmt.Sprintf("[%s][%s] %s", T("expiration"), system.Subscription.SubscriptionPlan.Name, T(system.UUID))
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", expireTextBody(system, period))
	m.SetHeader("Message-Id", fmt.Sprintf("%d@dartagnan.project", time.Now().UnixNano()))

	m.AddAlternative("text/html", expireHtmlBody(system, period))

	d := gomail.NewDialer(
		configuration.Config.Notifications.Email.SMTPHost,
		configuration.Config.Notifications.Email.SMTPPort,
		configuration.Config.Notifications.Email.SMTPUser,
		configuration.Config.Notifications.Email.SMTPPassword,
	)

	// send the email
	if err := d.DialAndSend(m); err != nil {
		status = false
	}

	return status
}
