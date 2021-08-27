package main

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
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

	applicationUserRepository         repository.ApplicationUserRepository         = repository.NewApplicationUserRepository(db)
	applicationMenuCategoryRepository repository.ApplicationMenuCategoryRepository = repository.NewApplicationMenuCategoryRepository(db)
	applicationMenuRepository         repository.ApplicationMenuRepository         = repository.NewApplicationMenuRepository(db)
	membershipRepository              repository.MembershipRepository              = repository.NewMembershipRepository(db)
	filemasterRepository              repository.FilemasterRepository              = repository.NewFilemasterRepository(db)
	emitenRepository                  repository.EmitenRepository                  = repository.NewEmitenRepository(db)
	emitenCategoryRepository          repository.EmitenCategoryRepository          = repository.NewEmitenCategoryRepository(db)
	webinarCategoryRepository         repository.WebinarCategoryRepository         = repository.NewWebinarCategoryRepository(db)
	sectorRepository                  repository.SectorRepository                  = repository.NewSectorRepository(db)
	tagRepository                     repository.TagRepository                     = repository.NewTagRepository(db)
	technicalAnalysisRepository       repository.TechnicalAnalysisRepository       = repository.NewTechnicalAnalysisRepository(db)
	fundamentalAnalysisRepository     repository.FundamentalAnalysisRepository     = repository.NewFundamentalAnalysisRepository(db)
	fundamentalAnalysisTagRepository  repository.FundamentalAnalysisTagRepository  = repository.NewFundamentalAnalysisTagRepository(db)

	authService                    service.AuthService                    = service.NewAuthService(applicationUserRepository)
	jwtService                     service.JWTService                     = service.NewJWTService()
	applicationUserService         service.ApplicationUserService         = service.NewApplicationUserService(applicationUserRepository)
	applicationMenuCategoryService service.ApplicationMenuCategoryService = service.NewApplicationMenuCategoryService(applicationMenuCategoryRepository)
	applicationMenuService         service.ApplicationMenuService         = service.NewApplicationMenuService(applicationMenuRepository)
	membershipService              service.MembershipService              = service.NewMembershipService(membershipRepository)
	filemasterService              service.FilemasterService              = service.NewFilemasterService(filemasterRepository)
	emitenService                  service.EmitenService                  = service.NewEmitenService(emitenRepository)
	emitenCategoryService          service.EmitenCategoryService          = service.NewEmitenCategoryService(emitenCategoryRepository)
	webinarCategoryService         service.WebinarCategoryService         = service.NewWebinarCategoryService(webinarCategoryRepository)
	sectorService                  service.SectorService                  = service.NewSectorService(sectorRepository)
	tagService                     service.TagService                     = service.NewTagService(tagRepository)
	technicalAnalysisService       service.TechnicalAnalysisService       = service.NewTechnicalAnalysisService(technicalAnalysisRepository)
	fundamentalAnalysisService     service.FundamentalAnalysisService     = service.NewFundamentalAnalysisService(fundamentalAnalysisRepository)
	fundamentalAnalysisTagService  service.FundamentalAnalysisTagService  = service.NewFundamentalAnalysisTagService(fundamentalAnalysisTagRepository)

	authController                    controllers.AuthController                    = controllers.NewAuthController(authService, jwtService)
	applicationUserController         controllers.ApplicationUserController         = controllers.NewApplicationUserController(applicationUserService, jwtService)
	applicationMenuCategoryController controllers.ApplicationMenuCategoryController = controllers.NewApplicationMenuCategoryController(applicationMenuCategoryService, jwtService)
	applicationMenuController         controllers.ApplicationMenuController         = controllers.NewApplicationMenuController(applicationMenuService, jwtService)
	membershipController              controllers.MembershipController              = controllers.NewMembershipController(membershipService, jwtService)
	filemasterController              controllers.FilemasterController              = controllers.NewFilemasterController(filemasterService, jwtService)
	emitenController                  controllers.EmitenController                  = controllers.NewEmitenController(emitenService, jwtService)
	emitenCategoryController          controllers.EmitenCategoryController          = controllers.NewEmitenCategoryController(emitenCategoryService, jwtService)
	webinarCategoryController         controllers.WebinarCategoryController         = controllers.NewWebinarCategoryController(webinarCategoryService, jwtService)
	sectorController                  controllers.SectorController                  = controllers.NewSectorController(sectorService, jwtService)
	tagController                     controllers.TagController                     = controllers.NewTagController(tagService, jwtService)
	technicalAnalysisController       controllers.TechnicalAnalysisController       = controllers.NewTechnicalAnalysisController(technicalAnalysisService, jwtService)
	fundamentalAnalysisController     controllers.FundamentalAnalysisController     = controllers.NewFundamentalAnalysisController(fundamentalAnalysisService, jwtService)
	fundamentalAnalysisTagController  controllers.FundamentalAnalysisTagController  = controllers.NewFundamentalAnalysisTagController(fundamentalAnalysisTagService, jwtService)
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

	r.Use(location.New(location.Config{
		Scheme:  "http",
		Host:    "foo.com",
		Base:    "/base",
		Headers: location.Headers{Scheme: "X-Forwarded-Proto", Host: "X-Forwarded-For"},
	}))

	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store))

	r.Static("/assets", "./assets")
	r.Static("/script", "./templates/js")
	r.HTMLRender = createMyRender("templates/views/")

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

	r.POST("/login", func(c *gin.Context) {
		c.Request.ParseForm()
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
			isVerified, _isAdmin, token, message := authController.LoginForm(c, email, password)

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

		// c.HTML(
		// 	http.StatusOK, "logout",
		// 	gin.H{
		// 		"title": "logout",
		// 	},
		// )
		authController.Logout(c)
		c.Redirect(http.StatusMovedPermanently, "/")
	})

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

	authRoutes := r.Group("api/auth")
	{
		authRoutes.POST("/login", authController.Login)
		authRoutes.GET("/logout", authController.Logout)
	}

	emitenApiRoutes := r.Group("api/emiten")
	{
		emitenApiRoutes.POST("/getDatatables", emitenController.GetDatatables)
		emitenApiRoutes.GET("/lookup", emitenController.Lookup)
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
		applicationMenuApiRoutes.GET("/lookup", applicationMenuController.Lookup)
		applicationMenuApiRoutes.POST("/save", applicationMenuController.Save)
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

	webinarCategoryApiRoutes := r.Group("api/webinar_category")
	{
		webinarCategoryApiRoutes.POST("/getDatatables", webinarCategoryController.GetDatatables)
		webinarCategoryApiRoutes.GET("/getTree", webinarCategoryController.GetTree)
		webinarCategoryApiRoutes.GET("/lookup", webinarCategoryController.Lookup)
		webinarCategoryApiRoutes.POST("/save", webinarCategoryController.Save)
		webinarCategoryApiRoutes.GET("/getById/:id", webinarCategoryController.GetById)
		webinarCategoryApiRoutes.DELETE("/deleteById/:id", webinarCategoryController.DeleteById)
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
		membershipApiRoutes.POST("/save", membershipController.Save)
		membershipApiRoutes.GET("/getById/:id", membershipController.GetById)
		membershipApiRoutes.DELETE("/deleteById/:id", membershipController.DeleteById)
	}

	filemasterApiRoutes := r.Group("api/filemaster")
	{
		filemasterApiRoutes.POST("/upload/:id", filemasterController.Insert)
		filemasterApiRoutes.GET("/getAll", filemasterController.GetAll)
		filemasterApiRoutes.DELETE("/deleteByRecordId/:recordId", filemasterController.DeleteByRecordId)
	}

	applicationUserApiRoutes := r.Group("api/application_user")
	{
		applicationUserApiRoutes.POST("/getDatatables", applicationUserController.GetDatatables)
	}

	technicalAnalysisApiRoutes := r.Group("api/technical_analysis")
	{
		technicalAnalysisApiRoutes.POST("/getDatatables", technicalAnalysisController.GetDatatables)
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

	r.Run(":10000")
}
