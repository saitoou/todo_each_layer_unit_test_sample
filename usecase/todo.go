//go:generate mockgen -source=$GOFILE -package=mock_service -destination=../mock/service/todo.go
package usecase

import (
	"context"
	"local-golang-prac/domain/entity"
	"time"
)

type TodoService interface {
	GetTodoByID(ctx context.Context, id uint) (*entity.Todo, error)
	CreateTodo(ctx context.Context, todo *entity.Todo) error
	UpdateTodo(ctx context.Context, todo *entity.Todo) error
	DeleteTodo(ctx context.Context, id uint) error
}

type TodoUsecaseInput struct {
	ID      uint   `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	// CreatedAt time.Time `json:"created_at"`
	// UpdatedAt time.Time `json:"updated_at,omitempty"`
}

type TodoUsecaseOutput struct {
	ID        uint      `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

type TodoUsecase struct {
	todoSvc TodoService
}

func NewTodoUseCase(todoSvc TodoService) *TodoUsecase {
	return &TodoUsecase{todoSvc: todoSvc}
}

func (uc *TodoUsecase) GetTodoByID(ctx context.Context, id uint) (*TodoUsecaseOutput, error) {
	ret, err := uc.todoSvc.GetTodoByID(ctx, id)
	if err != nil {
		return nil, err
	}

	todo := &TodoUsecaseOutput{
		ID:        ret.ID,
		Title:     ret.Title,
		Content:   ret.Content,
		CreatedAt: ret.CreatedAt,
	}

	return todo, nil
}

func (uc *TodoUsecase) CreateTodo(ctx context.Context, todo TodoUsecaseInput) error {

	insertTodo := &entity.Todo{
		Title:   todo.Title,
		Content: todo.Content,
	}

	if err := uc.todoSvc.CreateTodo(ctx, insertTodo); err != nil {
		return err
	}

	return nil
}

func (uc *TodoUsecase) UpdateTodo(ctx context.Context, todo TodoUsecaseInput) error {

	updateTodo := &entity.Todo{
		ID:      todo.ID,
		Title:   todo.Title,
		Content: todo.Content,
	}

	if err := uc.todoSvc.UpdateTodo(ctx, updateTodo); err != nil {
		return err
	}

	return nil
}

func (uc *TodoUsecase) DeleteTodo(ctx context.Context, id uint) error {

	if err := uc.todoSvc.DeleteTodo(ctx, id); err != nil {
		return err
	}

	return nil
}
