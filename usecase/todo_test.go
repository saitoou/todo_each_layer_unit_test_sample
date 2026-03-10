package usecase

import (
	"context"
	"errors"
	"local-golang-prac/domain/entity"
	mock_service "local-golang-prac/mock/service"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_Usecase_GetTodoByID(t *testing.T) {
	ctx := context.Background()
	now := time.Now().UTC().Truncate(time.Microsecond)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSvc := mock_service.NewMockTodoService(ctrl)
	uc := NewTodoUseCase(mockSvc)
	resultTodo := &entity.Todo{
		ID:        1,
		Title:     "title1",
		Content:   "content1",
		CreatedAt: now,
	}
	outputTodo := &TodoUsecaseOutput{
		ID:        1,
		Title:     "title1",
		Content:   "content1",
		CreatedAt: now,
	}

	tests := []struct {
		name     string
		id       uint
		mockFunc func()
		want     *TodoUsecaseOutput
		wantErr  error
	}{
		{
			name: "特定のIDでTODOを取得した場合",
			id:   uint(1),
			mockFunc: func() {
				mockSvc.EXPECT().GetTodoByID(ctx, uint(1)).Return(resultTodo, nil)
			},
			want:    outputTodo,
			wantErr: nil,
		},
		{
			name: "特定のIDでTODOを0件取得した場合",
			id:   uint(10),
			mockFunc: func() {
				mockSvc.EXPECT().GetTodoByID(ctx, uint(10)).Return(&entity.Todo{}, nil)
			},
			want:    &TodoUsecaseOutput{},
			wantErr: nil,
		},
		{
			name: "特定のIDでTODOを取得できずエラー",
			id:   uint(1),
			mockFunc: func() {
				mockSvc.EXPECT().GetTodoByID(ctx, uint(1)).Return(nil, errors.New("failed to get todo by id"))
			},
			want:    nil,
			wantErr: errors.New("failed to get todo by id"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()
			got, err := uc.GetTodoByID(ctx, tt.id)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err)
		})
	}

}

func Test_Usecase_CreateTodo(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSvc := mock_service.NewMockTodoService(ctrl)
	uc := NewTodoUseCase(mockSvc)
	createTodo := &entity.Todo{
		Title:   "title1",
		Content: "content1",
	}

	tests := []struct {
		name     string
		todo     TodoUsecaseInput
		mockFunc func()
		wantErr  error
	}{
		{
			name: "Todoの生成に成功した場合",
			todo: TodoUsecaseInput{
				Title:   "title1",
				Content: "content1",
			},
			mockFunc: func() {
				mockSvc.EXPECT().CreateTodo(ctx, createTodo).Return(nil)
			},
			wantErr: nil,
		},
		{
			name: "TODOの生成が失敗した場合",
			todo: TodoUsecaseInput{
				Title:   "title1",
				Content: "content1",
			},
			mockFunc: func() {
				mockSvc.EXPECT().CreateTodo(ctx, createTodo).Return(errors.New("failed to create todo"))
			},
			wantErr: errors.New("failed to create todo"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()
			err := uc.CreateTodo(ctx, tt.todo)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func Test_Usecase_UpdateTodo(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSvc := mock_service.NewMockTodoService(ctrl)
	uc := NewTodoUseCase(mockSvc)
	updateTodo := &entity.Todo{
		ID:      1,
		Title:   "title1",
		Content: "content1",
	}

	tests := []struct {
		name     string
		todo     TodoUsecaseInput
		mockFunc func()
		wantErr  error
	}{
		{
			name: "Updateに成功した場合",
			todo: TodoUsecaseInput{
				ID:      1,
				Title:   "title1",
				Content: "content1",
			},
			mockFunc: func() {
				mockSvc.EXPECT().UpdateTodo(ctx, updateTodo).Return(nil)
			},
			wantErr: nil,
		},
		{
			name: "Updateに失敗した場合",
			todo: TodoUsecaseInput{
				ID:      1,
				Title:   "title1",
				Content: "content1",
			},
			mockFunc: func() {
				mockSvc.EXPECT().UpdateTodo(ctx, updateTodo).Return(errors.New("failed to update"))
			},
			wantErr: errors.New("failed to update"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()
			err := uc.UpdateTodo(ctx, tt.todo)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func Test_Usecase_DeleteTodo(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSvc := mock_service.NewMockTodoService(ctrl)
	uc := NewTodoUseCase(mockSvc)

	tests := []struct {
		name     string
		id       uint
		mockFunc func()
		wantErr  error
	}{
		{
			name: "Deleteが成功した場合",
			id:   uint(1),
			mockFunc: func() {
				mockSvc.EXPECT().DeleteTodo(ctx, uint(1)).Return(nil)
			},
			wantErr: nil,
		},
		{
			name: "Deleteが失敗した場合",
			id:   uint(1),
			mockFunc: func() {
				mockSvc.EXPECT().DeleteTodo(ctx, uint(1)).Return(errors.New("failed to delete"))
			},
			wantErr: errors.New("failed to delete"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()
			err := uc.DeleteTodo(ctx, tt.id)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
