package main

import (
	"fmt"
	"net/http"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

var (
	db *gorm.DB
	v  *Validator
)

type Validator struct {
	validate *validator.Validate
}

type Task struct {
	ID        uint           `json:"id"`
	Objective string         `validate:"required" json:"objective"`
	IsDone    bool           `gorm:"default:false" json:"is_done"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-"`
}

func (v *Validator) ValidateStruct(i interface{}) error {
	if err := v.validate.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return nil
}

func (v *Validator) ValidateVar(field interface{}, tag string) error {
	if err := v.validate.Var(field, tag); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return nil
}

func main() {
	connectToDatabase()
	v = &Validator{validate: validator.New()}

	e := echo.New()

	e.GET("/tasks", getTasks)
	e.GET("/tasks/:id", getTaskById)
	e.POST("/tasks", createTask)
	e.DELETE("/tasks/:id", deleteTaskById)

	e.Logger.Fatal(e.Start(":8000"))
}

func connectToDatabase() {
	var err error
	db, err = gorm.Open(sqlite.Open("test.sqlite"), &gorm.Config{})

	db.AutoMigrate(&Task{})

	if err != nil {
		panic(err)
	}

	fmt.Println("Connected to Database")
}

func getTasks(c echo.Context) error {
	var tasks []Task
	db.Find(&tasks)

	return c.JSON(http.StatusOK, tasks)
}

func getTaskById(c echo.Context) error {
	var task Task

	id := c.Param("id")
	if err := v.ValidateVar(id, "required,number"); err != nil {
		return err
	}

	db.Find(&task, id)

	return c.JSON(http.StatusOK, task)
}

func createTask(c echo.Context) error {
	var task Task

	if err := c.Bind(&task); err != nil {
		return err
	}

	if err := v.ValidateStruct(&task); err != nil {
		return err
	}

	db.Create(&task)

	return c.JSON(http.StatusCreated, task)
}

func deleteTaskById(c echo.Context) error {
	id := c.Param("id")
	if err := v.ValidateVar(id, "required,number"); err != nil {
		return err
	}

	db.Delete(&Task{}, id)

	return c.NoContent(http.StatusNoContent)
}

func editTaskById(c echo.Context) error {
	var task Task

	id := c.Param("id")
	if err := v.ValidateVar(id, "required,number"); err != nil {
		return err
	}

	db.Find(&task, id)
}
