package service

import (
	"github.com/gin-gonic/gin"
	"leave_message/models"
	"log"
)
func SendMessage(c *gin.Context)(res int) {
	var msg models.LeaveMessage
	err := c.ShouldBind(&msg)
	if err != nil{
		log.Println("接受留言失败!",err)
		res = 3
		return res
	}
	res = models.SendMessage(&msg)
	return res
}

func ViewMessage(c *gin.Context) []models.MessageText {
	var msg models.LeaveMessage
	msg.Sender = c.PostForm("sender")
	msg.Receiver = c.PostForm("receiver")
	message := models.ViewMessage(&msg)
	return message
}

func DeliverComment(c *gin.Context) (ok bool) {
	var cmt models.Comment
	cmt.Sender,_ = c.Cookie("username")
	cmt.Comment = c.PostForm("comment")
	cmt.MessageId = c.PostForm("message_id")
	anonymous := c.PostForm("anonymous")
	if anonymous == "yes"{
		cmt.Sender = "匿名用户"
	}
	err := models.DeliverComment(&cmt)
	if err != nil{
		return false
	}
	return true
}

func DeliverReply(c *gin.Context) (ok bool){
	var rep models.Reply
	rep.Sender,_ = c.Cookie("username")
	rep.Receiver = c.PostForm("receiver")
	rep.CommentId = c.PostForm("comment_id")
	rep.Reply = c.PostForm("reply")
	anonymous := c.PostForm("anonymous")
	if anonymous == "yes"{
		rep.Sender = "匿名用户"
	}
	err := models.DeliverReply(&rep)
	if err != nil{
		return false
	}
	return true
}
func MessageLikes(c *gin.Context) bool {
	messageId := c.PostForm("message_id")
	id,_ := c.Cookie("id")
	ok := models.MessageLikes(id,messageId)
	return ok
}
func CommentLikes(c *gin.Context) bool {
	commentId := c.PostForm("comment_id")
	id,_ := c.Cookie("id")
	ok := models.CommentLikes(id,commentId)
	return ok
}
func ReplyLikes(c *gin.Context) bool {
	replyId := c.PostForm("reply_id")
	id,_ := c.Cookie("id")
	ok := models.ReplyLikes(id,replyId)
	return ok
}
func DeleteMessage(c *gin.Context) bool {
	messageId := c.PostForm("message_id")
	username,_ := c.Cookie("username")
	ok := models.DeleteMessage(username,messageId)
	return ok
}
func DeleteComment(c *gin.Context) bool {
	commentId := c.PostForm("comment_id")
	username,_ := c.Cookie("username")
	ok := models.DeleteComment(username,commentId)
	return ok
}
func DeleteReply(c *gin.Context) bool {
	replyId := c.PostForm("reply_id")
	username,_ := c.Cookie("username")
	ok := models.DeleteReply(username,replyId)
	return ok
}