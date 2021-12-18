package main

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	// _ "github.com/swaggo/gin-swagger/example/basic/docs"

	"joranvest/docs"

	"gorm.io/gorm"

	"joranvest/config"
	"joranvest/controllers"
	"joranvest/repository"
	"joranvest/service"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/multitemplate"
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
	roleRepository                    repository.RoleRepository                    = repository.NewRoleRepository(db)
	roleMemberRepository              repository.RoleMemberRepository              = repository.NewRoleMemberRepository(db)
	roleMenuRepository                repository.RoleMenuRepository                = repository.NewRoleMenuRepository(db)
	organizationRepository            repository.OrganizationRepository            = repository.NewOrganizationRepository(db)
	ratingMasterRepository            repository.RatingMasterRepository            = repository.NewRatingMasterRepository(db)
	paymentRepository                 repository.PaymentRepository                 = repository.NewPaymentRepository(db)
	emailRepository                   repository.EmailRepository                   = repository.NewEmailRepository(db)

	authService                    service.AuthService                    = service.NewAuthService(applicationUserRepository)
	jwtService                     service.JWTService                     = service.NewJWTService()
	applicationUserService         service.ApplicationUserService         = service.NewApplicationUserService(applicationUserRepository)
	applicationMenuCategoryService service.ApplicationMenuCategoryService = service.NewApplicationMenuCategoryService(applicationMenuCategoryRepository)
	applicationMenuService         service.ApplicationMenuService         = service.NewApplicationMenuService(applicationMenuRepository)
	membershipService              service.MembershipService              = service.NewMembershipService(membershipRepository)
	membershipUserService          service.MembershipUserService          = service.NewMembershipUserService(membershipUserRepository)
	filemasterService              service.FilemasterService              = service.NewFilemasterService(filemasterRepository)
	emitenService                  service.EmitenService                  = service.NewEmitenService(emitenRepository)
	emitenCategoryService          service.EmitenCategoryService          = service.NewEmitenCategoryService(emitenCategoryRepository)
	articleCategoryService         service.ArticleCategoryService         = service.NewArticleCategoryService(articleCategoryRepository)
	articleService                 service.ArticleService                 = service.NewArticleService(articleRepository)
	articleTagService              service.ArticleTagService              = service.NewArticleTagService(articleTagRepository)
	webinarCategoryService         service.WebinarCategoryService         = service.NewWebinarCategoryService(webinarCategoryRepository)
	webinarService                 service.WebinarService                 = service.NewWebinarService(webinarRepository)
	webinarSpeakerService          service.WebinarSpeakerService          = service.NewWebinarSpeakerService(webinarSpeakerRepository)
	webinarRegistrationService     service.WebinarRegistrationService     = service.NewWebinarRegistrationService(webinarRegistrationRepository)
	sectorService                  service.SectorService                  = service.NewSectorService(sectorRepository)
	tagService                     service.TagService                     = service.NewTagService(tagRepository)
	technicalAnalysisService       service.TechnicalAnalysisService       = service.NewTechnicalAnalysisService(technicalAnalysisRepository)
	fundamentalAnalysisService     service.FundamentalAnalysisService     = service.NewFundamentalAnalysisService(fundamentalAnalysisRepository)
	fundamentalAnalysisTagService  service.FundamentalAnalysisTagService  = service.NewFundamentalAnalysisTagService(fundamentalAnalysisTagRepository)
	roleService                    service.RoleService                    = service.NewRoleService(roleRepository)
	roleMemberService              service.RoleMemberService              = service.NewRoleMemberService(roleMemberRepository)
	roleMenuService                service.RoleMenuService                = service.NewRoleMenuService(roleMenuRepository)
	organizationService            service.OrganizationService            = service.NewOrganizationService(organizationRepository)
	ratingMasterService            service.RatingMasterService            = service.NewRatingMasterService(ratingMasterRepository)
	paymentService                 service.PaymentService                 = service.NewPaymentService(paymentRepository)
	emailService                   service.EmailService                   = service.NewEmailService(emailRepository)

	authController                    controllers.AuthController                    = controllers.NewAuthController(authService, emailService, jwtService)
	applicationUserController         controllers.ApplicationUserController         = controllers.NewApplicationUserController(applicationUserService, jwtService)
	applicationMenuCategoryController controllers.ApplicationMenuCategoryController = controllers.NewApplicationMenuCategoryController(applicationMenuCategoryService, jwtService)
	applicationMenuController         controllers.ApplicationMenuController         = controllers.NewApplicationMenuController(applicationMenuService, jwtService)
	membershipController              controllers.MembershipController              = controllers.NewMembershipController(membershipService, jwtService)
	membershipUserController          controllers.MembershipUserController          = controllers.NewMembershipUserController(membershipUserService, jwtService)
	filemasterController              controllers.FilemasterController              = controllers.NewFilemasterController(filemasterService, jwtService)
	emitenController                  controllers.EmitenController                  = controllers.NewEmitenController(emitenService, jwtService)
	emitenCategoryController          controllers.EmitenCategoryController          = controllers.NewEmitenCategoryController(emitenCategoryService, jwtService)
	articleCategoryController         controllers.ArticleCategoryController         = controllers.NewArticleCategoryController(articleCategoryService, jwtService)
	articleController                 controllers.ArticleController                 = controllers.NewArticleController(articleService, jwtService)
	articleTagController              controllers.ArticleTagController              = controllers.NewArticleTagController(articleTagService, jwtService)
	webinarCategoryController         controllers.WebinarCategoryController         = controllers.NewWebinarCategoryController(webinarCategoryService, jwtService)
	webinarController                 controllers.WebinarController                 = controllers.NewWebinarController(webinarService, jwtService)
	webinarSpeakerController          controllers.WebinarSpeakerController          = controllers.NewWebinarSpeakerController(webinarSpeakerService, jwtService)
	webinarRegistrationController     controllers.WebinarRegistrationController     = controllers.NewWebinarRegistrationController(webinarRegistrationService, jwtService)
	sectorController                  controllers.SectorController                  = controllers.NewSectorController(sectorService, jwtService)
	tagController                     controllers.TagController                     = controllers.NewTagController(tagService, jwtService)
	technicalAnalysisController       controllers.TechnicalAnalysisController       = controllers.NewTechnicalAnalysisController(technicalAnalysisService, jwtService)
	fundamentalAnalysisController     controllers.FundamentalAnalysisController     = controllers.NewFundamentalAnalysisController(fundamentalAnalysisService, jwtService)
	fundamentalAnalysisTagController  controllers.FundamentalAnalysisTagController  = controllers.NewFundamentalAnalysisTagController(fundamentalAnalysisTagService, jwtService)
	roleController                    controllers.RoleController                    = controllers.NewRoleController(roleService, jwtService)
	roleMemberController              controllers.RoleMemberController              = controllers.NewRoleMemberController(roleMemberService, jwtService)
	roleMenuController                controllers.RoleMenuController                = controllers.NewRoleMenuController(roleMenuService, jwtService)
	organizationController            controllers.OrganizationController            = controllers.NewOrganizationController(organizationService, jwtService)
	ratingMasterController            controllers.RatingMasterController            = controllers.NewRatingMasterController(ratingMasterService, jwtService)
	paymentController                 controllers.PaymentController                 = controllers.NewPaymentController(paymentService, jwtService)
	emailController                   controllers.EmailController                   = controllers.NewEmailController(emailService, jwtService)

	//clientest coreapi.ClientTest = coreapi.NewClientTest()
	// #endregion
)

