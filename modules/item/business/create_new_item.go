package business

import (
	"context"
	"golang_todo_list_api/modules/item/model"
	"strings"
)

type CreateItemStorage interface {
	CreateItem(ctx context.Context, data *model.TodoItemCreation) error
}

type createItemBusines struct {
	store CreateItemStorage
}

func NewCreateItemBussines(store CreateItemStorage) *createItemBusines {
	return &createItemBusines{store: store}
}

func (business *createItemBusines) CreateNewItem(ctx context.Context, data *model.TodoItemCreation) error {
	title := strings.TrimSpace(data.Title)
	if title == "" {
		return model.ErrTitleIsBlank
	}

	if err := business.store.CreateItem(ctx, data); err != nil {
		return err
	}

	return nil
}
