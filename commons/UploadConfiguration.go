package commons

import "github.com/gin-gonic/gin"

type FileConfig struct {
	Path      string   `valid:"required,length(1|20)"`
	MaxSize   int      `valid:"required"`
	AllowType []string `valid:"required,length(0|100)"`
}

type ImageConfig struct {
	Path      string `valid:"required,length(1|20)"`
	MaxSize   int    `valid:"required"`
	Thumbnail ThumbnailConfig
}

type ThumbnailConfig struct {
	Path      string `valid:"required,length(1|20)"`
	MaxWidth  int    `valid:"required"`
	MaxHeight int    `valid:"required"`
}

type TConfig struct {
	Path      string `valid:"required,length(1|20)"`
	UrlPrefix string `valid:"required,length(0|20)"`
	File      FileConfig
	Image     ImageConfig
}

type Uploader struct {
	Upload   *gin.RouterGroup
	Download *gin.RouterGroup
	Config   TConfig
	Engine   *gin.Engine
}
