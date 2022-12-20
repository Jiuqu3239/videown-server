package ginlet

import (
	// "videown/internal/middleware/limiter"
	// "log"
	"net/http"
	"strings"
	"time"

	"videown-server/internal/ginlet/api"
	"videown-server/internal/ginlet/middleware/logger"
	"videown-server/internal/ginlet/resp"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())          // Custom logging middleware
	router.Use(logger.LoggerToFile()) // Logs are redirected to files
	router.Use(cors.Default())
	router.Use(TrimGetSuffix())
	router.Use(gin.CustomRecovery(func(c *gin.Context, err any) {
		resp.ErrorWithHttpStatus(c, err.(error), http.StatusInternalServerError)
		c.Abort()
	}))

	// if global.AppSetting.Limiter.IsOpen {
	// 	limiterObj, err := limiter.NewLimiter(global.RedisCli)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	router.Use(limiter.RateMiddleware(limiterObj, global.AppSetting.Limiter.Gap, global.AppSetting.Limiter.Count))
	// }

	registerEndpointsForSys(router)
	registerEndpointsForVideo(router)
	registerEndpointsForNFT(router)
	router.StaticFile("/favicon.ico", "./static/favicon.ico")
	return router
}

func TrimGetSuffix() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == http.MethodGet {
			req := c.Request.RequestURI
			idx := strings.LastIndex(req, "&")
			if idx > 0 {
				c.Request.RequestURI = req[0:idx]
			}
		}
		c.Next()
	}
}

func registerEndpointsForVideo(router *gin.Engine) {
	v := api.NewVideoAPI()
	g := router.Group("/video")
	g.PUT("/list", v.QueryVideos)
	g.PUT("/search", v.SearchVideos)
	g.PUT("/cover", v.UploadVideoCoverImg)
	g.GET("/cover", v.DownloadVideoCoverImg)
	g.GET("/views", v.AddVideoViews)
}

func registerEndpointsForNFT(router *gin.Engine) {
	n := api.NewNftAPI()
	g := router.Group("/nft")
	//g.Use(auth.AuthRequiredForUser)
	{
		g.PUT("/create", n.CreateVideoMetadata)
		g.PUT("/mint/:act", n.MintNFT)
		g.PUT("/purchase/:act", n.BuyNFT)
		g.PUT("/transfer/:act", n.TransferNFT)
		g.PUT("/change/status/:act", n.ChangeSellingStatus)
		g.PUT("/change/price/:act", n.ChangeSellingPrice)
		g.PUT("/activity/list", n.QueryActivities)
	}
}

func registerEndpointsForSys(router *gin.Engine) {
	s := api.NewSysAPI()
	g := router.Group("/sys")
	{
		g.GET("/server-ts", func(c *gin.Context) { resp.Ok(c, time.Now().Unix()) })
		g.PUT("/login", s.UserLogin)
		//g.POST("/change-pwd", auth.AuthRequiredForUser, s.ChangePassword)
	}
}
