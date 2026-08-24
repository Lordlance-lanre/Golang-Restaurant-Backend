package controller

import (
	"context"
	// "fmt"
	"time"
	"net/http"
	"github.com/Lordlance-lanre/Golang-Restaurant-Backend/database"
	"github.com/Lordlance-lanre/Golang-Restaurant-Backend/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	// "go.mongodb.org/mongo-driver/v2/bson/primitive"
)

type OrderItemPack struct{
	Table_Id *string
	Order_items []models.OrderItems
}

var orderItemsCollection*mongo.Collection = database.OpenCollection(database.Client, "orderItems")

func GetOrderItems() gin.HandlerFunc{
	return func(c *gin.Context){
		// c.Writer.Write([]byte("Get Order Items"))
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		result, err := orderItemsCollection.Find(context.TODO(), bson.M{})
		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		var allOrderItems []bson.M
		if err := result.All(ctx, &allOrderItems); err != nil{
			msg := "Unable to get order items"
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		c.JSON(http.StatusOK, allOrderItems)
	}
}

func GetOrderItemsByOrder() gin.HandlerFunc{
	return func(c *gin.Context){
		// c.Writer.Write([]byte("Get Order Item By Order"))
		OrderId := c.Param("order_id")

		allOrderItems, err := ItemsByOrder(OrderId)
		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, allOrderItems)
	}
}

func ItemsByOrder(id string) (OrderItems []primitive.M, err error){
	// fmt.Println("Items By Order")

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	matchStage := bson.D{{Key: "$match", Value: bson.D{{Key: "order_id", Value: id}}}}

	lookupStage := bson.D{{Key: "$lookup", Value: bson.D{{Key: "from", Value: "food"}, {Key: "localField", Value: "food_id"}, {Key: "foreignField", Value: "food_id"}, {Key: "as", Value: "food"}}}}

	unwindStage := bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$food"}, {Key: "preserveNullAndEmptyArrays", Value: true}}}}

	lookupOrderStage := bson.D{{Key: "$lookup", Value: bson.D{{Key: "from", Value: "order"}, {Key: "localField", Value: "order_id"}, {Key: "foreignField", Value: "order_id"}, {Key: "as", Value: "order"}}}}
	unwindOrderStage := bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$order"}, {Key: "preserveNullAndEmptyArrays", Value: true}}}}

	lookupTableStage := bson.D{{Key: "$lookup", Value: bson.D{{Key: "from", Value: "table"}, {Key: "localField", Value: "order.table_id"}, {Key: "foreignField", Value: "table_id"}, {Key: "as", Value: "table"}}}}
	unwindTableStage := bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$table"}, {Key: "preserveNullAndEmptyArrays", Value: true}}}}


	projectStage := bson.D{
		{Key: "$project", Value: bson.D{
			{Key: "id", Value: 0},
			{Key: "amount", Value: "$food.price"},
			{Key: "total_count", Value: 1},
			{Key: "food_name", Value: "$food.name"},
			{Key: "food_image", Value: "$food.food_image"},
			{Key: "table_number", Value: "$table.table_number"},
			{Key: "table_id", Value: "$table.table_id"},
			{Key: "order_id", Value: "$order.order_id"},
			{Key: "price", Value: "$food.price"},
			{Key: "quantity", Value: 1},
		}},
	}

	groupStage := bson.D{{Key:"$group", Value: bson.D{{Key: "_id", Value: bson.D{{Key:"order_id", Value: "$order_id"}, {Key:"table_id", Value: "$table_id"}, {Key: "table_number",Value:"$table_number"}}}, {Key: "payment_due", Value: bson.D{{Key: "$sum", Value: "$amount"}}}, {Key:"total_count", Value: bson.D{{Key: "$sum", Value:1}}}, {Key:"order_items", Value: bson.D{{Key:"$push", Value: bson.D{{Key: "food_id", Value:"$food_id"}, {Key: "quantity",Value:"$quantity"}, {Key: "unit_price", Value: "$unit_price"}}}}}}}}

	projectStage2 := bson.D{ 
		{Key: "$project", Value: bson.D{
			{Key: "id", Value: 0},
			{Key: "payment_due", Value: 1},
			{Key: "total_count", Value: 1},
			{Key: "table_number", Value: "$_id.table_number"},
			{Key: "order_items", Value: 1},
		}},
	}

	result, err := orderItemsCollection.Aggregate(ctx, mongo.Pipeline{
		matchStage,
		lookupStage,
		unwindStage,
		lookupOrderStage,
		unwindOrderStage,
		lookupTableStage,
		unwindTableStage,
		projectStage,
		groupStage,
		projectStage2,
	})
	if err != nil {
		return OrderItems, err
	}
	if err = result.All(ctx, &OrderItems); err != nil {
		return OrderItems, err
	}
	return OrderItems, err
}

