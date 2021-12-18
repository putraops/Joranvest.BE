package service

import (
	"joranvest/dto"
	"joranvest/helper"
	"joranvest/models"
	entity_view_models "joranvest/models/view_models"
	"joranvest/repository"
	"log"

	"github.com/mashingan/smapping"
	"golang.org/x/crypto/bcrypt"
)

// var (
// 	userSession helper.UserSession = helper.NewSession(ctx)
// )

//AuthService is a contract about something that this service can do
type AuthService interface {
	VerifyCredential(username string, email string, password string) helper.Response
	CreateUser(user dto.ApplicationUserRegisterDto) (models.ApplicationUser, error)
	GetByEmail(email string) models.ApplicationUser
	IsDuplicateEmail(email string) bool
}

type authService struct {
	appUserRepo repository.ApplicationUserRepository
}

//NewAuthService creates a new instance of AuthService
func NewAuthService(appUserRepo repository.ApplicationUserRepository) AuthService {
	return &authService{
		appUserRepo: appUserRepo,
	}
}

func (service *authService) VerifyCredential(username string, email string, password string) helper.Response {
	res := service.appUserRepo.GetViewUserByEmail(username, email)
	if res == nil {
		return helper.ServerResponse(false, "Email yang Anda masukkan tidak terdaftar.", "", helper.EmptyObj{})
	}
	if v, ok := res.(entity_view_models.EntityApplicationUserView); ok {
		if !v.IsEmailVerified {
			return helper.ServerResponse(false, "Akun anda belum aktif. Silahkan periksa kembali Email untuk melakukan verifikasi.", "", helper.EmptyObj{})
		}

		comparedPassword := comparePassword(v.Password, []byte(password))
		if (v.Email == email || v.Username == username) && !comparedPassword {
			return helper.ServerResponse(false, "Password yang Anda masukkan salah", "", helper.EmptyObj{})
		}
	}
	return helper.ServerResponse(true, "Ok", "", helper.EmptyObj{})
}

func (service *authService) CreateUser(user dto.ApplicationUserRegisterDto) (models.ApplicationUser, error) {
	userToCreate := models.ApplicationUser{}
	errMap := smapping.FillStruct(&userToCreate, smapping.MapFields(&user))
	if errMap != nil {
		log.Fatalf("Failed map %v", errMap)
	}
	return service.appUserRepo.Insert(userToCreate)
}

func (service *authService) GetByEmail(email string) models.ApplicationUser {
	return service.appUserRepo.GetByEmail(email)
}

func (service *authService) IsDuplicateEmail(email string) bool {
	res := service.appUserRepo.IsDuplicateEmail(email)
	return !(res.Error == nil)
}

func comparePassword(hashedPwd string, plainPassword []byte) bool {
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPassword)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}
