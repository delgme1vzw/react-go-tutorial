package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

// custom data structure to give different fields with different data types
type Todo struct {
	ID        int    `json:"id"`        //0
	Completed bool   `json:"completed"` //false
	Body      string `json:"body"`      //empty string
}

func main() {
	fmt.Println("Hello, World!")

	app := fiber.New()

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error Loading .env file")
	}

	PORT := os.Getenv("PORT")

	//Create a slice of Todo structs to store our todos in memory.
	todos := []Todo{}

	//Context and error here where context is a pointer
	app.Get("/api/todos", func(c *fiber.Ctx) error {
		// id := c.Params("id") //get id from the query parameters, which is a string
		// for _, todo := range todos {
		// 	if fmt.Sprint(todo.ID) == id {
		// 		return c.Status(200).JSON(todo)
		// 	}
		// }
		// return c.Status(404).JSON(fiber.Map{"msg": "Todo not found"})
		return c.Status(200).JSON(todos)
	})

	//Create a todo
	app.Post("/api/todos", func(c *fiber.Ctx) error {
		//todo is the memory address of the Todo struct, so we need to use the & operator to get the memory address of the struct and assign it to the todo variable.
		todo := &Todo{}
		//Tries to parse the body of the request into the todo struct.
		if err := c.BodyParser(todo); err != nil {
			//return c.Status(400).JSON(fiber.Map{"error": "Cannot parse JSON"})
			return err
		}

		//todo here is a pointer, so we need to dereference it to access the Body field
		//it does not need a * because we are using the dot notation to access the field, which automatically dereferences the pointer
		if todo.Body == "" {
			return c.Status(400).JSON(fiber.Map{"error": "To do body is required"})
		}

		todo.ID = len(todos) + 1
		//Add the todo to the slice of todos. We need to dereference the pointer to get the value of the todo struct and append it to the slice.
		//We are getting the value of the todo struct by dereferencing the pointer using the * operator, which gives us the value of the struct that the pointer is pointing to.
		todos = append(todos, *todo)

		//201 = resource created with json string of the todo struct
		return c.Status(201).JSON(todo)
	})

	//Update a todo, always set todo = true
	app.Patch("/api/todos/:id", func(c *fiber.Ctx) error {
		id := c.Params("id") //get id from the url parameters, which is a string

		//loop through slice of todos
		for i, todo := range todos {
			//if we found the todo with that id
			if fmt.Sprint(todo.ID) == id {
				todos[i].Completed = true
				return c.Status(200).JSON(todos[i])
			}
		}
		//else we return a 404 error with a json string of an error message
		return c.Status(404).JSON(fiber.Map{"error": "To do not found"})
	})

	//Delete todo
	app.Delete("/api/todos/:id", func(c *fiber.Ctx) error {
		id := c.Params("id") //get id from the url parameters, which is a string

		//loop through slice of todos
		for i, todo := range todos {
			//if we found the todo with that id
			if fmt.Sprint(todo.ID) == id {
				//remove the todo from the slice of todos by slicing out the todo at index i and concatenating the two slices together
				//append from the start of the array up to but not including the index we found id at
				//and append the values after index i by unpacking them with ...
				todos = append(todos[:i], todos[i+1:]...)
				return c.Status(200).JSON(fiber.Map{"success": true})
			}
		}
		//else we return a 404 error with a json string of an error message
		return c.Status(404).JSON(fiber.Map{"error": "To do not found"})
	})
	log.Fatal(app.Listen(":" + PORT))
}
