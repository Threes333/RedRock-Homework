package service

import (
	"github.com/gin-gonic/gin"
	"leave_message/models"
	"log"
)

func RegisterUser(c *gin.Context)(res int) {
	var user models.User
	err := c.ShouldBind(&user)
	if err != nil{
		log.Println("接受注册数据失败!")
		res = 4
		return res
	}
	res = models.RegisterUser(&user)
	return res
}

func LoginUser(c *gin.Context)(res int){
	var user models.User
	user.Id = c.PostForm("id")
	user.Password = c.PostForm("password")
	if user.Id == "" || user.Password == ""{
		res = 3
		return res
	}
	res,user.Username = models.LoginUser(&user)
	if res == 0{
		c.SetCookie("id",c.PostForm("id"),240,"/","localhost",false,true)
		c.SetCookie("password",c.PostForm("password"),240,"/","localhost",false,true)
		c.SetCookie("username",user.Username,240,"/","localhost",false,true)
	}
	return res
}

