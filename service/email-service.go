package service

import (
	"bytes"
	"database/sql"
	"fmt"
	"html/template"
	"joranvest/commons"
	"joranvest/dto"
	"joranvest/helper"
	"joranvest/models"
	"joranvest/models/entity_view_models"
	"joranvest/repository"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
	"gorm.io/gorm"
)

type EmailService interface {
	NewPayment(record entity_view_models.EntityPaymentView) helper.Result
	SendEmailVerification(email string, userId string) helper.Response
	SendEmailVerified(email string) helper.Response
	ResetPassword(user models.ApplicationUser) helper.Response
	SendWebinarInformationToParticipants(dto dto.SendWebinarInformationDto, participant entity_view_models.EntityWebinarRegistrationView)
}

type emailService struct {
	emailRepository        repository.EmailRepository
	emailLoggingRepository repository.EmailLoggingRepository
	helper.AppSession
	FRONTEND_URL  string
	DASHBOARD_URL string
	smtpHost      string
	smtpPort      string
	smtpUsername  string
	smtpPassword  string
	DB            *gorm.DB
}

func NewEmailService(db *gorm.DB) EmailService {
	return &emailService{
		DB:                     db,
		emailRepository:        repository.NewEmailRepository(db),
		emailLoggingRepository: repository.NewEmailLoggingRepository(db),
		FRONTEND_URL:           os.Getenv("FRONTEND_URL"),
		DASHBOARD_URL:          os.Getenv("DASHBOARD_URL"),
		smtpHost:               os.Getenv("CONFIG_SMTP_HOST"),
		smtpPort:               os.Getenv("CONFIG_SMTP_PORT"),
		smtpUsername:           os.Getenv("CONFIG_SMTP_USERNAME"),
		smtpPassword:           os.Getenv("CONFIG_SMTP_PASSWORD"),
	}
}

type MailInfo struct {
	Title      *string
	Recipient  *string
	ActionUrl  *string
	Product    *Product
	MailFooter MailFooter
}

type Product struct {
	ProductType *string
	ProductName *string
	TotalPrice  *float64
	Duration    *float64
	Status      *int
}

type MailFooter struct {
	Year int
}

type AddressHeader struct {
	Email string
	Name  string
}

