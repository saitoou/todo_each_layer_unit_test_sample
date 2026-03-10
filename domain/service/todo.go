package service

import (
	"context"
	"local-golang-prac/domain/entity"
	"local-golang-prac/domain/repository"
	"local-golang-prac/utils"
	"time"
)

type TodoService struct {
	todoRepo repository.TodoRepository
}

func NewTodoService(todoRepo repository.TodoRepository) *TodoService {
	return &TodoService{
		todoRepo: todoRepo,
	}
}

func (svc *TodoService) GetTodoByID(ctx context.Context, id uint) ([]*entity.Todo, error) {

	ret, err := svc.todoRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (svc *TodoService) CreateTodo(ctx context.Context, todo *entity.Todo) error {

	todo.CreatedAt = time.Now().In(utils.JstLocation())
	err := svc.todoRepo.Create(ctx, todo)
	if err != nil {
		return err
	}
	return nil
}

func (svc *TodoService) UpdateTodo(ctx context.Context, todo *entity.Todo) error {

	todo.UpdatedAt = time.Now().In(utils.JstLocation())
	if err := svc.todoRepo.Update(ctx, todo); err != nil {
		return err
	}

	return nil
}

func (svc *TodoService) DeleteTodo(ctx context.Context, id uint) error {
	if err := svc.todoRepo.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}
