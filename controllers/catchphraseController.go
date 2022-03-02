package controllers

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/xvbnm48/mongodb-fiber/config"
	"github.com/xvbnm48/mongodb-fiber/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
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

func GetCatchphrase(c *fiber.Ctx) error {
	catchphraseCollection := config.MI.DB.Collection("catchphrases")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	var catchphrase models.Catchphrase
	objId, err := primitive.ObjectIDFromHex(c.Params("id"))
	findResult := catchphraseCollection.FindOne(ctx, bson.M{"_id": objId})
	if err = findResult.Err(); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  false,
			"message": "Catchphrase not found",
			"Error":   err,
		})
	}
	err = findResult.Decode(&catchphrase)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  false,
			"message": "Catchphrase not found",
			"Error":   err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    catchphrase,
		"success": true,
	})
}

func AddCatchphrase(c *fiber.Ctx) error {
	catchphraseCollection := config.MI.DB.Collection("catchphrases")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	catchphrase := new(models.Catchphrase)
	if err := c.BodyParser(catchphrase); err != nil {
		log.Println(err)
		return c.Status(400).JSON(fiber.Map{
			"status":  false,
			"message": "Invalid request body",
			"Error":   err,
		})
	}
	result, err := catchphraseCollection.InsertOne(ctx, catchphrase)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  false,
			"message": "Catchphrase failed to insert",
			"Error":   err,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data":    result,
		"success": true,
		"message": "catchphrase inserted successfully",
	})
}

func UpdateCatchphrase(c *fiber.Ctx) error {
	catchphraseCollection := config.MI.DB.Collection("catchphrases")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	catchphrase := new(models.Catchphrase)

	if err := c.BodyParser(catchphrase); err != nil {
		log.Println(err)
		return c.Status(400).JSON(fiber.Map{
			"status":  false,
			"message": "Invalid request body",
			"Error":   err,
		})
	}

	objId, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  false,
			"message": "Catchphrase not found",
			"Error":   err,
		})
	}

	update := bson.M{
		"$set": catchphrase,
	}
	_, err = catchphraseCollection.UpdateOne(ctx, bson.M{"_id": objId}, update)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  false,
			"message": "Catchphrase failed to update",
			"Error":   err,
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data":    catchphrase,
		"success": true,
		"message": "catchphrase updated successfully",
	})
}

func DeleteCatchphrase(c *fiber.Ctx) error {
	catchphraseCollection := config.MI.DB.Collection("catchphrases")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	objId, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  false,
			"message": "Catchphrase not found",
			"Error":   err,
		})
	}

	_, err = catchphraseCollection.DeleteOne(ctx, bson.M{"_id": objId})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  false,
			"message": "Catchphrase failed to delete",
			"Error":   err,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "catchphrase deleted successfully",
	})
}
