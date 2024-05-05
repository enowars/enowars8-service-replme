package main

import (
	"fmt"
	"log"
	"net/http"
	"slices"

	"cafedodo/client"
	"cafedodo/lib"
	"cafedodo/orchestrator"
	"cafedodo/renderer"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
)

func main() {
	imageTag := "ptwhy"

	o := orchestrator.Instance()
	o.BuildImage("../ptwhy/", imageTag)

	engine := gin.Default()

	engine.Static("/static", "./static")

	store := memstore.NewStore([]byte("secret"))
	engine.Use(sessions.Sessions("session", store))

	ginHtmlRenderer := engine.HTMLRender
	engine.HTMLRender = &renderer.HTMLTemplRenderer {
		FallbackHtmlRenderer: ginHtmlRenderer,
	}

	engine.GET("/", func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		var ptySessions []string
		s := session.Get("ptySessions")
		if s == nil {
			ptySessions = []string{}
		} else {
			ptySessions = s.([]string)
		}
		e := ctx.Request.URL.Query().Get("error")
		fmt.Println(e);
		ctx.HTML(http.StatusOK, "", client.Home(ptySessions, e))
	})

	engine.GET("/session/prepare", func(ctx *gin.Context) {
		ptySession := ctx.Request.URL.Query().Get("ptySession")
		session := sessions.Default(ctx)
		s := session.Get("ptySessions")
		if s == nil || !slices.Contains(s.([]string), ptySession) {
			ctx.Redirect(http.StatusTemporaryRedirect, "/?error=session_not_found")
			return
		}

		c := o.GetContainer(ptySession)
		if c == nil {
			ctx.Redirect(http.StatusTemporaryRedirect, "/?error=container_not_found")
			return
		}

		err := o.StartContainer(c.ID)
		if err != nil {
			log.Print(err, " :failed to start container ", c.ID)
			ctx.Redirect(http.StatusTemporaryRedirect, "/?error=could_not_start_container")
		} else {
			c = o.GetContainer(ptySession)
			if len(c.Ports) == 0 {
				ctx.Redirect(http.StatusTemporaryRedirect, "/?error=no_port_for_container")
				return
			}
			ctx.Redirect(
				http.StatusTemporaryRedirect,
				fmt.Sprintf(
					"/session?ptySession=%s&port=%d",
					ptySession,
					c.Ports[0].PublicPort,
				),
			)
		}

	})

	engine.GET("/session", func(ctx *gin.Context) {
		ptySession := ctx.Request.URL.Query().Get("ptySession")
		port := ctx.Request.URL.Query().Get("port")
		ctx.HTML(http.StatusOK, "", client.Session(ptySession, port))
	})

	engine.GET("/api/session/private", func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		var ptySessions []string
		s := session.Get("ptySessions")
		ptySession := lib.RandomString(32)
		if s == nil {
			ptySessions = []string{}
		} else {
			ptySessions = s.([]string)
		}

		c := o.GetContainer(ptySession)

		if c != nil {
			log.Print(c.ID, ": container exists")
		} else {
			response, err := o.CreateContainer(orchestrator.RunContainerOptions {
				ImageTag: "ptwhy",
				ContainerName: ptySession,
				Ports: nat.PortMap {
					nat.Port("3000/tcp"): []nat.PortBinding {
						{
							HostIP: "127.0.0.1",
							HostPort: "0",
						},
					},
				},
			})
			log.Print(response.ID, ": container created")
			if err != nil {
				log.Print(err, ": container creation failed")
				ctx.Redirect(http.StatusMovedPermanently, "/?error=container_creation_failed")
			} 
		}
		
		ptySessions = append(ptySessions, ptySession)
		session.Set("ptySessions", ptySessions)
		session.Save()

		ctx.Redirect(http.StatusTemporaryRedirect, "/")
	})

	engine.GET("/api/session/private/:ptySession/delete", func(ctx *gin.Context) {
		ptySession := ctx.Param("ptySession")
		session := sessions.Default(ctx)
		var ptySessions []string
		s := session.Get("ptySessions")
		if s == nil {
			ptySessions = []string{}
		} else {
			ptySessions = s.([]string)
		}

		c := o.GetContainer(ptySession)

		if c != nil {
			err := o.Client.ContainerKill(o.Context, c.ID, "")

			if err != nil {
				log.Print(err);
				ctx.Redirect(http.StatusTemporaryRedirect, "/?error=could_not_kill_container")
				return
			}

			err = o.Client.ContainerRemove(o.Context, c.ID, container.RemoveOptions{})

			if err != nil {
				log.Print(err);
				ctx.Redirect(http.StatusTemporaryRedirect, "/?error=could_not_remove_container")
				return
			}
		}

		for i, currPtySession := range ptySessions {
			if currPtySession == ptySession {
				ptySessions = append(ptySessions[:i], ptySessions[i+1:]...)
				break
			}
		}
		session.Set("ptySessions", ptySessions)
		session.Save()
		ctx.Redirect(http.StatusTemporaryRedirect, "/")
	})

	engine.Run(":8080")
}

