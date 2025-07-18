package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/vandannandwana/students-api/internal/storage"
	"github.com/vandannandwana/students-api/internal/types"
	"github.com/vandannandwana/students-api/internal/utils/response"
)

func New(storage storage.Storage) http.HandlerFunc {

	return func(writer http.ResponseWriter, req *http.Request) {
		slog.Info("Creating a student")

		var student types.Student
		err := json.NewDecoder(req.Body).Decode(&student)

		if errors.Is(err, io.EOF) {
			response.WriteJson(writer, http.StatusBadRequest, response.GeneralError(fmt.Errorf("empty body")))
			return
			// slog.Error("JSON body is empty")
		}

		if err != nil {
			response.WriteJson(writer, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		//Request Validation
		if err := validator.New().Struct(student); err != nil {
			validateErrs := err.(validator.ValidationErrors)
			response.WriteJson(writer, http.StatusBadRequest, response.ValidationError(validateErrs))
			return
		}

		// writer.Write([]byte("Welcome, to Student's API"))

		lastId, err := storage.CreateStudent(student.Name, student.Email, student.Age)

		if err != nil {
			response.WriteJson(writer, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		slog.Info("Student Created Successfully", slog.String("StudentId", fmt.Sprintf("%d", lastId)))

		response.WriteJson(writer, http.StatusCreated, map[string]int64{"id": lastId})

	}
}
