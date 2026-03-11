//go:generate mockgen -source=$GOFILE -package=mock_usecase -destination=../mock/usecase/todo.go

package handler

import (
	"context"
	"local-golang-prac/gen/openapi/common"
	"local-golang-prac/usecase"
	"net/http"

	"github.com/labstack/echo/v4"
)

type TodoUseCase interface {
	GetTodoByID(ctx context.Context, id uint) (*usecase.TodoUsecaseOutput, error)
	CreateTodo(ctx context.Context, todo usecase.TodoUsecaseInput) error
	UpdateTodo(ctx context.Context, todo usecase.TodoUsecaseInput) error
	DeleteTodo(ctx context.Context, id uint) error
}

type TodoHandler struct {
	todoUseCase TodoUseCase
}

func NewTodoHandler(todoUseCase TodoUseCase) *TodoHandler {
	return &TodoHandler{
		todoUseCase: todoUseCase,
	}
}

func (h *TodoHandler) GetTodoByID(c echo.Context, id uint) error {
	ctx := c.Request().Context()

	ret, err := h.todoUseCase.GetTodoByID(ctx, id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, ret)
}

func (h *TodoHandler) CreateTodo(c echo.Context) error {
	ctx := c.Request().Context()
	req := common.TodoCreateRequest{}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	createTodo := usecase.TodoUsecaseInput{
		Title:   req.Title,
		Content: req.Content,
	}

	err := h.todoUseCase.CreateTodo(ctx, createTodo)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, nil)
}

func (h *TodoHandler) UpdateTodo(c echo.Context, id uint) error {
	ctx := c.Request().Context()
	req := common.TodoUpdateRequest{}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	updateTodo := usecase.TodoUsecaseInput{
		ID:      id,
		Title:   req.Title,
		Content: req.Content,
	}

	err := h.todoUseCase.UpdateTodo(ctx, updateTodo)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, nil)
}

func (h *TodoHandler) DeleteTodo(c echo.Context, id uint) error {
	ctx := c.Request().Context()

	if err := h.todoUseCase.DeleteTodo(ctx, id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.NoContent(http.StatusNoContent)
}
