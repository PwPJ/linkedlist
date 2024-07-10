package v2

import (
	"fmt"
	"linkedlist/linkedlist"
	"net/http"
	"strconv"

	"github.com/go-playground/validator"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type ListEntity struct {
	Index uint `json:"index"`
	Value int  `json:"value" validate:"required"`
}

type server struct {
	list *linkedlist.LinkedList
}

type customValidator struct {
	validator *validator.Validate
}

func (cv *customValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

func V2() (*echo.Echo, error) {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Validator = &customValidator{validator: validator.New()}

	s := &server{}
	l := linkedlist.NewLinkedList()
	s.list = l
	e.POST("/numbers/:index/:value", s.Insert)
	e.DELETE("/numbers/:index", s.Remove)
	e.GET("/numbers/value/:value", s.Find)
	e.GET("/numbers/index/:index", s.Get)

	return e, nil
}

func (s *server) Insert(c echo.Context) error {
	data := ListEntity{}
	if err := c.Bind(&data); err != nil {
		return err
	}
	if err := c.Validate(&data); err != nil {
		return err
	}

	ok := s.list.Insert(data.Index, data.Value)
	if !ok {
		return echo.NewHTTPError(echo.ErrBadRequest.Code, "Invalid index")
	}
	c.JSON(http.StatusCreated, data)
	return nil
}

func (s *server) Remove(c echo.Context) error {
	indexStr := c.Param("index")
	index, err := strconv.ParseUint(indexStr, 10, 32)
	if err != nil {
		return echo.NewHTTPError(echo.ErrBadRequest.Code, "Invalid index")
	}

	ok := s.list.Remove(uint(index))
	if !ok {
		return echo.NewHTTPError(echo.ErrNotFound.Code, "Index not found")
	}

	c.NoContent(http.StatusOK)
	return nil

}

func (s *server) Find(c echo.Context) error {
	valueStr := c.Param("value")
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		fmt.Println(err)
		return echo.NewHTTPError(echo.ErrBadRequest.Code, "Invalid value")
	}

	index, ok := s.list.Find(value)
	if !ok {
		return echo.NewHTTPError(echo.ErrNotFound.Code, "Value not found")
	}

	data := ListEntity{
		Index: index,
		Value: value,
	}
	c.JSON(http.StatusOK, data)
	return nil
}

func (s *server) Get(c echo.Context) error {

	indexStr := c.Param("index")
	index, err := strconv.ParseUint(indexStr, 10, 32)
	if err != nil {
		return echo.NewHTTPError(echo.ErrBadRequest.Code, "Invalid index")
	}

	value, ok := s.list.Get(uint(index))
	if !ok {
		return echo.NewHTTPError(echo.ErrNotFound.Code, "Index not found")
	}
	data := ListEntity{
		Index: uint(index),
		Value: value,
	}

	c.JSON(http.StatusOK, data)
	return nil
}
