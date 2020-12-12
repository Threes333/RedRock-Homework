package controller

import (
	"github.com/gin-gonic/gin"
	"leave_message/models"
	"leave_message/service"
	"net/http"
)

func SendMessage(c *gin.Context){
	switch service.SendMessage(c) {
	case 0:
		c.JSON(http.StatusOK,gin.H{
		"code" : "10000",
		"message" : "success",
	})
	case 1:
		c.JSON(http.StatusOK,gin.H{
		"code" : "20001",
		"message" : "接受用户不存在!",
		})
	case 2:
		c.JSON(http.StatusOK,gin.H{
			"code" : "20002",
			"message" : "插入留言数据失败!",
		})
	case 3:
		c.JSON(http.StatusOK,gin.H{
			"code" : "20003",
			"message" : "留言数据有为空",
		})
	}
}

func RegisterUser(c *gin.Context){
	switch service.RegisterUser(c) {
	case 0:
		c.JSON(http.StatusOK,gin.H{
			"code" : "10000",
			"message" : "success register",
		})
	case 1:
		c.JSON(http.StatusOK,gin.H{
			"code" : "20001",
			"message" : "Id已被注册",
		})
	case 2:
		c.JSON(http.StatusOK,gin.H{
			"code" : "20002",
			"message" : "用户名已被注册",
		})
	case 3:
		c.JSON(http.StatusOK,gin.H{
			"code" : "20003",
			"message" : "Id和用户名都已被注册",
		})
	case 4:
		c.JSON(http.StatusOK,gin.H{
			"code" : "20004",
			"message" : "Id或用户名或密码为空",
		})
	}
}

func LoginUser(c *gin.Context){
	switch service.LoginUser(c) {
	case 0:
		c.JSON(http.StatusOK,gin.H{
			"code" : "10000",
			"message" : "success login",
		})
	case 1:
		c.JSON(http.StatusOK,gin.H{
			"code" : "20001",
			"message" : "用户不存在",
		})
	case 2:
		c.JSON(http.StatusOK,gin.H{
			"code" : "20002",
			"message" : "密码输入错误",
		})
	case 3:
		c.JSON(http.StatusOK,gin.H{
			"code" : "20003",
			"message" : "用户名或密码为空",
		})
	}
}

func ViewMessage(c *gin.Context) {
	msg := service.ViewMessage(c)
	if len(msg) == 0{
		c.JSON(http.StatusOK,gin.H{
			"code" : "20000",
			"msg" : "无信息",
		})
		return
	}
	c.JSON(http.StatusOK,gin.H{
		"code" : "10000",
		"msg" :  msg,
	})
}

func MiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err1 := c.Cookie("id")
		password,err2 := c.Cookie("password")
		if err1 != nil || err2 != nil{
			c.JSON(http.StatusUnauthorized,gin.H{
				"err" : "Not login!",
			})
			c.Abort()
			return
		}
		user :=models.User{Id: id,Password: password}
		if res,_ :=models.LoginUser(&user); res != 0{
			c.JSON(http.StatusUnauthorized,gin.H{
				"err" : "Not login!",
			})
			c.Abort()
			return
		}
		c.Next()
		return
	}
}

func DeliverComment(c *gin.Context) {
	ok := service.DeliverComment(c)
	if !ok{
		c.JSON(http.StatusOK,gin.H{
			"code" : "20000",
			"msg" : "评论保存失败",
		})
	}else {
		c.JSON(http.StatusOK,gin.H{
			"code" : "10000",
			"msg" : "评论保存成功",
		})
	}
}

func DeliverReply(c *gin.Context) {
	ok := service.DeliverReply(c)
	if !ok{
		c.JSON(http.StatusOK,gin.H{
			"code" : "20000",
			"msg" : "回复保存失败",
		})
	}else{
		c.JSON(http.StatusOK,gin.H{
			"code" : "10000",
			"msg" : "回复保存成功",
		})
	}
}

func CancelLogin(c *gin.Context){
	c.SetCookie("id","",-1,"/","localhost",false,true)
	c.SetCookie("password","",-1,"/","localhost",false,true)
	c.SetCookie("username","",-1,"/","localhost",false,true)
	c.JSON(http.StatusOK,gin.H{
		"code" : "10000",
		"msg" : "退出登陆成功!",
	})
}
func MessageLikes(c *gin.Context) {
	ok := service.MessageLikes(c)
	if !ok{
		c.JSON(http.StatusOK,gin.H{
			"code" : "20000",
			"msg" : "留言点赞失败",
		})
	}else{
		c.JSON(http.StatusOK,gin.H{
			"code" : "10000",
			"msg" : "留言点赞成功",
		})
	}
}
func CommentLikes(c *gin.Context) {
	ok := service.CommentLikes(c)
	if !ok{
		c.JSON(http.StatusOK,gin.H{
			"code" : "20000",
			"msg" : "评论点赞失败",
		})
	}else{
		c.JSON(http.StatusOK,gin.H{
			"code" : "10000",
			"msg" : "评论点赞成功",
		})
	}
}
func ReplyLikes(c *gin.Context) {
	ok := service.ReplyLikes(c)
	if !ok{
		c.JSON(http.StatusOK,gin.H{
			"code" : "20000",
			"msg" : "回复点赞失败",
		})
	}else{
		c.JSON(http.StatusOK,gin.H{
			"code" : "10000",
			"msg" : "回复点赞成功",
		})
	}
}
func DeleteMessage(c *gin.Context) {
	ok := service.DeleteMessage(c)
	if !ok{
		c.JSON(http.StatusOK,gin.H{
			"code" : "20000",
			"msg" : "留言删除失败",
		})
	}else{
		c.JSON(http.StatusOK,gin.H{
			"code" : "10000",
			"msg" : "留言删除成功",
		})
	}
}
func DeleteComment(c *gin.Context) {
	ok := service.DeleteComment(c)
	if !ok{
		c.JSON(http.StatusOK,gin.H{
			"code" : "20000",
			"msg" : "评论删除失败",
		})
	}else{
		c.JSON(http.StatusOK,gin.H{
			"code" : "10000",
			"msg" : "评论删除成功",
		})
	}
}
func DeleteReply(c *gin.Context) {
	ok := service.DeleteReply(c)
	if !ok{
		c.JSON(http.StatusOK,gin.H{
			"code" : "20000",
			"msg" : "回复删除失败",
		})
	}else{
		c.JSON(http.StatusOK,gin.H{
			"code" : "10000",
			"msg" : "回复删除成功",
		})
	}
}


