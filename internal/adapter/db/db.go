package db

import "github.com/mokan-r/place-for-your-thoughts/internal/model"

type Storager interface {
	AddPost(model model.Post) error
	GetEntriesCount() (res int, err error)
	GetEntry(id string) (model model.Post, err error)
	GetEntriesWithOffset(offset int, limit int) (res []model.Post, err error)
}
