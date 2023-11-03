package repository

import (
	"reflect"
	"testing"

	"github.com/nurcholisnanda/online-bookstore/domain/book"
	"gorm.io/gorm"
)

func TestNewBookRepositoryImpl(t *testing.T) {
	type args struct {
		db *gorm.DB
	}
	db, err := DBConn()
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}
	tests := []struct {
		name string
		args args
		want book.Repository
	}{
		{
			name: "implemented",
			args: args{
				db: db,
			},
			want: NewBookRepositoryImpl(db),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewBookRepositoryImpl(tt.args.db); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBookRepositoryImpl() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_bookRepositoryImpl_FindAllBook(t *testing.T) {
	db, err := DBConn()
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}
	books := seedBook(db)

	tests := []struct {
		name      string
		r         *bookRepositoryImpl
		wantBooks []book.Book
		wantErr   bool
	}{
		{
			name: "success fetching all books",
			r: &bookRepositoryImpl{
				db: db,
			},
			wantBooks: books,
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBooks, err := tt.r.FindAllBook()
			if (err != nil) != tt.wantErr {
				t.Errorf("bookRepositoryImpl.FindAllBook() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotBooks, tt.wantBooks) {
				t.Errorf("bookRepositoryImpl.FindAllBook() = %v, want %v", gotBooks, tt.wantBooks)
			}
		})
	}
}

func Test_bookRepositoryImpl_GetBooksByIDs(t *testing.T) {
	type args struct {
		ids []uint
	}
	db, err := DBConn()
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}
	books := seedBook(db)

	tests := []struct {
		name      string
		r         *bookRepositoryImpl
		args      args
		wantBooks []book.Book
		wantErr   bool
	}{
		{
			name: "success fetching books by ids",
			r: &bookRepositoryImpl{
				db: db,
			},
			args: args{
				ids: []uint{1, 2},
			},
			wantBooks: books,
			wantErr:   false,
		},
		{
			name: "fail fetching book",
			r: &bookRepositoryImpl{
				db: db,
			},
			args: args{
				ids: []uint{},
			},
			wantBooks: nil,
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBooks, err := tt.r.GetBooksByIDs(tt.args.ids)
			if (err != nil) != tt.wantErr {
				t.Errorf("bookRepositoryImpl.GetBooksByIDs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotBooks, tt.wantBooks) {
				t.Errorf("bookRepositoryImpl.GetBooksByIDs() = %v, want %v", gotBooks, tt.wantBooks)
			}
		})
	}
}
