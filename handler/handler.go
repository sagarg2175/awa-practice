package handler

import (
	"main/domain"
	"main/repo"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AwsProjHandler struct {
	svc *repo.AwsProjRepo
}

func NewAwsProjHandler(svc *repo.AwsProjRepo) *AwsProjHandler {
	return &AwsProjHandler{
		svc,
	}
}

type CreateRequest struct {
	Name   string `json:"name" validate:"required,customName"`
	Domain string `json:"domain" validate:"required,min=8"`
}

func (ah *AwsProjHandler) AwaProjCreate(ctx *gin.Context) {

	var req CreateRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		return
	}

	batch := domain.Master{
		Name:   &req.Name,
		Domain: &req.Domain,
	}

	err := ah.svc.CreateAwsProj(ctx, &batch)
	if err != nil {
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": []string{"User created successfully"},
	})

}
