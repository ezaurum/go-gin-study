package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/sessions"
	"net/http"
	"log"
	"time"
)
func setupRouter() *gin.Engine {
	r := gin.Default()

	store := sessions.NewCookieStore([]byte("secret"))

	handlerFunc := sessions.Sessions("WEB_SESSION", store)
	r.Use(handlerFunc)

	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	r.GET("login2", Hn(func(context *gin.Context, session sessions.Session) {
		var count int
		//v := session.Values["count"]
		v := session.Get("count")

		if v == nil {
			count = 0
			//session.Values["foo"] = "bar"
			//session.Values[42] = 43
		} else {
			count = v.(int)
			count++
		}
		session.Set("count", count)
		session.Save()
		context.JSON(200, gin.H{"count": count})
	}))
	return r
}

type SessionHandler func(context *gin.Context, session sessions.Session)

func Hn (handler SessionHandler) gin.HandlerFunc {
	return func(context *gin.Context) {
		session := sessions.Default(context)
		handler(context, session)
		session.Save()
	}
}

func main() {
	r := setupRouter()
	//TODO r.LoadHTMLGlob("templates/**/*")
	r.GET("/testing", startPage)


	r.GET("login", func(context *gin.Context) {
		session := sessions.Default(context)
		session.Options(sessions.Options{
			MaxAge:10,
		})
		session.Set("sessionID", 1)
		session.Save()
		context.JSON(200, gin.H{"sessionid": 1})
	})

	r.GET("check", func(context *gin.Context) {
		session := sessions.Default(context)

		sessionID := session.Get("sessionID")

		if nil == sessionID {
			context.JSON(http.StatusForbidden, gin.H{"error":"login-required"})
		} else {
			context.JSON(http.StatusOK, gin.H{"sessionID":sessionID})
		}
	})

	r.GET("sessions", func(context *gin.Context) {
		session := sessions.Default(context)

		// Get a session. We're ignoring the error resulted from decoding an
		// existing session: Get() always returns a session, even if empty.
		//session, _ := store.Get(context.Request, "session-name")
		// Set some session values.

		var count int
		//v := session.Values["count"]
		v := session.Get("count")

		if v == nil {
			count = 0
			//session.Values["foo"] = "bar"
			//session.Values[42] = 43
		} else {
			count = v.(int)
			count++
		}
		session.Set("count", count)
		//session.Values["count"] = count

		// Save it before we write to the response/return from the handler.
		//session.Save(context.Request, context.Writer)
		session.Save()
		context.JSON(200, gin.H{"count": count})
	})

	//TODO html := template.Must(template.ParseFiles("file1", "file2"))
	//TODO router.SetHTMLTemplate(html)

	// This handler will match /user/john but will not match neither /user/ or /user
	r.GET("/user/:name", func(c *gin.Context) {
		name := c.Param("name")
		c.String(http.StatusOK, "Hello %s", name)
	})

	// However, this one will match /user/john/ and also /user/john/send
	// If no other routers match /user/john, it will redirect to /user/john/
	r.GET("/user/:name/*action", func(c *gin.Context) {
		name := c.Param("name")
		action := c.Param("action")
		message := name + " is " + action
		c.String(http.StatusOK, message)
	})

	// Query string parameters are parsed using the existing underlying request object.
	// The request responds to a url matching:  /welcome?firstname=Jane&lastname=Doe
	r.GET("/welcome", func(c *gin.Context) {
		firstname := c.DefaultQuery("firstname", "Guest")
		lastname := c.Query("lastname") // shortcut for c.Request.URL.Query().Get("lastname")

		c.String(http.StatusOK, "Hello %s %s", firstname, lastname)

	})

	r.Run() // listen and serve on 0.0.0.0:8080
}
type Person struct {
	Name     string    `form:"name"`
	Address  string    `form:"address"`
	Birthday time.Time `form:"birthday" time_format:"2006-01-02" time_utc:"1"`
	Colors []string `form:"colors[]"`
}

func startPage(c *gin.Context) {
	var person Person
	// If `GET`, only `Form` binding engine (`query`) used.
	// If `POST`, first checks the `content-type` for `JSON` or `XML`, then uses `Form` (`form-data`).
	// See more at https://github.com/gin-gonic/gin/blob/master/binding/binding.go#L48
	if c.ShouldBind (&person) == nil {
		log.Println(person.Name)
		log.Println(person.Address)
		log.Println(person.Birthday)
		log.Println(person.Colors)
	}

	c.JSON(http.StatusOK, gin.H{"color": person.Colors})
}