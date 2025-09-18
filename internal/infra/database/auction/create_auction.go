package auction

import (
	"context"
	"os"
	"time"

	"goexpert-leilao/configuration/logger"
	"goexpert-leilao/internal/entity/auction_entity"
	"goexpert-leilao/internal/internal_error"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type AuctionEntityMongo struct {
	Id          string                          `bson:"_id"`
	ProductName string                          `bson:"product_name"`
	Category    string                          `bson:"category"`
	Description string                          `bson:"description"`
	Condition   auction_entity.ProductCondition `bson:"condition"`
	Status      auction_entity.AuctionStatus    `bson:"status"`
	Timestamp   int64                           `bson:"timestamp"`
}
type AuctionRepository struct {
	Collection *mongo.Collection
}

func NewAuctionRepository(database *mongo.Database) *AuctionRepository {
	return &AuctionRepository{
		Collection: database.Collection("auctions"),
	}
}

func (ar *AuctionRepository) closeAuction(
	ctx context.Context,
	auctionId string) error {
	update := bson.M{"$set": bson.M{"status": auction_entity.Completed}}
	filter := bson.M{"_id": auctionId}

	_, err := ar.Collection.UpdateOne(ctx, filter, update)

	if err != nil {
		logger.Error("Error trying to update auction", err, zap.String("auctionId", auctionId))
		return err
	}
	logger.Info("Auction closed successfully", zap.String("auctionId", auctionId))
	return nil
}

func (ar *AuctionRepository) CreateAuction(
	ctx context.Context,
	auctionEntity *auction_entity.Auction) *internal_error.InternalError {
	auctionEntityMongo := &AuctionEntityMongo{
		Id:          auctionEntity.Id,
		ProductName: auctionEntity.ProductName,
		Category:    auctionEntity.Category,
		Description: auctionEntity.Description,
		Condition:   auctionEntity.Condition,
		Status:      auctionEntity.Status,
		Timestamp:   auctionEntity.Timestamp.Unix(),
	}
	_, err := ar.Collection.InsertOne(ctx, auctionEntityMongo)
	if err != nil {
		logger.Error("Error trying to insert auction", err)
		return internal_error.NewInternalServerError("Error trying to insert auction")
	}

	go func() {

		select {
		case <-time.After(getAuctionInterval()):
			ar.closeAuction(ctx, auctionEntityMongo.Id)
		case <-ctx.Done():
			logger.Info("Context cancelled while trying to close auction", zap.String("auctionId", auctionEntityMongo.Id))
		}

	}()

	return nil
}

func getAuctionInterval() time.Duration {
	auctionInterval := os.Getenv("AUCTION_INTERVAL")
	duration, err := time.ParseDuration(auctionInterval)
	if err != nil {
		return time.Minute * 4
	}
	return duration
}
