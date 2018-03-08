/*
 * Copyright (C) 2018 Nethesis S.r.l.
 * http://www.nethesis.it - info@nethesis.it
 *
 * This file is part of Dartagnan project.
 *
 * Dartagnan is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License,
 * or any later version.
 *
 * Dartagnan is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with Dartagnan.  If not, see COPYING.
 *
 */

package notifications

import (
	"fmt"
	"time"

	gomail "gopkg.in/gomail.v2"

	"github.com/nicksnyder/go-i18n/i18n"
	"github.com/nethesis/dartagnan/athos/models"
	"github.com/nethesis/dartagnan/athos/configuration"
	"github.com/cbroglie/mustache"
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

	values["labelName"] =  T("Name")
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
	values["alertId"] =  alert.NameI18n
	values["alertStatus"] =  alert.Status
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

func MailNotification(address string, alert models.Alert, isNew bool) bool {
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

