package application

import (
	"github.com/jessicatarra/greenlight/internal/config"
	"github.com/jessicatarra/greenlight/internal/helpers"
	"github.com/jessicatarra/greenlight/internal/jsonlog"
	"github.com/jessicatarra/greenlight/internal/mailer"
	"github.com/jessicatarra/greenlight/internal/validator"
	"github.com/jessicatarra/greenlight/ms/auth/entity"
	"github.com/jessicatarra/greenlight/ms/auth/repositories"
	"sync"
	"time"
)

func ValidateUser(v *validator.Validator, user *entity.User) {
	v.Check(user.Name != "", "name", "must be provided")
	v.Check(len(user.Name) <= 500, "name", "must not be more than 500 bytes long")

	ValidateEmail(v, user.Email)

	if user.Password.Plaintext != nil {
		ValidatePasswordPlaintext(v, *user.Password.Plaintext)
	}

	if user.Password.Hash == nil {
		panic("missing password hash for user")
	}
}

func ValidateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "must be provided")
	v.Check(validator.Matches(email, validator.EmailRX), "email", "must be a valid email address")
}

func ValidatePasswordPlaintext(v *validator.Validator, password string) {
	v.Check(password != "", "password", "must be provided")
	v.Check(len(password) >= 8, "password", "must be at least 8 bytes long")
	v.Check(len(password) <= 72, "password", "must not be more than 72 bytes long")
}

type Appl interface {
	CreateUseCase(input entity.CreateUserRequest) (*entity.User, error)
}

type appl struct {
	userRepo  repositories.UserRepository
	tokenRepo repositories.TokenRepository
	helpers   helpers.Resource
	logger    *jsonlog.Logger
	wg        *sync.WaitGroup
	mailer    mailer.Mailer
}

func NewAppl(userRepo repositories.UserRepository, tokenRepo repositories.TokenRepository, logger *jsonlog.Logger,
	wg *sync.WaitGroup, cfg config.Config) Appl {
	return &appl{
		userRepo:  userRepo,
		tokenRepo: tokenRepo,
		helpers:   helpers.NewBackgroundTask(wg, logger),
		logger:    logger,
		wg:        wg,
		mailer:    mailer.New(cfg.Smtp.Host, cfg.Smtp.Port, cfg.Smtp.Username, cfg.Smtp.Password, cfg.Smtp.Sender),
	}
}

func (a *appl) CreateUseCase(input entity.CreateUserRequest) (*entity.User, error) {
	user := &entity.User{Name: input.Name, Email: input.Email, Activated: false}

	err := user.Password.Set(input.Password)
	if err != nil {
		return nil, err
	}

	v := validator.New()

	ValidateUser(v, user)

	if !v.Valid() {
		var err error
		return nil, err
	}

	err = a.userRepo.InsertNewUser(user)

	if err != nil {
		return nil, err
	}

	token, err := a.tokenRepo.New(user.ID, 3*24*time.Hour, repositories.ScopeActivation)
	if err != nil {
		return nil, err
	}

	fn := func() {
		data := map[string]interface{}{
			"activationToken": token.Plaintext,
			"userID":          user.ID,
		}
		print(token.Plaintext)

		err = a.mailer.Send(user.Email, "user_welcome.gohtml", data)
		if err != nil {
			a.logger.PrintError(err, nil)
		}
	}

	a.helpers.Background(fn)

	return user, err
}
