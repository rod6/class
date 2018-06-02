package main

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
	"github.com/rod6/class/ctrl"
)

func main() {
	e := echo.New()
	e.Logger.SetLevel(log.INFO)

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.RequestIDWithConfig(middleware.RequestIDConfig{Generator: func() string {
		id := uuid.New()
		return id.String()
	}}))

	// Routes
	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"message": "hello, world",
		})
	})

	// url: ip:1323/login
	// method: POST
	// header: "Content-type": "application/json"
	// body: '{"username":"ziang","password":"ziang"}'
	// 返回数据中包含一个token，这个token表示登陆成功，后续所有操作需要这个token
	// curl localhost:1323/login -X POST -H "Content-type: application/json" -d '{"username":"ziang","password":"ziang"}'
	e.POST("/login", ctrl.Login)

	// Restricted group
	r := e.Group("/r")
	r.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		Skipper:       middleware.DefaultSkipper,
		SigningMethod: middleware.AlgorithmHS256,
		ContextKey:    "token",
		TokenLookup:   "header:" + echo.HeaderAuthorization,
		AuthScheme:    "Bearer",
		Claims:        jwt.MapClaims{},
		SigningKey:    []byte("Ziang Qiu"),
	}))

	// url: ip:1323/r/course/list
	// method: POST
	// header: "Authorization: Bearer xxxxx" (xxxxx是上面login操作返回的token)
	// 返回值为一个json，格式是 {"course1":"teacher1", "course2":"teacher2"}, 可用于现实老师列表、课程列表等
	// curl localhost:1323/r/course/list -X POST -H "Content-type: application/json" -H "Authorization: Bearer xxxxxx"
	r.POST("/course/list", ctrl.CourseList)

	// url: ip:1323/r/course/add
	// method: POST
	// header: "Content-type": "application/json", "Authorization: Bearer xxxxx" (xxxxx是上面login操作返回的token)
	// body: '{"teacher":"ziang","course":"history"}'
	// curl localhost:1323/r/course/add -X POST -H "Content-type: application/json" -H "Authorization: Bearer xxxxxx" -d '{"teacher":"ziang", "course":"history"}'
	r.POST("/course/add", ctrl.CourseAdd)

	// url: ip:1323/r/course/student/list
	// method: POST
	// header: "Content-type": "application/json", "Authorization: Bearer xxxxx" (xxxxx是上面login操作返回的token)
	// body: '{"course":"history"}'
	// curl localhost:1323/r/course/student/list -X POST -H "Content-type: application/json" -H "Authorization: Bearer xxxxxx" -d '{"course":"history"}'
	r.POST("/course/student/list", ctrl.StudentList)

	// url: ip:1323/r/course/student/add
	// method: POST
	// header: "Content-type": "application/json", "Authorization: Bearer xxxxx" (xxxxx是上面login操作返回的token)
	// body: '{"course":"history", "student":"daoguo"}'
	// curl localhost:1323/r/course/student/add -X POST -H "Content-type: application/json" -H "Authorization: Bearer xxxxxx" -d '{"course":"history", "student":"ziang"}'
	r.POST("/course/student/add", ctrl.StudentAdd)

	// url: ip:1323/r/course/student/absent
	// method: POST
	// header: "Content-type": "application/json", "Authorization: Bearer xxxxx" (xxxxx是上面login操作返回的token)
	// body: '{"course":"history", "student":"daoguo", "memo":"此处写你提到的备注"}'
	// curl localhost:1323/r/course/student/absent -X POST -H "Content-type: application/json" -H "Authorization: Bearer xxxxxx" -d '{"course":"history", "student":"ziang", "memo":"i donot know"}'
	r.POST("/course/student/absent", ctrl.StudentAbsent)

	// url: ip:1323/r/course/student/absent
	// method: POST
	// header: "Content-type": "application/json", "Authorization: Bearer xxxxx" (xxxxx是上面login操作返回的token)
	// body: '{"course":"history"}'
	// 返回值是一个json，格式是 {"student1":["memo1-1", "memo1-2"], "student2":["memo2-1", "memo2-2", "memo2-3"]}，数组可用于统计个数，现实详细每次缺席的备注
	// curl localhost:1323/r/course/absent/list -X POST -H "Content-type: application/json" -H "Authorization: Bearer xxxxxx" -d '{"course":"history"}'
	r.POST("/course/absent/list", ctrl.AbsentList)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}
