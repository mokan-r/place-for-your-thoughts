package postgresdb

import (
	"github.com/jmoiron/sqlx"
	"github.com/mokan-r/place-for-your-thoughts/internal/model"
)

type Postgres struct {
	Client *sqlx.DB
}

func New(client *sqlx.DB) *Postgres {
	return &Postgres{client}
}

func (p *Postgres) AddPost(model model.Post) error {
	query := `insert into posts (name, text) values ($1, $2)`
	if _, err := p.Client.Query(query, model.Name, model.Text); err != nil {
		return err
	}
	return nil
}

func (p *Postgres) GetEntriesCount() (res int, err error) {
	query := `select count(*) from posts`
	err = p.Client.Get(&res, query)
	return
}

func (p *Postgres) GetEntriesWithOffset(offset int, limit int) (res []model.Post, err error) {
	query := `select name, id from posts limit $1 offset $2`
	err = p.Client.Select(&res, query, offset, limit)
	return
}

func (p *Postgres) GetEntry(id string) (model model.Post, err error) {
	query := `select name, text from posts where id = $1`
	err = p.Client.Get(&model, query, id)
	return
}
