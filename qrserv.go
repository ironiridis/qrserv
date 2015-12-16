package main

import (
	"bytes"
	"encoding/base64"

	"github.com/gin-gonic/gin"
	"github.com/ironiridis/qrserv/vendor/rsc/qr"
)

func encodableString(c *gin.Context) string {
	// note that we tried to use c.Request.URL.Path and c.Request.URL.RawPath
	// here, but using c.Param ends up being (almost) good enough... the only
	// problem is that our *Parameter callout still includes the leading slash

	if len(c.Request.URL.RawQuery) > 0 {
		return c.Param("URL")[1:] + "?" + c.Request.URL.RawQuery
	}
	return c.Param("URL")[1:]
}

// PNGRequest renders a PNG file
func PNGRequest(c *gin.Context) {
	c.Header("X-Attribution", "Uses code.google.com/p/rsc/qr")
	qrc, err := qr.Encode(encodableString(c), qr.L)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	// TODO: caching directives
	c.Data(200, "image/png", qrc.PNG())
}

// HTMLRequest renders an HTML document with an embedded PNG <img>
func HTMLRequest(c *gin.Context) {
	c.Header("X-Attribution", "Uses code.google.com/p/rsc/qr")
	qrc, err := qr.Encode(encodableString(c), qr.L)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	// TODO: caching directives
	var out bytes.Buffer
	out.Write([]byte("<img src='data:image/png;base64,"))
	pngenc := base64.NewEncoder(base64.StdEncoding, &out)
	pngenc.Write(qrc.PNG())
	pngenc.Close()
	out.Write([]byte("'>"))
	c.Data(200, "text/html", out.Bytes())
}

func main() {
	r := gin.Default()
	r.GET("/png/*URL", PNGRequest)
	r.GET("/html/*URL", HTMLRequest)
	r.Run("127.0.0.1:8080")
}
