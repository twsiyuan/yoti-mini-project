package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	yoti "github.com/getyoti/yoti-go-sdk"
)

type comment struct {
	UserID    string
	FullName  string
	Email     string
	Phone     string
	Text      string
	HTML      template.HTML
	Anonymous bool
}

func recoveryHandler(outputErr bool, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				stack := make([]byte, 2048)
				stack = stack[:runtime.Stack(stack, false)]
				displayErr := "Internal server error"
				if outputErr {
					displayErr = fmt.Sprintf("Unexpected error: %v, in %s", err, stack)
				}

				fmt.Fprintf(os.Stderr, "Unexpected error: %v, in %s\n", err, stack)
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, displayErr)

			}
		}()

		next.ServeHTTP(w, req)
	})
}

func imageHandler(db *sql.DB, anonymousFile string) http.Handler {
	anonymousRaws, err := ioutil.ReadFile(anonymousFile)
	if err != nil {
		panic(err)
	}
	anonymousExt := filepath.Ext(anonymousFile)
	switch anonymousExt {
	case ".jpg", ".jpeg":
		anonymousExt = "jpg"
	case ".png":
		anonymousExt = "png"
	}

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		id := req.URL.Path
		if len(id) <= 0 {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		var raws []byte
		var format string

		if id == "anonymous" {
			raws = anonymousRaws
			format = anonymousExt
		} else {
			if err := db.QueryRowContext(req.Context(), "SELECT `selfie`, `selfie_format` FROM `users` WHERE `user_id` = ?", id).Scan(&raws, &format); err == sql.ErrNoRows {
				w.WriteHeader(http.StatusNotFound)
				return
			} else if err != nil {
				panic(err)
			}
		}

		switch format {
		case "jpg":
			w.Header().Set("CONTENT-TYPE", "image/jpeg")
			break
		case "png":
			w.Header().Set("CONTENT-TYPE", "image/png")
			break
		}
		w.Write(raws)
	})
}

func postHandler(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		session, err := store.Get(req, sessionID)
		if err != nil {
			panic(err)
		}

		id, login := session.Values[sessionProfileKey]
		if login {
			text := req.FormValue("comment")
			anonymous := req.FormValue("anonymous") != ""

			if _, err := db.ExecContext(req.Context(), "INSERT INTO `comments`(`user_id`, `comment`, `anonymous`)VALUES(?, ?, ?)", id, text, anonymous); err != nil {
				panic(err)
			}
			session.AddFlash("Posted comment")
		} else {
			session.AddFlash("Need login.")

		}

		session.Save(req, w)
		http.Redirect(w, req, "/", http.StatusTemporaryRedirect)
	})
}

func mainHandler(db *sql.DB) http.Handler {
	tmpl := template.Must(template.ParseFiles(htmlTemplte))
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// TODO: Use client render, build app API
		rows, err := db.QueryContext(req.Context(), "SELECT `c`.`user_id`, `c`.`anonymous`, `c`.`comment`, `u`.`full_name`, `u`.`email_address`, `u`.`phone_number` FROM `comments` AS `c` LEFT JOIN `users` AS `u` ON `c`.`user_id`=`u`.`user_id` ORDER BY `c`.`create_time` DESC ")
		if err != nil {
			panic(err)
		}

		comments := make([]comment, 0)
		for rows.Next() {
			c := comment{}
			if err := rows.Scan(&c.UserID, &c.Anonymous, &c.Text, &c.FullName, &c.Email, &c.Phone); err != nil {
				panic(err)
			}
			c.HTML = template.HTML(strings.Replace(template.HTMLEscapeString(c.Text), "\n", "<br/>", -1))
			comments = append(comments, c)
		}

		options := struct {
			YotiAppID    string
			Login        bool
			Comments     []comment
			FlashMessage string
			AnyFlash     bool
		}{
			YotiAppID: appID,
			Comments:  comments,
		}

		session, err := store.Get(req, sessionID)
		if err != nil {
			panic(err)
		}
		_, options.Login = session.Values[sessionProfileKey]

		if flashes := session.Flashes(); flashes != nil {
			options.FlashMessage = flashes[0].(string)
			options.AnyFlash = true
		}
		tmpl.Execute(w, options)
	})
}

func callbackHandler(c *yoti.Client, db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		// Query
		token := req.URL.Query().Get("token")
		profile, err := c.GetUserProfile(token)
		if err != nil {
			panic(err)
		}

		id := profile.ID
		email := profile.EmailAddress
		mobile := profile.MobileNumber
		name := profile.FullName
		selfie := ([]byte)(nil)
		format := ""

		if pselfie := profile.Selfie; pselfie != nil {
			selfie = pselfie.Data
			switch pselfie.Type {
			case yoti.ImageTypeJpeg:
				format = "jpg"
			case yoti.ImageTypePng:
				format = "png"
			}
		}

		// Insert database
		tx, err := db.BeginTx(req.Context(), nil)
		defer tx.Rollback()
		if err != nil {
			panic(err)
		}

		var userID string
		if err := tx.QueryRowContext(req.Context(), "SELECT `user_id` FROM `users` WHERE `user_id`=?", profile.ID).Scan(&userID); err == sql.ErrNoRows {
			if _, err := tx.ExecContext(req.Context(), "INSERT INTO `users`(`user_id`, `email_address`, `phone_number`, `selfie`, `selfie_format`, `full_name`)VALUES(?, ?, ?, ?, ?, ?)", id, email, mobile, selfie, format, name); err != nil {
				panic(err)
			}
		} else if err != nil {
			panic(err)
		}

		if err := tx.Commit(); err != nil {
			panic(err)
		}

		// Session
		session, err := store.Get(req, sessionID)
		if err != nil {
			panic(err)
		}
		session.Values[sessionProfileKey] = id
		if err := session.Save(req, w); err != nil {
			panic(err)
		}

		http.Redirect(w, req, "/", http.StatusTemporaryRedirect)
	})
}
