package service

import (
	"auth/internal/models"
	"auth/internal/repository"
	"log/slog"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type TokenManager struct {
	repo       *repository.Repository
	jwtSecret  string
	accessTTL  time.Duration
	refreshTTL time.Duration
}

func NewTokenManager(repo *repository.Repository) (*TokenManager, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	accessTTL, err := time.ParseDuration(os.Getenv("ACCESS_TTL"))
	if err != nil {
		slog.Debug("Failed to parse ACCESS_TTL")
		return nil, err
	}

	refreshTTL, err := time.ParseDuration(os.Getenv("REFRESH_TTL"))
	if err != nil {
		slog.Debug("Failed to parse REFRESH_TTL")
		return nil, err
	}

	TokenManager := &TokenManager{
		repo:       repo,
		jwtSecret:  jwtSecret,
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL}
	return TokenManager, nil
}

func (s *TokenManager) GetToken(guid string, ip string) (*models.JWT, error) {
	JWT, err := s.GenerateToken(guid, ip)
	if err != nil {
		slog.Debug("Failed to generate JWT")
		return nil, err
	}

	hashedRefreshToken, err := bcrypt.GenerateFromPassword([]byte(JWT.RefreshToken), bcrypt.DefaultCost)
	if err != nil {
		slog.Debug("Failed to hash token")
		return nil, err
	}

	tx, err := s.repo.Storage.DB.Begin()
	if err != nil {
		slog.Debug("Failed to start transaction")
		return nil, err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	err = s.repo.UpdateToken(tx, guid, string(hashedRefreshToken))
	if err != nil {
		slog.Debug("Failed to update token in database")
		slog.Error(err.Error())
		return nil, err
	}

	return JWT, nil
}

func (s *TokenManager) CheckToken(refreshToken string, ip string) (*models.JWT, error) {
	parsedToken, err := jwt.ParseWithClaims(refreshToken, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.jwtSecret), nil
	})
	if err != nil {
		slog.Debug("Failed to parse refresh token")
		return nil, err
	}

	if claims, ok := parsedToken.Claims.(*models.Claims); ok && parsedToken.Valid {
		if claims.IP != ip {
			err := s.Warning(claims.GUID)
			if err != nil {
				slog.Debug("Warning failed")
				return nil, err
			}
		}
	} else {
		slog.Debug("Not Valid Token")
		return nil, err
	}

	hashedRefreshToken, err := bcrypt.GenerateFromPassword([]byte(refreshToken), bcrypt.DefaultCost)
	if err != nil {
		slog.Debug("Failed to hash token")
		return nil, err
	}

	userGUID, err := s.repo.GetGUID(string(hashedRefreshToken))
	if err != nil {
		slog.Debug("Failed to get user ID by refresh token")
		return nil, err
	}

	if userGUID == "" {
		slog.Debug("Not correct UserID")
		return nil, err
	}
	return s.GetToken(userGUID, ip)
}

func (s *TokenManager) GenerateToken(guid string, ip string) (*models.JWT, error) {
	claims := &models.Claims{
		GUID: guid,
		IP:   ip,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.accessTTL)),
		},
	}

	access := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(s.refreshTTL))
	refresh := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	accessToken, err := access.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return nil, err
	}
	refreshToken, err := refresh.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return nil, err
	}

	JWT := &models.JWT{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	return JWT, nil
}

func (s *TokenManager) Warning(guid string) error {
	email, err := s.repo.GetEmail(guid)
	if err != nil {
		slog.Debug("Failed to get Email")
		return err
	}
	slog.Info("Warning on Email", "email", email)
	return nil
}

func (s *TokenManager) InsertData(users []models.User) ([]models.User, error) {
	for _, user := range users {
		user.GUID = uuid.New().String()
		jwt, err := s.GenerateToken(user.GUID, user.IP)
		if err != nil {
			slog.Debug("Failed to generate Token")
			return nil, err
		}
		HashedRefreshToken, err := bcrypt.GenerateFromPassword([]byte(jwt.RefreshToken), bcrypt.DefaultCost)
		if err != nil {
			slog.Debug("Failed to hash RefreshToken")
			return nil, err
		}
		err = s.repo.InsertGUID(user.GUID, user.Email, string(HashedRefreshToken))
		if err != nil {
			slog.Debug("Failed to add GUID")
			return nil, err
		}
	}
	return  users, nil
}
