package controllers

import (
	"blog/models"
	"github.com/beego/beego/v2/adapter/logs"
	beego "github.com/beego/beego/v2/server/web"
	"strings"
)

type UserController struct {
	beego.Controller
}

type LoginInfo struct {
	Username string `form:"username"`
	Password string `form:"password"`
}

type registerInfo struct {
	LoginInfo
	Confirm_pwd string `form:"confirm_password"`
}

func (c *UserController) Get() {
	var requestPath string
	c.Data["hidden"] = "hidden"
	c.Data["message"] = ""
	if strings.Contains(c.Ctx.Request.RequestURI, "?") {
		requestPath = strings.Split(c.Ctx.Request.RequestURI, "?")[0]
	} else {
		requestPath = c.Ctx.Request.RequestURI
	}
	switch requestPath {
	case "/register":
		c.TplName = "register.html"
	case "/login":
		c.TplName = "login.html"
	}
}

func (c *UserController) RegistUser() {
	var info registerInfo
	var user models.User
	if err := c.ParseForm(&info); err != nil {
		logs.Error(err)
		c.Data["hidden"] = ""
		c.Data["message"] = "服务器解析错误，注册失败"
		c.TplName = "register.html"
		return
	}
	models.MGRDB.Getuser(info.Username, &user)
	if user.ID != 0 {
		c.Data["hidden"] = ""
		c.Data["message"] = "用户已经存在"
		c.TplName = "register.html"
		return
	}
	if strings.Compare(info.Password, info.Confirm_pwd) != 0 {
		c.Data["hidden"] = ""
		c.Data["message"] = "两次输入的密码不一致，重新输入"
		c.TplName = "register.html"
		return
	}
	user = models.User{Username: info.Username, Password: info.Password}
	models.MGRDB.Adduser(&user)
	c.TplName = "register_success.html"
}

func (c *UserController) Login() {
	var info LoginInfo
	var user models.User
	if err := c.ParseForm(&info); err != nil {
		logs.Error(err)
		c.Data["hidden"] = ""
		c.Data["message"] = "服务器解析错误，登录失败"
		c.TplName = "login.html"
		return
	}
	models.MGRDB.Getuser(info.Username, &user)
	if user.ID == 0 {
		c.Data["hidden"] = ""
		c.Data["message"] = "用户不存在，请先注册"
		c.TplName = "login.html"
		return
	}
	if strings.Compare(user.Password, info.Password) != 0 {
		c.Data["hidden"] = ""
		c.Data["message"] = "密码错误"
		c.TplName = "login.html"
		return
	}
	c.Data["form_hidden"], c.Data["p_hidden"] = "hidden", ""
	c.Data["username"] = info.Username
	c.TplName = "header.html"
	c.Layout = "index.html"
}
