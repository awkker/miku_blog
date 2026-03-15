package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"nanamiku-blog/backend/query"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrTokenExpired       = errors.New("token expired or revoked")
	ErrTokenInvalid       = errors.New("invalid token")
)

type JWTConfig struct {
	Secret     string
	AccessTTL  time.Duration
	RefreshTTL time.Duration
}

type AuthService struct {
	q   *query.Queries
	db  *pgxpool.Pool
	cfg JWTConfig
}

func NewAuthService(db *pgxpool.Pool, cfg JWTConfig) *AuthService {
	return &AuthService{
		q:   query.New(db),
		db:  db,
		cfg: cfg,
	}
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    int64  `json:"expires_at"`
}

type AdminClaims struct {
	jwt.RegisteredClaims
	AdminID  uuid.UUID `json:"admin_id"`
	Username string    `json:"username"`
	Role     string    `json:"role"`
}

func (s *AuthService) Login(ctx context.Context, username, password string) (*TokenPair, error) {
	admin, err := s.q.GetAdminByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrInvalidCredentials
		}
		return nil, fmt.Errorf("query admin: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	pair, err := s.generateTokenPair(ctx, admin.ID, admin.Username, admin.Role)
	if err != nil {
		return nil, err
	}

	_ = s.q.UpdateAdminLastLogin(ctx, admin.ID)
	return pair, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, rawRefresh string) (*TokenPair, error) {
	hash := hashToken(rawRefresh)

	rt, err := s.q.GetRefreshTokenByHash(ctx, hash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrTokenExpired
		}
		return nil, fmt.Errorf("query refresh token: %w", err)
	}

	admin, err := s.q.GetAdminByID(ctx, rt.AdminUserID)
	if err != nil {
		return nil, fmt.Errorf("query admin: %w", err)
	}

	_ = s.q.RevokeRefreshToken(ctx, rt.ID)

	pair, err := s.generateTokenPair(ctx, admin.ID, admin.Username, admin.Role)
	if err != nil {
		return nil, err
	}

	return pair, nil
}

func (s *AuthService) Logout(ctx context.Context, rawRefresh string) error {
	if rawRefresh == "" {
		return nil
	}
	hash := hashToken(rawRefresh)
	rt, err := s.q.GetRefreshTokenByHash(ctx, hash)
	if err != nil {
		return nil
	}
	return s.q.RevokeRefreshToken(ctx, rt.ID)
}

func (s *AuthService) LogoutAll(ctx context.Context, adminID uuid.UUID) error {
	return s.q.RevokeAllUserTokens(ctx, adminID)
}

func (s *AuthService) GetAdminInfo(ctx context.Context, adminID uuid.UUID) (*query.GetAdminByIDRow, error) {
	admin, err := s.q.GetAdminByID(ctx, adminID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}
	return &admin, nil
}

func (s *AuthService) ValidateAccessToken(tokenStr string) (*AdminClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &AdminClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(s.cfg.Secret), nil
	})
	if err != nil {
		return nil, ErrTokenInvalid
	}

	claims, ok := token.Claims.(*AdminClaims)
	if !ok || !token.Valid {
		return nil, ErrTokenInvalid
	}

	return claims, nil
}

func (s *AuthService) generateTokenPair(ctx context.Context, adminID uuid.UUID, username, role string) (*TokenPair, error) {
	now := time.Now()
	accessExp := now.Add(s.cfg.AccessTTL)

	claims := &AdminClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessExp),
			IssuedAt:  jwt.NewNumericDate(now),
			ID:        uuid.New().String(),
		},
		AdminID:  adminID,
		Username: username,
		Role:     role,
	}

	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(s.cfg.Secret))
	if err != nil {
		return nil, fmt.Errorf("sign access token: %w", err)
	}

	rawRefresh := uuid.New().String()
	refreshHash := hashToken(rawRefresh)
	refreshExp := now.Add(s.cfg.RefreshTTL)

	_, err = s.q.CreateRefreshToken(ctx, query.CreateRefreshTokenParams{
		AdminUserID: adminID,
		TokenHash:   refreshHash,
		ExpiresAt:   refreshExp,
	})
	if err != nil {
		return nil, fmt.Errorf("store refresh token: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: rawRefresh,
		ExpiresAt:    accessExp.Unix(),
	}, nil
}

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func hashToken(raw string) string {
	h := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(h[:])
}
