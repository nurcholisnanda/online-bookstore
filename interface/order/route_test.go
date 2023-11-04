package order

import (
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func TestAddOrderRoutes(t *testing.T) {
	type args struct {
		rg *gin.RouterGroup
		db *gorm.DB
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "ïmplemented",
			args: args{
				rg: &gin.Default().RouterGroup,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			AddOrderRoutes(tt.args.rg, tt.args.db)
		})
	}
}
