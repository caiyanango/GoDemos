package controllers

import (
	beego "github.com/beego/beego/v2/server/web"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	c.Data["form_hidden"] = ""
	c.Data["login"] = "登录"
	c.Data["register"] = "注册"
	c.Data["p_hidden"] = "hidden"
	c.Layout = "index.html"
	c.TplName = "header.html"
}
