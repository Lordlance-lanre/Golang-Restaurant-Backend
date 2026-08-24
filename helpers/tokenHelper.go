package helpers

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/Lordlance-lanre/Golang-Restaurant-Backend/database"
	"github.com/dgrijalva/jwt-go"
	// "go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	bson "go.mongodb.org/mongo-driver/v2/bson"
)

type SignedDetails struct {
	Email 				string
	First_name 			string
	Last_name 			string
	Uid 				string
	jwt.StandardClaims
}

var SECRET_KEY string = os.Getenv("SECRET_KEY")

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")

func GenerateAllTokens(email string, userId string, firstName string, lastName string) (signedToken string, signedRefreshToken string, err error){
	claim := &SignedDetails{
		First_name: firstName,
		Last_name: lastName,
		Uid: userId,
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Minute * 60).Unix(),
			IssuedAt: time.Now().Unix(),
		},
	}
	refreshClaims := &SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * 24 * 7).Unix(),
			IssuedAt: time.Now().Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claim).SignedString([]byte(SECRET_KEY))
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SECRET_KEY))

	if err != nil{
		log.Panic(err)
		return
	}

	return token, refreshToken, err
}


func UpdateAllTokens(signedToken string, signedRefreshToken string, userId string){
	ctx, cancel := context.WithTimeout(context.Background(), 100 * time.Second)
	defer cancel()

	var updateObj []bson.E

	updateObj = append(updateObj, bson.E{Key: "token", Value: signedToken})
	updateObj = append(updateObj, bson.E{Key: "refresh_token", Value: signedRefreshToken})
	updateObj = append(updateObj, bson.E{Key: "updated_at", Value: time.Now()})

	upsert := true
	filter := bson.M{"user_id": userId}

	opt := options.UpdateOne().SetUpsert(upsert)

	_, err := userCollection.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: updateObj}}, opt)

	if err != nil {
		log.Panic(err)
		return
	}

}

func ValidateAllTokens(signedToken string) (claims *SignedDetails, msg string) {
	//Validate token authenticity
	token, err := jwt.ParseWithClaims(
		signedToken,
		&SignedDetails{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(SECRET_KEY), nil
		},
	)
	if err != nil {
		msg = err.Error()
		return nil, msg
	}

	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		msg = "the token is invalid"
		return nil, msg
	}

	// Vlaidate token if it has expired
	if claims.ExpiresAt < time.Now().Local().Unix() {
		msg = "token is expired"
		return nil, msg
	}

	return claims, ""
}
