package main

import (
	"database/sql"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"net/http"
	"os"
)

var DB *sql.DB

func SetupRoutes(db *sql.DB, fiber *fiber.App) {
	DB = db

	todos := fiber.Group("todos")
	todos.Get("/", getAllTodos)
	todos.Post("/", createTodo)
	todos.Delete("/:todo", deleteTodo)
	todos.Post("/search", searchTodo)

	SetupWeb(fiber)
}

func SetupWeb(f *fiber.App) {

	dir := "./web/dist"

	//	check if dir exists
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		f.Get("/", func(c *fiber.Ctx) error {
			return c.SendString("Web directory not found")
		})

	} else {
		f.Use(filesystem.New(filesystem.Config{
			Root:   http.Dir(dir),
			Browse: false,
			Index:  "index.html",
			MaxAge: 3600,
		}))
	}

}

func searchTodo(c *fiber.Ctx) error {
	var q string
	err := c.BodyParser(&q)
	if err != nil {
		return err
	}
	rows, err := DB.Query("SELECT id, content from todos where content like ?;", "%"+q+"%")
	if err != nil {
		return err
	}
	defer rows.Close()

	var todos []Todo
	for rows.Next() {
		var todo Todo
		err := rows.Scan(&todo.Id, &todo.Content)
		if err != nil {
			return err
		}
		todos = append(todos, todo)
	}
	return c.JSON(todos)
}

func deleteTodo(c *fiber.Ctx) error {
	todo := c.Params("todo")
	_, err := DB.Exec("DELETE FROM todos WHERE id = ?;", todo)
	if err != nil {
		return err
	}
	return c.SendStatus(200)
}
func createTodo(c *fiber.Ctx) error {
	todo := new(string)
	if err := c.BodyParser(todo); err != nil {
		return err
	}
	data, err := DB.Query("INSERT INTO todos (content) VALUES (?) returning id;", todo)
	if err != nil {
		return err
	}
	var id int
	for data.Next() {
		err := data.Scan(&id)
		if err != nil {
			return err
		}
	}
	return c.JSON(Todo{
		Id:      id,
		Content: *todo,
	})
}

func getAllTodos(c *fiber.Ctx) error {
	q, err := DB.Query("SELECT id, content FROM todos;")
	if err != nil {
		return err
	}
	defer q.Close()
	var todos []Todo
	for q.Next() {
		var todo Todo
		err := q.Scan(&todo.Id, &todo.Content)
		if err != nil {
			return err
		}
		todos = append(todos, todo)
	}
	if todos == nil {
		return c.JSON([]string{})
	}
	return c.JSON(todos)

}
