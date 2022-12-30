package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

func UpdateExpense(c echo.Context) error {
	id := c.Param("id")
	ex := Expenses{}
	err := c.Bind(&ex)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}
	upex, err := DB.Prepare("UPDATE expenses SET title=$2, amount=$3, note=$4, tags=$5 WHERE id=$1;")
	if err != nil {
		log.Fatal("Can't prepare expenses for update", err)
	}
	if _, err := upex.Exec(id, ex.Title, ex.Amount, ex.Note, pq.Array(&ex.Tags)); err != nil {
		log.Fatal("error execute update ", err)
	}
	ex.ID, _ = strconv.Atoi(id)
	return c.JSON(http.StatusOK, ex)
}
