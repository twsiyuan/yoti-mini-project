package main

import (
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

func newCookieStore() *sessions.CookieStore {
	codecs := securecookie.CodecsFromPairs(([]byte)("DCfkw328jdcnslcnslpwoe3SDSF3er9fi"))

	cs := &sessions.CookieStore{
		Codecs: codecs,
		Options: &sessions.Options{
			Path:   "/",
			MaxAge: 86400 * 30,
		},
	}

	cs.MaxAge(cs.Options.MaxAge)
	return cs
}
