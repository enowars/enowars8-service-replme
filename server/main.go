package main

import (
	"cafedodo/components"
	"cafedodo/lib"
	"cafedodo/orchestrator"
	"cafedodo/renderer"
	"fmt"
	"log"
	"net/http"
	"slices"

	"github.com/docker/go-connections/nat"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
)

func main() {

	port := 6969

	imageTag := "ptwhy"

	o := orchestrator.Instance()
	o.BuildImage("../ptwhy/", imageTag)

	o.GetContainer()

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
		e := c.Request.URL.Query().Get("error")
		fmt.Println(e);
		c.HTML(http.StatusOK, "", components.Home(ptySessions, e))
	})

	engine.GET("/session/:ptySession", func(c *gin.Context) {
		ptySession := c.Param("ptySession")
		session := sessions.Default(c)
		s := session.Get("ptySessions")
		fmt.Println(c.Request.URL.Query().Get("error"));
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
		ptySession := lib.RandomString(32)
		if s == nil {
			ptySessions = []string{}
		} else {
			ptySessions = s.([]string)
		}

		container := o.GetContainer("asdf")

		if container != nil {
			log.Print(container.ID, ": container exists")
		} else {
			response, err := o.CreateContainer(orchestrator.RunContainerOptions {
				ImageTag: "ptwhy",
				ContainerName: ptySession,
				Ports: nat.PortMap {
					"3000": []nat.PortBinding {
						{
							HostIP: "127.0.0.1",
							HostPort: fmt.Sprint(port),
						},
					},
				},
			})
			port++
			log.Print(response.ID, ": container created")
			if err != nil {
				log.Print(err, ": container creation failed")
				c.Redirect(http.StatusMovedPermanently, "/?error=container_creation_failed")
			} 
		}
		
		ptySessions = append(ptySessions, ptySession)
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

