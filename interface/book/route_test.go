package book

import (
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func TestAddBookRoutes(t *testing.T) {
	type args struct {
		rg *gin.RouterGroup
		db *gorm.DB
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "integrated",
			args: args{
				rg: &gin.Default().RouterGroup,
				db: &gorm.DB{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			AddBookRoutes(tt.args.rg, tt.args.db)
		})
	}
}
