package models

import (
	"database/sql"
	"fmt"
	"leave_message/dao"
	"log"
	"strings"
)
type LeaveMessage struct {
	Sender string `json:"sender" form:"sender"`
	Message  string `json:"message" form:"message"`
	Receiver string `json:"receiver" form:"receiver"`
}
type Comment struct {
	Sender string `json:"sender" form:"sender"`
	Comment string `json:"comment" form:"comment"`
	MessageId string `json:"message_id" form:"message_id"`
}
type Reply struct {
	Sender string `json:"sender" form:"sender"`
	Reply string `json:"reply" form:"reply"`
	CommentId string `json:"comment_id" form:"comment_id"`
	Receiver string `json:"receiver" form:"receiver"`
}
type ReplyText struct {
	Reply string
	Likes int
}
type CommentText struct {
	Comment string
	Likes int
	CmtReply []ReplyText
}
type MessageText struct {
	Message string
	Likes int
	MsgComment []CommentText
}
type Liked struct {
	id string
	MessageId string
	CommentId string
	ReplyId string
}
func SendMessage(msg *LeaveMessage)(res int) {
	sqlstr := "Select username from user where username = ?"
	stmt, _ := dao.DB.Prepare(sqlstr)
	var name interface{}
	err := stmt.QueryRow(msg.Receiver).Scan(name)
	if err == sql.ErrNoRows{
		log.Println("接收用户不存在!")
		res = 1
		return res
	}
	sqlstr = "insert into message (sender,message,receiver) values (?,?,?)"
	stmt, _ = dao.DB.Prepare(sqlstr)
	_, err = stmt.Exec(msg.Sender, msg.Message, msg.Receiver)
	if err != nil {
		log.Println("插入留言数据失败!")
		res = 2
		return res
	}
	res = 0
	return res
}

func DeliverComment(cmt *Comment)(err error){
	sqlstr := "Insert into comment (sender,comment,message_id) values (?,?,?)"
	stmt, _ := dao.DB.Prepare(sqlstr)
	_, err = stmt.Exec(cmt.Sender, cmt.Comment, cmt.MessageId)
	if err != nil{
		log.Println("保存评论失败!",err)
	}
	return err
}

func DeliverReply(rep *Reply)(err error){
	sqlstr := "Insert into reply (sender,reply,comment_id,receiver) values (?,?,?,?)"
	stmt, _ := dao.DB.Prepare(sqlstr)
	_,err = stmt.Exec(rep.Sender,rep.Reply,rep.CommentId,rep.Receiver)
	if err != nil{
		log.Println("保存回复失败!",err)
	}
	return err
}

func ViewMessage(msg *LeaveMessage) []MessageText {
	message := make([]MessageText,0)
	var rows *sql.Rows
	if msg.Receiver == ""{
		sqlstr := "Select message,id,likes from message where sender = ? order by id"
		stmt, _ := dao.DB.Prepare(sqlstr)
		rows, _= stmt.Query(msg.Sender)
	}else{
		sqlstr := "Select message,id,likes from message where sender = ? and receiver = ? order by id"
		stmt, _ := dao.DB.Prepare(sqlstr)
		rows, _= stmt.Query(msg.Sender,msg.Receiver)
	}
	if rows == nil{
		return nil
	}
	var i = 1
	for rows.Next(){
		var likes int
		var meg MessageText
		var id string
		_ = rows.Scan(&meg.Message,&id,&likes)
		meg.Message = fmt.Sprintf("%d. %s 给 %s 留言: '%s' ",i,msg.Sender,msg.Receiver,meg.Message)
		meg.Likes = likes
		meg.MsgComment = FindComment(id)
		message = append(message,meg)
		i++
	}
	return message
}

func FindComment(messageId string) []CommentText {
	comment := make([]CommentText,0)
	sqlstr := "Select sender,comment,id,likes from comment where message_id = ? order by likes"
	stmt, _ := dao.DB.Prepare(sqlstr)
	rows, _:= stmt.Query(messageId)
	var i = 1
	for rows.Next(){
		var likes int
		var cmt CommentText
		var sender,comm,id string
		_ = rows.Scan(&sender,&comm,&id,&likes)
		cmt.Comment = fmt.Sprintf("%d. %s 评论道 '%s' ",i,sender,comm)
		cmt.Likes = likes
		cmt.CmtReply = FindReply(id)
		comment = append(comment,cmt)
		i++
	}
	return comment
}
func FindReply(commentId string) []ReplyText {
	reply := make([]ReplyText,0)
	sqlstr := "Select sender,reply,receiver,likes from reply where comment_id = ? order by id"
	stmt, _ := dao.DB.Prepare(sqlstr)
	rows, _:= stmt.Query(commentId)
	var i = 1
	for rows.Next(){
		var rep ReplyText
		var likes int
		var sender, answer, receiver string
		_ = rows.Scan(&sender,&answer,&receiver,&likes)
		rep.Reply = fmt.Sprintf("%d. %s 回复 %s : '%s' ",i, sender, receiver, answer)
		rep.Likes = likes
		reply = append(reply,rep)
		i++
	}
	return reply
}

