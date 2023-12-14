package main

import (
	"database/sql"
	"fmt"
	"log"

	D "kwizz/functions"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

var db *sql.DB

type Categorie struct {
	Cat_Id          int
	Cat_ShortName   string
	Cat_Name        string
	Cat_Description string
	Cat_Image       string
}

type Quizz struct {
	Quizz_Id          int
	Quizz_Name        string
	Quizz_Description string
	Created_at        string
	Quizz_CatID       int
}

type Question struct {
	Q_Question_ID int
	Quizz_ID      int
	Question      string
	// Response        string
	Order_questions *int
}

type Response struct {
	Response_ID   int
	R_Question_ID int
	// Session_ID  int
	Answer    string
	isCorrect *bool
}

type HappyCouples struct {
	AnswerOne   string
	AnswerTwo   string
	AnswerThree *string
	Question    string
}

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
		var Categories []Categorie
		var pseudonym string
		err := db.QueryRow("select pseudonym from users").Scan(&pseudonym)
		check(err)

		rows, err := db.Query("select * from categories")
		check(err)
		defer rows.Close()

		for rows.Next() {
			var cat Categorie
			err := rows.Scan(&cat.Cat_Id, &cat.Cat_Name, &cat.Cat_ShortName, &cat.Cat_Description, &cat.Cat_Image)
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
			err := rows.Scan(&quizz.Quizz_Id, &quizz.Quizz_Name, &quizz.Quizz_Description, &quizz.Created_at, &quizz.Quizz_CatID)
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
		myQuizz := c.Params("id")

		id := 0
		err := db.QueryRow(`Insert into quizz_sessions(quizz_id) values ($1) returning session_id`, myQuizz).Scan(&id)
		check(err)

		url := fmt.Sprintf("/session/%s/%d", myQuizz, id)

		return c.Redirect(url)

	})

	// ACTIVE SESSION PAGE
	app.Get("/session/:quizz_id/:session_id", func(c *fiber.Ctx) error {
		myQuizz := c.Params("quizz_id")
		// mySession := c.Params("session_id")

		var Questions []Question

		rows, err := db.Query("select * from questions where quizz_id = $1", myQuizz)
		check(err)
		// defer rows.Close()

		for rows.Next() {
			var question Question
			err := rows.Scan(&question.Q_Question_ID, &question.Quizz_ID, &question.Question, &question.Order_questions)
			check(err)
			Questions = append(Questions, question)
		}

		var Responses []Response

		otherRows, err := db.Query("select * from responses")
		check(err)
		defer otherRows.Close()

		for otherRows.Next() {
			var response Response
			err := otherRows.Scan(&response.Response_ID, &response.R_Question_ID, &response.Answer, &response.isCorrect)
			check(err)
			Responses = append(Responses, response)
		}

		var Happy_couples []HappyCouples

		for _, q := range Questions {
			var Happy_couple HappyCouples
			for _, r := range Responses {
				if q.Q_Question_ID == r.R_Question_ID {
					Happy_couple = append(Happy_couple, r.Answer)
				}
			}
			Happy_couples = append(Happy_couples, q.Question)
		}

		fmt.Print(Happy_couples)

		// otherRows.Close()

		return c.Render("pages/session", fiber.Map{
			// "Answers":   Split_answers,
			"Happy_couples": Happy_couples,
			"Questions":     Questions,
			"Responses":     Responses,
		})
	})

	// Renders CSS
	app.Static("/assets/", "./assets")
	// Renders pictures
	app.Static("/public/", "./public")

	// Open port to listen to
	log.Fatal(app.Listen(":3000"))

}
