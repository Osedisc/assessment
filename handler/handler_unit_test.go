package handler_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Osedisc/assessment/handler"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestExp01PostExpenses(t *testing.T) {
	mockData := []byte(`{
		"title": "test1",
	  "amount": 2,
	  "note": "test3",
	  "tags": ["test4", "test5"]
		}`)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/expenses", bytes.NewBuffer(mockData))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	expected := `{"id":1,"title":"test1","amount":2,"note":"test3","tags":["test4","test5"]}`
	exmockRows := sqlmock.NewRows([]string{"id"}).AddRow("1")

	db, mock, err := sqlmock.New()
	mock.ExpectQuery("INSERT INTO expenses .*").WithArgs().WillReturnRows(exmockRows)

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	handler.DB = db
	c := e.NewContext(req, rec)

	if assert.NoError(t, handler.PostExpenses(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Equal(t, expected, strings.TrimSpace(rec.Body.String()))
	}
}

func TestExp02GetExpensesbyID(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/expenses/1", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	expected := `{"id":1,"title":"test1","amount":2,"note":"test3","tags":["test4","test5"]}`
	exmockRows := sqlmock.NewRows([]string{"id", "title", "amount", "note", "tags"}).
	AddRow(1, "test1", 2, "test3", `{test4,test5}`)

	db, mock, err := sqlmock.New()
	mock.ExpectPrepare("SELECT (.+) FROM expenses WHERE id=?").
		ExpectQuery().
		WithArgs("1").
		WillReturnRows(exmockRows)

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	handler.DB = db
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	if assert.NoError(t, handler.GetExpensebyid(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, expected, strings.TrimSpace(rec.Body.String()))
	}
}

//I don't Know How to test EXP03 UpdateExpense T_T

func TestExp04GetAllExpenses(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/expenses", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	
	expected := `[{"id":1,"title":"test1","amount":2,"note":"test3","tags":["test4","test5"]},{"id":2,"title":"test1","amount":3,"note":"test3","tags":["test4","test5"]}]`
	exmockRows := sqlmock.NewRows([]string{"id", "title", "amount", "note", "tags"}).
	AddRow(1, "test1", 2, "test3", `{test4,test5}`).
	AddRow(2, "test1", 3, "test3", `{test4,test5}`)

	db, mock, err := sqlmock.New()
	mock.ExpectPrepare("SELECT (.+) FROM expenses").
		ExpectQuery().
		WillReturnRows(exmockRows)

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	handler.DB = db
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")
	if assert.NoError(t, handler.GetAllExpenses(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, expected, strings.TrimSpace(rec.Body.String()))
	}
}