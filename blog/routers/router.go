package routers

import (
	"blog/controllers"
	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	beego.Router("/", &controllers.MainController{})

	beego.Router("/register", &controllers.UserController{}, "post:RegistUser")
	beego.Router("/register", &controllers.UserController{})
	beego.Router("/login", &controllers.UserController{}, "post:Login")
	beego.Router("/login", &controllers.UserController{})
}
