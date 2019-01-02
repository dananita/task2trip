package backend

import (
	"github.com/itimofeev/task2trip/rest/models"
	"github.com/itimofeev/task2trip/rest/restapi/operations/tasks"
	"time"
)

type User struct {
	ID       string
	Email    string `sql:",unique,notnull"`
	Password string `sql:",notnull"`
}

type Category struct {
	ID           string
	Key          string `sql:",unique,notnull"`
	DefaultValue string `sql:",notnull"`
}

type Task struct {
	ID             string
	UserID         string
	Name           string `sql:",notnull"`
	Description    string `sql:",notnull"`
	CategoryID     string `sql:",notnull"`
	Category       *Category
	CreateTime     time.Time `sql:",notnull"`
	BudgetEstimate int64     `sql:",notnull"`
}

type Store interface {
	GetUserByID(token string) (*User, error)
	GetUserByEmailAndPassword(email, password string) (user *User, err error)
	CreateUser(email, password string) (user *User, err error)
	CreateTask(user *User, params *models.TaskCreateParams) (task *Task, err error)
	ListCategories() (categories []*Category, err error)
	SearchTasks(user *User, params tasks.SearchTasksParams) (tasks []*Task, total int64, err error)
}
