//go:generate mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=../../mock/$GOPACKAGE/$GOFILE
package repository

import (
	"context"
	"local-golang-prac/domain/entity"
)

type TodoRepository interface {
	FindByID(ctx context.Context, id uint) ([]*entity.Todo, error)
	Create(ctx context.Context, todo *entity.Todo) error
	Update(ctx context.Context, todo *entity.Todo) error
	Delete(ctx context.Context, id uint) error
}
