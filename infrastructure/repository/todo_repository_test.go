package infrastructure

import (
	"context"
	"errors"
	"local-golang-prac/domain/entity"
	"local-golang-prac/testutils"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func Test_NewTodoRepository(t *testing.T) {
	mockdb, _ := testutils.Mock(t)
	repo := NewTodoRepository(mockdb)
	assert.Equal(t, mockdb, repo.db)
}

func Test_TodoRepository_FindByID(t *testing.T) {
	ctx := context.Background()
	now := time.Now().UTC().Truncate(time.Millisecond)
	tests := []struct {
		name     string
		id       uint
		mockFunc func(mock sqlmock.Sqlmock)
		want     []*entity.Todo
		wantErr  error
	}{
		{
			name: "IDでTodoが見つかる",
			id:   1,
			mockFunc: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "title", "content", "created_at"})
				rows.AddRow(1, "title1", "content1", now)
				rows.AddRow(2, "title2", "content2", now)
				rows.AddRow(3, "title3", "content3", now)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "todos" WHERE id = $1`)).WithArgs(1).WillReturnRows(rows)
			},
			want: []*entity.Todo{
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
			},
			wantErr: nil,
		},
		{
			name: "UserIDで対象のTodoが見つからない",
			id:   10,
			mockFunc: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "title", "content", "created_at"})
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "todos" WHERE id = $1`)).WithArgs(10).WillReturnRows(rows)
			},
			want:    []*entity.Todo{},
			wantErr: nil,
		},
		{
			name: "DBエラーが発生する",
			id:   1,
			mockFunc: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(
					regexp.QuoteMeta(`SELECT * FROM "todos" WHERE id = $1`),
				).WithArgs(1).WillReturnError(errors.New("db error"))
			},
			want:    nil,
			wantErr: errors.New("db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockdb, gormMock := testutils.Mock(t)
			repo := NewTodoRepository(mockdb)
			tt.mockFunc(gormMock)
			got, err := repo.FindByID(ctx, tt.id)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func Test_TodoRepository_Create(t *testing.T) {
	ctx := context.Background()
	now := time.Now().UTC().Truncate(time.Millisecond)
	wantTodoData := &entity.Todo{ID: 1, Title: "title1", Content: "content1", CreatedAt: now}
	tests := []struct {
		name     string
		todo     *entity.Todo
		mockFunc func(mock sqlmock.Sqlmock)
		wantErr  error
	}{
		{
			name: "Createが成功した場合",
			todo: wantTodoData,
			mockFunc: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "todos"`)).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				mock.ExpectCommit()
			},
			wantErr: nil,
		},
		{
			name: "Createが失敗した場合",
			todo: wantTodoData,
			mockFunc: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "todos"`)).
					WillReturnError(errors.New("failed to create"))
				mock.ExpectRollback()
			},
			wantErr: errors.New("failed to create"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockdb, gormMock := testutils.Mock(t)
			repo := NewTodoRepository(mockdb)
			tt.mockFunc(gormMock)
			err := repo.Create(ctx, tt.todo)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func Test_TodoRepository_Update(t *testing.T) {
	ctx := context.Background()
	now := time.Now().UTC().Truncate(time.Millisecond)
	wantTodoData := &entity.Todo{ID: 1, Title: "title1", Content: "content1", CreatedAt: now, UpdatedAt: now}
	tests := []struct {
		name     string
		todo     *entity.Todo
		mockFunc func(mock sqlmock.Sqlmock)
		wantErr  error
	}{
		{
			name: "Updateが成功した場合",
			todo: wantTodoData,
			mockFunc: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "todos"`)).
					WithArgs("title1",
						"content1",
						now,
						sqlmock.AnyArg(),
						1).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
			},
			wantErr: nil,
		},
		{
			name: "Updateが失敗した場合",
			todo: wantTodoData,
			mockFunc: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "todos"`)).
					WillReturnError(errors.New("failed to create"))
				mock.ExpectRollback()
			},
			wantErr: errors.New("failed to create"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockdb, gormMock := testutils.Mock(t)
			repo := NewTodoRepository(mockdb)
			tt.mockFunc(gormMock)
			err := repo.Update(ctx, tt.todo)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func Test_TodoRepository_Delete(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name     string
		id       uint
		mockFunc func(mock sqlmock.Sqlmock)
		wantErr  error
	}{
		{
			name: "Deleteが成功した場合",
			id:   1,
			mockFunc: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "todos" WHERE "todos"."id" = $1`)).
					WithArgs(1).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
			},
			wantErr: nil,
		},
		{
			name: "Updateが失敗した場合",
			id:   1,
			mockFunc: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "todos" WHERE "todos"."id" = $1`)).
					WithArgs(1).
					WillReturnError(errors.New("failed to create"))
				mock.ExpectRollback()
			},
			wantErr: errors.New("failed to create"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockdb, gormMock := testutils.Mock(t)
			repo := NewTodoRepository(mockdb)
			tt.mockFunc(gormMock)
			err := repo.Delete(ctx, tt.id)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
