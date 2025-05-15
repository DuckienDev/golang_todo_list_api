package storage

import (
	"context"
	"golang_todo_list_api/modules/item/model"
)

func (s *sqlStore) CreateItem(ctx context.Context, data *model.TodoItemCreation) error {
	if err := s.db.Create(&data).Error; err != nil {

		return err
	}
	return nil

}
