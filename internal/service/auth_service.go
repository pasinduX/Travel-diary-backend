package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"travel-diary-backend/internal/config"
	"travel-diary-backend/internal/dao"
	"travel-diary-backend/internal/dto"
	"travel-diary-backend/internal/integrations"
	"travel-diary-backend/internal/models"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
)

type AuthService struct {
	cfg      config.Config
	users    dao.UserDAO
	sessions dao.SessionDAO
	oauth    *oauth2.Config
}

func NewAuthService(cfg config.Config, users dao.UserDAO, sessions dao.SessionDAO) *AuthService {
	return &AuthService{
		cfg:      cfg,
		users:    users,
		sessions: sessions,
		oauth:    integrations.GoogleOAuthConfig(cfg.GoogleClientID, cfg.GoogleClientSecret, cfg.GoogleRedirectURL),
	}
}

func (s *AuthService) Register(ctx context.Context, req dto.RegisterRequest) (dto.AuthResponse, error) {
	req.Username = strings.TrimSpace(req.Username)
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	if req.Username == "" || req.Email == "" || req.Password == "" {
		return dto.AuthResponse{}, errors.New("username, email, and password are required")
	}

	if _, err := s.users.FindByUsername(ctx, req.Username); err == nil {
		return dto.AuthResponse{}, errors.New("username already exists")
	} else if !errors.Is(err, dao.ErrUserNotFound) {
		return dto.AuthResponse{}, err
	}
	if _, err := s.users.FindByEmail(ctx, req.Email); err == nil {
		return dto.AuthResponse{}, errors.New("email already exists")
	} else if !errors.Is(err, dao.ErrUserNotFound) {
		return dto.AuthResponse{}, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return dto.AuthResponse{}, err
	}

	user := models.User{
		ID:           uuid.NewString(),
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hash),
		AuthProvider: "local",
		Name:         req.Name,
	}
	user, err = s.users.Create(ctx, user)
	if err != nil {
		return dto.AuthResponse{}, err
	}

	token, refreshToken, err := s.issueTokens(ctx, user)
	if err != nil {
		return dto.AuthResponse{}, err
	}

	return dto.AuthResponse{Token: token, RefreshToken: refreshToken, User: toUserResponse(user)}, nil
}

func (s *AuthService) Login(ctx context.Context, req dto.LoginRequest) (dto.AuthResponse, error) {
	req.Username = strings.TrimSpace(req.Username)
	if req.Username == "" || req.Password == "" {
		return dto.AuthResponse{}, errors.New("username and password are required")
	}

	user, err := s.users.FindByUsername(ctx, req.Username)
	if err != nil {
		return dto.AuthResponse{}, errors.New("invalid username or password")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return dto.AuthResponse{}, errors.New("invalid username or password")
	}

	token, refreshToken, err := s.issueTokens(ctx, user)
	if err != nil {
		return dto.AuthResponse{}, err
	}
	return dto.AuthResponse{Token: token, RefreshToken: refreshToken, User: toUserResponse(user)}, nil
}

func (s *AuthService) GoogleLoginURL() string {
	return s.oauth.AuthCodeURL(s.cfg.GoogleAuthStateSecret)
}

func (s *AuthService) GoogleCallback(ctx context.Context, code, state string) (dto.AuthResponse, error) {
	if state != s.cfg.GoogleAuthStateSecret {
		return dto.AuthResponse{}, errors.New("invalid oauth state")
	}

	token, err := s.oauth.Exchange(ctx, code)
	if err != nil {
		return dto.AuthResponse{}, err
	}

	info, err := integrations.FetchGoogleUserInfo(ctx, token)
	if err != nil {
		return dto.AuthResponse{}, err
	}

	username := strings.ToLower(strings.ReplaceAll(info.Email, "@", "_"))
	user := models.User{
		ID:           uuid.NewString(),
		Username:     username,
		Email:        strings.ToLower(info.Email),
		PasswordHash: "",
		GoogleID:     info.Sub,
		AuthProvider: "google",
		Name:         info.Name,
		PictureURL:   info.Picture,
	}

	user, err = s.users.UpsertGoogleUser(ctx, user)
	if err != nil {
		return dto.AuthResponse{}, err
	}

	appToken, refreshToken, err := s.issueTokens(ctx, user)
	if err != nil {
		return dto.AuthResponse{}, err
	}

	return dto.AuthResponse{Token: appToken, RefreshToken: refreshToken, User: toUserResponse(user)}, nil
}

func (s *AuthService) Refresh(ctx context.Context, refreshToken string) (dto.AuthResponse, error) {
	_, claims, err := integrations.ParseRefreshToken(s.cfg.JWTSecret, refreshToken)
	if err != nil {
		return dto.AuthResponse{}, errors.New("invalid refresh token")
	}

	if tokenType, _ := claims["type"].(string); tokenType != "refresh" {
		return dto.AuthResponse{}, errors.New("invalid refresh token")
	}

	tokenID, _ := claims["jti"].(string)
	userID, _ := claims["sub"].(string)
	if tokenID == "" || userID == "" {
		return dto.AuthResponse{}, errors.New("invalid refresh token")
	}

	session, err := s.sessions.FindActiveByTokenID(ctx, tokenID)
	if err != nil {
		return dto.AuthResponse{}, errors.New("invalid refresh token")
	}
	if session.TokenHash != integrations.HashToken(refreshToken) {
		return dto.AuthResponse{}, errors.New("invalid refresh token")
	}

	user, err := s.users.FindByID(ctx, userID)
	if err != nil {
		return dto.AuthResponse{}, errors.New("invalid refresh token")
	}

	if err := s.sessions.RevokeByTokenID(ctx, tokenID); err != nil {
		return dto.AuthResponse{}, err
	}

	accessToken, newRefreshToken, err := s.issueTokens(ctx, user)
	if err != nil {
		return dto.AuthResponse{}, err
	}

	return dto.AuthResponse{Token: accessToken, RefreshToken: newRefreshToken, User: toUserResponse(user)}, nil
}

func (s *AuthService) issueTokens(ctx context.Context, user models.User) (string, string, error) {
	accessToken, err := integrations.CreateAccessToken(s.cfg.JWTSecret, user.ID, user.Username, user.Email)
	if err != nil {
		return "", "", err
	}
	refreshToken, tokenID, err := integrations.CreateRefreshToken(s.cfg.JWTSecret, user.ID)
	if err != nil {
		return "", "", err
	}
	if _, err := s.sessions.Create(ctx, models.RefreshSession{
		ID:        tokenID,
		UserID:    user.ID,
		TokenID:   tokenID,
		TokenHash: integrations.HashToken(refreshToken),
		ExpiresAt: time.Now().UTC().Add(integrations.RefreshTokenTTL),
	}); err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
}

func toUserResponse(u models.User) dto.UserResponse {
	return dto.UserResponse{
		ID:       u.ID,
		Username: u.Username,
		Email:    u.Email,
		Name:     u.Name,
		Provider: u.AuthProvider,
	}
}
