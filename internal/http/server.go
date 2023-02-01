package http

import (
	"github.com/gin-gonic/gin"
	"html/template"
	"io/fs"
	"net/http"
	ws "smshog/internal/websocket"
	"strings"
	"time"
)

func NewHttpServer(files fs.FS, hub *ws.Hub, messages *[]Message) *gin.Engine {
	tmpl := template.Must(template.New("").ParseFS(files, "templates/*.html"))
	js, _ := fs.Sub(files, "assets/js")
	css, _ := fs.Sub(files, "assets/css")

	r := gin.Default()
	//r.LoadHTMLGlob("templates/*")
	r.SetHTMLTemplate(tmpl)
	r.StaticFS("/js", http.FS(js))
	r.StaticFS("/css", http.FS(css))
	r.StaticFileFS("/favicon.ico", "assets/favicon.ico", http.FS(files))
	//r.StaticFileFS("/favicon.ico", "/assets/favicon.ico", httpFS)

	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
	sendHandler := getSendHandler(messages, hub)
	r.GET("/sys/send.php", sendHandler)
	r.POST("/sys/send.php", sendHandler)
	r.GET("/api/messages", func(c *gin.Context) {
		var req ListMessagesRequest
		if err := c.BindQuery(&req); err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}

		var filtered []Message

		if req.Search == "" {
			filtered = *messages
		} else {
			for _, msg := range *messages {
				if strings.Contains(msg.Text, req.Search) || strings.Contains(msg.Phone, req.Search) {
					filtered = append(filtered, msg)
				}
			}
		}

		from := req.Offset
		to := min(req.Offset+req.Limit, len(filtered))

		var slice []Message

		if from >= len(filtered) {
			slice = make([]Message, 0)
		} else {
			slice = filtered[from:to]
		}

		c.JSON(http.StatusOK, gin.H{
			"items":  slice,
			"offset": from,
			"count":  max(0, to-from),
			"total":  len(filtered),
		})
	})
	r.DELETE("/api/messages", func(c *gin.Context) {
		*messages = nil
		c.String(http.StatusOK, "success")
	})
	r.GET("/api/websocket", func(c *gin.Context) {
		ws.ServeWS(hub, c.Writer, c.Request)
	})
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})

	return r
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func getSendHandler(messages *[]Message, hub *ws.Hub) func(c *gin.Context) {
	return func(c *gin.Context) {
		var req MessageRequest
		if err := c.ShouldBind(&req); err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}

		for _, phone := range req.GetPhones() {
			message := Message{
				Sender:    req.Sender,
				Phone:     phone,
				Text:      req.Mes,
				CreatedAt: time.Now(),
			}
			*messages = append(*messages, message)
		}

		c.String(http.StatusOK, "success")

		hub.Broadcast <- []byte("{\"event\":\"updated\"}")
	}
}
