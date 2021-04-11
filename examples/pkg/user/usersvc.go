package user

import (
	"easyquery"
	"easyquery/examples/pkg/db"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"golang.org/x/crypto/bcrypt"
)

var UserCrudService *UserService

type UserService struct {
	easyquery.Crud
}

func init() {
	UserCrudService = NewDefaultUserService()
}

func NewUserService(crud easyquery.Crud) *UserService {
	return &UserService{crud}
}

func NewDefaultUserService() *UserService {
	return NewUserService(easyquery.NewCrudService(Model, Models, retreiveGormDB))
}

func retreiveGormDB() *gorm.DB {
	return db.Postgres
}

func CurrentUser(ctx *gin.Context) (*User, bool) {
	var model *User
	current, ok := ctx.Get("user")
	if !ok {
		return model, false
	}
	model, ok = current.(*User)
	if !ok {
		return model, false
	}
	return model, true
}

func (service *UserService) CheckUserPwd(user *User, password *string) (*User, bool) {
	if pErr := bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(*password)); pErr != nil {
		return user, false
	}
	return user, true
}

func (service *UserService) EncryptPassword(user *User) (*User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(*user.Password), bcrypt.DefaultCost)
	if err != nil {
		return user, err
	}
	pwd := string(hash)
	user.Password = &pwd
	return user, nil
}
func (service *UserService) FindByLoginName(user *User) (*User, error) {
	var model User
	currentScope := retreiveGormDB().Where("login_name = ?", user.LoginName).First(&model)
	return &model, currentScope.Error
}

func (service *UserService) FindByToken(token string) (*User, error) {
	var model User
	currentScope := retreiveGormDB().Where("password = ?", token).First(&model)
	return &model, currentScope.Error
}
