package main

import (
	"database/sql"
	"log"

	D "kwizz/functions"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

var db *sql.DB

type Categorie struct {
	Id          int
	ShortName   string
	Name        string
	Description string
	Image       string
}

var Categories []Categorie

type Quizz struct {
	Id               int
	QuizzName        string
	QuizzDescription string
	Quizz_CatID      int
}

var Quizzes []Categorie

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {

	// Database connexion
	db = D.Connect()
	defer db.Close()

	// Start app
	engine := html.New("./views", ".html")

	// Add unescape function to allow innerHTML
	// engine.AddFunc(
	// 	"unescape", func(s string) template.HTML {
	// 		return template.HTML(s)
	// 	},
	// )

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// HOME PAGE
	app.Get("/", func(c *fiber.Ctx) error {
		var pseudonym string
		err := db.QueryRow("select pseudonym from users").Scan(&pseudonym)
		check(err)

		rows, err := db.Query("select * from categories")
		check(err)
		defer rows.Close()

		for rows.Next() {
			var cat Categorie
			err := rows.Scan(&cat.Id, &cat.Name, &cat.ShortName, &cat.Description, &cat.Image)
			check(err)
			Categories = append(Categories, cat)
		}

		err = rows.Err()
		check(err)

		// Render index template
		return c.Render("main_page", fiber.Map{
			"Pseudonym":  pseudonym,
			"Categories": Categories,
		})
	})

	//CATEGORIES

	// Famille

	// for cat := range Categories {
	// 	if cat.Id =
	// }

	// app.Get("/famille", func(c *fiber.Ctx) error {

	// 	rows, err := db.Query("select * from quizzes")
	// 	check(err)
	// 	defer rows.Close()

	// 	for rows.Next() {
	// 		var cat Quizzes
	// 		err := rows.Scan(&cat.Id, &cat.Name, &cat.ShortName, &cat.Description, &cat.Image)
	// 		check(err)
	// 		Categories = append(Categories, cat)
	// 	}

	// 	return c.Render("categories/famille", fiber.Map{
	// 		"Title": "Hello world",
	// 	})
	// })

	// Culture Générale
	// app.Get("/culturegenerale", func(c *fiber.Ctx) error {
	// 	return c.Render("categories/culturegenerale", fiber.Map{
	// 		"Title": "Hello world",
	// 	})
	// })

	// // Mathématiques
	// app.Get("/mathematiques", func(c *fiber.Ctx) error {
	// 	return c.Render("categories/mathematiques", fiber.Map{
	// 		"Title": "Hello world",
	// 	})
	// })

	// // Français
	// app.Get("/francais", func(c *fiber.Ctx) error {
	// 	return c.Render("categories/francais", fiber.Map{
	// 		"Title": "Hello world",
	// 	})
	// })

	// Renders CSS
	app.Static("/assets/", "./assets")
	// Renders pictures
	app.Static("/public/", "./public")

	// Open port to listen to
	log.Fatal(app.Listen(":3000"))

}
