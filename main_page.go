package main

import (
	"database/sql"
	"fmt"
	"log"
	"runtime"

	D "kwizz/functions"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/lib/pq"
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

type ViewSelectedQuestion struct {
	Session_ID    int
	V_Question_ID int
	V_Question    string
	V_Response_ID []int32
	V_Answer      []string
	V_isCorrrect  []bool
	Selected      []bool
}

func check(err error) (b bool) {
	if err != nil {
		_, filename, line, _ := runtime.Caller(1)
		log.Printf("[error] %s:%d %v", filename, line, err)
		b = true
	}
	return
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

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// HOME PAGE
	app.Get("/", func(c *fiber.Ctx) error {
		var Categories []Categorie
		var pseudonym string
		err := db.QueryRow("select pseudonym from users").Scan(&pseudonym)
		check(err)

		rows, err := db.Query("select cat_id, cat_name, cat_short_name, cat_description, cat_image from categories")
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

		var Score string
		err = db.QueryRow("select coalesce((select sum(score) from quizz_sessions where user_id = 1), 0)").Scan(&Score)
		check(err)

		// Render index template
		return c.Render("main_page", fiber.Map{
			"Pseudonym":  pseudonym,
			"Categories": Categories,
			"Score":      Score,
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

		rows, err := db.Query("select quizz_id, quizz_title, quizz_description, created_at, cat_id from quizzes where cat_id = $1", myCat)
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

		var finished bool
		errerr := db.QueryRow(`select finished from quizz_sessions where session_id = $1`, mySession).Scan(&finished)
		check(errerr)

		// CASE 1 : quizz was already taken
		if finished {
			var Score int
			errerr := db.QueryRow(`select score from quizz_sessions where session_id = $1`, mySession).Scan(&Score)
			check(errerr)

			var ViewSelectedQuestions []ViewSelectedQuestion
			rows, err := db.Query("select session_id, question_id, question, response_ids, answers, isCorrects, selected from v_selected_questions where session_id = $1", mySession)
			check(err)
			defer rows.Close()

			for rows.Next() {
				var view ViewSelectedQuestion
				err := rows.Scan(&view.Session_ID, &view.V_Question_ID, &view.V_Question, (*pq.Int32Array)(&view.V_Response_ID), (*pq.StringArray)(&view.V_Answer), (*pq.BoolArray)(&view.V_isCorrrect), (*pq.BoolArray)(&view.Selected))
				check(err)
				ViewSelectedQuestions = append(ViewSelectedQuestions, view)
			}

			// fmt.Println(ViewSelectedQuestions)

			return c.Render("pages/finished_session", fiber.Map{
				"Score": Score,
				"View":  ViewSelectedQuestions,
			})

		}

		// CASE 2 : first time taking quizz
		var Questions []Question

		rows, err := db.Query("select question_id, quizz_id, question, order_questions from questions where quizz_id = $1", myQuizz)
		check(err)
		defer rows.Close()

		for rows.Next() {
			var question Question
			err := rows.Scan(&question.Q_Question_ID, &question.Quizz_ID, &question.Question, &question.Order_questions)
			check(err)
			Questions = append(Questions, question)
		}

		var HappyCouples []HappyCouple

		for _, q := range Questions {
			var Responses []Response

			otherRows, err := db.Query("select response_id, question_id, answer, isCorrect  from responses where question_id = $1", q.Q_Question_ID)
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
			err := db.QueryRow(`select question_id from questions where quizz_id = $1 order by question_id offset $2 limit 1`, myQuizz, i).Scan(&QuestionID)
			check(err)

			value := c.FormValue(QuestionID)
			// fmt.Println(value)

			var isCorrect bool
			errerr := db.QueryRow(`select iscorrect from responses where response_id = $1`, value).Scan(&isCorrect)
			check(errerr)

			//Incremente session score if answers are correct
			if isCorrect {
				_, err := db.Exec(`update quizz_sessions set score = score + 1  where session_id = $1`, mySession)
				check(err)
			}

			//Store answers
			_, errStore := db.Exec(`insert into sess_resp(session_id, response_id) values ($1, $2)`, mySession, value)
			check(errStore)

		}

		//Set session to finish
		_, err = db.Exec(`update quizz_sessions set finished = true where session_id = $1`, mySession)
		check(err)

		url := fmt.Sprintf("/session/%s/%s", myQuizz, mySession)
		return c.Redirect(url)

	})

	// Renders CSS
	app.Static("/assets/", "./assets")
	// Renders pictures
	app.Static("/public/", "./public")

	// Open port to listen to
	log.Fatal(app.Listen(":19000"))

}
