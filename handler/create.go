package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

func PostExpenses(c echo.Context) error {
	ex := Expenses{}
	err := c.Bind(&ex)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	row := DB.QueryRow("INSERT INTO expenses (title, amount, note, tags) VALUES ($1, $2, $3,$4) RETURNING id", ex.Title, ex.Amount, ex.Note, pq.Array(&ex.Tags))
	err = row.Scan(&ex.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, ex)
}