func GetOrderItemById() gin.HandlerFunc{
	return func(c *gin.Context){
		// c.Writer.Write([]byte("Get Order Item By Id"))
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		 orderItemId := c.Param("order_item_id")
		 var orderItem models.OrderItems

		err:= orderItemsCollection.FindOne(ctx, bson.M{
			"order_item_id": orderItemId,
		 }).Decode(&orderItem)

		if err != nil{
			msg := "Unable to get order item by id"
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		c.JSON(http.StatusOK, orderItem)		
	}
}

func CreateOrderItem() gin.HandlerFunc{
	return func(c *gin.Context){
		// c.Writer.Write([]byte("Create Order Item"))
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var orderItemPack OrderItemPack
		var order models.Order

		if err := c.BindJSON(&orderItemPack); err != nil{
			msg := "Bad request made"
			c.JSON(http.StatusBadRequest, gin.H{"error": msg})
			return
		}
		order.Order_date, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		orderItemsToBeInserted := []interface{}{}
		order.Table_id = orderItemPack.Table_Id
		order_id := OrderItemOrderCreator(order)

		for _, orderItem:= range orderItemPack.Order_items{
			orderItem.Order_id = &order_id
			validationErr := validate.Struct(orderItem)

			if validationErr != nil{
				c.JSON(http.StatusInternalServerError, gin.H{"error": validationErr.Error()})
				return
			}

			orderItem.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
			orderItem.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
			orderItem.ID = primitive.NewObjectID()
			orderItem_id := orderItem.ID.Hex()
			orderItem.Order_item_id = &orderItem_id
			var num = toFixed(*orderItem.Unit_price, 2)
			orderItem.Unit_price = &num

			orderItemsToBeInserted = append(orderItemsToBeInserted, orderItem)
		}
		insertOrderItems, err := orderItemsCollection.InsertMany(ctx, orderItemsToBeInserted)

		if err != nil{
			msg := "Unable to add order items"
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		c.JSON(http.StatusOK, insertOrderItems)
	}
}

func UpdateOrderItem() gin.HandlerFunc{
	return func(c *gin.Context){
		// c.Writer.Write([]byte("Update Order Item"))

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var orderItem models.OrderItems
		orderItemId := c.Param("order_item_id")

		filter := bson.M{"order_item_id": orderItemId}

		var updateObj bson.D

		if orderItem.Unit_price != nil {
			updateObj = append(updateObj, bson.E{Key: "unit_price", Value: orderItem.Unit_price})
		}
		if orderItem.Quantity != nil{
			updateObj = append(updateObj, bson.E{Key: "quantity", Value: orderItem.Quantity})
		}
		if orderItem.Food_id != nil {
			updateObj = append(updateObj, bson.E{Key: "food_id", Value: orderItem.Food_id})
		}

		orderItem.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		updateObj = append(updateObj, bson.E{Key: "updated_at", Value: orderItem.Updated_at})

		upsert := true
		opt := options.UpdateOne().SetUpsert(upsert)

		result, err := orderItemsCollection.UpdateOne(ctx, filter, bson.D{
			bson.E{Key: "$set", Value: updateObj},
		},opt)

		if err != nil{
			msg := "Error occured while updating order item"
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		c.JSON(http.StatusOK, result)

	}
}

func DeleteOrderItem() gin.HandlerFunc{
	return func(c *gin.Context){
		// c.Writer.Write([]byte("Delete Order Item"))
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		
		orderItemId := c.Param("order_item_id")
		if err := orderItemsCollection.FindOne(ctx, bson.M{"order_item_id": orderItemId}).Err(); err != nil{
			msg := "Unable to get order item by id"
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		result, err := orderItemsCollection.DeleteOne(ctx, bson.M{"order_item_id": orderItemId})
		if err != nil{
			msg := "Error occured while deleting order item"
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		c.JSON(http.StatusOK, result)
	}
}