func (service *emailService) SendEmailVerification(email string, userId string) helper.Response {
	var total = service.emailLoggingRepository.GetLastIntervalLogging(email, commons.MailTypeAccountVerification, commons.MaxSendEmailOneInterval)
	if total > commons.MaxSendEmailOneInterval {
		return helper.ServerResponse(false, "Email Already Sent", "", helper.EmptyObj{})
	}

	var to = []string{email}
	commons.Logger()
	err := godotenv.Load()
	if err != nil {
		log.Error(service.getCurrentFuncName())
		log.Error("Failed to get SMTP Configuration")
		return helper.ServerResponse(false, "Failed to get SMTP Configuration", "", helper.EmptyObj{})
	}

	smtpPort, err := strconv.Atoi(service.smtpPort)
	if err != nil {
		log.Error(service.getCurrentFuncName())
		log.Error("Failed to Convert Port")
		return helper.ServerResponse(false, "Failed to Convert Port", "", helper.EmptyObj{})
	}

	smtpSenderName := os.Getenv("CONFIG_SENDER_NAME_NO_REPLY")
	mailer := gomail.NewMessage()
	mailer.SetHeader("From", smtpSenderName)
	mailer.SetHeader("To", to...)
	mailer.SetHeader("Subject", "Verifikasi Email")
	mailer.SetBody("text/html", `<!doctype html>
        <html>
            <head>
                <meta name="viewport" content="width=device-width" />
                <meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
                <title>Joranvest - Verifikasi Email</title>
                <style>
                    body {
                        background-color:#f6f6f6;
                        font-family:sans-serif;
                        -webkit-font-smoothing:antialiased;
                        font-size:14px;
                        line-height:1.4;
                        letter-spacing: 2px;
                        margin:0;
                        padding:0;
                        -ms-text-size-adjust:100%;
                        -webkit-text-size-adjust:100%;
                    }
                    table {
                        border-collapse:separate;
                        mso-table-lspace:0pt;
                        mso-table-rspace:0pt;
                        width:100%;
                    }
                    table td {
                        font-family:sans-serif;
                        font-size:14px;
                        vertical-align:top; 
                    }
                    .body{
                        background-color:#f6f6f6;
                        width:100%; 
                    }
                    .container{
                        display:block;
                        margin:0 auto !important;
                        /* makes it centered */
                        max-width: 630px;
                        width: 630px;
                        padding: 10px;
                    }
                    .content{
                        box-sizing:border-box;
                        margin: 30px auto 30px auto;
                        max-width: 630px;
                        /*padding: 20px 50px 20px 50px;*/
                        display:flex;
                        flex-direction: column;
                        justify-content: center;
                        align-items: center;
                    }

                    #tbl-button {
                        margin-left: auto;
                        margin-right: auto;
                    }

                    .main{
                        background:#ffffff;
                        border-radius:5px;
                        width:100%;
                        border: 1px solid #dee2e6;
                    }
                    .body-wrapper{
                        box-sizing:border-box;
                        padding: 60px 30px 5px 30px; 
                    }

                    /* Typograpghy */
                    h1, h2, h3, h4{
                        color: #000000;
                        font-family: sans-serif;
                        font-weight: 400;
                        margin: 0;
                        margin-bottom: 30px; 
                    }
                    h1 {
                        font-size: 35px;
                        font-weight: 300;
                        text-align: center;
                    }
                    .no-reply {
                        font-size: 13px;    
                        color: #6c757d!important
                    }
                    
                    .text-bold {
                        font-weight: bold;
                    }
                    .text-center {
                        text-align: center !important;
                    }
                    .text-white {
                        color: white !important;
                    }

                    p,ul,ol{
                        font-family:sans-serif;
                        font-size:14px;
                        font-weight:normal;
                        margin:0;
                        margin-bottom:15px; 
                    }

                    p li,ul li,ol li{
                        list-style-position:inside;
                        margin-left:5px; 
                    }
                    a{
                        color:#3498db;
                        text-decoration:underline; 
                    }
                    

                    .btn {
                        background-color:#ffffff;
                        border:solid 1px #3498db;
                        border-radius:5px;
                        box-sizing:border-box;
                        color:#3498db;
                        cursor:pointer;
                        display:inline-block;
                        font-size:14px;
                        font-weight:bold;
                        margin:0;
                        padding:12px 25px;
                        text-align:center; 
                        text-decoration:none;
                        text-transform:capitalize; 
                    }
                
                    .btn-primary {
                        background-color:#3498db;
                        border-color:#3498db;
                        color: #ffffff; 
                    }
                    .btn-primary:hover{
                        background-color:#34495e !important; 
                    }
                    .btn-primary:hover{
                        background-color:#34495e !important;
                        border-color:#34495e !important; 
                    } 

                    .clear{
                        clear:both; 
                    }
                    hr{
                        border:0;
                        border-bottom:1px solid #f6f6f6;
                        margin:20px 0; 
                    }
                    .shadow{
                        box-shadow:0 2px 4px rgba(0,0,0,.075);
                    }
                    .joranvest-logo {
                        margin-right: auto;
                        margin-left: auto;
                        margin-bottom: 10px;
                        width: 45%;
                    }

                    @media only screen and (max-width:620px){
                        h1 {
                            font-size: 20px;
                        }
                        .no-reply {
                            font-size: 10px;    
                        }
                        .joranvest-logo {
                            width: 70%;
                        }

                        table[class=body] p,
                        table[class=body] ul,
                        table[class=body] ol,
                        table[class=body] td,
                        table[class=body] span,
                        table[class=body] a{
                            font-size:16px !important; 
                        }
                        table[class=body] .body-wrapper,
                        table[class=body] .article{
                            padding:10px !important; 
                        }
                        table[class=body] .content{
                            padding:0 !important; 
                        }
                        table[class=body] .container{
                            padding:0 !important;
                            width:100% !important; 
                        }
                        table[class=body] .main{
                            border-left-width:0 !important;
                            border-radius:0 !important;
                            border-right-width:0 !important; 
                        }
                        table[class=body] .btn table{
                            width:100% !important; 
                        }
                        table[class=body] .btn a{
                            width:100% !important; 
                        }
                        table[class=body] .img-responsive{
                            height:auto !important;
                            max-width:100% !important;
                            width:auto !important; 
                        }
                    }

                </style>
                </head>
                <body class="">
                <table role="presentation" border="0" cellpadding="0" cellspacing="0" class="body">
                <tr>
                    <td>&nbsp;</td>
                    <td class="container">
                    <div class="content">
                    <table role="presentation" class="main shadow">
                        <tr>
                            <td class="body-wrapper">
                                <table role="presentation" border="0" cellpadding="0" cellspacing="0">
                                    <tr>
                                    <td>
                                        <p class="text-center">
                                            <img class="joranvest-logo" src="https://joranvest.com/assets/img/logo.png" alt="Joranvest"/>
                                        </p>
                                        <h1 class="text-center text-bold">Selamat datang di Joranvest</h1>
                                        <p class="text-center">Untuk menyelesaikan Registrasi akun Anda, Silahkan Verifikasi Email Anda dengan cara menekan tombol di bawah.</p>

                                        <p class="text-center">
                                            <a class="btn btn-primary text-white" href="`+service.FRONTEND_URL+`/register-verification/`+userId+`" target="_blank">Verifikasi Email</a>
                                        </p>
                                    
                                        <hr style="margin-top: 30px;" />
                                        <p class="no-reply text-center">Email ini adalah email otomatis. Mohon untuk tidak membalas email ini.</p>
                                    </td>
                                    </tr>
                                </table>
                            </td>
                        </tr>
                        </table>
                        
                    </div>
                    </td>
                    <td>&nbsp;</td>
                </tr>
                </table>
            </body>
        </html>`)

	dialer := gomail.NewDialer(
		service.smtpHost,
		smtpPort,
		service.smtpUsername,
		service.smtpPassword,
	)

	var emailLoggingRecord models.EmailLogging
	emailLoggingRecord.Email = email
	emailLoggingRecord.MailType = commons.MailTypeAccountVerification
	emailLoggingRecord.LastSent = sql.NullTime{Time: time.Now(), Valid: true}
	res := service.emailLoggingRepository.Insert(emailLoggingRecord)

	if res.Status {
		errSend := dialer.DialAndSend(mailer)
		if err != nil {
			log.Error(service.getCurrentFuncName())
			log.Error(fmt.Sprintf("%v,", errSend))
			return helper.ServerResponse(false, fmt.Sprintf("%v,", errSend), fmt.Sprintf("%v,", errSend), helper.EmptyObj{})
		}
	}

	return helper.ServerResponse(true, "Email Sent", "", helper.EmptyObj{})
}

