package student

import "net/http"

func New() http.HandlerFunc {
	return func(writer http.ResponseWriter, req *http.Request) {
		writer.Write([]byte("Welcome, to Student's API"))
	}
}
