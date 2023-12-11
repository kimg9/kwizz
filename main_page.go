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

type Quizz struct {
	Id               int
	QuizzName        string
	QuizzDescription string
	Created_at       string
	Quizz_CatID      int
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
		// log.Fatalf("query error: %v\n", err)
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
		var Categories []Categorie
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

	// CATEGORIES PAGE
	app.Get("/categorie/:id", func(c *fiber.Ctx) error {
		myCat := c.Params("id")

		var queryValue string
		err := db.QueryRow("select count(*) != 0 from categories where cat_id = $1", myCat).Scan(&queryValue)
		check(err)

		if queryValue != "true" {
			return c.SendStatus(404)
		}

		var Quizzes []Quizz

		rows, err := db.Query("select * from quizzes where cat_id = $1", myCat)
		check(err)
		defer rows.Close()

		for rows.Next() {
			var quizz Quizz
			err := rows.Scan(&quizz.Id, &quizz.QuizzName, &quizz.QuizzDescription, &quizz.Created_at, &quizz.Quizz_CatID)
			check(err)
			Quizzes = append(Quizzes, quizz)
		}

		err = rows.Err()
		check(err)

		var CurrentCat string
		errerr := db.QueryRow("select cat_name from categories where cat_id = $1", myCat).Scan(&CurrentCat)
		check(errerr)

		return c.Render("pages/categorie", fiber.Map{
			"CurrentCat": CurrentCat,
			"Quizzes":    Quizzes,
		})
	})

	// QUIZZES PAGE
	app.Get("/quizz/:id", func(c *fiber.Ctx) error {

		return c.Render("pages/quizz", fiber.Map{})
	})

	// Renders CSS
	app.Static("/assets/", "./assets")
	// Renders pictures
	app.Static("/public/", "./public")

	// Open port to listen to
	log.Fatal(app.Listen(":3000"))

}
