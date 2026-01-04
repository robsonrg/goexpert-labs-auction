package bid_usecase

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/robsonrg/goexpert-labs-auction/configuration/logger"
	"github.com/robsonrg/goexpert-labs-auction/internal/entity/bid_entity"
	"github.com/robsonrg/goexpert-labs-auction/internal/internal_error"
)

type BidInputDTO struct {
	UserId    string  `json:"user_id"`
	AuctionId string  `json:"auction_id"`
	Amount    float64 `json:"amount"`
}

type BidOutputDTO struct {
	Id        string    `json:"id"`
	UserId    string    `json:"user_id"`
	AuctionId string    `json:"auction_id"`
	Amount    float64   `json:"amount"`
	Timestamp time.Time `json:"timestamp" time_format:"2006-01-02 15:04:05"`
}

type BidUseCase struct {
	BidRepository bid_entity.BidEntityRepository

	timer               *time.Timer
	maxBatchSize        int
	batchInsertInterval time.Duration
	bidChannel          chan bid_entity.Bid
}

func NewBidUseCase(bidRepository bid_entity.BidEntityRepository) BidUseCaseInterface {
	maxSizeInterval := getMaxBatchSizeInterval()
	maxBatchSize := getMaxBatchSize()

	bidUseCase := &BidUseCase{
		BidRepository:       bidRepository,
		maxBatchSize:        maxBatchSize,
		batchInsertInterval: maxSizeInterval,
		timer:               time.NewTimer(maxSizeInterval),
		bidChannel:          make(chan bid_entity.Bid, maxBatchSize),
	}

	bidUseCase.triggerCreateRoutine(context.Background())

	return bidUseCase
}

var bidBatch []bid_entity.Bid

type BidUseCaseInterface interface {
	CreateBid(
		ctx context.Context,
		bidInputDTO BidInputDTO) (string, *internal_error.InternalError)

	FindWinningBidByAuctionId(
		ctx context.Context, auctionId string) (*BidOutputDTO, *internal_error.InternalError)

	FindBidByAuctionId(
		ctx context.Context, auctionId string) ([]BidOutputDTO, *internal_error.InternalError)
}

func (bu *BidUseCase) triggerCreateRoutine(ctx context.Context) {
	go func() {
		defer close(bu.bidChannel)

		for {
			select {
			case bidEntity, ok := <-bu.bidChannel:
				if !ok {
					if len(bidBatch) > 0 {
						if err := bu.BidRepository.CreateBid(ctx, bidBatch); err != nil {
							logger.Error("error trying to process bid batch list", err)
						}
					}
					fmt.Println("[channel closed] bid batch created successfully")
					return
				}

				bidBatch = append(bidBatch, bidEntity)

				if len(bidBatch) >= bu.maxBatchSize {
					if err := bu.BidRepository.CreateBid(ctx, bidBatch); err != nil {
						logger.Error("error trying to process bid batch list", err)
					}

					bidBatch = nil
					bu.timer.Reset(bu.batchInsertInterval)
					fmt.Println("[max batch size] bid batch created successfully")
				}
			case <-bu.timer.C:
				if err := bu.BidRepository.CreateBid(ctx, bidBatch); err != nil {
					logger.Error("error trying to process bid batch list", err)
				}
				bidBatch = nil
				bu.timer.Reset(bu.batchInsertInterval)
				fmt.Println("[timer] bid batch created successfully")
			}
		}
	}()
}

func (bu *BidUseCase) CreateBid(
	ctx context.Context,
	bidInputDTO BidInputDTO) (string, *internal_error.InternalError) {

	bidEntity, err := bid_entity.CreateBid(bidInputDTO.UserId, bidInputDTO.AuctionId, bidInputDTO.Amount)
	if err != nil {
		return "", err
	}

	bu.bidChannel <- *bidEntity

	return bidEntity.Id, nil
}

func getMaxBatchSizeInterval() time.Duration {
	batchInsertInterval := os.Getenv("BATCH_INSERT_INTERVAL")
	duration, err := time.ParseDuration(batchInsertInterval)
	if err != nil {
		return 3 * time.Minute
	}

	return duration
}

func getMaxBatchSize() int {
	value, err := strconv.Atoi(os.Getenv("MAX_BATCH_SIZE"))
	if err != nil {
		return 5
	}

	return value
}
