package application

import (
	"github.com/jessicatarra/greenlight/internal/concurrent"
	"github.com/jessicatarra/greenlight/internal/config"
	"github.com/jessicatarra/greenlight/internal/mailer"
	"github.com/jessicatarra/greenlight/ms/auth/domain"
	"github.com/jessicatarra/greenlight/ms/auth/repositories"
	"github.com/pascaldekloe/jwt"
	"strconv"
	"sync"
	"time"
)

type appl struct {
	userRepo       domain.UserRepository
	tokenRepo      domain.TokenRepository
	permissionRepo domain.PermissionRepository
	concurrent     concurrent.Resource
	mailer         mailer.Mailer
	cfg            config.Config
}

func NewAppl(userRepo domain.UserRepository, tokenRepo domain.TokenRepository, permissionRepo domain.PermissionRepository, cfg config.Config, wg *sync.WaitGroup) domain.Appl {
	return &appl{
		userRepo:       userRepo,
		tokenRepo:      tokenRepo,
		permissionRepo: permissionRepo,
		concurrent:     concurrent.NewBackgroundTask(wg),
		mailer:         mailer.New(cfg.Smtp.Host, cfg.Smtp.Port, cfg.Smtp.Username, cfg.Smtp.Password, cfg.Smtp.From),
		cfg:            cfg,
	}
}

func (a *appl) CreateUseCase(input domain.CreateUserRequest, hashedPassword string) (*domain.User, error) {
	user := &domain.User{Name: input.Name, Email: input.Email, Activated: false}

	err := a.userRepo.InsertNewUser(user, hashedPassword)

	if err != nil {
		return nil, err
	}

	err = a.permissionRepo.AddForUser(user.ID, "movies:read")
	if err != nil {
		return nil, err
	}

	token, err := a.tokenRepo.New(user.ID, 3*24*time.Hour, repositories.ScopeActivation)
	if err != nil {
		return nil, err
	}

	fn := func() error {
		data := map[string]interface{}{
			"activationToken": token.Plaintext,
			"userID":          user.ID,
		}
		//print(token.Plaintext)

		err = a.mailer.Send(user.Email, "user_welcome.gohtml", data)
		if err != nil {
			return err
		}
		return nil
	}

	a.concurrent.BackgroundTask(fn)

	return user, err
}

func (a *appl) ActivateUseCase(tokenPlainText string) (*domain.User, error) {
	user, err := a.userRepo.GetForToken(repositories.ScopeActivation, tokenPlainText)
	if err != nil {
		return nil, err
	}

	user.Activated = true

	err = a.userRepo.UpdateUser(user)
	if err != nil {
		return nil, err
	}

	err = a.tokenRepo.DeleteAllForUser(repositories.ScopeActivation, user.ID)
	if err != nil {
		return nil, err
	}

	return user, err
}

func (a *appl) GetByEmailUseCase(email string) (*domain.User, error) {
	existingUser, err := a.userRepo.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}

	return existingUser, nil
}

func (a *appl) CreateAuthTokenUseCase(userID int64) ([]byte, error) {
	var claims jwt.Claims
	claims.Subject = strconv.FormatInt(userID, 10)
	claims.Issued = jwt.NewNumericTime(time.Now())
	claims.NotBefore = jwt.NewNumericTime(time.Now())
	claims.Expires = jwt.NewNumericTime(time.Now().Add(24 * time.Hour))
	claims.Issuer = a.cfg.BaseURL
	claims.Audiences = []string{a.cfg.BaseURL}

	jwtBytes, err := claims.HMACSign(jwt.HS256, []byte(a.cfg.Jwt.Secret))
	if err != nil {
		return nil, err
	}

	return jwtBytes, nil
}
