package db

import (
	"bookstore/model"
	"bookstore/pb"
	"context"
	"errors"
	"time"

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
	result := DB.WithContext(ctx).Delete(&model.Shelf{}, shelfID)
	return result.Error
}

// 查询某个书架上的书的信息（可选分页)
func (b *BookStore) ListBooks(ctx context.Context, shelfID int64, nextID string, pageSize int64) ([]*model.Book, error) {
	// 查询数据库中是否有这个书架
	// 如果没有则返回错误
	if _, err := b.GetShelf(ctx, shelfID); errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, gorm.ErrRecordNotFound
	}
	// 条件查询
	var books []*model.Book
	res := DB.WithContext(ctx).Where("shelf_id = ? and id > ? ", shelfID, nextID).Order("id asc").Limit(int(pageSize)).Find(&books)
	return books, res.Error
}

// 创建一本书，加入书架
func (b *BookStore) CreateBook(ctx context.Context, req *pb.CreateBookRequest) (*model.Book, error) {
	// 查询数据库中是否有这个书架
	// 如果没有则返回错误
	if _, err := b.GetShelf(ctx, req.GetShelf()); errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, gorm.ErrRecordNotFound
	}
	// 查询这本书是否在书架
	if !errors.Is(DB.WithContext(ctx).First(&model.Book{ID: req.GetBook().Id}).Error, gorm.ErrRecordNotFound) {
		return nil, gorm.ErrRegistered
	}
	// 加入书架
	book := &model.Book{
		ID:        req.Book.GetId(),
		Author:    req.Book.GetAuthor(),
		Title:     req.Book.GetTitle(),
		ShelfID:   req.Book.GetShelfId(),
		CreateAt:  time.Now(),
		UpdatedAt: time.Now(),
	}
	result := DB.WithContext(ctx).Create(book)
	return book, result.Error
}