func createMyRender(view_path string) multitemplate.Renderer {
	r := multitemplate.NewRenderer()
	// _ = []string{
	// 	shared_path + "_header.html",
	// 	shared_path + "_nav.html",
	// 	shared_path + "_logout.html",
	// 	shared_path + "_baseScript.html",
	// 	// ...
	// }

	shared_path := view_path + "shared/"
	// #region Template Route
	r.AddFromFiles("index", "templates/views/index.html", shared_path+"_header.html", shared_path+"_baseScript.html")
	// #endregion
	return r
}

func Setup(c *gin.Context, title string, header string, subheader string, nav string, subnav string) map[string]string {
	data := make(map[string]string)
	session := sessions.Default(c)
	isAdmin := fmt.Sprint(session.Get("IsAdmin"))
	entityid := fmt.Sprint(session.Get("EntityId"))
	data["title"] = title
	data["header"] = header
	data["subheader"] = subheader
	data["nav"] = nav
	data["subnav"] = subnav
	data["timenow"] = time.Now().Format("2006-01-02 15:04:05.000000")
	data["token"] = fmt.Sprint(session.Get("mytoken"))
	if isAdmin == "true" {
		data["isadmin"] = "1"
	} else {
		data["isadmin"] = "0"
	}
	data["entityid"] = entityid
	data["userLoginName"] = fmt.Sprint(session.Get("userLoginName"))
	data["userFirstName"] = fmt.Sprint(session.Get("userFirstName"))

	if session.Get("userFirstName") == nil {
		data["userFirstName"] = ""
	}

	//url := location.Get(c)
	fmt.Println("=======================================")
	// fmt.Println("Scheme: ", url.Scheme)
	// fmt.Println("Host: ", url.Host)
	// fmt.Println("Path: ", url.Path)
	fmt.Println("=======================================")
	//data["hostName"] = url.Scheme + "://" + url.Host

	return data
}

