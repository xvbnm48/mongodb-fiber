package controllers

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/xvbnm48/mongodb-fiber/config"
	"github.com/xvbnm48/mongodb-fiber/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"math"
	"strconv"
	"time"
)

func GetAllCatchphrases(c *fiber.Ctx) error {
	catchphraseCollection := config.MI.DB.Collection("catchphrases")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	var catchphrases []models.Catchphrase

	filter := bson.M{}
	findOptions := options.Find()

	if s := c.Query("s"); s != "" {
		filter = bson.M{
			"$or": []bson.M{
				{
					"movieName": primitive.Regex{
						Pattern: s,
						Options: "i",
					},
				},
				{
					"catchphrase": bson.M{
						"$regex": primitive.Regex{
							Pattern: s,
							Options: "i",
						},
					},
				},
			},
		}
	}
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limitVal, _ := strconv.Atoi(c.Query("limit", "10"))
	var limit int64 = int64(limitVal)

	total, _ := catchphraseCollection.CountDocuments(ctx, filter)

	findOptions.SetSkip(int64(page) - 1*limit)
	findOptions.SetLimit(limit)
	cursor, err := catchphraseCollection.Find(ctx, filter, findOptions)
	defer cursor.Close(ctx)

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  false,
			"message": "Catchphrases not found",
			"Error":   err,
		})
	}
	for cursor.Next(ctx) {
		var catchphrase models.Catchphrase
		cursor.Decode(&catchphrase)
		catchphrases = append(catchphrases, catchphrase)
	}
	last := math.Ceil(float64(total / limit))
	if last < 1 && total > 0 {
		last = 1
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":      catchphrases,
		"total":     total,
		"page":      page,
		"last_page": last,
		"limit":     limit,
	})
}
