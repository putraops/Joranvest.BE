package service

import (
	"joranvest/dto"
	"joranvest/models"
	entity_view_models "joranvest/models/view_models"
	"joranvest/repository"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/mashingan/smapping"
	"golang.org/x/crypto/bcrypt"
)

// var (
// 	userSession helper.UserSession = helper.NewSession(ctx)
// )

//AuthService is a contract about something that this service can do
type AuthService interface {
	VerifyCredential(username string, email string, password string) interface{}
	CreateUser(user dto.ApplicationUserRegisterDto) (models.ApplicationUser, error)
	GetByEmail(email string) models.ApplicationUser
	IsDuplicateEmail(email string) bool
}

type authService struct {
	appUserRepo repository.ApplicationUserRepository
	ctx         *gin.Context
}

//NewAuthService creates a new instance of AuthService
func NewAuthService(appUserRepo repository.ApplicationUserRepository) AuthService {
	return &authService{
		appUserRepo: appUserRepo,
	}
}

func (service *authService) VerifyCredential(username string, email string, password string) interface{} {
	res := service.appUserRepo.GetUserByUsernameOrEmail(username, email)
	if res == nil {
		return nil
	}
	if v, ok := res.(entity_view_models.EntityApplicationUserView); ok {
		comparedPassword := comparePassword(v.Password, []byte(password))
		if (v.Email == email || v.Username == username) && comparedPassword {
			return res
		}
		return false
	}
	return nil
}

func (service *authService) CreateUser(user dto.ApplicationUserRegisterDto) (models.ApplicationUser, error) {
	userToCreate := models.ApplicationUser{}
	errMap := smapping.FillStruct(&userToCreate, smapping.MapFields(&user))
	if errMap != nil {
		log.Fatalf("Failed map %v", errMap)
	}

	//
	// //asdasd := userSession.GetAppSession()
	// println("userSession.GetAppSession().EntityId")
	// session := sessions.Default(ctx)
	// println(fmt.Sprintf("%v", session.Get("EntityId")))

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
