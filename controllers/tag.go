package controllers

import (
	"github.com/astaxie/beego"
	"activist_api/models"
	"encoding/json"
	"log"
)

type TagController struct {
	beego.Controller
}

func (c *TagController) sendError(message string, code float64) {
	var response models.DefaultResponse
	response.Ok = false
	response.Error = &models.Error{ UserMessage: message, Code: code }
	c.Data["json"] = &response
	c.ServeJSON()
}

func (c *TagController) sendErrorWithStatus(message string, code float64, status int) {
	c.Ctx.Output.SetStatus(status)
	var response models.DefaultResponse
	response.Ok = false
	response.Error = &models.Error{ UserMessage: message, Code: code }
	c.Data["json"] = &response
	c.ServeJSON()
}

func (c *TagController) sendSuccess() {
	var response models.DefaultResponse
	response.Ok = true
	c.Data["json"] = &response
	c.ServeJSON()
}

func (c *TagController) QueryTags() {
	tag := c.Input().Get("query")
	tags := models.GetTags(tag)
	c.Data["json"] = &tags
	c.ServeJSON()
}

func (c *TagController) GetTagStatus() {
	var tagName string
	var userId, tagId int64
	tagName = c.Ctx.Input.Param(":tag")

	if payload, err := validateToken(c.Ctx.Input.Header("Authorization")); err != nil {
		log.Println(err)
		c.sendErrorWithStatus("Invalid token. Access denied", 401, 401)
		return
	} else {
		user := models.GetUserById(int64(payload["sub"].(float64)))
		userId = user.Id
	}

	if tag := models.GetTag(tagName); tag == nil {
		c.sendError("Tag not found", 14)
		return
	} else {
		tagId = tag.Id
	}

	var response models.GetTagStatusResponse

	if tagStatus := models.GetTagStatus(userId, tagId); tagStatus == nil {
		response.HasStatus = false
	} else {
		response.HasStatus = true
		response.Status = tagStatus.Status
	}

	response.Ok = true
	c.Data["json"] = &response
	c.ServeJSON()
}

func (c *TagController) DeleteTagStatus() {
	var tagName string
	var userId, tagId int64
	tagName = c.Ctx.Input.Param(":tag")

	if payload, err := validateToken(c.Ctx.Input.Header("Authorization")); err != nil {
		log.Println(err)
		c.sendErrorWithStatus("Invalid token. Access denied", 401, 401)
		return
	} else {
		user := models.GetUserById(int64(payload["sub"].(float64)))
		userId = user.Id
	}

	if tag := models.GetTag(tagName); tag == nil {
		c.sendError("Tag not found", 14)
		return
	} else {
		tagId = tag.Id
	}

	if err := models.DeleteTagStatus(tagId, userId); err != nil {
		c.sendError("Couldn't delete tag status", 14)
		return
	}

	c.sendSuccess()
}

func (c *TagController) AddTagStatus() {
	var tag string
	var userId int64
	var status bool
	tag = c.Ctx.Input.Param(":tag")

	if payload, err := validateToken(c.Ctx.Input.Header("Authorization")); err != nil {
		log.Println(err)
		c.sendErrorWithStatus("Invalid token. Access denied", 401, 401)
		return
	} else {
		user := models.GetUserById(int64(payload["sub"].(float64)))
		userId = user.Id
	}

	var request models.AddFavHideTagRequest
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &request); err == nil {
		status = request.Status
	}

	if statusId := models.AddTagStatus(tag, userId, status); statusId == 0 {
		c.sendError("Couldn't add tag status", 14)
		return
	}
	c.sendSuccess()
}
