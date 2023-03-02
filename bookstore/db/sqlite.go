package db

import (
	"bookstore/model"
	"context"
	"errors"

	"google.golang.org/genproto/googleapis/rpc/status"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func init() {
	// 连接到数据库
	db, err := gorm.Open(sqlite.Open("bookstore.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	// 迁移模型
	models := []interface{}{
		&model.Shelf{},
		&model.Book{},
	}
	db.AutoMigrate(models...)
	DB = db
}

type BookStore struct{}

// 查询所有书架信息
func (b *BookStore) ListShelves(ctx context.Context) ([]*model.Shelf, error) {
	var shelves []*model.Shelf
	result := DB.WithContext(ctx).Find(&shelves)
	return shelves, result.Error
}

// 创建一个书架
func (b *BookStore) CreateShelf(ctx context.Context, shelf *model.Shelf) (*model.Shelf, error) {
	if len(shelf.Theme) <= 0 || shelf.Size <= 0 {
		return nil, errors.New("invalid args")
	}
	result := DB.WithContext(ctx).Omit("DeletedAt").Create(&shelf)
	return shelf, result.Error
}

// 获取一个书架
func (b *BookStore) GetShelf(ctx context.Context, shelfID int64) (*model.Shelf, error) {
	if shelfID <= 0 {
		return nil, gorm.ErrInvalidValue
	}
	shelf := &model.Shelf{
		ID: shelfID,
	}
	result := DB.WithContext(ctx).First(shelf)
	return shelf, result.Error
}

// 删除一个书架
func (b *BookStore) DeleteShelf(ctx context.Context, shelfID int64) error {
	result := DB.Delete(&model.Shelf{}, shelfID)
	return result.Error
}
