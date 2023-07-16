package doc

import (
	"os"

	"github.com/gin-gonic/gin"
)

func GetAPIDoc(source string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		w := ctx.Writer
		ctx.Header("Content-Type", "text/html; charset=utf-8")

		path := "../../static/templates/doc/" + source + ".html"

		htmlTemplate, err := os.ReadFile(path)

		if err != nil {
			panic(err)
		}

		_, err = w.Write(htmlTemplate)
		if err != nil {
			panic(err)
		}

		ctx.AbortWithStatus(200)
	}
}
