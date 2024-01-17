package main

import (
	"os"
	"fmt"
	"log"
	"database/sql"
	//	"net/http"
	"strconv"
	_"github.com/lib/pq"
//	"github.com/gin-gonic/gin"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"

	"github.com/joho/godotenv"
//	"gorm.io/gorm"
//	"gorm.io/driver/postgres"
)

type todo struct {
	Item string
}

func indexHandler(c *fiber.Ctx, db *sql.DB) error {
	var res string
	var todos []string
	rows, err := db.Query("SELECT * FROM todos")
	fmt.Printf("STUFF: %v\n", rows)
	defer rows.Close()

	if err != nil {
		log.Fatalln(err)
		c.JSON("An error has occured!")
	}

	for rows.Next() {
		rows.Scan(&res)
		fmt.Printf("STUFF IN FOR LOOP: %v\n", rows)
		todos = append(todos, res)
	}
	
	return c.Render("index", fiber.Map{"Todos": todos,})
}
		
func postHandler(c *fiber.Ctx, db *sql.DB) error {
	newTodo := todo{}
	if err := c.BodyParser(&newTodo); err != nil {
		log.Printf("An error has occured trying to post: %v", err)
		return c.SendString(err.Error())
	}

	fmt.Printf("%v\n", newTodo)
	if newTodo.Item != "" {
		_, err := db.Exec("INSERT into todos VALUES ($1)", newTodo.Item)
		if err != nil {
			log.Fatalf("An erro occured while executing POST: %v", err)
		}
	}
	return c.Redirect("/")
}

func putHandler(c *fiber.Ctx, db *sql.DB) error {
	olditem := c.Query("olditem")
	newitem := c.Query("newitem")
	db.Exec("UPDATE todos SET item=$1 WHERE item=$2", newitem, olditem)
	return c.Redirect("/")
}

func deleteHandler(c *fiber.Ctx, db *sql.DB) error {
	todoToDelete := c.Query("item")
	db.Exec("DELETE from todos WHERE item=$1", todoToDelete)
	return c.SendString("Deleted")
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var postUser string = os.Getenv("PG_USER")
	var postPassword string = os.Getenv("PG_PASSSWORD")
	var postHost string = os.Getenv("PG_HOST")
	postPort, _ := strconv.Atoi(os.Getenv("PG_PORT"))
	var dbName string = os.Getenv("DB_NAME")

	//var connStr string = fmt.Sprintf("postgresql://%s:%s@%s/%s?sslmode=disable",
	var connStr string = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		postHost,
		postPort,
		postUser,
		postPassword,
		dbName )
	
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(fmt.Sprintf("Error connecting to postgres: %s", err))
	}

	engine := html.New("./../client/views", ".html")
	app := fiber.New(fiber.Config{Views: engine,})
	
	app.Get("/", func(c *fiber.Ctx) error {
		return indexHandler(c, db)
	})

	app.Post("/",func(c *fiber.Ctx) error {
		return postHandler(c ,db)
	})

	app.Put("/update",func(c *fiber.Ctx) error {
		return putHandler(c, db)
	})

	app.Delete("/delete",func(c *fiber.Ctx) error {
		return deleteHandler(c,db)
	})
	
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	app.Static("/", "./../client/public")
	log.Fatalln(app.Listen(fmt.Sprintf(":%v", port)))
}

