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
	Q_Question_ID   int
	Quizz_ID        int
	Question        string
	Order_questions *int
}

type Response struct {
	Response_ID   int
	R_Question_ID int
	Answer        string
	isCorrect     *bool
}

type HappyCouple struct {
	TheQuestion  Question
	TheResponses []Response
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func newSession(myQuizz string) int {
	id := 0
	err := db.QueryRow(`Insert into quizz_sessions(quizz_id, user_id) values ($1, '1') returning session_id`, myQuizz).Scan(&id)
	check(err)
	return id
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

		// Check for already done sessions
		var finished bool
		errerr := db.QueryRow(`select exists(select 1 from quizz_sessions where quizz_id = $1 and finished = true)`, myQuizz).Scan(&finished)
		check(errerr)

		if finished {
			lastSession := 0
			err := db.QueryRow(`select session_id from quizz_sessions where finished = true and quizz_id = $1 order by created_at desc limit 1`, myQuizz).Scan(&lastSession)
			check(err)

			url := fmt.Sprintf("/session/%s/%d", myQuizz, lastSession)

			return c.Redirect(url)
		}

		//no previous done sessions found
		id := newSession(myQuizz)

		url := fmt.Sprintf("/session/%s/%d", myQuizz, id)

		return c.Redirect(url)

	})

	// ACTIVE SESSION PAGE
	app.Get("/session/:quizz_id/:session_id", func(c *fiber.Ctx) error {
		myQuizz := c.Params("quizz_id")
		mySession := c.Params("session_id")

		// Get questions that will be rendered in HTML
		var Questions []Question

		rows, err := db.Query("select * from questions where quizz_id = $1", myQuizz)
		check(err)
		defer rows.Close()

		for rows.Next() {
			var question Question
			err := rows.Scan(&question.Q_Question_ID, &question.Quizz_ID, &question.Question, &question.Order_questions)
			check(err)
			Questions = append(Questions, question)
		}

		// Get answers : two different cases whether the quizz was already completed or not
		var finished bool
		errerr := db.QueryRow(`select finished from quizz_sessions where session_id = $1`, mySession).Scan(&finished)
		check(errerr)

		// CASE 1 : quizz was already taken
		if finished {
			var Score int
			errerr := db.QueryRow(`select score from quizz_sessions where session_id = $1`, mySession).Scan(&Score)
			check(errerr)

			// for _, q := range Questions {
			// 	var respID int
			// 	var questID int
			// 	var answer string
			// 	var isCorrect bool

			// 	otherRows, err := db.Query("select * from responses where question_id = $1", q.Q_Question_ID)
			// 	check(err)
			// 	defer otherRows.Close()

			// 	for otherRows.Next() {
			// 		var response Response
			// 		err := otherRows.Scan(&response.Response_ID, &response.R_Question_ID, &response.Answer, &response.isCorrect)
			// 		check(err)
			// 		Responses = append(Responses, response)
			// 	}

			// 	var happyCouple HappyCouple
			// 	happyCouple.TheQuestion = q
			// 	happyCouple.TheResponses = Responses

			// 	HappyCouples = append(HappyCouples, happyCouple)

			// }

			return c.Render("pages/finished_session", fiber.Map{
				"Score": Score,
			})

		}

		// CASE 2 : first time taking quizz
		var HappyCouples []HappyCouple

		for _, q := range Questions {
			var Responses []Response

			otherRows, err := db.Query("select * from responses where question_id = $1", q.Q_Question_ID)
			check(err)
			defer otherRows.Close()

			for otherRows.Next() {
				var response Response
				err := otherRows.Scan(&response.Response_ID, &response.R_Question_ID, &response.Answer, &response.isCorrect)
				check(err)
				Responses = append(Responses, response)
			}

			var happyCouple HappyCouple
			happyCouple.TheQuestion = q
			happyCouple.TheResponses = Responses

			HappyCouples = append(HappyCouples, happyCouple)

		}

		return c.Render("pages/session", fiber.Map{
			"HappyCouples": HappyCouples,
		})
	})

	// FORM ANSWERS
	app.Post("/session/:quizz_id/:session_id", func(c *fiber.Ctx) error {
		myQuizz := c.Params("quizz_id")
		mySession := c.Params("session_id")

		// If user asks to redo the quizz, new session is created
		if c.FormValue("redo") != "" {
			// id := 0
			// err := db.QueryRow(`Insert into quizz_sessions(quizz_id) values ($1) returning session_id`, myQuizz).Scan(&id)
			// check(err)
			id := newSession(myQuizz)

			url := fmt.Sprintf("/session/%s/%d", myQuizz, id)

			return c.Redirect(url)
		}

		//Get responses from form
		var count int
		err := db.QueryRow(`select count(*) from questions where quizz_id = $1`, myQuizz).Scan(&count)
		check(err)

		for i := 0; i < count; i++ {
			var QuestionID string
			err := db.QueryRow(`select question_id from questions where quizz_id = 1 order by question_id offset $1 limit 1`, i).Scan(&QuestionID)
			check(err)

			value := c.FormValue(QuestionID)

			var isCorrect bool
			errerr := db.QueryRow(`select iscorrect from responses where response_id = $1`, value).Scan(&isCorrect)
			check(errerr)

			//Incremente session score if answers are correct
			if isCorrect {
				_, err := db.Exec(`update quizz_sessions set score = score + 1  where session_id = $1`, mySession)
				check(err)
			}

			//Store answers
			_, errStore := db.Exec(`insert into sess_resp(session_id, response_id) values ($1, $2)`, mySession, QuestionID)
			check(errStore)

		}

		//Set session to finish
		_, err = db.Exec(`update quizz_sessions set finished = true where session_id = $1`, mySession)
		check(err)

		//Update general score
		// _, err = db.Exec(`insert into user_score(user_id, total_score) values quizz_sessions('1', score) where session_id = $1`, mySession)
		// check(err)

		//TODO:update general score in leader board

		url := fmt.Sprintf("/session/%s/%s", myQuizz, mySession)
		return c.Redirect(url)

	})

	// Renders CSS
	app.Static("/assets/", "./assets")
	// Renders pictures
	app.Static("/public/", "./public")

	// Open port to listen to
	log.Fatal(app.Listen(":3000"))

}
