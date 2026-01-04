package bid_controller

import (
	"context"
	"net/http"

	"github.com/robsonrg/goexpert-labs-auction/configuration/rest_err"
	"github.com/robsonrg/goexpert-labs-auction/internal/infra/api/web/validation"
	"github.com/robsonrg/goexpert-labs-auction/internal/usecase/bid_usecase"

	"github.com/gin-gonic/gin"
)

type BidController struct {
	bidUseCase bid_usecase.BidUseCaseInterface
}

func NewBidController(bidUseCase bid_usecase.BidUseCaseInterface) *BidController {
	return &BidController{
		bidUseCase: bidUseCase,
	}
}

func (u *BidController) CreateBid(c *gin.Context) {
	var bidInputDTO bid_usecase.BidInputDTO

	if err := c.ShouldBindJSON(&bidInputDTO); err != nil {
		restErr := validation.ValidateErr(err)

		c.JSON(restErr.Code, restErr)
		return
	}

	bidId, err := u.bidUseCase.CreateBid(context.Background(), bidInputDTO)
	if err != nil {
		restErr := rest_err.ConvertError(err)

		c.JSON(restErr.Code, restErr)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"bid_id": bidId})
}
