package config

import (
	"fmt"
	"joranvest/helper"
	"joranvest/models"
	entity_view_models "joranvest/models/entity_view_models"
	"joranvest/models/view_models"
	"joranvest/repository"
	"os"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type configService struct {
	applicationUserRepository repository.ApplicationUserRepository
}

func (db *configService) setupAdministrator() (bool, error) {
	var superAdmin = models.ApplicationUser{
		FirstName: "Administrator",
		LastName:  "",
		Username:  "sys_admin",
		Password:  "joranvest",
		Address:   "",
		Phone:     "",
		IsAdmin:   true,
		IsDefault: true,
		Email:     "admin@joranvest.com",
	}
	newRecord, err := db.applicationUserRepository.Insert(superAdmin)
	if err != nil {
		return false, err
	} else {
		return newRecord.IsActive, nil
	}
}

//-- Setup and Open Database Connection
func SetupDatabaseConnection() *gorm.DB {
	fmt.Println("=======================================================")
	fmt.Println("Open Connection....")
	fmt.Println("=======================================================")
	fmt.Println("Open env file....")
	fmt.Println("=======================================================")
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Failed to open env file!")
		fmt.Println("=======================================================")
		panic("")
	}

	//dbDial := os.Getenv("DB_DIAL")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s",
		dbHost,
		dbUser,
		dbPassword,
		dbName,
		dbPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	fmt.Println("Connecting....")
	if err != nil {
		fmt.Println("Failed to connect to database!")
		panic("=======================================================")
	} else {
		fmt.Println("Connected...")
	}

	//-- This function to generate model to database table
	db.AutoMigrate(
		&models.Entity{},
		&models.ApplicationUser{},
		&models.Team{},
		&models.Role{},
		&models.TeamRole{},
		&models.TeamMember{},
		&models.RoleMember{},
		&models.ApplicationMenuCategory{},
		&models.ApplicationMenu{},
		&models.RoleMenu{},
		&models.Organization{},

		&models.EmitenCategory{},
		&models.Emiten{},
		&models.Sector{},
		&models.Tag{},
		&models.ArticleCategory{},
		&models.Article{},
		&models.ArticleTag{},
		&models.Membership{},
		&models.MembershipUser{},
		&models.Payment{},
		&models.Filemaster{},

		&models.WebinarCategory{},
		&models.Webinar{},
		&models.WebinarSpeaker{},
		&models.WebinarRegistration{},
		&models.WebinarRecording{},

		&models.TechnicalAnalysis{},
		&models.FundamentalAnalysis{},
		&models.FundamentalAnalysisTag{},

		&models.RatingMaster{},

		&models.EducationCategory{},
		&models.Education{},
		&models.EducationPlaylist{},
		&models.EducationPlaylistUser{},

		&models.EmailLogging{},
		&models.EmailBlacklist{},

		&models.Order{},
	)

	viewList := make(map[string]map[string]string)

	var vw_application_user = entity_view_models.EntityApplicationUserView{}
	viewList[vw_application_user.TableName()] = vw_application_user.Migration()

	var vw_team = entity_view_models.EntityTeamView{}
	viewList[vw_team.TableName()] = vw_team.Migration()
	var vw_team_member = entity_view_models.EntityTeamMemberView{}
	viewList[vw_team_member.TableName()] = vw_team_member.Migration()
	var vw_role = entity_view_models.EntityRoleView{}
	viewList[vw_role.TableName()] = vw_role.Migration()
	var vw_team_role = entity_view_models.EntityTeamRoleView{}
	viewList[vw_team_role.TableName()] = vw_team_role.Migration()

	var vw_role_member = entity_view_models.EntityRoleMemberView{}
	viewList[vw_role_member.TableName()] = vw_role_member.Migration()

	var vw_organization = entity_view_models.EntityOrganizationView{}
	viewList[vw_organization.TableName()] = vw_organization.Migration()

	var vw_application_menu_category = entity_view_models.EntityApplicationMenuCategoryView{}
	viewList[vw_application_menu_category.TableName()] = vw_application_menu_category.Migration()

	var vw_application_menu = entity_view_models.EntityApplicationMenuView{}
	viewList[vw_application_menu.TableName()] = vw_application_menu.Migration()

	var vw_role_menu = entity_view_models.EntityRoleMenuView{}
	viewList[vw_role_menu.TableName()] = vw_role_menu.Migration()

	var vw_order = entity_view_models.EntityMembershipView{}
	viewList[vw_order.TableName()] = vw_order.Migration()

	//-- Article, Article Category, Article Tag
	var vw_article_category = entity_view_models.EntityArticleCategoryView{}
	viewList[vw_article_category.TableName()] = vw_article_category.Migration()
	var vw_article = entity_view_models.EntityArticleView{}
	viewList[vw_article.TableName()] = vw_article.Migration()
	var vw_article_tag = entity_view_models.EntityArticleTagView{}
	viewList[vw_article_tag.TableName()] = vw_article_tag.Migration()

	var vw_membership = entity_view_models.EntityMembershipView{}
	viewList[vw_membership.TableName()] = vw_membership.Migration()
	var vw_membership_user = entity_view_models.EntityMembershipUserView{}
	viewList[vw_membership_user.TableName()] = vw_membership_user.Migration()
	var vw_payment = entity_view_models.EntityPaymentView{}
	viewList[vw_payment.TableName()] = vw_payment.Migration()

	var vw_emiten_category = entity_view_models.EntityEmitenCategoryView{}
	viewList[vw_emiten_category.TableName()] = vw_emiten_category.Migration()

	var vw_emiten = entity_view_models.EntityEmitenView{}
	viewList[vw_emiten.TableName()] = vw_emiten.Migration()

	var vw_sector = entity_view_models.EntitySectorView{}
	viewList[vw_sector.TableName()] = vw_sector.Migration()

	var vw_filemaster = entity_view_models.EntityFilemasterView{}
	viewList[vw_filemaster.TableName()] = vw_filemaster.Migration()

	var vw_webinar_category = entity_view_models.EntityWebinarCategoryView{}
	viewList[vw_webinar_category.TableName()] = vw_webinar_category.Migration()
	var vw_webinar = entity_view_models.EntityWebinarView{}
	viewList[vw_webinar.TableName()] = vw_webinar.Migration()
	var vw_webinar_speaker = entity_view_models.EntityWebinarSpeakerView{}
	viewList[vw_webinar_speaker.TableName()] = vw_webinar_speaker.Migration()
	var vw_webinar_speaker_review = view_models.WebinarSpeakerReviewViewModel{}
	viewList[vw_webinar_speaker_review.TableName()] = vw_webinar_speaker_review.Migration()
	var vw_webinar_registration = entity_view_models.EntityWebinarRegistrationView{}
	viewList[vw_webinar_registration.TableName()] = vw_webinar_registration.Migration()
	var vw_webinar_recording = entity_view_models.EntityWebinarRecordingView{}
	viewList[vw_webinar_recording.TableName()] = vw_webinar_recording.Migration()

	var vw_tag = entity_view_models.EntityTagView{}
	viewList[vw_tag.TableName()] = vw_tag.Migration()

	var vw_technical_analysis = entity_view_models.EntityTechnicalAnalysisView{}
	viewList[vw_technical_analysis.TableName()] = vw_technical_analysis.Migration()

	var vw_fundamental_analysis = entity_view_models.EntityFundamentalAnalysisView{}
	viewList[vw_fundamental_analysis.TableName()] = vw_fundamental_analysis.Migration()
	var vw_fundamental_analysis_tag = entity_view_models.EntityFundamentalAnalysisTagView{}
	viewList[vw_fundamental_analysis_tag.TableName()] = vw_fundamental_analysis_tag.Migration()

	//-- Education:: Module Pembelajaran & Recording
	var vw_education_category = entity_view_models.EntityEducationCategoryView{}
	viewList[vw_education_category.TableName()] = vw_education_category.Migration()
	var vw_education = entity_view_models.EntityEducationView{}
	viewList[vw_education.TableName()] = vw_education.Migration()
	var vw_education_playlist = entity_view_models.EntityEducationPlaylistView{}
	viewList[vw_education_playlist.TableName()] = vw_education_playlist.Migration()
	var vw_education_playlist_user = entity_view_models.EntityEducationPlaylistUserView{}
	viewList[vw_education_playlist_user.TableName()] = vw_education_playlist_user.Migration()

	var vw_email_logging = entity_view_models.EntityEmailLoggingView{}
	viewList[vw_email_logging.TableName()] = vw_email_logging.Migration()
	var vw_email_blacklist = entity_view_models.EntityEmailBlacklistView{}
	viewList[vw_email_blacklist.TableName()] = vw_email_blacklist.Migration()

	var vw_rating_master = entity_view_models.EntityRatingMasterView{}
	viewList[vw_rating_master.TableName()] = vw_rating_master.Migration()

	if len(viewList) > 0 {
		for _, detail := range viewList {
			db.Exec(fmt.Sprintf("CREATE OR REPLACE VIEW %s AS %s", detail["view_name"], detail["query"]))
		}
	}

	var entityRecord = models.Entity{}
	if err := db.Where("Name = ?", "Joranvest").First(&entityRecord).Error; err != nil {
		fmt.Println("Creating Administrator Started")
		entityRecord.Id = uuid.New().String()
		entityRecord = models.Entity{
			Name:        "Joranvest",
			Description: "Joranvest",
		}
		entityRecord.Id = uuid.New().String()
		db.Create(&entityRecord)
		fmt.Println("New Entity has been created.")
	}

	var superAdmin = models.ApplicationUser{}
	if err := db.Where("Username = ?", os.Getenv("DEFAULT_ADMINISTRATOR_USERNAME")).First(&superAdmin).Error; err != nil {
		fmt.Println("Creating Administrator Started")
		superAdmin.Id = uuid.New().String()
		superAdmin = models.ApplicationUser{
			FirstName: "System",
			LastName:  "Administrator",
			Username:  os.Getenv("DEFAULT_ADMINISTRATOR_USERNAME"),
			Address:   "",
			Phone:     "",
			IsAdmin:   true,
			Email:     "admin@joranvest.com",
		}
		superAdmin.Id = uuid.New().String()
		superAdmin.CreatedBy = superAdmin.Id
		superAdmin.Password = helper.HashAndSalt([]byte("joranvest"))
		superAdmin.EntityId = entityRecord.Id
		db.Create(&superAdmin)
		fmt.Println("Finished and Enjoy.")
	}
	return db
}

//-- Close Database Connection
func CloseDatabaseConnection(db *gorm.DB) {
	dbSQL, err := db.DB()
	if err != nil {
		panic("Failed to close connection")
	}
	dbSQL.Close()
}