func CommentLikes(id string, commentId string,)bool{
	sqlstr := "SELECT comment_id from liked where user_id = ?"
	stmt, _ := dao.DB.Prepare(sqlstr)
	rows, _ := stmt.Query(id)
	for rows.Next(){
		var cmtId string
		_ = rows.Scan(&cmtId)
		if strings.Compare(strings.ToLower(commentId),strings.ToLower(cmtId)) == 0{
			log.Println("该用户已点赞")
			return false
		}
	}
	sqlstr = "Update comment set likes = likes + 1 where id = ?"
	tx,_ := dao.DB.Begin()
	stmt, _ = tx.Prepare(sqlstr)
	_, err := stmt.Exec(commentId)
	if err != nil{
		log.Println("评论点赞数增加失败",err)
		return false
	}
	sqlstr = "Insert into liked (user_id,comment_id) values (?,?)"
	stmt,_ =tx.Prepare(sqlstr)
	_, err = stmt.Exec(id, commentId)
	if err != nil{
		log.Println("用户点赞失败",err)
		_ = tx.Rollback()
		_ = tx.Commit()
		return false
	}
	_ = tx.Commit()
	return true
}
func MessageLikes(id string, messageId string,)bool{
	sqlstr := "SELECT message_id from liked where user_id = ?"
	stmt, _ := dao.DB.Prepare(sqlstr)
	rows, _ := stmt.Query(id)
	for rows.Next(){
		var msgId string
		_ = rows.Scan(&msgId)
		if strings.Compare(strings.ToLower(messageId),strings.ToLower(msgId)) == 0{
			log.Println("该用户已点赞")
			return false
		}
	}
	sqlstr = "Update message set likes = likes + 1 where id = ?"
	tx,_ := dao.DB.Begin()
	stmt, _ = tx.Prepare(sqlstr)
	_, err := stmt.Exec(messageId)
	if err != nil{
		log.Println("留言点赞数增加失败",err)
		return false
	}
	sqlstr = "Insert into liked (user_id,message_id) values (?,?)"
	stmt,_ =tx.Prepare(sqlstr)
	_, err = stmt.Exec(id, messageId)
	if err != nil{
		log.Println("用户点赞失败",err)
		_ = tx.Rollback()
		_ = tx.Commit()
		return false
	}
	_ = tx.Commit()
	return true
}
func ReplyLikes(id string, replyId string,)bool{
	sqlstr := "SELECT reply_id from liked where user_id = ?"
	stmt, _ := dao.DB.Prepare(sqlstr)
	rows, _ := stmt.Query(id)
	for rows.Next(){
		var repId string
		_ = rows.Scan(&repId)
		if strings.Compare(strings.ToLower(replyId),strings.ToLower(repId)) == 0{
			log.Println("该用户已点赞")
			return false
		}
	}
	sqlstr = "Update reply set likes = likes + 1 where id = ?"
	tx,_ := dao.DB.Begin()
	stmt, _ = tx.Prepare(sqlstr)
	_, err := stmt.Exec(replyId)
	if err != nil{
		log.Println("回复点赞数增加失败",err)
		return false
	}
	sqlstr = "Insert into liked (user_id,reply_id) values (?,?)"
	stmt,_ =tx.Prepare(sqlstr)
	_, err = stmt.Exec(id, replyId)
	if err != nil{
		log.Println("用户点赞失败",err)
		_ = tx.Rollback()
		_ = tx.Commit()
		return false
	}
	_ = tx.Commit()
	return true
}
func DeleteMessage(username string, messageId string) bool {
	sqlstr := "SELECT sender from message where id = ?"
	stmt,_ := dao.DB.Prepare(sqlstr)
	var sender string
	_ = stmt.QueryRow(messageId).Scan(&sender)
	if strings.Compare(strings.ToLower(sender),strings.ToLower(username)) != 0{
		log.Println("无法删除他人留言")
		return false
	}
	sqlstr = "Update message set message = '该留言已被删除' where id = ?"
	stmt,_ = dao.DB.Prepare(sqlstr)
	_,err := stmt.Exec(messageId)
	if err != nil{
		log.Println("留言删除失败",err)
		return false
	}
	return true
}
func DeleteComment(username string, commentId string) bool {
	sqlstr := "SELECT sender from comment where id = ?"
	stmt,_ := dao.DB.Prepare(sqlstr)
	var sender string
	_ = stmt.QueryRow(commentId).Scan(&sender)
	if strings.Compare(strings.ToLower(sender),strings.ToLower(username)) != 0{
		log.Println("无法删除他人评论")
		return false
	}
	sqlstr = "Update comment set comment = '该评论已被删除' where id = ?"
	stmt,_ = dao.DB.Prepare(sqlstr)
	_,err := stmt.Exec(commentId)
	if err != nil{
		log.Println("评论删除失败",err)
		return false
	}
	return true
}
func DeleteReply(username string, replyId string) bool {
	sqlstr := "SELECT sender from reply where id = ?"
	stmt,_ := dao.DB.Prepare(sqlstr)
	var sender string
	_ = stmt.QueryRow(replyId).Scan(&sender)
	if strings.Compare(strings.ToLower(sender),strings.ToLower(username)) != 0{
		log.Println("无法删除他人回复:")
		return false
	}
	sqlstr = "Update reply set reply = '该回复已被删除' where id = ?"
	stmt,_ = dao.DB.Prepare(sqlstr)
	_,err := stmt.Exec(replyId)
	if err != nil{
		log.Println("回复删除失败",err)
		return false
	}
	return true
}