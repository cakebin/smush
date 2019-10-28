package email

import (
	"errors"
	"fmt"
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// ResetPWInfo is a convenience data structure for holding
// all relevant information for resetting a user's password
type ResetPWInfo struct {
	UserEmail string `json:"userEmail"`
	ResetURL  string `json:"resetUrl"`
}

/*---------------------------------
            Interface
----------------------------------*/

// PasswordEmailer describes all of the methods
// used for sending emails related to a user's password
type PasswordEmailer interface {
	SendResetPWEmail(resetPWInfo *ResetPWInfo) (bool, error)
}

// SendResetPWEmail sends an email to a user
// allowing them to reset their password
func (e *Email) SendResetPWEmail(resetPWInfo *ResetPWInfo) (bool, error) {
	from := mail.NewEmail("Cakebin", "cae@cakeforge.co")
	subject := "Reset your password for smush-tracker"
	to := mail.NewEmail("Smusher", resetPWInfo.UserEmail)

	resetPWBody := fmt.Sprintf(`
   <p>Hallo friend,</p>

   <p>We received a request to reset your password. You can do so by visiting the following link:</p>
	 
	 <a href="%s">%s</a>

   <p>If you did not initate this request, you can safely ignore it as it will expire shortly.</p>

   <p>Keep smushing! :)</p>
   `, resetPWInfo.ResetURL, resetPWInfo.ResetURL)
	content := mail.NewContent("text/html", resetPWBody)

	m := mail.NewV3MailInit(from, subject, to, content)

	request := sendgrid.GetRequest(
		os.Getenv("SENDGRID_API_KEY"),
		"/v3/mail/send",
		"https://api.sendgrid.com",
	)
	request.Method = "POST"
	request.Body = mail.GetRequestBody(m)

	response, err := sendgrid.API(request)

	if response.StatusCode != 202 {
		return false, errors.New(response.Body)
	}

	if err != nil {
		return false, err
	}

	return true, nil
}
