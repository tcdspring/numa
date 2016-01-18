package admin

import (
	"encoding/base64"
	"fmt"
	"numa/models"
	"numa/toolkit"
	"os"
	"runtime"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

// IndexController 后台首页
type IndexController struct {
	BaseController
}

// Index 后台首页
func (m *IndexController) Index() {
	m.Data["title"] = "首页"
	m.Data["hostname"], _ = os.Hostname()
	m.Data["goversion"] = runtime.Version()
	m.Data["os"] = runtime.GOOS
	m.Data["cpunum"] = runtime.NumCPU()
	m.Data["arch"] = runtime.GOARCH
	m.display()
}

// Login 登录处理
func (m *IndexController) Login() {

	if m.Ctx.Request.Method == "GET" {

	} else if m.Ctx.Request.Method == "POST" {

		account := strings.TrimSpace(m.GetString("account"))
		password := strings.TrimSpace(m.GetString("password"))
		remember := m.GetString("remember")

		// fmt.Println("account = " + account)
		// fmt.Println("password = " + password)

		if account != "" && password != "" {
			var admin models.Admin
			admin.Account = account

			// if account == "admin" && password == "123456" {
			// 	admin.NickName = "我是管理员"
			// 	admin.LastIP = m.GetClientIP()
			// 	admin.Password = toolkit.SHA256([]byte(password))
			//     admin.LastTime = toolkit.GetTime()
			// 	admin.Token = toolkit.GenUID()
			// 	admin.Email = "lurma@qq.com"
			// 	admin.Status = 1
            //     admin.Head = "default.png"
			// 	admin.Insert()
			// 	m.Redirect(beego.AppConfig.String("adminurl"), 302)
			// 	return
			// }

			if admin.Read("account") != nil || admin.Password != toolkit.SHA256([]byte(password)) {
				// 用户名或密码输入错误
				fmt.Println("账号或密码填写错误")
				m.Data["errmsg"] = "账号或密码填写错误"
			} else {
				// 登录成功
				token := toolkit.GenUID()

				admin.LastIP = m.GetClientIP()
				admin.LastTime = toolkit.GetTime()
				admin.Token = token
				admin.Update()

				key := []byte(beego.AppConfig.String("aeskey"))
				result, err := toolkit.AesEncrypt([]byte(m.GetClientIP()+"|"+token), key)
				if err != nil {
					return
				}

				auth := base64.StdEncoding.EncodeToString(result)
				fmt.Println("auth = " + auth)

				if remember == "yes" {
					m.Ctx.SetCookie("auth", auth, 7*86400)
				} else {
					m.Ctx.SetCookie("auth", auth)
				}

				m.Redirect(beego.AppConfig.String("adminurl"), 302)
			}
		} else {
			// 用户名或密码为空
			fmt.Println("账号或密码为空")
			m.Data["errmsg"] = "账号或密码不能为空"
		}
	}

	m.TplNames = beego.AppConfig.String("adminpath") + "/login.html"
}

// Logout 注销登录
func (m *IndexController) Logout() {
	m.Ctx.SetCookie("auth", "")
	m.Redirect(beego.AppConfig.String("adminurl")+"/login", 302)
}

// Profile 查看并修改管理员信息
func (m *IndexController) Profile() {

	admin := models.Admin{Account: m.userName}
	o := orm.NewOrm()
	if err := o.Read(&admin, "Account"); err != nil {
		m.error(err.Error())
	}

	if m.Ctx.Request.Method == "POST" {
		fmt.Println("收到POST请求")

		nickname := strings.TrimSpace(m.GetString("nickname"))
		email := strings.TrimSpace(m.GetString("email"))

		fmt.Printf("nickname=%s, email=%s\n", nickname, email)

		admin.NickName = nickname
		admin.Email = email

		if num, err := o.Update(&admin, "nick_name", "email"); err == nil {
			if num > 0 {
				m.Data["result"] = true
			} else {
				m.Data["errmsg"] = "更新失败"
			}
		} else {
			m.error(err.Error())
			m.Data["errmsg"] = "更新失败"
		}

	}
    
    format := admin.LastTime.Format("2006-01-02 15:04:05")
    m.Data["time"] = format
    m.Data["head"] = beego.AppConfig.String("headpath") + "/" + admin.Head
    
    fmt.Printf("head = %s\n", admin.Head)
    
    m.Data["admin"] = admin
    
    m.display("profile")
}

// Password 修改登录密码
func (m *IndexController) Password() {

	admin := models.Admin{Account: m.userName}
	if err := orm.NewOrm().Read(&admin, "Account"); err != nil {
		m.error(err.Error())
	}

	if m.Ctx.Request.Method == "POST" {

		var errmsg string
		oldPassword := strings.TrimSpace(m.GetString("old-password"))
		newPassword := strings.TrimSpace(m.GetString("new-password"))
		confirmPassword := strings.TrimSpace(m.GetString("confirm-password"))

		if newPassword != "" {
			if oldPassword == "" || toolkit.SHA256([]byte(oldPassword)) != admin.Password {
				errmsg = "原密码输入错误"
			} else if len(newPassword) < 6 {
				errmsg = "密码长度不能少于6个字符"
			} else if newPassword != confirmPassword {
				errmsg = "两次输入的密码不一致"
			}

			if len(errmsg) == 0 {
				admin.Password = toolkit.SHA256([]byte(newPassword))
				admin.Update("password")
				m.Data["result"] = true
			} else {
				m.Data["errmsg"] = errmsg
			}
		}

	}

	m.Data["admin"] = admin
	m.display("password")
}

// Head 上传头像
func (m *IndexController) Head() {
    
    admin := models.Admin{Account: m.userName}
	o := orm.NewOrm()
	if err := o.Read(&admin, "Account"); err != nil {
		m.error(err.Error())
	}
    
    if m.Ctx.Request.Method == "POST" {
        fmt.Printf("upload head...\n")
    }
    
    m.Data["admin"] = admin
    m.display("profile")
}