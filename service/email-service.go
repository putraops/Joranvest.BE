package service

import (
	"fmt"
	"joranvest/commons"
	"joranvest/helper"
	"joranvest/repository"
	"os"
	"runtime"
	"strconv"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
)

type EmailService interface {
	SendEmailVerification(to []string, userId string) helper.Response
}

type emailService struct {
	emailRepository repository.EmailRepository
	helper.AppSession
}

func NewEmailService(emailRepository repository.EmailRepository) EmailService {
	return &emailService{
		emailRepository: emailRepository,
	}
}

func (service *emailService) SendEmailVerification(to []string, userId string) helper.Response {
	commons.Logger()
	err := godotenv.Load()
	log.Info(service.getCurrentFuncName())
	if err != nil {
		log.Error(service.getCurrentFuncName())
		log.Error("Failed to get SMTP Configuration")
		return helper.ServerResponse(false, "Failed to get SMTP Configuration", "", helper.EmptyObj{})
	}

	smtpHost := os.Getenv("CONFIG_SMTP_HOST")
	smtpPort, err := strconv.Atoi(os.Getenv("CONFIG_SMTP_PORT"))
	if err != nil {
		log.Error(service.getCurrentFuncName())
		log.Error("Failed to Convert Port")
		return helper.ServerResponse(false, "Failed to Convert Port", "", helper.EmptyObj{})
	}
	smtpSenderName := os.Getenv("CONFIG_SENDER_NAME_NO_REPLY")
	smtpEmail := os.Getenv("CONFIG_AUTH_EMAIL_NO_REPLY")
	smtpPassword := os.Getenv("CONFIG_AUTH_PASSWORD_NO_REPLY")
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
                                            <img class="joranvest-logo" src="https://joranvest.com/assets/img/logo-joranvest.png" alt="Joranvest"/>
                                        </p>
                                        <h1 class="text-center text-bold">Selamat datang di Joranvest</h1>
                                        <p class="text-center">Untuk menyelesaikan Registrasi akun Anda, Silahkan Verifikasi Email Anda dengan cara menekan tombol di bawah.</p>

                                        <p class="text-center">
                                            <a class="btn btn-primary text-white" href="`+os.Getenv("FRONTEND_URL")+`/register-verification/`+userId+`" target="_blank">Verifikasi Email</a>
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
		smtpHost,
		smtpPort,
		smtpEmail,
		smtpPassword,
	)

	errSend := dialer.DialAndSend(mailer)
	if err != nil {
		log.Error(service.getCurrentFuncName())
		log.Error(fmt.Sprintf("%v,", errSend))
		return helper.ServerResponse(false, fmt.Sprintf("%v,", errSend), fmt.Sprintf("%v,", errSend), helper.EmptyObj{})
	}
	return helper.ServerResponse(true, "Email Sent", "", helper.EmptyObj{})
}

func (service *emailService) getCurrentFuncName() string {
	pc, _, _, _ := runtime.Caller(1)
	return runtime.FuncForPC(pc).Name()
}
