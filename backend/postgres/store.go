package postgres

import (
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/itimofeev/task2trip/backend"
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
	//db.Exec(`ALTER TABLE goods ADD CONSTRAINT goods_unique UNIQUE (user_id, basket_id, goods_id)`)
	//db.Exec(`CREATE INDEX tbl_col_text_pattern_ops_idx ON products(name text_pattern_ops)`)
	//
	_, _ = store.CreateUser("user1@gmail.com", "123")
	//store.CreateUser("user2@gmail.com", "123")

	return store
}

type Store struct {
	db *pg.DB
}

func (s *Store) GetUserByEmailAndPassword(email, password string) (user *backend.User, err error) {
	user = &backend.User{}
	return user, s.db.Model(user).Where("email = ? AND password = ?", email, password).Select()
}

func (s *Store) GetUserByID(id string) (*backend.User, error) {
	user := &backend.User{ID: id}
	err := s.db.Select(user)
	return user, err
}

func (s *Store) CreateUser(email, password string) (*backend.User, error) {
	user := &backend.User{
		ID:       xid.New().String(),
		Email:    email,
		Password: password,
	}
	return user, s.db.Insert(user)
}

func createSchema(db *pg.DB) error {
	for _, mdl := range []interface{}{
		(*backend.User)(nil),
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
