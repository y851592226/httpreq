package httpreq

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	ts *httptest.Server
)

type FooStruct struct {
	Foo string `json:"foo" xml:"foo"`
}

type EFooStruct struct {
	Foo int `json:"foo" xml:"foo"`
}

func TestHttpreq(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Httpreq Suite")
}

var _ = BeforeSuite(func() {
	setupRoute()
})

func setupRoute() {
	gin.SetMode(gin.ReleaseMode)
	handler := gin.New()
	handler.GET("/test/timeout", func(c *gin.Context) {
		time.Sleep(time.Second)
		c.String(200, "hello")
	})
	handler.GET("/test/get", func(c *gin.Context) {
		c.String(200, "this is a get request")
	})
	handler.POST("/test/post", func(c *gin.Context) {
		c.String(200, "this is a post request")
	})
	handler.Any("/test/response/error", func(c *gin.Context) {
		c.Header("Content-Length", "100")
		c.String(200, "hello world")
	})
	handler.Any("/test/response/success", func(c *gin.Context) {
		c.Header("server", "gin")
		c.SetCookie("package", "github.com/y851592226/httpreq", 1000, "", "", false, true)
		c.String(200, "welcome")
	})
	handler.Any("/test/response/bind/json", func(c *gin.Context) {
		c.Header("server", "gin")
		c.SetCookie("package", "github.com/y851592226/httpreq", 1000, "", "", false, true)
		c.JSON(200, FooStruct{"bar-json"})
	})
	handler.Any("/test/response/bind/xml", func(c *gin.Context) {
		c.Header("server", "gin")
		c.SetCookie("package", "github.com/y851592226/httpreq", 1000, "", "", false, true)
		c.XML(200, FooStruct{"bar-xml"})
	})
	ts = httptest.NewServer(handler)
}