func (service *emailService) SendEmailVerified(email string) helper.Response {
	var to = []string{email}

	commons.Logger()
	err := godotenv.Load()
	if err != nil {
		log.Error(service.getCurrentFuncName())
		log.Error("Failed to get SMTP Configuration")
		return helper.ServerResponse(false, "Failed to get SMTP Configuration", "", helper.EmptyObj{})
	}

	smtpPort, err := strconv.Atoi(service.smtpPort)
	if err != nil {
		log.Error(service.getCurrentFuncName())
		log.Error("Failed to Convert Port")
		return helper.ServerResponse(false, "Failed to Convert Port", "", helper.EmptyObj{})
	}

	smtpSenderName := os.Getenv("CONFIG_SENDER_NAME_NO_REPLY")
	mailer := gomail.NewMessage()
	mailer.SetHeader("From", smtpSenderName)
	mailer.SetHeader("To", to...)
	// mailer.SetAddressHeader("Cc", "tralalala@gmail.com", "Tra Lala La")
	mailer.SetHeader("Subject", "Verifikasi Email")
	mailer.SetBody("text/html", `<!doctype html>
        <html>
            <head>
                <meta name="viewport" content="width=device-width" />
                <meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
                <title>Joranvest - Verifikasi Email</title>
                <style>
                    body {
                        background-color:#f6f6f6;
                        font-family:sans-serif;
                        -webkit-font-smoothing:antialiased;
                        font-size:14px;
                        line-height:1.4;
                        letter-spacing: 2px;
                        margin:0;
                        padding:0;
                        -ms-text-size-adjust:100%;
                        -webkit-text-size-adjust:100%;
                    }
                    table {
                        border-collapse:separate;
                        mso-table-lspace:0pt;
                        mso-table-rspace:0pt;
                        width:100%;
                    }
                    table td {
                        font-family:sans-serif;
                        font-size:14px;
                        vertical-align:top; 
                    }
                    .body{
                        background-color:#f6f6f6;
                        width:100%; 
                    }
                    .container{
                        display:block;
                        margin:0 auto !important;
                        /* makes it centered */
                        max-width: 700px;
                        width: 700px;
                        padding: 10px;
                    }
                    .content{
                        box-sizing:border-box;
                        margin: 30px auto 30px auto;
                        max-width: 700px;
                        /*padding: 20px 50px 20px 50px;*/
                        min-height: 50vh;
                        display:flex;
                        flex-direction: column;
                        justify-content: center;
                        align-items: center;
                    }

                    #tbl-button {
                        margin-left: auto;
                        margin-right: auto;
                    }

                    .main{
                        background:#ffffff;
                        border-radius:5px;
                        width:100%;
                        border: 1px solid #dee2e6;
                    }
                    .body-wrapper{
                        box-sizing:border-box;
                        padding: 60px 30px 5px 30px; 
                    }

                    /* Typograpghy */
                    h1, h2, h3, h4{
                        color: #000000;
                        font-family: sans-serif;
                        font-weight: 400;
                        margin: 0;
                        margin-bottom: 30px; 
                    }
                    h1 {
                        font-size: 35px;
                        font-weight: 300;
                        text-align: center;
                    }
                    .no-reply {
                        font-size: 13px;    
                        color: #6c757d!important
                    }
                    
                    .text-bold {
                        font-weight: bold;
                    }
                    .text-center {
                        text-align: center !important;
                    }
                    .text-white {
                        color: white !important;
                    }

                    p,ul,ol{
                        font-family:sans-serif;
                        font-size:14px;
                        font-weight:normal;
                        margin:0;
                        margin-bottom:15px; 
                    }

                    p li,ul li,ol li{
                        list-style-position:inside;
                        margin-left:5px; 
                    }
                    a{
                        color:#3498db;
                        text-decoration:underline; 
                    }
                    

                    .btn {
                        background-color:#ffffff;
                        border:solid 1px #3498db;
                        border-radius:5px;
                        box-sizing:border-box;
                        color:#3498db;
                        cursor:pointer;
                        display:inline-block;
                        font-size:14px;
                        font-weight:bold;
                        margin:0;
                        padding:12px 25px;
                        text-align:center; 
                        text-decoration:none;
                        text-transform:capitalize; 
                    }
                
                    .btn-primary {
                        background-color:#3498db;
                        border-color:#3498db;
                        color: #ffffff; 
                    }
                    .btn-primary:hover{
                        background-color:#34495e !important; 
                    }
                    .btn-primary:hover{
                        background-color:#34495e !important;
                        border-color:#34495e !important; 
                    } 

                    .clear{
                        clear:both; 
                    }
                    hr{
                        border:0;
                        border-bottom:1px solid #f6f6f6;
                        margin:20px 0; 
                    }
                    .shadow{
                        box-shadow:0 2px 4px rgba(0,0,0,.075);
                    }
                    .joranvest-logo {
                        margin-right: auto;
                        margin-left: auto;
                        margin-bottom: 10px;
                        width: 45%;
                    }

                    @media only screen and (max-width:620px){
                        h1 {
                            font-size: 20px;
                        }
                        .no-reply {
                            font-size: 10px;    
                        }
                        .joranvest-logo {
                            width: 70%;
                        }

                        table[class=body] p,
                        table[class=body] ul,
                        table[class=body] ol,
                        table[class=body] td,
                        table[class=body] span,
                        table[class=body] a{
                            font-size:16px !important; 
                        }
                        table[class=body] .body-wrapper,
                        table[class=body] .article{
                            padding:10px !important; 
                        }
                        table[class=body] .content{
                            padding:0 !important; 
                        }
                        table[class=body] .container{
                            padding:0 !important;
                            width:100% !important; 
                        }
                        table[class=body] .main{
                            border-left-width:0 !important;
                            border-radius:0 !important;
                            border-right-width:0 !important; 
                        }
                        table[class=body] .btn table{
                            width:100% !important; 
                        }
                        table[class=body] .btn a{
                            width:100% !important; 
                        }
                        table[class=body] .img-responsive{
                            height:auto !important;
                            max-width:100% !important;
                            width:auto !important; 
                        }
                    }

                </style>
                </head>
                <body class="">
                <table role="presentation" border="0" cellpadding="0" cellspacing="0" class="body">
                <tr>
                    <td>&nbsp;</td>
                    <td class="container">
                    <div class="content">
                    <table role="presentation" class="main shadow">
                        <tr>
                            <td class="body-wrapper">
                                <table role="presentation" border="0" cellpadding="0" cellspacing="0">
                                    <tr>
                                    <td>
                                        <p class="text-center">
                                            <img class="joranvest-logo" src="https://joranvest.com/assets/img/logo.png" alt="Joranvest"/>
                                        </p>
                                        <h1 class="text-center text-bold">Selamat Akun Kamu Telah Terverifikasi</h1>
                                        <p class="text-center">Tekan tombol dibawah ini untuk login ke dalam Aplikasi.</p>

                                        <p class="text-center">
                                            <a class="btn btn-primary text-white" href="`+service.FRONTEND_URL+`/login" target="_blank">Login</a>
                                        </p>
                                    
                                        <hr style="margin-top: 30px;" />
                                        <p class="no-reply text-center">Email ini adalah email otomatis. Mohon untuk tidak membalas email ini.</p>
                                    </td>
                                    </tr>
                                </table>
                            </td>
                        </tr>
                        </table>
                        
                    </div>
                    </td>
                    <td>&nbsp;</td>
                </tr>
                </table>
            </body>
        </html>`)

	dialer := gomail.NewDialer(
		service.smtpHost,
		smtpPort,
		service.smtpUsername,
		service.smtpPassword,
	)

	var emailLoggingRecord models.EmailLogging
	emailLoggingRecord.Email = email
	emailLoggingRecord.MailType = commons.MailTypeEmailVerified
	emailLoggingRecord.LastSent = sql.NullTime{Time: time.Now(), Valid: true}
	res := service.emailLoggingRepository.Insert(emailLoggingRecord)

	if res.Status {
		errSend := dialer.DialAndSend(mailer)
		if err != nil {
			log.Error(service.getCurrentFuncName())
			log.Error(fmt.Sprintf("%v,", errSend))
			return helper.ServerResponse(false, fmt.Sprintf("%v,", errSend), fmt.Sprintf("%v,", errSend), helper.EmptyObj{})
		}
	}

	return helper.ServerResponse(true, "Email Sent", "", helper.EmptyObj{})
}

