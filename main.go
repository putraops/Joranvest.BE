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

	"github.com/gin-contrib/location"
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
	filemasterRepository              repository.FilemasterRepository              = repository.NewFilemasterRepository(db)
	emitenRepository                  repository.EmitenRepository                  = repository.NewEmitenRepository(db)
	emitenCategoryRepository          repository.EmitenCategoryRepository          = repository.NewEmitenCategoryRepository(db)
	articleCategoryRepository         repository.ArticleCategoryRepository         = repository.NewArticleCategoryRepository(db)
	articleRepository                 repository.ArticleRepository                 = repository.NewArticleRepository(db)
	articleTagRepository              repository.ArticleTagRepository              = repository.NewArticleTagRepository(db)
	webinarCategoryRepository         repository.WebinarCategoryRepository         = repository.NewWebinarCategoryRepository(db)
	webinarRepository                 repository.WebinarRepository                 = repository.NewWebinarRepository(db)
	webinarSpeakerRepository          repository.WebinarSpeakerRepository          = repository.NewWebinarSpeakerRepository(db)
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

	authService                    service.AuthService                    = service.NewAuthService(applicationUserRepository)
	jwtService                     service.JWTService                     = service.NewJWTService()
	applicationUserService         service.ApplicationUserService         = service.NewApplicationUserService(applicationUserRepository)
	applicationMenuCategoryService service.ApplicationMenuCategoryService = service.NewApplicationMenuCategoryService(applicationMenuCategoryRepository)
	applicationMenuService         service.ApplicationMenuService         = service.NewApplicationMenuService(applicationMenuRepository)
	membershipService              service.MembershipService              = service.NewMembershipService(membershipRepository)
	filemasterService              service.FilemasterService              = service.NewFilemasterService(filemasterRepository)
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
	roleService                    service.RoleService                    = service.NewRoleService(roleRepository)
	roleMemberService              service.RoleMemberService              = service.NewRoleMemberService(roleMemberRepository)
	roleMenuService                service.RoleMenuService                = service.NewRoleMenuService(roleMenuRepository)
	organizationService            service.OrganizationService            = service.NewOrganizationService(organizationRepository)
	ratingMasterService            service.RatingMasterService            = service.NewRatingMasterService(ratingMasterRepository)

	authController                    controllers.AuthController                    = controllers.NewAuthController(authService, jwtService)
	applicationUserController         controllers.ApplicationUserController         = controllers.NewApplicationUserController(applicationUserService, jwtService)
	applicationMenuCategoryController controllers.ApplicationMenuCategoryController = controllers.NewApplicationMenuCategoryController(applicationMenuCategoryService, jwtService)
	applicationMenuController         controllers.ApplicationMenuController         = controllers.NewApplicationMenuController(applicationMenuService, jwtService)
	membershipController              controllers.MembershipController              = controllers.NewMembershipController(membershipService, jwtService)
	filemasterController              controllers.FilemasterController              = controllers.NewFilemasterController(filemasterService, jwtService)
	emitenController                  controllers.EmitenController                  = controllers.NewEmitenController(emitenService, jwtService)
	emitenCategoryController          controllers.EmitenCategoryController          = controllers.NewEmitenCategoryController(emitenCategoryService, jwtService)
	articleCategoryController         controllers.ArticleCategoryController         = controllers.NewArticleCategoryController(articleCategoryService, jwtService)
	articleController                 controllers.ArticleController                 = controllers.NewArticleController(articleService, jwtService)
	articleTagController              controllers.ArticleTagController              = controllers.NewArticleTagController(articleTagService, jwtService)
	webinarCategoryController         controllers.WebinarCategoryController         = controllers.NewWebinarCategoryController(webinarCategoryService, jwtService)
	webinarController                 controllers.WebinarController                 = controllers.NewWebinarController(webinarService, jwtService)
	webinarSpeakerController          controllers.WebinarSpeakerController          = controllers.NewWebinarSpeakerController(webinarSpeakerService, jwtService)
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
	admin_shared_path := view_path + "admin/shared/"

	// #region Template Route
	r.AddFromFiles("index", "templates/views/index.html", shared_path+"_header.html", shared_path+"_baseScript.html")
	r.AddFromFiles("login", "templates/views/login.html", shared_path+"_header.html", shared_path+"_baseScript.html")
	r.AddFromFiles("logout", "templates/views/logout.html", shared_path+"_header.html", shared_path+"_baseScript.html")
	r.AddFromFiles("register", "templates/views/register/index.html", shared_path+"_header.html", shared_path+"_baseScript.html")
	r.AddFromFiles("register/success", "templates/views/register/success.html", shared_path+"_header.html", shared_path+"_baseScript.html")
	r.AddFromFiles("browse",
		view_path+"_base.html", view_path+"order/order.index.html",
		shared_path+"_header.html", shared_path+"_nav.html", shared_path+"_topNav.html",
		shared_path+"_logout.html", shared_path+"_footer.html", shared_path+"_baseScript.html")

	r.AddFromFiles("browse/technical", view_path+"browse/technical.html", shared_path+"_header.html", shared_path+"_baseScript.html")
	r.AddFromFiles("article", view_path+"article/index.html", shared_path+"_header.html", shared_path+"_baseScript.html")

	r.AddFromFiles("dashboard",
		view_path+"_base.html", view_path+"admin/dashboard.html",
		admin_shared_path+"_header.html", admin_shared_path+"_nav.html", admin_shared_path+"_topNav.html",
		admin_shared_path+"_logout.html", admin_shared_path+"_footer.html", admin_shared_path+"_baseScript.html")

	r.AddFromFiles("emiten_category",
		view_path+"_base.html", view_path+"admin/emiten_category/emiten_category.index.html",
		admin_shared_path+"_header.html", admin_shared_path+"_nav.html", admin_shared_path+"_topNav.html",
		admin_shared_path+"_logout.html", admin_shared_path+"_footer.html", admin_shared_path+"_baseScript.html")

	r.AddFromFiles("webinar_category",
		view_path+"_base.html", view_path+"admin/webinar_category/webinar_category.index.html",
		admin_shared_path+"_header.html", admin_shared_path+"_nav.html", admin_shared_path+"_topNav.html",
		admin_shared_path+"_logout.html", admin_shared_path+"_footer.html", admin_shared_path+"_baseScript.html")
	r.AddFromFiles("webinar",
		view_path+"_base.html", view_path+"admin/webinar/webinar.index.html",
		admin_shared_path+"_header.html", admin_shared_path+"_nav.html", admin_shared_path+"_topNav.html",
		admin_shared_path+"_logout.html", admin_shared_path+"_footer.html", admin_shared_path+"_baseScript.html")
	r.AddFromFiles("webinar_detail",
		view_path+"_base.html", view_path+"admin/webinar/webinar.detail.html",
		admin_shared_path+"_header.html", admin_shared_path+"_nav.html", admin_shared_path+"_topNav.html",
		admin_shared_path+"_logout.html", admin_shared_path+"_footer.html", admin_shared_path+"_baseScript.html")

	r.AddFromFiles("article_category",
		view_path+"_base.html", view_path+"admin/article_category/article_category.index.html",
		admin_shared_path+"_header.html", admin_shared_path+"_nav.html", admin_shared_path+"_topNav.html",
		admin_shared_path+"_logout.html", admin_shared_path+"_footer.html", admin_shared_path+"_baseScript.html")
	r.AddFromFiles("article",
		view_path+"_base.html", view_path+"admin/article/article.index.html",
		admin_shared_path+"_header.html", admin_shared_path+"_nav.html", admin_shared_path+"_topNav.html",
		admin_shared_path+"_logout.html", admin_shared_path+"_footer.html", admin_shared_path+"_baseScript.html")
	r.AddFromFiles("article_detail",
		view_path+"_base.html", view_path+"admin/article/article.detail.html",
		admin_shared_path+"_header.html", admin_shared_path+"_nav.html", admin_shared_path+"_topNav.html",
		admin_shared_path+"_logout.html", admin_shared_path+"_footer.html", admin_shared_path+"_baseScript.html")

	r.AddFromFiles("application_menu",
		view_path+"_base.html", view_path+"admin/application_menu/application_menu.index.html",
		admin_shared_path+"_header.html", admin_shared_path+"_nav.html", admin_shared_path+"_topNav.html",
		admin_shared_path+"_logout.html", admin_shared_path+"_footer.html", admin_shared_path+"_baseScript.html")
	r.AddFromFiles("application_menu_detail",
		view_path+"_base.html", view_path+"admin/application_menu/application_menu.detail.html",
		admin_shared_path+"_header.html", admin_shared_path+"_nav.html", admin_shared_path+"_topNav.html",
		admin_shared_path+"_logout.html", admin_shared_path+"_footer.html", admin_shared_path+"_baseScript.html")

	r.AddFromFiles("emiten",
		view_path+"_base.html", view_path+"admin/emiten/emiten.index.html",
		admin_shared_path+"_header.html", admin_shared_path+"_nav.html", admin_shared_path+"_topNav.html",
		admin_shared_path+"_logout.html", admin_shared_path+"_footer.html", admin_shared_path+"_baseScript.html")

	r.AddFromFiles("sector",
		view_path+"_base.html", view_path+"admin/sector/sector.index.html",
		admin_shared_path+"_header.html", admin_shared_path+"_nav.html", admin_shared_path+"_topNav.html",
		admin_shared_path+"_logout.html", admin_shared_path+"_footer.html", admin_shared_path+"_baseScript.html")

	r.AddFromFiles("tag",
		view_path+"_base.html", view_path+"admin/tag/tag.index.html",
		admin_shared_path+"_header.html", admin_shared_path+"_nav.html", admin_shared_path+"_topNav.html",
		admin_shared_path+"_logout.html", admin_shared_path+"_footer.html", admin_shared_path+"_baseScript.html")

	r.AddFromFiles("membership",
		view_path+"_base.html", view_path+"admin/membership/membership.index.html",
		admin_shared_path+"_header.html", admin_shared_path+"_nav.html", admin_shared_path+"_topNav.html",
		admin_shared_path+"_logout.html", admin_shared_path+"_footer.html", admin_shared_path+"_baseScript.html")

	r.AddFromFiles("application_user",
		view_path+"_base.html", view_path+"admin/application_user/application_user.index.html",
		admin_shared_path+"_header.html", admin_shared_path+"_nav.html", admin_shared_path+"_topNav.html",
		admin_shared_path+"_logout.html", admin_shared_path+"_footer.html", admin_shared_path+"_baseScript.html")

	r.AddFromFiles("technical_analysis_index",
		view_path+"_base.html", view_path+"admin/analysis/technical_analysis.index.html",
		admin_shared_path+"_header.html", admin_shared_path+"_nav.html", admin_shared_path+"_topNav.html",
		admin_shared_path+"_logout.html", admin_shared_path+"_footer.html", admin_shared_path+"_baseScript.html")
	r.AddFromFiles("technical_analysis_detail",
		view_path+"_base.html", view_path+"admin/analysis/technical_analysis.detail.html",
		admin_shared_path+"_header.html", admin_shared_path+"_nav.html", admin_shared_path+"_topNav.html",
		admin_shared_path+"_logout.html", admin_shared_path+"_footer.html", admin_shared_path+"_baseScript.html")

	r.AddFromFiles("fundamental_analysis_index",
		view_path+"_base.html", view_path+"admin/analysis/fundamental_analysis.index.html",
		admin_shared_path+"_header.html", admin_shared_path+"_nav.html", admin_shared_path+"_topNav.html",
		admin_shared_path+"_logout.html", admin_shared_path+"_footer.html", admin_shared_path+"_baseScript.html")
	r.AddFromFiles("fundamental_analysis_detail",
		view_path+"_base.html", view_path+"admin/analysis/fundamental_analysis.detail.html",
		admin_shared_path+"_header.html", admin_shared_path+"_nav.html", admin_shared_path+"_topNav.html",
		admin_shared_path+"_logout.html", admin_shared_path+"_footer.html", admin_shared_path+"_baseScript.html")

	r.AddFromFiles("role",
		view_path+"_base.html", view_path+"admin/role/role.index.html",
		admin_shared_path+"_header.html", admin_shared_path+"_nav.html", admin_shared_path+"_topNav.html",
		admin_shared_path+"_logout.html", admin_shared_path+"_footer.html", admin_shared_path+"_baseScript.html")
	r.AddFromFiles("role_member",
		view_path+"_base.html", view_path+"admin/role/role.member.html",
		admin_shared_path+"_header.html", admin_shared_path+"_nav.html", admin_shared_path+"_topNav.html",
		admin_shared_path+"_logout.html", admin_shared_path+"_footer.html", admin_shared_path+"_baseScript.html")
	r.AddFromFiles("role_menu",
		view_path+"_base.html", view_path+"admin/role/role.menu.html",
		admin_shared_path+"_header.html", admin_shared_path+"_nav.html", admin_shared_path+"_topNav.html",
		admin_shared_path+"_logout.html", admin_shared_path+"_footer.html", admin_shared_path+"_baseScript.html")

	r.AddFromFiles("organization",
		view_path+"_base.html", view_path+"admin/organization/organization.index.html",
		admin_shared_path+"_header.html", admin_shared_path+"_nav.html", admin_shared_path+"_topNav.html",
		admin_shared_path+"_logout.html", admin_shared_path+"_footer.html", admin_shared_path+"_baseScript.html")
	// #endregion
	return r
}

