package service

import (
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/nurcholisnanda/online-bookstore/application/dto"
	"github.com/nurcholisnanda/online-bookstore/domain/book"
	"github.com/nurcholisnanda/online-bookstore/domain/book/mock"
)

func TestNewBookService(t *testing.T) {
	type args struct {
		bookRepo book.Repository
	}
	ctrl := gomock.NewController(t)
	mockRepo := mock.NewMockRepository(ctrl)

	tests := []struct {
		name string
		args args
		want BookService
	}{
		{
			name: "implemented",
			args: args{
				bookRepo: mockRepo,
			},
			want: NewBookService(mockRepo),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewBookService(tt.args.bookRepo); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBookService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_bookService_GetBooks(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mock.NewMockRepository(ctrl)
	res := seedBookRes()
	books := seedBook()

	mockRepoRecord := func(res []book.Book, err error) func(m *mock.MockRepository) {
		return func(m *mock.MockRepository) {
			m.EXPECT().FindAllBook().Return(res, err)
		}
	}

	tests := []struct {
		name    string
		s       *bookService
		want    []dto.BookResponse
		wantErr bool
		custom  any
	}{
		{
			name: "success get books",
			s: &bookService{
				bookRepo: mockRepo,
			},
			want:    res,
			wantErr: false,
		},
		{
			name: "fail get books",
			s: &bookService{
				bookRepo: mockRepo,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if strings.Contains(tt.name, "success") {
				mockRepoRecord(books, nil)(mockRepo)
			} else {
				mockRepoRecord(nil, errors.New("any error"))(mockRepo)
			}
			got, err := tt.s.GetBooks()
			if (err != nil) != tt.wantErr {
				t.Errorf("bookService.GetBooks() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("bookService.GetBooks() = %v, want %v", got, tt.want)
			}
		})
	}
}