func (service *emailService) ResetPassword(user models.ApplicationUser) helper.Response {
	title := "Reset Password"
	userFullname := fmt.Sprintf("%v %v", user.FirstName, user.LastName)
	actionUrl := fmt.Sprintf("%v/recover-password/%v/%v", service.FRONTEND_URL, user.Id, user.Email)

	result := SendEmail(
		"reset-password.html",
		"Reset Password",
		[]string{user.Email},
		nil,
		nil,
		MailInfo{
			Title:     &title,
			Recipient: &userFullname,
			ActionUrl: &actionUrl,
			MailFooter: MailFooter{
				Year: time.Now().Year(),
			},
		},
	)

	if result.Status {
		return helper.ServerResponse(true, "Email Sent", "", helper.EmptyObj{})
	} else {
		return helper.ServerResponse(false, result.Message, result.Message, nil)
	}
}

func (service *emailService) NewPayment(record entity_view_models.EntityPaymentView) helper.Result {
	title := "Pembayaran Baru"
	recipient := "Finance Team"
	actionUrl := fmt.Sprintf("%v/payment", service.DASHBOARD_URL)
	totalPrice := float64(record.Price) + float64(record.UniqueNumber)
	var duration *float64
	productType := ""
	productName := ""
	if record.WebinarId != "" {
		productType = "Webinar"
		productName = record.WebinarTitle
		duration = nil
	} else if record.MembershipId != "" {
		productType = "Membership"
		productName = record.MembershipName
		duration = record.MembershipDuration
	} else if record.ProductId != "" {
		productType = "JCS"
		productName = "JCS"
		duration = record.ProductDuration
	}

	var to []string

	// func (r roleMemberRepository) GetViewById(recordId string) helper.Result {
	var roleNotifications []entity_view_models.EntityRoleNotificationView
	service.DB.Find(&roleNotifications, "has_payment_notification = ?", true)

	if len(roleNotifications) > 0 {
		for _, item := range roleNotifications {
			var roleMembers []entity_view_models.EntityRoleMemberView
			service.DB.Find(&roleMembers, "role_id = ?", item.RoleId)

			if len(roleMembers) > 0 {
				for _, member := range roleMembers {
					to = append(to, member.ApplicationUserEmail)
				}
			}
		}
	}
	to = append(to, "putraops@gmail.com")

	res := SendEmail(
		"new-payment.html",
		title,
		to,
		nil,
		nil,
		MailInfo{
			Title:     &title,
			Recipient: &recipient,
			ActionUrl: &actionUrl,
			Product: &Product{
				ProductType: &productType,
				ProductName: &productName,
				TotalPrice:  &totalPrice,
				Duration:    duration,
			},
			MailFooter: MailFooter{
				Year: time.Now().Year(),
			},
		},
	)

	return res
}