type App struct {
	c *gin.Context
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

	url := location.Get(c)
	fmt.Println("=======================================")
	fmt.Println("Scheme: ", url.Scheme)
	fmt.Println("Host: ", url.Host)
	fmt.Println("Path: ", url.Path)
	fmt.Println("=======================================")
	data["hostName"] = url.Scheme + "://" + url.Host

	return data
}

func main() {
	defer config.CloseDatabaseConnection(db)
	r := gin.Default()

	// programatically set swagger info
	docs.SwaggerInfo.Title = "Swagger Example API"
	docs.SwaggerInfo.Description = "This is a sample server Petstore server."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:10000"
	docs.SwaggerInfo.BasePath = "/api"

	r.Use(location.New(location.Config{
		Scheme:  "http",
		Host:    "foo.com",
		Base:    "/base",
		Headers: location.Headers{Scheme: "X-Forwarded-Proto", Host: "X-Forwarded-For"},
	}))

	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store))

	r.Use(CORSMiddleware())
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
	r.GET("/logout", func(c *gin.Context) {
		data := Setup(c, "Logout", "", "", "", "")
		c.HTML(
			http.StatusOK,
			"logout",
			gin.H{
				"title": "Logout",
				"err":   "",
				"data":  data,
			},
		)
	})

	r.GET("/login", func(c *gin.Context) {
		c.HTML(
			http.StatusOK,
			"login",
			gin.H{
				"title": "Login",
				"err":   "",
			},
		)
	})
	r.GET("/register", func(c *gin.Context) {
		data := Setup(c, "Dashboard", "", "", "", "")
		c.HTML(
			http.StatusOK,
			"register",
			gin.H{
				"title": "Register",
				"err":   "",
				"data":  data,
			},
		)
	})
	// r.GET("/register", func(c *gin.Context) {
	// 	c.HTML(
	// 		http.StatusOK,
	// 		"register",
	// 		gin.H{
	// 			"title": "Register",
	// 			"err":   "",
	// 		},
	// 	)
	// })

	applicationUserRoutes := r.Group("user")
	{
		applicationUserRoutes.POST("/register", authController.RegisterForm)
		applicationUserRoutes.GET("/register/success", func(c *gin.Context) {
			data := Setup(c, "Dashboard", "", "", "", "")
			c.HTML(
				http.StatusOK,
				"register/success",
				gin.H{
					"title": "Register",
					"err":   "",
					"data":  data,
				},
			)
		})
	}
	// #endregion

	// #region Auth Route
	r.POST("/login", func(c *gin.Context) {
		c.Request.ParseForm()
		username := c.PostFormArray("username")[0]
		email := c.PostFormArray("email")[0]
		password := c.PostFormArray("password")[0]
		var err = ""
		var isValid = false
		var isAdmin = false

		if email == "" && password == "" {
			err = "Email dan Password tidak boleh kosong."
		} else if email == "" {
			err = "Email tidak boleh kosong."
		} else if password == "" {
			err = "Password tidak boleh kosong."
		} else {
			err = "Nice"
			isVerified, _isAdmin, token, message := authController.LoginForm(c, username, email, password)

			if isVerified {
				session := sessions.Default(c)
				session.Set("mytoken", token)
				session.Save()
				err = "Login Successfully."
				isValid = true
				isAdmin = _isAdmin

			} else {
				err = message
			}
		}

		if isValid {
			if isAdmin {
				c.Redirect(http.StatusMovedPermanently, "/admin/dashboard")
			} else {
				fmt.Println("end")
				c.Redirect(http.StatusMovedPermanently, "/")
			}
		} else {
			c.HTML(
				http.StatusOK, "login",
				gin.H{
					"title": "Login",
					"err":   err,
				},
			)
		}
	})
	r.POST("/logout", func(c *gin.Context) {
		c.Request.ParseForm()
		authController.Logout(c)
		c.Redirect(http.StatusMovedPermanently, "/")
	})
	// #endregion

	// #region Web Admin Route
	adminRoutes := r.Group("admin")
	{
		adminRoutes.GET("/dashboard", func(c *gin.Context) {
			data := Setup(c, "Dashboard", "", "", "", "")
			c.HTML(
				http.StatusOK,
				"dashboard",
				gin.H{
					"data": data,
				},
			)
		})

		adminRoutes.GET("/application_menu", func(c *gin.Context) {
			data := Setup(c, "Application Menu", "Application Menu", "Application Menu", "Application Menu", "Application Menu")
			c.HTML(
				http.StatusOK,
				"application_menu",
				gin.H{
					"data": data,
				},
			)
		})

		adminRoutes.GET("/application_menu/detail", func(c *gin.Context) {
			data := Setup(c, "Application Menu", "Application Menu", "Application Menu", "Application Menu", "Application Menu")
			qry := c.Request.URL.Query()
			if _, found := qry["id"]; found {
				data["id"] = fmt.Sprint(qry["id"][0])
			}
			c.HTML(
				http.StatusOK,
				"application_menu_detail",
				gin.H{
					"data": data,
				},
			)
		})

		adminRoutes.GET("/membership", func(c *gin.Context) {
			data := Setup(c, "Membership", "", "", "", "")
			c.HTML(
				http.StatusOK,
				"membership",
				gin.H{
					"data": data,
				},
			)
		})

		adminRoutes.GET("/emiten", func(c *gin.Context) {
			data := Setup(c, "Emiten", "", "", "", "")
			c.HTML(
				http.StatusOK,
				"emiten",
				gin.H{
					"data": data,
				},
			)
		})

		adminRoutes.GET("/emiten_category", func(c *gin.Context) {
			data := Setup(c, "Emiten Category", "Emiten Category", "Emiten Category", "Emiten Category", "Emiten Category")
			c.HTML(
				http.StatusOK,
				"emiten_category",
				gin.H{
					"data": data,
				},
			)
		})

		adminRoutes.GET("/article_category", func(c *gin.Context) {
			data := Setup(c, "Article Category", "Article Category", "Article Category", "Article Category", "Article Category")
			c.HTML(
				http.StatusOK,
				"article_category",
				gin.H{
					"data": data,
				},
			)
		})

		adminRoutes.GET("/webinar_category", func(c *gin.Context) {
			data := Setup(c, "Webinar Category", "Webinar Category", "Webinar Category", "Webinar Category", "Webinar Category")
			c.HTML(
				http.StatusOK,
				"webinar_category",
				gin.H{
					"data": data,
				},
			)
		})

		adminRoutes.GET("/webinar", func(c *gin.Context) {
			data := Setup(c, "Webinar", "Webinar", "Webinar", "Webinar", "Webinar")
			c.HTML(
				http.StatusOK,
				"webinar",
				gin.H{
					"data": data,
				},
			)
		})
		adminRoutes.GET("/webinar/detail", func(c *gin.Context) {
			data := Setup(c, "Webinar Detail", "Webinar Detail", "Webinar Detail", "Webinar Detail", "Webinar Detail")
			qry := c.Request.URL.Query()
			if _, found := qry["id"]; found {
				data["id"] = fmt.Sprint(qry["id"][0])
			}
			c.HTML(
				http.StatusOK,
				"webinar_detail",
				gin.H{
					"data": data,
				},
			)
		})

		adminRoutes.GET("/article", func(c *gin.Context) {
			data := Setup(c, "Article", "Article", "Article", "Article", "Article")
			c.HTML(
				http.StatusOK,
				"article",
				gin.H{
					"data": data,
				},
			)
		})
		adminRoutes.GET("/article/detail", func(c *gin.Context) {
			data := Setup(c, "Article Detail", "Article Detail", "Article Detail", "Article Detail", "Article Detail")
			qry := c.Request.URL.Query()
			if _, found := qry["id"]; found {
				data["id"] = fmt.Sprint(qry["id"][0])
			}
			c.HTML(
				http.StatusOK,
				"article_detail",
				gin.H{
					"data": data,
				},
			)
		})

		adminRoutes.GET("/sector", func(c *gin.Context) {
			data := Setup(c, "Sector", "Sector", "Sector", "Sector", "Sector")
			c.HTML(
				http.StatusOK,
				"sector",
				gin.H{
					"data": data,
				},
			)
		})

		adminRoutes.GET("/tag", func(c *gin.Context) {
			data := Setup(c, "Tag", "", "", "", "")
			c.HTML(
				http.StatusOK,
				"tag",
				gin.H{
					"data": data,
				},
			)
		})

		adminRoutes.GET("/technical_analysis", func(c *gin.Context) {
			data := Setup(c, "Analisa Teknikal", "Analisa Teknikal", "Analisa Teknikal", "Analisa Teknikal", "Analisa Teknikal")
			c.HTML(
				http.StatusOK,
				"technical_analysis_index",
				gin.H{
					"data": data,
				},
			)
		})
		adminRoutes.GET("/technical_analysis/detail", func(c *gin.Context) {
			data := Setup(c, "Analisa Teknikal", "Analisa Teknikal", "Analisa Teknikal", "Analisa Teknikal", "Analisa Teknikal")
			qry := c.Request.URL.Query()
			if _, found := qry["id"]; found {
				data["id"] = fmt.Sprint(qry["id"][0])
			}
			c.HTML(
				http.StatusOK,
				"technical_analysis_detail",
				gin.H{
					"data": data,
				},
			)
		})

		adminRoutes.GET("/fundamental_analysis", func(c *gin.Context) {
			data := Setup(c, "Analisa Fundamental", "Analisa Fundamental", "Analisa Fundamental", "Analisa Fundamental", "Analisa Fundamental")
			c.HTML(
				http.StatusOK,
				"fundamental_analysis_index",
				gin.H{
					"data": data,
				},
			)
		})
		adminRoutes.GET("/fundamental_analysis/detail", func(c *gin.Context) {
			data := Setup(c, "Analisa Fundamental", "Analisa Fundamental", "Analisa Fundamental", "Analisa Fundamental", "Analisa Fundamental")
			qry := c.Request.URL.Query()
			if _, found := qry["id"]; found {
				data["id"] = fmt.Sprint(qry["id"][0])
			}
			c.HTML(
				http.StatusOK,
				"fundamental_analysis_detail",
				gin.H{
					"data": data,
				},
			)
		})

		adminRoutes.GET("/application_user", func(c *gin.Context) {
			data := Setup(c, "Application User", "", "", "", "")
			c.HTML(
				http.StatusOK,
				"application_user",
				gin.H{
					"data": data,
				},
			)
		})

		adminRoutes.GET("/role", func(c *gin.Context) {
			data := Setup(c, "Role", "Role", "Role", "Role", "Role")
			c.HTML(
				http.StatusOK,
				"role",
				gin.H{
					"data": data,
				},
			)
		})
		adminRoutes.GET("/role/member", func(c *gin.Context) {
			data := Setup(c, "Role Member", "Role Member", "Role Member", "Role Member", "Role Member")
			qry := c.Request.URL.Query()
			if _, found := qry["id"]; found {
				data["id"] = fmt.Sprint(qry["id"][0])
			}
			c.HTML(
				http.StatusOK,
				"role_member",
				gin.H{
					"data": data,
				},
			)
		})
		adminRoutes.GET("/role/menu", func(c *gin.Context) {
			data := Setup(c, "Role Menu", "Role Menu", "Role Menu", "Role Menu", "Role Menu")
			qry := c.Request.URL.Query()
			if _, found := qry["id"]; found {
				data["id"] = fmt.Sprint(qry["id"][0])
			}
			c.HTML(
				http.StatusOK,
				"role_menu",
				gin.H{
					"data": data,
				},
			)
		})

		adminRoutes.GET("/organization", func(c *gin.Context) {
			data := Setup(c, "Organization", "Organization", "Organization", "Organization", "Organization")
			c.HTML(
				http.StatusOK,
				"organization",
				gin.H{
					"data": data,
				},
			)
		})
	}

	browseRoutes := r.Group("browse")
	{
		browseRoutes.GET("", func(c *gin.Context) {

			data := Setup(c, "Order", "Daftar Order", "Daftar Order", "", "")
			c.HTML(
				http.StatusOK, "order",
				gin.H{
					"data": data,
				},
			)
		})

		browseRoutes.GET("/technical", func(c *gin.Context) {
			data := Setup(c, "Technical", "Technical", "Technical", "", "")
			//qry := c.Request.URL.Query()
			// if _, found := qry["id"]; found {
			// 	data["id"] = fmt.Sprint(qry["id"][0])
			// }

			c.HTML(
				http.StatusOK, "browse/technical",
				gin.H{
					"data": data,
				},
			)
		})
	}

	articleRoutes := r.Group("article")
	{
		articleRoutes.GET("", func(c *gin.Context) {

			data := Setup(c, "Article", "List Article", "List Article", "", "")
			c.HTML(
				http.StatusOK, "article",
				gin.H{
					"data": data,
				},
			)
		})
	}

	entityRoutes := r.Group("entity")
	{
		entityRoutes.GET("", func(c *gin.Context) {
			data := Setup(c, "Entity", "Daftar Entity", "Daftar Entity", "", "")
			c.HTML(
				http.StatusOK, "entity",
				gin.H{
					"data": data,
				},
			)
		})
	}
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
		emitenApiRoutes.GET("/lookup", emitenController.Lookup)
		emitenApiRoutes.POST("/emitenLookup", emitenController.EmitenLookup)
		emitenApiRoutes.POST("/save", emitenController.Save)
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
		webinarCategoryApiRoutes.GET("/getTree", webinarCategoryController.GetTree)
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
		webinarSpeakerApiRoutes.GET("/getById/:id", webinarSpeakerController.GetById)
		webinarSpeakerApiRoutes.GET("/getAll", webinarSpeakerController.GetAll)
	}

	tagApiRoutes := r.Group("api/tag")
	{
		tagApiRoutes.POST("/getDatatables", tagController.GetDatatables)
		tagApiRoutes.GET("/lookup", tagController.Lookup)
		tagApiRoutes.POST("/save", tagController.Save)
		tagApiRoutes.GET("/getById/:id", tagController.GetById)
		tagApiRoutes.DELETE("/deleteById/:id", tagController.DeleteById)
	}

	sectorApiRoutes := r.Group("api/sector")
	{
		sectorApiRoutes.POST("/getDatatables", sectorController.GetDatatables)
		sectorApiRoutes.GET("/lookup", sectorController.Lookup)
		sectorApiRoutes.POST("/save", sectorController.Save)
		sectorApiRoutes.GET("/getById/:id", sectorController.GetById)
		sectorApiRoutes.DELETE("/deleteById/:id", sectorController.DeleteById)
	}

	membershipApiRoutes := r.Group("api/membership")
	{
		membershipApiRoutes.POST("/getDatatables", membershipController.GetDatatables)
		membershipApiRoutes.GET("/getAll", membershipController.GetAll)
		membershipApiRoutes.POST("/save", membershipController.Save)
		membershipApiRoutes.POST("/setRecommendation", membershipController.SetRecommendation)
		membershipApiRoutes.GET("/getById/:id", membershipController.GetById)
		membershipApiRoutes.DELETE("/deleteById/:id", membershipController.DeleteById)
	}

	filemasterApiRoutes := r.Group("api/filemaster")
	{
		filemasterApiRoutes.POST("/single_upload/:id", filemasterController.SingleUpload)
		filemasterApiRoutes.POST("/uploadByType/:module/:filetype/:id", filemasterController.UploadByType)
		filemasterApiRoutes.POST("/upload/:id", filemasterController.Insert)
		filemasterApiRoutes.GET("/getAll", filemasterController.GetAll)
		filemasterApiRoutes.DELETE("/deleteByRecordId/:recordId", filemasterController.DeleteByRecordId)
	}

	applicationUserApiRoutes := r.Group("api/application_user")
	{
		applicationUserApiRoutes.POST("/getDatatables", applicationUserController.GetDatatables)
		applicationUserApiRoutes.GET("/lookup", applicationUserController.Lookup)
		applicationUserApiRoutes.POST("/register", authController.RegisterForm)
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
		fundamentalAnalysisApiRoutes.POST("/save", fundamentalAnalysisController.Save)
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
		organizationApiRoutes.POST("/getDatatables", organizationController.GetDatatables)
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
	// #endregion

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Run(":10000")
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", c.Request.Header.Get("Origin"))
		// c.Writer.Header().Set("Access-Control-Allow-Origin", c.Request.Header.Get("Origin"))
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Access-Control-Allow-Origin, Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		fmt.Println(c.Request.Method)
		if c.Request.Method == "OPTIONS" {
			fmt.Println("OPTIONS")
			c.AbortWithStatus(200)
		} else {
			c.Next()
		}
	}
}
