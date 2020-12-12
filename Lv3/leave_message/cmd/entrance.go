package cmd

import (
	"github.com/gin-gonic/gin"
	"leave_message/controller"
)
func Entrance() {
	r := gin.Default()
	r.GET("/login",controller.LoginUser)
	r.POST("/register",controller.RegisterUser)
	r.POST("/message",controller.MiddleWare(),controller.SendMessage)
	r.GET("/message",controller.ViewMessage)
	r.POST("/comment",controller.MiddleWare(),controller.DeliverComment)
	r.POST("/reply",controller.MiddleWare(),controller.DeliverReply)
	r.DELETE("/login",controller.MiddleWare(),controller.CancelLogin)
	r.PUT("/message",controller.MiddleWare(),controller.MessageLikes)
	r.PUT("/comment",controller.MiddleWare(),controller.CommentLikes)
	r.PUT("/reply",controller.MiddleWare(),controller.ReplyLikes)
	r.DELETE("/message",controller.MiddleWare(),controller.DeleteMessage)
	r.DELETE("/comment",controller.MiddleWare(),controller.DeleteComment)
	r.DELETE("/reply",controller.MiddleWare(),controller.DeleteReply)
	r.Run(":8080")
}