func (service *emailService) SendWebinarInformationToParticipants(dto dto.SendWebinarInformationDto, participant entity_view_models.EntityWebinarRegistrationView) {

	//--- Participants
	var to = []string{participant.UserEmail}
	var subject = "Webinar #" + participant.WebinarTitle

	var webinarDate string
	webinarStartDate := strconv.Itoa(participant.WebinarStartDate.Time.Day()) + " " + helper.ConvertMonthNameENGtoID(participant.WebinarStartDate.Time.Month().String()) + " " + strconv.Itoa(participant.WebinarStartDate.Time.Year())
	webinarEndDate := strconv.Itoa(participant.WebinarEndDate.Time.Day()) + " " + helper.ConvertMonthNameENGtoID(participant.WebinarEndDate.Time.Month().String()) + " " + strconv.Itoa(participant.WebinarEndDate.Time.Year())

	webinarStartTime := strconv.Itoa(participant.WebinarStartDate.Time.Hour()) + "." + strconv.Itoa(participant.WebinarStartDate.Time.Minute())
	webinarEndTime := strconv.Itoa(participant.WebinarEndDate.Time.Hour()) + "." + strconv.Itoa(participant.WebinarEndDate.Time.Minute())

	if webinarStartDate == webinarEndDate {
		webinarDate = webinarStartDate
	} else if webinarStartDate == webinarEndDate {
		webinarDate = webinarStartDate + " - " + webinarEndDate
	}

	webinarDate += " | Pukul "
	if webinarStartTime == webinarEndTime {
		webinarDate += webinarStartTime
	} else {
		webinarDate += webinarStartTime + " - " + webinarEndTime
	}
	webinarDate += " WIB"

	commons.Logger()
	err := godotenv.Load()
	if err != nil {
		log.Error(service.getCurrentFuncName())
		log.Error(participant.WebinarTitle)
		log.Error("Failed to get SMTP Configuration")
	}

	smtpPort, err := strconv.Atoi(service.smtpPort)
	if err != nil {
		log.Error(service.getCurrentFuncName())
		log.Error(participant.WebinarTitle)
		log.Error("Failed to Convert Port")
	}
	smtpSenderName := os.Getenv("CONFIG_SENDER_NAME_NO_REPLY")

	mailer := gomail.NewMessage()
	mailer.SetHeader("From", smtpSenderName)
	mailer.SetHeader("To", to...)
	mailer.SetHeader("Subject", subject)
	mailer.SetBody("text/html", `<!doctype html>
        <html>
            <head>
                <meta name="viewport" content="width=device-width" />
                <meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
                <title>Joranvest - Webinar Information</title>
                <style>
                    body {
                        background-color:#f6f6f6;
                        font-family:sans-serif;
                        -webkit-font-smoothing:antialiased;
                        font-size:14px;
                        line-height:1.4;
                        letter-spacing: 2px;
                        margin:0;
                        padding:0;
                        -ms-text-size-adjust:100%;
                        -webkit-text-size-adjust:100%;
                    }
                    table {
                        border-collapse:separate;
                        mso-table-lspace:0pt;
                        mso-table-rspace:0pt;
                        width:100%;
                    }
                    table td {
                        font-family:sans-serif;
                        font-size:14px;
                        vertical-align:top; 
                    }
                    .body{
                        background-color:#f6f6f6;
                        width:100%; 
                    }
                    .container{
                        display:block;
                        margin:0 auto !important;
                        /* makes it centered */
                        max-width: 630px;
                        width: 630px;
                        padding: 10px;
                    }
                    .content-main{
                        box-sizing:border-box;
                        margin: 40px auto 0px auto;
                        max-width: 630px;
                        /*padding: 20px 50px 20px 50px;*/
                        min-height: 50vh;
                        display:flex;
                        flex-direction: column;
                        justify-content: center;
                        align-items: center;
                    }
                    .content-footer{
                        box-sizing: border-box;
                        margin-left: auto;
                        margin-right: auto;
                        max-width: 630px;
                        margin-bottom: 50px;
                    }

                    #tbl-button {
                        margin-left: auto;
                        margin-right: auto;
                    }

                    .main{
                        background:#ffffff;
                        border-radius: 0px;
                        width:100%;
                        border: 1px solid #dee2e6;
                    }
                    .body-wrapper{
                        box-sizing:border-box;
                        padding: 40px 30px 20px 30px; 
                    }

                    .footer {
                        background: #111111!important;
                        color: white !important;
                        border-radius: 0px;
                        width:100%;
                    }
                    .footer-wrapper{
                        box-sizing:border-box;
                        padding: 25px 30px 25px 30px; 
                    }

                    /* Typograpghy */
                    h1, h2, h3, h4{
                        color: #000000;
                        font-family: sans-serif;
                        font-weight: 400;
                        margin: 0;
                        margin-bottom: 30px; 
                    }
                    h1 {
                        font-size: 35px;
                        font-weight: 300;
                        text-align: center;
                    }
                    .no-reply {
                        font-size: 13px;    
                        color: #6c757d!important
                    }
                    
                    .text-bold {
                        font-weight: bold;
                    }
                    .text-center {
                        text-align: center !important;
                    }
                    .text-white {
                        color: white !important;
                    }
                    
                    .mb-0 {
                        margin-bottom: 0px !important;
                    }

                    p,ul,ol{
                        font-family:sans-serif;
                        font-size:14px;
                        font-weight:normal;
                        margin:0;
                        margin-bottom:15px; 
                    }

                    p li,ul li,ol li{
                        list-style-position:inside;
                        margin-left:5px; 
                    }
                    a{
                        color:#3498db;
                        text-decoration:underline; 
                    }
                    

                    .btn {
                        background-color:#ffffff;
                        border:solid 1px #3498db;
                        border-radius:5px;
                        box-sizing:border-box;
                        color:#3498db;
                        cursor:pointer;
                        display:inline-block;
                        font-size:14px;
                        font-weight:bold;
                        margin:0;
                        padding:12px 25px;
                        text-align:center; 
                        text-decoration:none;
                        text-transform:capitalize; 
                    }
                
                    .btn-primary {
                        background-color:#3498db;
                        border-color:#3498db;
                        color: #ffffff; 
                    }
                    .btn-primary:hover{
                        background-color:#34495e !important; 
                    }
                    .btn-primary:hover{
                        background-color:#34495e !important;
                        border-color:#34495e !important; 
                    } 

                    .clear{
                        clear:both; 
                    }
                    hr{
                        border:0;
                        border-bottom:1px solid #f6f6f6;
                        margin:20px 0; 
                    }
                    .shadow{
                        box-shadow:0 2px 4px rgba(0,0,0,.075);
                    }
                    .joranvest-logo {
                        margin-right: auto;
                        margin-left: auto;
                        margin-bottom: 5px;
                        width: 35%;
                    }

                    @media only screen and (max-width:620px){
                        h1 {
                            font-size: 20px;
                        }
                        .no-reply {
                            font-size: 10px;    
                        }
                        .joranvest-logo {
                            width: 70%;
                        }

                        table[class=body] p,
                        table[class=body] ul,
                        table[class=body] ol,
                        table[class=body] td,
                        table[class=body] span,
                        table[class=body] a{
                            font-size:16px !important; 
                        }
                        table[class=body] .body-wrapper,
                        table[class=body] .article{
                            padding:10px !important; 
                        }
                        table[class=body] .content{
                            padding:0 !important; 
                        }
                        table[class=body] .container{
                            padding:0 !important;
                            width:100% !important; 
                        }
                        table[class=body] .main{
                            border-left-width:0 !important;
                            border-radius:0 !important;
                            border-right-width:0 !important; 
                        }
                        table[class=body] .btn table{
                            width:100% !important; 
                        }
                        table[class=body] .btn a{
                            width:100% !important; 
                        }
                        table[class=body] .img-responsive{
                            height:auto !important;
                            max-width:100% !important;
                            width:auto !important; 
                        }
                    }

                </style>
                </head>
                <body class="">
                    <table role="presentation" border="0" cellpadding="0" cellspacing="0" class="body">
                        <tr>
                            <td>&nbsp;</td>
                            <td class="container" style="padding-bottom: 0px !important;">
                                <div class="content-main">
                                <table role="presentation" class="main shadow">
                                    <tr>
                                        <td class="body-wrapper">
                                            <table role="presentation" border="0" cellpadding="0" cellspacing="0">
                                                <tr>
                                                <td>
                                                    <p class="text-center">
                                                        <a href="https://joranvest.com">
                                                            <img class="joranvest-logo" src="https://joranvest.com/assets/img/logo.png" alt="Joranvest"/>
                                                        </a>
                                                    </p>
                                                    <hr />
                                                    
                                                    <p class="mb-0">Halo Kak `+participant.UserFirstName+`,</p>
                                                    <p>Terima Kasih telah bergabung bersama Webinar.</p>

                                                    <p><a href="#">
                                                        <img style="width: 100%"
                                                            src="https://api.joranvest.com/`+participant.WebinarFilepath+`" />
                                                    </a><p>

                                                    <p class="mb-0"><strong>Topic: </strong>LIVE WEBINAR: #`+participant.WebinarTitle+`</p>
                                                    <p><strong>Time: </strong> `+webinarDate+`</p>
                                                    
                                                    <p>Catat Informasinya!</p>
                                                    <p class="mb-0">Join Zoom Meeting</p>
                                                    <p>`+dto.MeetingUrl+`</p>
                                                        
                                                    <p class="mb-0">Meeting ID: `+dto.MeetingId+`</p>
                                                    <p>Passcode: `+dto.Passcode+` <span style="font-size: font-size: 12px; color: red; font-style: italic;">(Jangan lupa password ini)</span></p>

                                                    <p>Room mulai dibuka 30 menit sebelum acara dimulai dan pastikan menggunakan Nama yang sama dengan Akun di JORANVEST.<p>

                                                    <p>Sampai bertemu dan Terima Kasih banyak atas partisipasinya</p>
                                                    
                                                    <p class="mb-0">Salam Hangat,</p>
                                                    <span>Admin JORANVEST</span>
                                                </td>
                                                </tr>
                                            </table>
                                        </td>
                                    </tr>
                                    </table>
                                </div>
                            </td>
                            <td>&nbsp;</td>
                        </tr>
                        <tr>
                            <td>&nbsp;</td>
                            <td class="container"  style="padding-top: 0px !important;">
                                <div class="content-footer text-white">
                                <table role="presentation" class="footer shadow">
                                    <tr>
                                        <td class="footer-wrapper">
                                            <table role="presentation" border="0" cellpadding="0" cellspacing="0">
                                                <tr>
                                                    <td>
                                                        <div>
                                                            <p class="text-center">Temukan Kami</p>
                                                            <p class="text-center">
                                                                <a href="facebook.com/joranvest"><img style="margin-right: 2px; margin-left: 2px; width: 37px; height: 37px; border: 3px solid white; border-radius: 6px;" src="https://joranvest.com/assets/icons/icon-white-facebook.png" alt="Facebook" /></a>
                                                                <a href="instagram.com/joranvest"><img style="margin-right: 2px; margin-left: 2px; width: 37px; height: 37px; border: 3px solid white; border-radius: 6px" src="https://joranvest.com/assets/icons/icon-white-instagram.png" alt="Instagram" /></a>
                                                                <a href="twitter.com/joranvest"><img style="margin-right: 2px; margin-left: 2px; width: 37px; height: 37px; border: 3px solid white; border-radius: 6px" src="https://joranvest.com/assets/icons/icon-white-twitter.png" alt="Twitter" /></a>
                                                                <a href="api.whatsapp.com/send?phone=6281228822774"><img style="margin-right: 2px; margin-left: 2px; width: 37px; height: 37px; border: 3px solid white; border-radius: 6px" src="https://joranvest.com/assets/icons/icon-white-whatsapp.png" alt="Whatsapp Business" /></a>
                                                                <a href="t.me/joranvest"><img style="margin-right: 2px; margin-left: 2px; width: 37px; height: 37px; border: 3px solid white; border-radius: 6px" src="https://joranvest.com/assets/icons/icon-white-telegram.png" alt="Telegram" /></a>
                                                            </p>
                                                            <p class="mb-0 text-center">Copyright © `+strconv.Itoa(time.Now().Year())+` Joranvest</p>
                                                            <p class="mb-0 text-center">All rights reserved</p>
                                                        </div>
                                                    </td>
                                                </tr>
                                            </table>
                                        </td>
                                    </tr>
                                    </table>
                                </div>
                            </td>
                            <td>&nbsp;</td>
                        </tr>
                    </table>
                </body>
        </html>`)

	dialer := gomail.NewDialer(
		service.smtpHost,
		smtpPort,
		service.smtpUsername,
		service.smtpPassword,
	)

	errSend := dialer.DialAndSend(mailer)
	if err != nil {
		log.Error("Error Send Email....")
		log.Error(participant.WebinarTitle)
		log.Error(service.getCurrentFuncName())
		log.Error(fmt.Sprintf("%v,", errSend))
	}
}

