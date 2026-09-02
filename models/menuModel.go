package models

import(
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Menu struct{
	ID				primitive.ObjectID  `bson:"_id"`
	Name			*string 			`json:"name" validate:"required" bson:"name"`
	Category		*string 			`json:"category" validate:"required" bson:"category"`
	Start_Date		time.Time 			`json:"start_date" bson:"start_date"`
	End_Date		time.Time 			`json:"end_date" bson:"end_date"`
	Created_at		time.Time 			`json:"created_at" bson:"created_at"`
	Updated_at		time.Time 			`json:"updated_at" bson:"updated_at"`
	Menu_id			*string 			`json:"menu_id" bson:"menu_id"`
}