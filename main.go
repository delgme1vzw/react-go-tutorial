package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// custom data structure to give different fields with different data types
// mongo db has its own int type for ids, so using that primitive.objectId here
// omitempty will ignore it if the value is 0
type Todo struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"` //bson here is mongo db column (binary json)
	Completed bool               `json:"completed"`                          //false
	Body      string             `json:"body"`                               //empty string
}

var collection *mongo.Collection

func main() {
	fmt.Println("Hello, World!")

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	MONGODB_URI := os.Getenv("MONGODB_URI")
	clientOptions := options.Client().ApplyURI(MONGODB_URI)
	//Usually when you create a connection you want some cancellations/timeouts ... here we don't so this is saying I am starting task but no special requirements that bind to it
	client, err := mongo.Connect(context.Background(), clientOptions)

	if err != nil {
		log.Fatal("Error connecting to DB: ", err)
	}

	//When main function done (main) connecting to db, we'd like to disconnect from db
	defer client.Disconnect(context.Background())

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal("Error connecting to DB cluster: ", err)
	}

	fmt.Println("Connected to MongoDB Atlas")

	//Here the database name is golang_db we will grab from dotenv, and the collection(table) is todo
	collection = client.Database("golang_db").Collection("todo")

	app := fiber.New()

	app.Get("/api/todos", getTodos)
	app.Post("/api/todos", createTodo)
	app.Patch("/api/todos/:id", updateTodo)
	app.Delete("/api/todos/:id", deleteTodo)

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	log.Fatal(app.Listen("0.0.0.0:" + port))

}

func getTodos(c *fiber.Ctx) error {
	var todos []Todo
	//trying to pass filter (no filter fetching all todos in collection)
	//cursor is a pointer to rSet so you can use it to iterate over documents (rows) returned from query
	cursor, err := collection.Find(context.Background(), bson.M{})

	if err != nil {
		return err
	}

	//postpone execution of function until surrounding function completes
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var todo Todo
		if err := cursor.Decode(&todo); err != nil {
			return err
		}
		todos = append(todos, todo)
	}

	return c.JSON(todos)
}

func createTodo(c *fiber.Ctx) error {
	//Will use c.bodyparser so this needs to be a pointer
	todo := new(Todo)
	//Parses response body into the struct you pass (todo)
	if err := c.BodyParser(todo); err != nil {
		return err
	}

	if todo.Body == "" {
		return c.Status(400).JSON(fiber.Map{"error": "To do body cannot be empty"})
	}

	//Context, document (row) to be inserted
	insertResult, err := collection.InsertOne(context.Background(), todo)
	if err != nil {
		return err
	}

	//At this point, a newly declared Todo has no initialized values so id=0...let's set id
	//This will check what ID was returned from the insert statement since it is sequential key
	todo.ID = insertResult.InsertedID.(primitive.ObjectID)

	return c.Status(201).JSON(todo)
}

func updateTodo(c *fiber.Ctx) error {
	id := c.Params("id")
	//Convert from string to primitive id for mongo db
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid todo ID"})
	}

	//First, what is the "where" clause for the update?
	filter := bson.M{"_id": objectId}
	//Now what do we want the "update" clause to be? completed=true
	update := bson.M{"$set": bson.M{"completed": true}}
	//Put it together, update one row to true where filter=whatever
	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err

	}
	return c.Status(200).JSON(fiber.Map{"success": true})

}

func deleteTodo(c *fiber.Ctx) error {
	id := c.Params("id")
	//convert to primitive
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid todo ID"})
	}

	//Where clause
	filter := bson.M{"_id": objectId}
	_, err = collection.DeleteOne(context.Background(), filter)

	if err != nil {
		return err
	}

	return c.Status(200).JSON(fiber.Map{"success": true})

}
