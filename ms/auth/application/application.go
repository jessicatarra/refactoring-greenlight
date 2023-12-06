package application

import (
	"github.com/jessicatarra/greenlight/internal/concurrent"
	"github.com/jessicatarra/greenlight/internal/config"
	"github.com/jessicatarra/greenlight/internal/jsonlog"
	"github.com/jessicatarra/greenlight/internal/mailer"
	"github.com/jessicatarra/greenlight/ms/auth/domain"
	"github.com/jessicatarra/greenlight/ms/auth/repositories"
	"sync"
	"time"
)

type appl struct {
	userRepo  domain.UserRepository
	tokenRepo domain.TokenRepository
	helpers   concurrent.Resource
	logger    *jsonlog.Logger
	wg        *sync.WaitGroup
	mailer    mailer.Mailer
}

func NewAppl(userRepo domain.UserRepository, tokenRepo domain.TokenRepository, logger *jsonlog.Logger,
	wg *sync.WaitGroup, cfg config.Config) domain.Appl {
	return &appl{
		userRepo:  userRepo,
		tokenRepo: tokenRepo,
		helpers:   concurrent.NewBackgroundTask(wg, logger),
		logger:    logger,
		wg:        wg,
		mailer:    mailer.New(cfg.Smtp.Host, cfg.Smtp.Port, cfg.Smtp.Username, cfg.Smtp.Password, cfg.Smtp.Sender),
	}
}

func (a *appl) CreateUseCase(input domain.CreateUserRequest) (*domain.User, error) {
	user := &domain.User{Name: input.Name, Email: input.Email, Activated: false}

	err := a.userRepo.InsertNewUser(user)

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

func (a *appl) ActivateUseCase(tokenPlainText string) (*domain.User, error) {
	//TODO implement me
	panic("implement me")
}

func (a *appl) GetByEmailUseCase(input domain.CreateUserRequest) (*domain.User, error) {
	email := input.Email

	existingUser, err := a.userRepo.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}

	return existingUser, nil
}
