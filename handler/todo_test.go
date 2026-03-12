package handler

import (
	"context"
	"errors"
	"fmt"
	mock_usecase "local-golang-prac/mock/usecase"
	"local-golang-prac/usecase"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"go.uber.org/mock/gomock"
)

func Test_handler_GetTodoByID(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUc := mock_usecase.NewMockTodoUseCase(ctrl)
	h := NewTodoHandler(mockUc)
	tests := []struct {
		name       string
		id         uint
		mockFunc   func()
		wantStatus int
	}{
		{
			name: "リクエストが正常に通る",
			id:   1,
			mockFunc: func() {
				mockUc.EXPECT().GetTodoByID(ctx, uint(1)).Return(
					&usecase.TodoUsecaseOutput{
						ID:      1,
						Title:   "title1",
						Content: "content1",
					}, nil)
			},

			wantStatus: http.StatusOK,
		},
		{
			name: "取得件数が０件の場合",
			id:   1,
			mockFunc: func() {
				mockUc.EXPECT().GetTodoByID(ctx, uint(1)).Return(
					&usecase.TodoUsecaseOutput{}, nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "リクエストが失敗した場合",
			id:   10,
			mockFunc: func() {
				mockUc.EXPECT().GetTodoByID(ctx, uint(10)).Return(
					nil, errors.New("todo not found"))
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/todo/%d", tt.id), nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/todo/:id")
			c.SetParamNames("id")
			c.SetParamValues(strconv.Itoa(int(tt.id)))
			tt.mockFunc()
			err := h.GetTodoByID(c, tt.id)
			if err != nil {
				e.HTTPErrorHandler(err, c)
			}
			if rec.Code != tt.wantStatus {
				t.Fatalf("expected status %d, got %d", tt.wantStatus, rec.Code)
			}
		})
	}
}

func Test_handler_CreateTodo(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUc := mock_usecase.NewMockTodoUseCase(ctrl)
	h := NewTodoHandler(mockUc)
	createTodo := &usecase.TodoUsecaseInput{
		Title:   "title1",
		Content: "content1",
	}
	tests := []struct {
		name       string
		body       string
		mockFunc   func()
		wantStatus int
	}{
		{
			name: "TODOの登録が成功する",
			body: `{"id":"0","Title":"title1","Content":"content1"}`,
			mockFunc: func() {
				mockUc.EXPECT().CreateTodo(ctx, createTodo).Return(nil)
			},
			wantStatus: http.StatusCreated,
		},
		{
			name:       "TODOの登録が失敗する",
			body:       `{"id":"0","Title":"title1","Content":`,
			mockFunc:   func() {},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "usecaseでエラーになる",
			body: `{"id":"0","Title":"title1","Content":"content1"}`,
			mockFunc: func() {
				mockUc.EXPECT().CreateTodo(ctx, &usecase.TodoUsecaseInput{
					Title:   "title1",
					Content: "content1",
				}).Return(errors.New("failed to create"))
			},
			wantStatus: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/todo", strings.NewReader(tt.body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()
			err := h.CreateTodo(c)
			if err != nil {
				e.HTTPErrorHandler(err, c)
			}
			if rec.Code != tt.wantStatus {
				t.Fatalf("expected status %d, got %d. body=%s", tt.wantStatus, rec.Code, rec.Body.String())
			}
		})
	}
}

func Test_handler_UpdateTodo(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUc := mock_usecase.NewMockTodoUseCase(ctrl)
	h := NewTodoHandler(mockUc)
	updateTodo := &usecase.TodoUsecaseInput{
		ID:      1,
		Title:   "title2",
		Content: "content2",
	}
	tests := []struct {
		name       string
		id         uint
		body       string
		todo       *usecase.TodoUsecaseInput
		mockFunc   func()
		wantStatus int
	}{
		{
			name: "正常に更新できる",
			id:   1,
			body: `{"id":"1","Title":"title2","Content":"content2"}`,
			todo: updateTodo,
			mockFunc: func() {
				mockUc.EXPECT().UpdateTodo(ctx, updateTodo).Return(nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "更新に失敗する",
			id:   1,
			body: `{"id":"1","Title":"title2","Content":"content2"}`,
			todo: updateTodo,
			mockFunc: func() {
				mockUc.EXPECT().UpdateTodo(ctx, updateTodo).Return(errors.New("failed to update"))
			},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name: "Json形式が壊れている",
			id:   1,
			body: `{"id":"1","Title":"title2","Content":`,
			todo: updateTodo,
			mockFunc: func() {
			},
			wantStatus: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/todo/%d", tt.id), strings.NewReader(tt.body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/todo/:id")
			c.SetParamNames("id")
			c.SetParamValues(strconv.Itoa(int(tt.id)))
			tt.mockFunc()
			err := h.UpdateTodo(c, tt.id)
			if err != nil {
				e.HTTPErrorHandler(err, c)
			}
			if rec.Code != tt.wantStatus {
				t.Fatalf("expected status %d, got %d", tt.wantStatus, rec.Code)
			}
		})
	}
}

func Test_handler_DeleteTodo(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUc := mock_usecase.NewMockTodoUseCase(ctrl)
	h := NewTodoHandler(mockUc)
	tests := []struct {
		name       string
		id         uint
		mockFunc   func()
		wantStatus int
	}{
		{
			name: "正常に削除できる",
			id:   1,
			mockFunc: func() {
				mockUc.EXPECT().DeleteTodo(ctx, uint(1)).Return(nil)
			},
			wantStatus: http.StatusNoContent,
		},
		{
			name: "削除に失敗する",
			id:   1,
			mockFunc: func() {
				mockUc.EXPECT().DeleteTodo(ctx, uint(1)).Return(errors.New("failed to delete"))
			},
			wantStatus: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/todo/%d", tt.id), nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/todo/:id")
			c.SetParamNames("id")
			c.SetParamValues(strconv.Itoa(int(tt.id)))
			tt.mockFunc()
			err := h.DeleteTodo(c, tt.id)
			if err != nil {
				e.HTTPErrorHandler(err, c)
			}
			if rec.Code != tt.wantStatus {
				t.Fatalf("expected status %d, got %d", tt.wantStatus, rec.Code)
			}
		})
	}
}
