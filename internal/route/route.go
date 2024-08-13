package route

import (
	"duval/internal/configuration"
	"duval/internal/route/docs"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

var engine *gin.Engine

func init() {
	if configuration.App.IsProd() {
		gin.SetMode(gin.ReleaseMode)
	}

	engine = gin.Default()

	// Global middleware
	// Logger middleware will write the logs to gin.DefaultWriter even if you set with GIN_MODE=release.
	// By default, gin.DefaultWriter = os.Stdout
	//engine.Use(gin.Logger())

	// Recovery middleware recovers from any panics and writes a 500 if there was one.

	//CORS middleware
	config := cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	engine.Use(cors.New(config))
	engine.Use(gin.Recovery())

}

func Serve() (err error) {
	err = attach(engine)
	if err != nil {
		panic(err)
	}

	err = engine.Run(configuration.App.Host + ":" + configuration.App.Port)
	if err != nil {
		panic(err)
	}

	return err
}

func attach(g *gin.Engine) (err error) {
	g.GET("/test", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	for i := 0; i < len(RootRoutesGroup); i++ {
		group := g.Group(RootRoutesGroup[i].Group)
		err = docs.GenerateDocumentation(group, RootRoutesGroup[i].Paths)
		if err != nil {
			panic(err)
		}
	}

	return err
}
