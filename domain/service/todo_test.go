package service

import (
	"context"
	"errors"
	"local-golang-prac/domain/entity"
	mock_repository "local-golang-prac/mock/repository"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_Service_GetTodoByID(t *testing.T) {
	ctx := context.Background()
	now := time.Now().UTC().Truncate(time.Millisecond)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mock_repository.NewMockTodoRepository(ctrl)
	svc := NewTodoService(mockRepo)
	wantTodos := []*entity.Todo{
		{
			ID:        1,
			Title:     "title1",
			Content:   "content1",
			CreatedAt: now,
		},
		{
			ID:        2,
			Title:     "title2",
			Content:   "content2",
			CreatedAt: now,
		},
		{
			ID:        3,
			Title:     "title3",
			Content:   "content3",
			CreatedAt: now,
		},
	}

	tests := []struct {
		name     string
		id       uint
		mockFunc func()
		want     []*entity.Todo
		wantErr  error
	}{
		{
			name: "Todoを取得した場合",
			id:   1,
			mockFunc: func() {
				mockRepo.EXPECT().FindByID(ctx, uint(1)).Return(wantTodos, nil)
			},
			want:    wantTodos,
			wantErr: nil,
		},
		{
			name: "Todoが１件もなかった場合",
			id:   10,
			mockFunc: func() {
				mockRepo.EXPECT().FindByID(ctx, uint(10)).Return([]*entity.Todo{}, nil)
			},
			want:    []*entity.Todo{},
			wantErr: nil,
		},
		{
			name: "DBエラーの場合",
			id:   1,
			mockFunc: func() {
				mockRepo.EXPECT().FindByID(ctx, uint(1)).Return(nil, errors.New("failed to connect DB"))
			},
			want:    nil,
			wantErr: errors.New("failed to connect DB"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()
			got, err := svc.GetTodoByID(ctx, tt.id)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func Test_Service_CreateTodo(t *testing.T) {
	ctx := context.Background()
	now := time.Now().UTC().Truncate(time.Millisecond)
	wantTodoData := &entity.Todo{ID: 1, Title: "title1", Content: "content1", CreatedAt: now}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mock_repository.NewMockTodoRepository(ctrl)
	svc := NewTodoService(mockRepo)

	tests := []struct {
		name     string
		todo     *entity.Todo
		mockFunc func()
		wantErr  error
	}{
		{
			name: "Createが成功した場合",
			todo: wantTodoData,
			mockFunc: func() {
				mockRepo.EXPECT().Create(ctx, wantTodoData).Return(nil)
			},
			wantErr: nil,
		},
		{
			name: "Create処理が失敗した場合",
			todo: wantTodoData,
			mockFunc: func() {
				mockRepo.EXPECT().Create(ctx, wantTodoData).Return(errors.New("failed to create"))
			},
			wantErr: errors.New("failed to create"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()
			err := svc.CreateTodo(ctx, tt.todo)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func Test_Service_UpdateTodo(t *testing.T) {
	ctx := context.Background()
	now := time.Now().UTC().Truncate(time.Millisecond)
	wantTodoData := &entity.Todo{ID: 1, Title: "title1", Content: "content1", CreatedAt: now}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mock_repository.NewMockTodoRepository(ctrl)
	svc := NewTodoService(mockRepo)

	tests := []struct {
		name     string
		todo     *entity.Todo
		mockFunc func()
		wantErr  error
	}{
		{
			name: "Updateが成功した場合",
			todo: wantTodoData,
			mockFunc: func() {
				mockRepo.EXPECT().Update(ctx, wantTodoData).Return(nil)
			},
			wantErr: nil,
		},
		{
			name: "Updateが失敗した場合",
			todo: wantTodoData,
			mockFunc: func() {
				mockRepo.EXPECT().Update(ctx, wantTodoData).Return(errors.New("failed to update"))
			},
			wantErr: errors.New("failed to update"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()
			err := svc.UpdateTodo(ctx, tt.todo)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func Test_Service_Delete(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mock_repository.NewMockTodoRepository(ctrl)
	svc := NewTodoService(mockRepo)

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
				mockRepo.EXPECT().Delete(ctx, uint(1)).Return(nil)
			},
			wantErr: nil,
		},
		{
			name: "Deleteが失敗した場合",
			id:   uint(1),
			mockFunc: func() {
				mockRepo.EXPECT().Delete(ctx, uint(1)).Return(errors.New("failed to delete"))
			},
			wantErr: errors.New("failed to delete"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()
			err := svc.DeleteTodo(ctx, tt.id)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