func SendEmail(templateName string, subject string, to []string, cc *AddressHeader, bcc *AddressHeader, data interface{}) helper.Result {
	fmt.Println("Sending....")

	commons.Logger()
	var err error

	smtpHost := os.Getenv("CONFIG_SMTP_HOST")
	smtpPort, err := strconv.Atoi(os.Getenv("CONFIG_SMTP_PORT"))
	if err != nil {
		log.Error("Failed to Convert Port")
		return helper.StandartResult(false, "Failed to Convert Port", nil)
	}
	smtpSenderName := os.Getenv("CONFIG_SENDER_NAME_NO_REPLY")
	smtpUsername := os.Getenv("CONFIG_SMTP_USERNAME")
	smtpPassword := os.Getenv("CONFIG_SMTP_PASSWORD")

	t, err := template.ParseGlob("templates/email/*")
	if err != nil {
		log.Error("SendEmail:ParseGlob")
		log.Error(fmt.Sprintf("%v,", err))
		return helper.StandartResult(false, fmt.Sprintf("%v,", err), nil)
	}

	var tpl bytes.Buffer
	err = t.ExecuteTemplate(&tpl, templateName, data)
	if err != nil {
		log.Error("SendEmail:ExecuteTemplate")
		log.Error(fmt.Sprintf("%v,", err))
		return helper.StandartResult(false, fmt.Sprintf("%v,", err), nil)
	}
	result := tpl.String()

	mailer := gomail.NewMessage()
	mailer.SetHeader("From", smtpSenderName)
	mailer.SetHeader("To", to...)
	if cc != nil {
		mailer.SetAddressHeader("Cc", cc.Email, cc.Name)
	}
	if bcc != nil {
		mailer.SetAddressHeader("Bcc", bcc.Email, bcc.Name)
	}
	mailer.SetHeader("Subject", subject)
	mailer.SetBody("text/html", result)
	// m.Attach(templateFile) // attach whatever you want

	dialer := gomail.NewDialer(
		smtpHost,
		smtpPort,
		smtpUsername,
		smtpPassword,
	)

	err = dialer.DialAndSend(mailer)
	if err != nil {
		log.Error("SendEmail:DialAndSend")
		log.Error(fmt.Sprintf("%v,", err))
		return helper.StandartResult(false, fmt.Sprintf("%v,", err), nil)
	}

	fmt.Println("Sent")
	return helper.StandartResult(true, "Ok", nil)
}

func (service *emailService) getCurrentFuncName() string {
	pc, _, _, _ := runtime.Caller(1)
	return runtime.FuncForPC(pc).Name()
}