func main() {
	defer config.CloseDatabaseConnection(db)
	//getCardToken()

	r := gin.New()
	// r := gin.Default()

	// programatically set swagger info
	docs.SwaggerInfo.Title = "Swagger Example API"
	docs.SwaggerInfo.Description = "This is a sample server Petstore server."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:10000"
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
	r.HTMLRender = createMyRender("templates/views/")

	// #region User Web View
	r.GET("/", func(c *gin.Context) {
		data := Setup(c, "Joranvest", "", "", "", "")
		c.HTML(
			http.StatusOK,
			"index",
			gin.H{
				"title": "Joranvest",
				"err":   "",
				"data":  data,
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
		webinarApiRoutes.POST("/save", webinarController.Save)
		webinarApiRoutes.POST("/submit/:id", webinarController.Submit)
		webinarApiRoutes.GET("/getById/:id", webinarController.GetById)
		webinarApiRoutes.DELETE("/deleteById/:id", webinarController.DeleteById)
	}

	webinarSpeakerApiRoutes := r.Group("api/webinar_speaker")
	{
		webinarSpeakerApiRoutes.POST("/save", webinarSpeakerController.Save)
		webinarSpeakerApiRoutes.GET("/getById/:id", webinarSpeakerController.GetById)
		webinarSpeakerApiRoutes.GET("/getAll", webinarSpeakerController.GetAll)
	}

	webinarRegistrationApiRoutes := r.Group("api/webinar_registration")
	{
		webinarRegistrationApiRoutes.POST("/getDatatables", webinarRegistrationController.GetDatatables)
		webinarRegistrationApiRoutes.POST("/getPagination", webinarRegistrationController.GetPagination)
		webinarRegistrationApiRoutes.GET("/getById/:id", webinarRegistrationController.GetById)
		webinarRegistrationApiRoutes.GET("/getViewById/:id", webinarRegistrationController.GetViewById)
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

	membershipApiRoutes := r.Group("api/membership")
	{
		membershipApiRoutes.POST("/getDatatables", membershipController.GetDatatables)
		membershipApiRoutes.POST("/getPagination", membershipController.GetPagination)
		membershipApiRoutes.GET("/getAll", membershipController.GetAll)
		membershipApiRoutes.POST("/save", membershipController.Save)
		membershipApiRoutes.POST("/setRecommendation", membershipController.SetRecommendation)
		membershipApiRoutes.GET("/getById/:id", membershipController.GetById)
		membershipApiRoutes.GET("/getViewById/:id", membershipController.GetViewById)
		membershipApiRoutes.DELETE("/deleteById/:id", membershipController.DeleteById)
	}

	membershipUserApiRoutes := r.Group("api/membershipUser")
	{
		membershipUserApiRoutes.POST("/getDatatables", membershipUserController.GetDatatables)
		membershipUserApiRoutes.GET("/getAll", membershipUserController.GetAll)
		membershipUserApiRoutes.POST("/save", membershipUserController.Save)
		membershipUserApiRoutes.GET("/getById/:id", membershipUserController.GetById)
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
		applicationUserApiRoutes.POST("/getDatatables", applicationUserController.GetDatatables)
		applicationUserApiRoutes.GET("/lookup", applicationUserController.Lookup)
		applicationUserApiRoutes.POST("/changePassword", applicationUserController.ChangePassword)
		applicationUserApiRoutes.POST("/recoverPassword", applicationUserController.RecoverPassword)
		applicationUserApiRoutes.POST("/register", authController.Register)
		applicationUserApiRoutes.GET("/getViewById/:id", applicationUserController.GetViewById)
		applicationUserApiRoutes.GET("/emailVerificationById/:id", applicationUserController.EmailVerificationById)
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

	roleApiRoutes := r.Group("api/role")
	{
		roleApiRoutes.POST("/getDatatables", roleController.GetDatatables)
		roleApiRoutes.GET("/lookup", roleController.Lookup)
		roleApiRoutes.POST("/save", roleController.Save)
		roleApiRoutes.GET("/getById/:id", roleController.GetById)
		roleApiRoutes.DELETE("/deleteById/:id", roleController.DeleteById)
	}
	roleMemberApiRoutes := r.Group("api/role_member")
	{
		roleMemberApiRoutes.POST("/getDatatables", roleMemberController.GetDatatables)
		roleMemberApiRoutes.POST("/save", roleMemberController.Save)
		roleMemberApiRoutes.GET("/getById/:id", roleMemberController.GetById)
		roleMemberApiRoutes.DELETE("/deleteById/:id", roleMemberController.DeleteById)
		roleMemberApiRoutes.GET("/getUsersInRole/:roleId", roleMemberController.GetUsersInRole)
		roleMemberApiRoutes.GET("/getUsersNotInRole/:roleId/:search", roleMemberController.GetUsersNotInRole)
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
		organizationApiRoutes.DELETE("/deleteById/:id", organizationController.DeleteById)
	}

	ratingMasterApiRoutes := r.Group("api/rating_master")
	{
		ratingMasterApiRoutes.POST("/getAll", ratingMasterController.GetAll)
		ratingMasterApiRoutes.POST("/save", ratingMasterController.Save)
		ratingMasterApiRoutes.GET("/getById/:id", ratingMasterController.GetById)
		ratingMasterApiRoutes.DELETE("/deleteById/:id", ratingMasterController.DeleteById)
	}

	paymentApiRoutes := r.Group("api/payment")
	{
		paymentApiRoutes.POST("/getPagination", paymentController.GetPagination)
		paymentApiRoutes.POST("/createTokenByCard", paymentController.CreateTokenIdByCard)
		paymentApiRoutes.GET("/getById/:id", paymentController.GetById)
		paymentApiRoutes.GET("/getUniqueNumber", paymentController.GetUniqueNumber)
		paymentApiRoutes.POST("/charge", paymentController.Charge)
		paymentApiRoutes.POST("/membershipPayment", paymentController.MembershipPayment)
		paymentApiRoutes.POST("/webinarPayment", paymentController.WebinarPayment)
		paymentApiRoutes.POST("/updatePaymentStatus", paymentController.UpdatePaymentStatus)
	}

	emailApiRoutes := r.Group("api/email")
	{
		emailApiRoutes.POST("/sendEmailVerification", emailController.SendEmailVerification)
	}
	// #endregion

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// url := location.Get()
	fmt.Println("=======================================")
	// fmt.Println("Scheme: ", url.Scheme)
	// fmt.Println("Host: ", url.Host)
	// fmt.Println("Path: ", url.Path)
	fmt.Println("=======================================")
	r.Run(":10000")
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
