package main

import (
	"bytes"
	"encoding/base64"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ironiridis/qrserv/vendor/rsc/qr"
)

// per RFC2616 sec 14.21: 1 year in the future is the maximum cache expiry
// pre-calculate this once at startup to avoid pointless CPU churn
var futureExpire = time.Now().AddDate(1, 0, 0).Format(time.RFC1123)

func encodableString(c *gin.Context) string {
	// note that we tried to use c.Request.URL.Path and c.Request.URL.RawPath
	// here, but using c.Param ends up being (almost) good enough... the only
	// problem is that our *Parameter callout still includes the leading slash

	if len(c.Request.URL.RawQuery) > 0 {
		return c.Param("URL")[1:] + "?" + c.Request.URL.RawQuery
	}
	return c.Param("URL")[1:]
}

func standardHeaders(c *gin.Context) {
	c.Header("X-Attribution", "Uses code.google.com/p/rsc/qr")
	c.Header("Cache-Control", "max-age=31536000") // 60*60*24*365 aka 1 year
	c.Header("Expires", futureExpire)
}

// PNGRequest renders a PNG file
func PNGRequest(c *gin.Context) {
	qrc, err := qr.Encode(encodableString(c), qr.L)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	standardHeaders(c)
	c.Data(200, "image/png", qrc.PNG())
}

// HTMLRequest renders an HTML document with an embedded PNG <img>
func HTMLRequest(c *gin.Context) {
	qrc, err := qr.Encode(encodableString(c), qr.L)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	standardHeaders(c)
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
