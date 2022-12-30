//go:build integration

package handler_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/Osedisc/assessment/handler"
	"github.com/stretchr/testify/assert"
)

// Test EXP01 post expenses

func TestCreateExpenses(t *testing.T) {
	body := bytes.NewBufferString(
		`{
		"title": "strawberry smoothie",
		"amount": 79,
		"note": "night market promotion discount 10 bath", 
		"tags": ["food", "beverage"]
		}`,
	)
	var e handler.Expenses

	res := request(http.MethodPost, uri("expenses"), body)
	err := res.Decode(&e)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, res.StatusCode)
	assert.NotEqual(t, 0, e.ID)
	assert.Equal(t, "strawberry smoothie", e.Title)
	assert.Equal(t, 79, e.Amount)
	assert.Equal(t, "night market promotion discount 10 bath", e.Note)
	assert.Equal(t, []string{"food", "beverage"}, e.Tags)
}

// Test EXP02 get expenses by ID
func TestGetExpensebyID(t *testing.T) {
	e := seedExpenses(t)
	var latest handler.Expenses
	res := request(http.MethodGet, uri("expenses", strconv.Itoa(e.ID)), nil)
	err := res.Decode(&latest)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, e.ID, latest.ID)
	assert.NotEmpty(t, latest.Title)
	assert.NotEmpty(t, latest.Amount)
	assert.NotEmpty(t, latest.Note)
	assert.NotEmpty(t, latest.Tags)
}

// Test EXP03 update expenses
func TestUpdateExpense(t *testing.T) {
	id := seedExpenses(t).ID
	e := handler.Expenses{
		ID:     id,
		Title:  "buy a new phone",
		Amount: 39000,
		Note:   "buy a new phone",
		Tags:   []string{"gadget", "shopping"},
	}
	turntojson, _ := json.Marshal(e)
	res := request(http.MethodPut, uri("expenses", strconv.Itoa(id)), bytes.NewBuffer(turntojson))
	var latest handler.Expenses
	err := res.Decode(&latest)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, latest.Title, e.Title)
	assert.Equal(t, latest.Amount, e.Amount)
	assert.Equal(t, latest.Note, e.Note)
	assert.Equal(t, latest.Tags, e.Tags)
}

// Test EXP04 get all expenses
func TestGetAllExpenses(t *testing.T) {
	seedExpenses(t)
	var ex []handler.Expenses

	res := request(http.MethodGet, uri("expenses"), nil)
	err := res.Decode(&ex)

	assert.Nil(t, err)
	assert.EqualValues(t, http.StatusOK, res.StatusCode)
	assert.Greater(t, len(ex), 0)
}

func seedExpenses(t *testing.T) handler.Expenses {
	var e handler.Expenses
	body := bytes.NewBufferString(
		`{
		"title": "strawberry smoothie",
		"amount": 79,
		"note": "night market promotion discount 10 bath", 
		"tags": ["food", "beverage"]
		}`,
	)
	err := request(http.MethodPost, uri("expenses"), body).Decode(&e)
	if err != nil {
		t.Fatal("can't create uomer:", err)
	}
	return e
}
func uri(paths ...string) string {
	host := "http://localhost:2565"
	if paths == nil {
		return host
	}

	url := append([]string{host}, paths...)
	return strings.Join(url, "/")
}

func request(method, url string, body io.Reader) *Response {
	req, _ := http.NewRequest(method, url, body)
	req.Header.Add("Content-Type", "application/json")
	client := http.Client{}
	res, err := client.Do(req)
	return &Response{res, err}
}

type Response struct {
	*http.Response
	err error
}

func (r *Response) Decode(v interface{}) error {
	if r.err != nil {
		return r.err
	}

	return json.NewDecoder(r.Body).Decode(v)
}
