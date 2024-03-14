package services

import (
	"context"
	"errors"
	"eshell/dao/entities"
	"eshell/dao/repo"
	"eshell/ioc"
	"eshell/services/models"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/universalmacro/common/fault"
	"github.com/universalmacro/common/singleton"
	"github.com/universalmacro/common/snowflake"
	"github.com/universalmacro/common/utils"
)

func GetAccountService() *AccountService {
	return accountServiceSingleton.Get()
}

var accountServiceSingleton = singleton.SingletonFactory(newAccountService, singleton.Eager)

func newAccountService() *AccountService {
	return &AccountService{accountRepo: repo.GetAccountRepo()}
}

type AccountService struct {
	accountRepo *repo.AccountRepository
}

func (s *AccountService) GetById(id uint) *models.Account {
	account, _ := s.accountRepo.GetById(id)
	return &models.Account{Account: account}
}

func (s *AccountService) GetByAccount(account string) *models.Account {
	ac, _ := s.accountRepo.FindOne("account = ?", account)
	return &models.Account{Account: ac}
}

func (s *AccountService) Create(account, password string) *models.Account {
	ac := &entities.Account{Account: account}
	ac.SetPassword(password)
	_, db := s.accountRepo.Create(ac)
	if db.Error != nil {
		return nil
	}
	return &models.Account{Account: ac}
}

func (s *AccountService) Login(account, password string) string {

	return ""
}

type Claims struct {
	jwt.StandardClaims
	ID        string `json:"id"`
	AccountID string `json:"accountId"`
	Account   string `json:"account"`
}

var sessionIdGenerator = snowflake.NewIdGenertor(0)

func (s *AccountService) CreateSession(account, password string) (string, error) {
	admin, _ := s.accountRepo.FindOne("account = ?", account)
	if admin == nil {
		return "", fault.ErrNotFound
	}
	if !admin.PasswordMatching(password) {
		return "", errors.New("password not match")
	}
	expired := time.Now().Add(time.Hour * 24 * 7).Unix()
	claims := Claims{
		ID:             sessionIdGenerator.String(),
		AccountID:      utils.UintToString(admin.ID),
		StandardClaims: jwt.StandardClaims{ExpiresAt: expired}}
	return ioc.GetJwtSigner().SignJwt(claims)
}

func (s *AccountService) VerifyToken(ctx context.Context, token string) (*models.Account, error) {
	claims, err := ioc.GetJwtSigner().VerifyJwt(token)
	if err != nil {
		return nil, err
	}
	account := s.GetById(utils.StringToUint(claims["accountId"].(string)))
	if account == nil {
		return nil, fault.ErrNotFound
	}
	return account, nil
}
