package postgres

import (
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/itimofeev/task2trip/backend"
	"github.com/itimofeev/task2trip/rest/models"
	"github.com/itimofeev/task2trip/rest/restapi/operations/offers"
	"github.com/itimofeev/task2trip/rest/restapi/operations/tasks"
	"github.com/itimofeev/task2trip/util"
	"github.com/rs/xid"
	"time"
)

func NewStore(connectURL string) backend.Store {
	opts, err := pg.ParseURL(connectURL)
	util.CheckErr(err, "pg.ParseURL")
	db := pg.Connect(opts)
	util.CheckErr(createSchema(db))

	db.OnQueryProcessed(func(event *pg.QueryProcessedEvent) {
		query, err := event.FormattedQuery()
		if err != nil {
			panic(err)
		}
		util.Log.Printf("%s %s", time.Since(event.StartTime), query)
	})

	store := &Store{
		db: db,
	}

	//db.Exec(`ALTER TABLE shares ADD CONSTRAINT shares_unique UNIQUE (user_id, basket_id)`)
	_, _ = db.Exec(`ALTER TABLE offers ADD CONSTRAINT offer__user_task__unique UNIQUE (user_id, task_id)`)
	//db.Exec(`CREATE INDEX tbl_col_text_pattern_ops_idx ON products(name text_pattern_ops)`)
	//
	_, _ = store.CreateUser("user1@gmail.com", "123")
	_, _ = store.CreateCategory("names.category.name.guide", "Гиды")
	_, _ = store.CreateCategory("names.category.name.route", "Маршруты")
	_, _ = store.CreateCategory("names.category.name.misc", "Остальное")

	return store
}

type Store struct {
	db *pg.DB
}

func (s *Store) ListOffers(user *backend.User, taskId string) (offers []*backend.Offer, err error) {
	err = s.db.RunInTransaction(func(tx *pg.Tx) error {
		task := &backend.Task{ID: taskId}
		if err := tx.Select(task); err != nil {
			return err
		}

		query := tx.Model(&offers).Where("task_id = ?", taskId)

		if task.UserID != user.ID {
			query = query.Where("user_id = ?", user.ID)
		}

		return query.Relation("User").Select()
	})

	return offers, err
}

func (s *Store) CreateOffer(user *backend.User, taskId string, params offers.CreateOfferBody) (offer *backend.Offer, err error) {
	err = s.db.RunInTransaction(func(tx *pg.Tx) error {
		task := &backend.Task{ID: taskId}
		if err := tx.Select(task); err != nil {
			return err
		}

		if task.UserID == user.ID {
			return util.NewBadRequestError("errors.offer.unable.create.offer.for.self.task", "unable to create offer for user's self task")
		}

		offer = &backend.Offer{
			ID:         xid.New().String(),
			UserID:     user.ID,
			CreateTime: time.Now(),
			User:       user,
			Comment:    params.Comment,
			Price:      *params.Price,
			TaskID:     task.ID,
		}
		return tx.Insert(offer)
	})

	return offer, err
}

func (s *Store) SearchTasks(user *backend.User, params tasks.SearchTasksParams) (tasks []*backend.Task, total int64, err error) {
	query := s.db.Model(&tasks)
	if params.UserID != nil && user != nil && user.ID == *params.UserID {
		query = query.Where("user_id = ?", user.ID)
	}

	if params.CategoryID != nil {
		query = query.Where("category_id = ?", *params.CategoryID)
	}

	if params.SearchString != nil {
		query = query.Where("name ilike ?", "%"+*params.SearchString+"%")
	}

	limit := 100
	if params.Limit != nil && *params.Limit > 0 && *params.Limit < 100 {
		limit = 100
	}

	skip := 0
	if params.Skip != nil && *params.Skip > 0 {
		skip = int(*params.Skip)
	}

	total_, err := query.
		Limit(limit).
		Offset(skip).
		Relation("Category").SelectAndCount()
	return tasks, int64(total_), err
}

func (s *Store) ListCategories() (categories []*backend.Category, err error) {
	return categories, s.db.Model(&categories).Select()
}

func (s *Store) CreateTask(user *backend.User, params *models.TaskCreateParams) (task *backend.Task, err error) {
	err = s.db.RunInTransaction(func(tx *pg.Tx) error {
		category := &backend.Category{ID: *params.CategoryID}
		if err := tx.Select(category); err != nil {
			return err
		}

		task = &backend.Task{
			ID:             xid.New().String(),
			UserID:         user.ID,
			CreateTime:     time.Now(),
			Name:           *params.Name,
			Description:    *params.Description,
			CategoryID:     category.ID,
			Category:       category,
			BudgetEstimate: *params.BudgetEstimate,
		}
		return tx.Insert(task)
	})
	return task, err
}

func (s *Store) GetUserByEmailAndPassword(email, password string) (user *backend.User, err error) {
	user = &backend.User{}
	return user, s.db.Model(user).Where("email = ? AND password = ?", email, password).Select()
}

func (s *Store) GetUserByID(id string) (*backend.User, error) {
	user := &backend.User{ID: id}
	return user, s.db.Select(user)
}

func (s *Store) CreateUser(email, password string) (*backend.User, error) {
	user := &backend.User{
		ID:       xid.New().String(),
		Email:    email,
		Password: password,
	}
	return user, s.db.Insert(user)
}

func (s *Store) CreateCategory(key, defaultValue string) (*backend.Category, error) {
	category := &backend.Category{
		ID:           xid.New().String(),
		Key:          key,
		DefaultValue: defaultValue,
	}
	return category, s.db.Insert(category)
}

func createSchema(db *pg.DB) error {
	for _, mdl := range []interface{}{
		(*backend.User)(nil),
		(*backend.Category)(nil),
		(*backend.Task)(nil),
		(*backend.Offer)(nil),
	} {
		err := db.CreateTable(mdl, &orm.CreateTableOptions{
			IfNotExists: true,
		})
		if err != nil {
			return err
		}
	}

	return nil
}
