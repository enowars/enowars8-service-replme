package main

import (
	"cafedodo/components"
	"cafedodo/renderer"
	"cafedodo/lib"
	"net/http"
	"slices"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
)

func main() {

	engine := gin.Default()

	engine.Static("/static", "./static")

	store := memstore.NewStore([]byte("secret"))
	engine.Use(sessions.Sessions("session", store))

	ginHtmlRenderer := engine.HTMLRender
	engine.HTMLRender = &renderer.HTMLTemplRenderer {
		FallbackHtmlRenderer: ginHtmlRenderer,
	}

	engine.GET("/", func(c *gin.Context) {
		session := sessions.Default(c)
		var ptySessions []string
		s := session.Get("ptySessions")
		if s == nil {
			ptySessions = []string{}
		} else {
			ptySessions = s.([]string)
		}
		c.HTML(http.StatusOK, "", components.Home(ptySessions))
	})

	engine.GET("/session/:ptySession", func(c *gin.Context) {
		ptySession := c.Param("ptySession")
		session := sessions.Default(c)
		s := session.Get("ptySessions")
		if s == nil || !slices.Contains(s.([]string), ptySession) {
			c.Redirect(http.StatusTemporaryRedirect, "/?error=session_not_found")
		} else {
			c.HTML(http.StatusOK, "", components.Session(ptySession))
		}
	})

	engine.GET("/api/session/private", func(c *gin.Context) {
		session := sessions.Default(c)
		var ptySessions []string
		s := session.Get("ptySessions")
		v := lib.RandomString(32)
		if s == nil {
			ptySessions = []string{}
		} else {
			ptySessions = s.([]string)
		}
		ptySessions = append(ptySessions, v)
		session.Set("ptySessions", ptySessions)
		session.Save()
		c.Redirect(http.StatusTemporaryRedirect, "/")
	})

	engine.GET("/api/session/private/:ptySession/delete", func(c *gin.Context) {
		ptySession := c.Param("ptySession")
		session := sessions.Default(c)
		var ptySessions []string
		s := session.Get("ptySessions")
		if s == nil {
			ptySessions = []string{}
		} else {
			ptySessions = s.([]string)
		}
		for i, currPtySession := range ptySessions {
			if currPtySession == ptySession {
				ptySessions = append(ptySessions[:i], ptySessions[i+1:]...)
				break
			}
		}
		session.Set("ptySessions", ptySessions)
		session.Save()
		c.Redirect(http.StatusTemporaryRedirect, "/")
	})

	engine.Run(":8080")
}

