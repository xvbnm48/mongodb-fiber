package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Catchphrase struct {
	id primitive.ObjectID `bson:"_id,omitempty" bson:"_id,omitempty"`
	//MovieName    string             `bson:"movie_name", json:"movieName, omitempty"`
	MovieName    string `json:"movieName, omitempty" bson:"movie_name,omitempty"`
	Catchphrase  string `json:"catchphrase, omitempty" bson:"catchphrase, omitempty"`
	MovieContext string `json:"MovieContext, omitempty" bson:"MovieContext, omitempty"`
}
