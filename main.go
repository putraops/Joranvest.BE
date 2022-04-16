package main

import (
	"os"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	// _ "github.com/swaggo/gin-swagger/example/basic/docs"

	"joranvest/docs"
	"joranvest/middleware"

	"gorm.io/gorm"

	"joranvest/config"
	"joranvest/controllers"
	"joranvest/repository"
	"joranvest/service"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
)

var (
	ctx *gin.Context
	db  *gorm.DB = config.SetupDatabaseConnection()

	// #region Configuration
	applicationUserRepository         repository.ApplicationUserRepository         = repository.NewApplicationUserRepository(db)
	applicationMenuCategoryRepository repository.ApplicationMenuCategoryRepository = repository.NewApplicationMenuCategoryRepository(db)
	applicationMenuRepository         repository.ApplicationMenuRepository         = repository.NewApplicationMenuRepository(db)
	membershipRepository              repository.MembershipRepository              = repository.NewMembershipRepository(db)
	membershipUserRepository          repository.MembershipUserRepository          = repository.NewMembershipUserRepository(db)
	filemasterRepository              repository.FilemasterRepository              = repository.NewFilemasterRepository(db)
	emitenRepository                  repository.EmitenRepository                  = repository.NewEmitenRepository(db)
	emitenCategoryRepository          repository.EmitenCategoryRepository          = repository.NewEmitenCategoryRepository(db)
	articleCategoryRepository         repository.ArticleCategoryRepository         = repository.NewArticleCategoryRepository(db)
	articleRepository                 repository.ArticleRepository                 = repository.NewArticleRepository(db)
	articleTagRepository              repository.ArticleTagRepository              = repository.NewArticleTagRepository(db)
	webinarCategoryRepository         repository.WebinarCategoryRepository         = repository.NewWebinarCategoryRepository(db)
	webinarRepository                 repository.WebinarRepository                 = repository.NewWebinarRepository(db)
	webinarSpeakerRepository          repository.WebinarSpeakerRepository          = repository.NewWebinarSpeakerRepository(db)
	webinarRegistrationRepository     repository.WebinarRegistrationRepository     = repository.NewWebinarRegistrationRepository(db)
	sectorRepository                  repository.SectorRepository                  = repository.NewSectorRepository(db)
	tagRepository                     repository.TagRepository                     = repository.NewTagRepository(db)
	technicalAnalysisRepository       repository.TechnicalAnalysisRepository       = repository.NewTechnicalAnalysisRepository(db)
	fundamentalAnalysisRepository     repository.FundamentalAnalysisRepository     = repository.NewFundamentalAnalysisRepository(db)
	fundamentalAnalysisTagRepository  repository.FundamentalAnalysisTagRepository  = repository.NewFundamentalAnalysisTagRepository(db)
	roleMenuRepository                repository.RoleMenuRepository                = repository.NewRoleMenuRepository(db)
	organizationRepository            repository.OrganizationRepository            = repository.NewOrganizationRepository(db)
	ratingMasterRepository            repository.RatingMasterRepository            = repository.NewRatingMasterRepository(db)
	paymentRepository                 repository.PaymentRepository                 = repository.NewPaymentRepository(db)

	jwtService                     service.JWTService                     = service.NewJWTService()
	authService                    service.AuthService                    = service.NewAuthService(applicationUserRepository)
	applicationMenuCategoryService service.ApplicationMenuCategoryService = service.NewApplicationMenuCategoryService(applicationMenuCategoryRepository)
	applicationMenuService         service.ApplicationMenuService         = service.NewApplicationMenuService(applicationMenuRepository)
	membershipService              service.MembershipService              = service.NewMembershipService(membershipRepository)
	membershipUserService          service.MembershipUserService          = service.NewMembershipUserService(membershipUserRepository)
	filemasterService              service.FilemasterService              = service.NewFilemasterService(filemasterRepository, webinarRepository, applicationUserRepository, organizationRepository)
	emitenService                  service.EmitenService                  = service.NewEmitenService(emitenRepository)
	emitenCategoryService          service.EmitenCategoryService          = service.NewEmitenCategoryService(emitenCategoryRepository)
	articleCategoryService         service.ArticleCategoryService         = service.NewArticleCategoryService(articleCategoryRepository)
	articleService                 service.ArticleService                 = service.NewArticleService(articleRepository)
	articleTagService              service.ArticleTagService              = service.NewArticleTagService(articleTagRepository)
	webinarCategoryService         service.WebinarCategoryService         = service.NewWebinarCategoryService(webinarCategoryRepository)
	webinarService                 service.WebinarService                 = service.NewWebinarService(webinarRepository)
	webinarSpeakerService          service.WebinarSpeakerService          = service.NewWebinarSpeakerService(webinarSpeakerRepository)
	sectorService                  service.SectorService                  = service.NewSectorService(sectorRepository)
	tagService                     service.TagService                     = service.NewTagService(tagRepository)
	technicalAnalysisService       service.TechnicalAnalysisService       = service.NewTechnicalAnalysisService(technicalAnalysisRepository)
	fundamentalAnalysisService     service.FundamentalAnalysisService     = service.NewFundamentalAnalysisService(fundamentalAnalysisRepository)
	fundamentalAnalysisTagService  service.FundamentalAnalysisTagService  = service.NewFundamentalAnalysisTagService(fundamentalAnalysisTagRepository)
	roleMenuService                service.RoleMenuService                = service.NewRoleMenuService(roleMenuRepository)
	organizationService            service.OrganizationService            = service.NewOrganizationService(organizationRepository)
	ratingMasterService            service.RatingMasterService            = service.NewRatingMasterService(ratingMasterRepository)

	authController                    controllers.AuthController                    = controllers.NewAuthController(db, authService, jwtService)
	applicationUserController         controllers.ApplicationUserController         = controllers.NewApplicationUserController(db, jwtService)
	applicationMenuCategoryController controllers.ApplicationMenuCategoryController = controllers.NewApplicationMenuCategoryController(applicationMenuCategoryService, jwtService)
	applicationMenuController         controllers.ApplicationMenuController         = controllers.NewApplicationMenuController(applicationMenuService, jwtService)
	membershipController              controllers.MembershipController              = controllers.NewMembershipController(membershipService, jwtService)
	membershipUserController          controllers.MembershipUserController          = controllers.NewMembershipUserController(membershipUserService, jwtService)
	filemasterController              controllers.FilemasterController              = controllers.NewFilemasterController(filemasterService, jwtService)
	emitenController                  controllers.EmitenController                  = controllers.NewEmitenController(emitenService, jwtService)
	emitenCategoryController          controllers.EmitenCategoryController          = controllers.NewEmitenCategoryController(emitenCategoryService, jwtService)
	educationCategoryController       controllers.EducationCategoryController       = controllers.NewEducationCategoryController(db, jwtService)
	educationController               controllers.EducationController               = controllers.NewEducationController(db, jwtService)
	articleCategoryController         controllers.ArticleCategoryController         = controllers.NewArticleCategoryController(articleCategoryService, jwtService)
	articleController                 controllers.ArticleController                 = controllers.NewArticleController(articleService, jwtService)
	articleTagController              controllers.ArticleTagController              = controllers.NewArticleTagController(articleTagService, jwtService)
	webinarCategoryController         controllers.WebinarCategoryController         = controllers.NewWebinarCategoryController(webinarCategoryService, jwtService)
	webinarController                 controllers.WebinarController                 = controllers.NewWebinarController(webinarService, jwtService)
	webinarSpeakerController          controllers.WebinarSpeakerController          = controllers.NewWebinarSpeakerController(webinarSpeakerService, jwtService)
	webinarRegistrationController     controllers.WebinarRegistrationController     = controllers.NewWebinarRegistrationController(db, jwtService)
	webinarRecordingController        controllers.WebinarRecordingController        = controllers.NewWebinarRecordingController(db, jwtService)
	sectorController                  controllers.SectorController                  = controllers.NewSectorController(sectorService, jwtService)
	tagController                     controllers.TagController                     = controllers.NewTagController(tagService, jwtService)
	technicalAnalysisController       controllers.TechnicalAnalysisController       = controllers.NewTechnicalAnalysisController(technicalAnalysisService, jwtService)
	fundamentalAnalysisController     controllers.FundamentalAnalysisController     = controllers.NewFundamentalAnalysisController(fundamentalAnalysisService, jwtService)
	fundamentalAnalysisTagController  controllers.FundamentalAnalysisTagController  = controllers.NewFundamentalAnalysisTagController(fundamentalAnalysisTagService, jwtService)
	roleController                    controllers.RoleController                    = controllers.NewRoleController(db, jwtService)
	roleMenuController                controllers.RoleMenuController                = controllers.NewRoleMenuController(roleMenuService, jwtService)
	roleMemberController              controllers.RoleMemberController              = controllers.NewRoleMemberController(db, jwtService)
	organizationController            controllers.OrganizationController            = controllers.NewOrganizationController(organizationService, jwtService)
	ratingMasterController            controllers.RatingMasterController            = controllers.NewRatingMasterController(ratingMasterService, jwtService)
	paymentController                 controllers.PaymentController                 = controllers.NewPaymentController(db, jwtService)
	productController                 controllers.ProductController                 = controllers.NewProductController(db, jwtService)
	emailController                   controllers.EmailController                   = controllers.NewEmailController(db, jwtService)

	//clientest coreapi.ClientTest = coreapi.NewClientTest()
	// #endregion
)

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	defer config.CloseDatabaseConnection(db)
	status, SERVER_NAME, SERVER_PORT := GetServerDefinition()
	if !status {
		panic("Mo Server_NAME OR SERVER_PORT found on .env")
	}
	//getCardToken()

	// r := gin.New()
	r := gin.Default()

	// programatically set swagger info
	docs.SwaggerInfo.Title = "Joranvest API Documentation"
	// docs.SwaggerInfo.Description = "This is a full API documentation of Joranvest."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = *SERVER_NAME + ":" + *SERVER_PORT
	docs.SwaggerInfo.BasePath = "/api"

	// r.Use(location.New(location.Config{
	// 	Scheme:  "http",
	// 	Host:    "foo.com",
	// 	Base:    "/base",
	// 	Headers: location.Headers{Scheme: "X-Forwarded-Proto", Host: "X-Forwarded-For"},
	// }))

	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store))

	//r.Use(CORSMiddleware())
	r.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		// AllowOrigins:  []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Accept", "Access-Control-Allow-Origin", "Origin", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		// AllowOriginFunc: func(origin string) bool {
		// 	return origin == "https://github.com"
		// },
		MaxAge: 12 * time.Hour,
	}))

	r.Static("/assets", "./assets")
	r.Static("/upload", "./upload")
	r.Static("/script", "./templates/js")
	// r.HTMLRender = createMyRender("templates/views/")

	// #region User Web View
	r.GET("/", func(c *gin.Context) {
		c.HTML(
			http.StatusOK,
			"index",
			gin.H{
				"title": "Joranvest",
				"err":   "",
				"data":  nil,
			},
		)
	})
	// #endregion

	// #region API
	authRoutes := r.Group("api/auth")
	{
		authRoutes.POST("/login", authController.Login)
		authRoutes.GET("/logout", authController.Logout)
	}

	emitenApiRoutes := r.Group("api/emiten")
	{
		emitenApiRoutes.POST("/getDatatables", emitenController.GetDatatables)
		emitenApiRoutes.POST("/getPagination", emitenController.GetPagination)
		emitenApiRoutes.GET("/lookup", emitenController.Lookup)
		emitenApiRoutes.POST("/emitenLookup", emitenController.EmitenLookup)
		emitenApiRoutes.POST("/save", emitenController.Save)
		emitenApiRoutes.POST("/patchingEmiten", emitenController.PatchingEmiten)
		emitenApiRoutes.GET("/getById/:id", emitenController.GetById)
		emitenApiRoutes.DELETE("/deleteById/:id", emitenController.DeleteById)
	}

	educationCategoryApiRoutes := r.Group("api")
	{
		educationCategoryWithAuthRoutes := educationCategoryApiRoutes.Group("/auth/education_category", middleware.AuthorizeJWT(jwtService))
		{
			educationCategoryWithAuthRoutes.POST("/save", educationCategoryController.Save)
			educationCategoryWithAuthRoutes.DELETE("/deleteById/:id", educationCategoryController.DeleteById)
		}
		educationCategoryWithoutAuthRoutes := educationCategoryApiRoutes.Group("/education_category")
		{
			educationCategoryWithoutAuthRoutes.POST("/getPagination", educationCategoryController.GetPagination)
			educationCategoryWithoutAuthRoutes.GET("/lookup", educationCategoryController.Lookup)
			educationCategoryWithoutAuthRoutes.GET("/getById/:id", educationCategoryController.GetById)
		}
	}

	educationApiRoutes := r.Group("api")
	{
		educationWithAuthRoutes := educationApiRoutes.Group("/auth/education", middleware.AuthorizeJWT(jwtService))
		{
			educationWithAuthRoutes.POST("/save", educationController.Save)
			educationWithAuthRoutes.GET("/submit/:id", educationController.Submit)
			educationWithAuthRoutes.POST("/markVideoAsWatched", educationController.MarkVideoAsWatched)
			educationWithAuthRoutes.POST("/uploadEducationCover/:id", educationController.UploadEducationCover)
			educationWithAuthRoutes.POST("/addToPlaylist", educationController.AddToPlaylist)
			educationWithAuthRoutes.DELETE("/deleteById/:id", educationController.DeleteById)
			educationWithAuthRoutes.DELETE("/removeFromPlaylistById/:id", educationController.RemoveFromPlaylistById)
		}
		educationWithoutAuthRoutes := educationApiRoutes.Group("/education")
		{
			educationWithoutAuthRoutes.POST("/getPagination", educationController.GetPagination)
			educationWithoutAuthRoutes.POST("/getPlaylist", educationController.GetPlaylist)
			educationWithoutAuthRoutes.GET("/getPlaylistByUserId/:education_id/:application_user_id", educationController.GetPlaylistByUserId)
			educationWithoutAuthRoutes.GET("/lookup", educationController.Lookup)
			educationWithoutAuthRoutes.GET("/getByPathUrl/:path_url", educationController.GetByPathUrl)
			educationWithoutAuthRoutes.GET("/getById/:id", educationController.GetById)
			educationWithoutAuthRoutes.GET("/getViewById/:id", educationController.GetViewById)
			educationWithoutAuthRoutes.GET("/getPlaylistById/:id", educationController.GetPlaylistById)
		}
	}

	applicationMenuCategoryApiRoutes := r.Group("api/application_menu_category")
	{
		applicationMenuCategoryApiRoutes.POST("/getDatatables", applicationMenuCategoryController.GetDatatables)
		applicationMenuCategoryApiRoutes.GET("/lookup", applicationMenuCategoryController.Lookup)
		applicationMenuCategoryApiRoutes.POST("/save", applicationMenuCategoryController.Save)
		applicationMenuCategoryApiRoutes.GET("/getById/:id", applicationMenuCategoryController.GetById)
		applicationMenuCategoryApiRoutes.DELETE("/deleteById/:id", applicationMenuCategoryController.DeleteById)
	}

	applicationMenuApiRoutes := r.Group("api/application_menu")
	{
		applicationMenuApiRoutes.GET("/getTree", applicationMenuController.GetTree)
		applicationMenuApiRoutes.GET("/getTreeByRoleId/:roleId", applicationMenuController.GetTreeByRoleId)
		applicationMenuApiRoutes.GET("/lookup", applicationMenuController.Lookup)
		applicationMenuApiRoutes.POST("/save", applicationMenuController.Save)
		applicationMenuApiRoutes.POST("/orderTree", applicationMenuController.OrderTree)
		applicationMenuApiRoutes.GET("/getById/:id", applicationMenuController.GetById)
		applicationMenuApiRoutes.DELETE("/deleteById/:id", applicationMenuController.DeleteById)
	}

	emitenCategoryApiRoutes := r.Group("api/emiten_category")
	{
		emitenCategoryApiRoutes.POST("/getDatatables", emitenCategoryController.GetDatatables)
		emitenCategoryApiRoutes.POST("/getPagination", emitenCategoryController.GetPagination)
		emitenCategoryApiRoutes.GET("/lookup", emitenCategoryController.Lookup)
		emitenCategoryApiRoutes.POST("/save", emitenCategoryController.Save)
		emitenCategoryApiRoutes.GET("/getById/:id", emitenCategoryController.GetById)
		emitenCategoryApiRoutes.DELETE("/deleteById/:id", emitenCategoryController.DeleteById)
	}

	articleCategoryApiRoutes := r.Group("api/article_category")
	{
		articleCategoryApiRoutes.POST("/getDatatables", articleCategoryController.GetDatatables)
		articleCategoryApiRoutes.GET("/getTree", articleCategoryController.GetTree)
		articleCategoryApiRoutes.GET("/lookup", articleCategoryController.Lookup)
		articleCategoryApiRoutes.POST("/save", articleCategoryController.Save)
		articleCategoryApiRoutes.GET("/getById/:id", articleCategoryController.GetById)
		articleCategoryApiRoutes.DELETE("/deleteById/:id", articleCategoryController.DeleteById)
	}

	articleApiRoutes := r.Group("api/article")
	{
		articleApiRoutes.POST("/getDatatables", articleController.GetDatatables)
		articleApiRoutes.POST("/getPagination", articleController.GetPagination)
		articleApiRoutes.POST("/save", articleController.Save)
		articleApiRoutes.POST("/submit/:id", articleController.Submit)
		articleApiRoutes.GET("/getById/:id", articleController.GetById)
		articleApiRoutes.GET("/getViewById/:id", articleController.GetViewById)
		articleApiRoutes.GET("/getArticleCoverById/:id", articleController.GetArticleCoverById)
		articleApiRoutes.DELETE("/deleteById/:id", articleController.DeleteById)
	}

	articleTagApiRoutes := r.Group("api/article_tag")
	{
		articleTagApiRoutes.GET("/getById/:id", articleTagController.GetById)
		articleTagApiRoutes.GET("/getAll", articleTagController.GetAll)
	}

	webinarCategoryApiRoutes := r.Group("api/webinar_category")
	{
		webinarCategoryApiRoutes.POST("/getDatatables", webinarCategoryController.GetDatatables)
		webinarCategoryApiRoutes.GET("/getTreeParent", webinarCategoryController.GetTreeParent)
		webinarCategoryApiRoutes.GET("/getTree", webinarCategoryController.GetTree)
		webinarCategoryApiRoutes.POST("/orderTree", webinarCategoryController.OrderTree)
		webinarCategoryApiRoutes.GET("/lookup", webinarCategoryController.Lookup)
		webinarCategoryApiRoutes.POST("/save", webinarCategoryController.Save)
		webinarCategoryApiRoutes.GET("/getById/:id", webinarCategoryController.GetById)
		webinarCategoryApiRoutes.DELETE("/deleteById/:id", webinarCategoryController.DeleteById)
	}

	webinarApiRoutes := r.Group("api/webinar")
	{
		webinarApiRoutes.POST("/getDatatables", webinarController.GetDatatables)
		webinarApiRoutes.POST("/getPagination", webinarController.GetPagination)
		webinarApiRoutes.POST("/getPaginationRegisteredByUser/:user_id", webinarController.GetPaginationRegisteredByUser)
		webinarApiRoutes.POST("/save", webinarController.Save)
		webinarApiRoutes.POST("/submit/:id", webinarController.Submit)
		webinarApiRoutes.POST("/uploadWebinarCover/:id", filemasterController.UploadWebinarCover)
		webinarApiRoutes.GET("/getById/:id", webinarController.GetById)
		webinarApiRoutes.GET("/getViewById/:id", webinarController.GetViewById)
		webinarApiRoutes.GET("/getWebinarWithRatingByUserId/:webinar_id/:user_id", webinarController.GetWebinarWithRatingByUserId)
		webinarApiRoutes.DELETE("/deleteById/:id", webinarController.DeleteById)
	}

	webinarRecordingApiRoutes := r.Group("api")
	{
		webinarRecordingWithAuthRoutes := webinarRecordingApiRoutes.Group("/webinar_recording", middleware.AuthorizeJWT(jwtService))
		{
			webinarRecordingWithAuthRoutes.POST("/save", webinarRecordingController.Save)
			webinarRecordingWithAuthRoutes.GET("/submit/:id", webinarRecordingController.Submit)
			webinarRecordingWithAuthRoutes.DELETE("/deleteById/:id", webinarRecordingController.DeleteById)
		}
		webinarRecordingWithoutAuthRoutes := webinarRecordingApiRoutes.Group("/webinar_recording")
		{
			webinarRecordingWithoutAuthRoutes.POST("/getPagination", webinarRecordingController.GetPagination)
			webinarRecordingWithoutAuthRoutes.GET("/getByPathUrl/:path_url", webinarRecordingController.GetByPathUrl)
			webinarRecordingWithoutAuthRoutes.GET("/getById/:id", webinarRecordingController.GetById)
			webinarRecordingWithoutAuthRoutes.GET("/getViewById/:id", webinarRecordingController.GetViewById)
			webinarRecordingWithoutAuthRoutes.GET("/getByWebinarId/:webinarId", webinarRecordingController.GetByWebinarId)
		}
	}

	webinarSpeakerApiRoutes := r.Group("api/webinar_speaker")
	{
		webinarSpeakerApiRoutes.POST("/save", webinarSpeakerController.Save)
		webinarSpeakerApiRoutes.GET("/getById/:id", webinarSpeakerController.GetById)
		webinarSpeakerApiRoutes.GET("/getAll", webinarSpeakerController.GetAll)
		webinarSpeakerApiRoutes.GET("/getSpeakersRatingByWebinarId/:webinarId", webinarSpeakerController.GetSpeakersRatingByWebinarId)
		webinarSpeakerApiRoutes.GET("/getSpeakerReviewById/:id", webinarSpeakerController.GetSpeakerReviewById)
	}

	webinarRegistrationApiRoutes := r.Group("api/webinar_registration")
	{
		webinarRegistrationApiRoutes.POST("/getDatatables", webinarRegistrationController.GetDatatables)
		webinarRegistrationApiRoutes.POST("/getPagination", webinarRegistrationController.GetPagination)
		webinarRegistrationApiRoutes.GET("/getById/:id", webinarRegistrationController.GetById)
		webinarRegistrationApiRoutes.GET("/getViewById/:id", webinarRegistrationController.GetViewById)
		webinarRegistrationApiRoutes.POST("/sendInvitation", webinarRegistrationController.SendInvitation)
		webinarRegistrationApiRoutes.POST("/sendMeetingInformation", webinarRegistrationController.SendMeetingInformation)
		webinarRegistrationApiRoutes.POST("/save", webinarRegistrationController.Save)
		webinarRegistrationApiRoutes.GET("/isWebinarRegistered/:id", webinarRegistrationController.IsWebinarRegistered)
		webinarRegistrationApiRoutes.DELETE("/deleteById/:id", webinarRegistrationController.DeleteById)
	}

	tagApiRoutes := r.Group("api/tag")
	{
		tagApiRoutes.POST("/getDatatables", tagController.GetDatatables)
		tagApiRoutes.POST("/getPagination", tagController.GetPagination)
		tagApiRoutes.GET("/lookup", tagController.Lookup)
		tagApiRoutes.POST("/save", tagController.Save)
		tagApiRoutes.GET("/getById/:id", tagController.GetById)
		tagApiRoutes.DELETE("/deleteById/:id", tagController.DeleteById)
	}

	sectorApiRoutes := r.Group("api/sector")
	{
		sectorApiRoutes.POST("/getPagination", sectorController.GetPagination)
		sectorApiRoutes.GET("/lookup", sectorController.Lookup)
		sectorApiRoutes.POST("/save", sectorController.Save)
		sectorApiRoutes.GET("/getById/:id", sectorController.GetById)
		sectorApiRoutes.DELETE("/deleteById/:id", sectorController.DeleteById)
	}

	membershipApiRoutes := r.Group("api")
	{
		membershipWithAuthRoutes := membershipApiRoutes.Group("/membership", middleware.AuthorizeJWT(jwtService))
		{
			membershipWithAuthRoutes.POST("/save", membershipController.Save)
			membershipWithAuthRoutes.DELETE("/deleteById/:id", membershipController.DeleteById)
		}
		membershipWithoutAuthRoutes := membershipApiRoutes.Group("/membership")
		{
			membershipWithoutAuthRoutes.GET("/lookup", membershipController.Lookup)
			membershipWithoutAuthRoutes.POST("/getDatatables", membershipController.GetDatatables)
			membershipWithoutAuthRoutes.POST("/getPagination", membershipController.GetPagination)
			membershipWithoutAuthRoutes.GET("/getAll", membershipController.GetAll)
			membershipWithoutAuthRoutes.POST("/setRecommendation", membershipController.SetRecommendation)
			membershipWithoutAuthRoutes.GET("/getById/:id", membershipController.GetById)
			membershipWithoutAuthRoutes.GET("/getViewById/:id", membershipController.GetViewById)
		}
	}

	membershipUserApiRoutes := r.Group("api/membership_user")
	{
		membershipUserApiRoutes.POST("/getPagination", membershipUserController.GetPagination)
		membershipUserApiRoutes.GET("/getAll", membershipUserController.GetAll)
		membershipUserApiRoutes.POST("/save", membershipUserController.Save)
		membershipUserApiRoutes.GET("/getById/:id", membershipUserController.GetById)
		membershipUserApiRoutes.GET("/getByUserLogin", membershipUserController.GetByUserLogin)
		membershipUserApiRoutes.DELETE("/deleteById/:id", membershipUserController.DeleteById)
	}

	filemasterApiRoutes := r.Group("api/filemaster")
	{
		filemasterApiRoutes.POST("/single_upload/:id", filemasterController.SingleUpload)
		filemasterApiRoutes.POST("/uploadByType/:module/:filetype/:id", filemasterController.UploadByType)
		filemasterApiRoutes.POST("/uploadPDFDocuments/:module/:id", filemasterController.UploadPDFDocuments)
		filemasterApiRoutes.POST("/uploadProfilePicture/:id", filemasterController.UploadProfilePicture)
		filemasterApiRoutes.POST("/upload/:id", filemasterController.Insert)
		filemasterApiRoutes.GET("/getAll", filemasterController.GetAll)
		filemasterApiRoutes.POST("/getAllByRecordIds", filemasterController.GetAllByRecordIds)
		filemasterApiRoutes.DELETE("/deleteById/:id", filemasterController.DeleteById)
		filemasterApiRoutes.DELETE("/deleteByRecordId/:recordId", filemasterController.DeleteByRecordId)
	}

	applicationUserApiRoutes := r.Group("api/application_user")
	{
		applicationUserApiRoutes.POST("/getPagination", applicationUserController.GetPagination)
		applicationUserApiRoutes.GET("/lookup", applicationUserController.Lookup)
		applicationUserApiRoutes.POST("/userLookup", applicationUserController.UserLookup)
		applicationUserApiRoutes.POST("/changePhone", applicationUserController.ChangePhone)
		applicationUserApiRoutes.POST("/changePassword", applicationUserController.ChangePassword)
		applicationUserApiRoutes.POST("/updateProfile", applicationUserController.UpdateProfile)
		applicationUserApiRoutes.POST("/recoverPassword", applicationUserController.RecoverPassword)
		applicationUserApiRoutes.POST("/resetPassword", applicationUserController.ResetPassword)
		applicationUserApiRoutes.POST("/register", authController.Register)
		applicationUserApiRoutes.GET("/getViewById/:id", applicationUserController.GetViewById)
		applicationUserApiRoutes.PATCH("/emailVerificationById/:id", applicationUserController.EmailVerificationById)
	}

	technicalAnalysisApiRoutes := r.Group("api/technical_analysis")
	{
		technicalAnalysisApiRoutes.POST("/getDatatables", technicalAnalysisController.GetDatatables)
		technicalAnalysisApiRoutes.POST("/getPagination", technicalAnalysisController.GetPagination)
		technicalAnalysisApiRoutes.POST("/save", technicalAnalysisController.Save)
		technicalAnalysisApiRoutes.GET("/getById/:id", technicalAnalysisController.GetById)
		technicalAnalysisApiRoutes.DELETE("/deleteById/:id", technicalAnalysisController.DeleteById)
	}

	fundamentalAnalysisApiRoutes := r.Group("api/fundamental_analysis")
	{
		fundamentalAnalysisApiRoutes.POST("/getDatatables", fundamentalAnalysisController.GetDatatables)
		fundamentalAnalysisApiRoutes.POST("/getPagination", fundamentalAnalysisController.GetPagination)
		fundamentalAnalysisApiRoutes.POST("/save", fundamentalAnalysisController.Save)
		fundamentalAnalysisApiRoutes.POST("/submit/:id", fundamentalAnalysisController.Submit)
		fundamentalAnalysisApiRoutes.GET("/getById/:id", fundamentalAnalysisController.GetById)
		fundamentalAnalysisApiRoutes.DELETE("/deleteById/:id", fundamentalAnalysisController.DeleteById)
	}

	fundamentalAnalysisTagApiRoutes := r.Group("api/fundamental_analysis_tag")
	{
		fundamentalAnalysisTagApiRoutes.GET("/getById/:id", fundamentalAnalysisTagController.GetById)
		fundamentalAnalysisTagApiRoutes.GET("/getAll", fundamentalAnalysisTagController.GetAll)
	}

	// roleApiRoutes := r.Group("api/role", middleware.DBTransactionMiddleware(db))
	// {
	// 	roleApiRoutes.POST("/getPagination", roleController.GetPagination)
	// 	roleApiRoutes.POST("/save", roleController.Save)
	// 	roleApiRoutes.GET("/getById/:id", roleController.GetById)
	// 	roleApiRoutes.DELETE("/deleteById/:id", roleController.DeleteById)
	// }

	roleApiRoutes := r.Group("api")
	{
		roleAuthRoutes := roleApiRoutes.Group("/role", middleware.AuthorizeJWT(jwtService))
		{
			roleAuthRoutes.POST("/save", roleController.Save)
			roleAuthRoutes.POST("/set/paymentNotification", roleController.SetPaymentNotification)
			roleAuthRoutes.POST("/set/dashboardAccess", roleController.SetDashboardAccess)
			roleAuthRoutes.POST("/set/fullAccess", roleController.SetFullAccess)
			roleAuthRoutes.GET("/notification/getByRoleId/:roleId", roleController.GetNotificationByRoleId)
			roleAuthRoutes.DELETE("/deleteById/:id", roleController.DeleteById)
		}
		roleNoAuthRoutes := roleApiRoutes.Group("/role")
		{
			// roleNoAuthRoutes.GET("/lookup", roleController.Lookup)
			roleNoAuthRoutes.POST("/getPagination", roleController.GetPagination)
			roleNoAuthRoutes.GET("/getAll", roleController.GetAll)
			roleNoAuthRoutes.GET("/getById/:id", roleController.GetById)
			roleNoAuthRoutes.GET("/getViewById/:id", roleController.GetViewById)
		}
	}

	roleMemberApiRoutes := r.Group("api")
	{
		roleMemberApiAuthRoutes := roleMemberApiRoutes.Group("/role_member", middleware.AuthorizeJWT(jwtService))
		{
			roleMemberApiAuthRoutes.POST("/save", roleMemberController.Save)
			roleMemberApiAuthRoutes.DELETE("/deleteById/:id", roleMemberController.DeleteById)
		}
		roleMemberNoAuthRoutes := roleMemberApiRoutes.Group("/role_member")
		{
			roleMemberNoAuthRoutes.POST("/getPagination", roleMemberController.GetPagination)
			roleMemberNoAuthRoutes.GET("/getById/:id", roleMemberController.GetById)
			roleMemberNoAuthRoutes.GET("/getByRoleId/:roleId", roleMemberController.GetByRoleId)
			roleMemberNoAuthRoutes.GET("/getViewById/:id", roleMemberController.GetViewById)
		}
	}

	roleMenuApiRoutes := r.Group("api/role_menu")
	{
		roleMenuApiRoutes.POST("/save", roleMenuController.Save)
		roleMenuApiRoutes.GET("/getById/:id", roleMenuController.GetById)
		roleMenuApiRoutes.DELETE("/deleteById/:id", roleMenuController.DeleteById)
		roleMenuApiRoutes.POST("/deleteByRoleAndMenuId", roleMenuController.DeleteByRoleAndMenuId)
	}

	organizationApiRoutes := r.Group("api/organization")
	{
		organizationApiRoutes.POST("/getPagination", organizationController.GetPagination)
		organizationApiRoutes.GET("/lookup", organizationController.Lookup)
		organizationApiRoutes.POST("/save", organizationController.Save)
		organizationApiRoutes.GET("/getById/:id", organizationController.GetById)
		organizationApiRoutes.GET("/getViewById/:id", organizationController.GetViewById)
		organizationApiRoutes.DELETE("/deleteById/:id", organizationController.DeleteById)
	}

	ratingMasterApiRoutes := r.Group("api/rating_master")
	{
		ratingMasterApiRoutes.POST("/getPagination", ratingMasterController.GetPagination)
		ratingMasterApiRoutes.POST("/getAll", ratingMasterController.GetAll)
		ratingMasterApiRoutes.POST("/save", ratingMasterController.Save)
		ratingMasterApiRoutes.GET("/getById/:id", ratingMasterController.GetById)
		ratingMasterApiRoutes.DELETE("/deleteById/:id", ratingMasterController.DeleteById)
	}

	paymentApiRoutes := r.Group("api/payment", middleware.AuthorizeJWT(jwtService))
	{
		paymentApiRoutes.POST("/getPagination", paymentController.GetPagination)
		paymentApiRoutes.POST("/createTokenByCard", paymentController.CreateTokenIdByCard)
		paymentApiRoutes.GET("/getById/:id", paymentController.GetById)
		paymentApiRoutes.GET("/getByProviderRecordId/:id", paymentController.GetByProviderRecordId)
		paymentApiRoutes.GET("/getByProviderReferenceId/:id", paymentController.GetByProviderReferenceId)
		paymentApiRoutes.GET("/getEWalletPaymentStatusByReferenceId/:reference_id", paymentController.GetEWalletPaymentStatusByReferenceId)
		paymentApiRoutes.POST("/updateWalletPaymentStatus", paymentController.UpdateWalletPaymentStatus)
		paymentApiRoutes.GET("/getUniqueNumber", paymentController.GetUniqueNumber)
		paymentApiRoutes.POST("/charge", paymentController.Charge)
		paymentApiRoutes.POST("/createTransferPayment", paymentController.CreateTransferPayment)
		paymentApiRoutes.POST("/createEWalletPayment", paymentController.CreateEWalletPayment)
		paymentApiRoutes.POST("/createQRCode", paymentController.CreateQRCode)
		paymentApiRoutes.POST("/membershipPayment", paymentController.MembershipPayment)
		paymentApiRoutes.POST("/webinarPayment", paymentController.WebinarPayment)
		paymentApiRoutes.POST("/updatePaymentStatus", paymentController.UpdatePaymentStatus)
	}

	productApiRoutes := r.Group("api")
	{
		productApiAuthRoutes := productApiRoutes.Group("/product", middleware.AuthorizeJWT(jwtService))
		{
			productApiAuthRoutes.POST("/save", productController.Save)
			productApiAuthRoutes.DELETE("/deleteById/:id", productController.DeleteById)
		}
		productApiNoAuthRoutes := productApiRoutes.Group("/product")
		{
			productApiNoAuthRoutes.GET("/getAll", productController.GetAll)
			productApiNoAuthRoutes.POST("/getPagination", productController.GetPagination)
			// webinarRecordingNoAuthRoutes.GET("/lookup", productController.Lookup)
			productApiNoAuthRoutes.GET("/getById/:id", productController.GetById)
			productApiNoAuthRoutes.GET("/getProductByRecordId/:record_id", productController.GetProductByRecordId)
			productApiNoAuthRoutes.GET("/getByProductType/:product_type", productController.GetByProductType)
			productApiNoAuthRoutes.GET("/getViewById/:id", productController.GetViewById)
		}
	}

	webhookRoutes := r.Group("webhook/payment")
	{
		webhookRoutes.POST("/gateway", paymentController.HookForXendit)
	}
	// #endregion

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Run(":" + *SERVER_PORT)
}

func GetServerDefinition() (status bool, SERVER_NAME *string, SERVER_PORT *string) {
	name := os.Getenv("SERVER_NAME")
	port := os.Getenv("SERVER_PORT")
	if name == "" || port == "" {
		return false, nil, nil
	}
	return true, &name, &port
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, DELETE, GET, PUT, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
