package user

// Interface that contains repository contract for user domain
//go:generate mockgen -source=repository.go -destination=mock/user.go -package=mock
type Repository interface {
	InsertUser(*User) error
	FindUserByID(uint) (*User, error)
	FindUserByEmail(string) (*User, error)
}
