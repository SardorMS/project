package spotify

import (
	"context"
	"errors"
	"log"
	"strings"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"

	"github.com/SardorMS/project/pkg/types"
)

var (
	ErrInternal            = errors.New("internal error")
	ErrEmptyRows           = errors.New("your data is already exists")
	ErrNotFound            = errors.New("not found")
	ErrTokenNotFound       = errors.New("token not found")
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
)

// Service - структура для сервиса spotify.
type Service struct {
	pool *pgxpool.Pool
}

// NewService - создание сервиса для spotify.
func NewService(pool *pgxpool.Pool) *Service {
	return &Service{pool: pool}
}

// RegisterApp - сервис регистрации приложения, точнее пользователя в момент использования приложения.
func (s *Service) RegisterApp(ctx context.Context, data *types.DataApp) (*types.Application, error) {

	var err error
	item := &types.Application{}

	sql := `INSERT INTO applications (client_id, redirect_uri, code_challenge, code_verifier, state, scope, code) 
		    VALUES ($1, $2, $3, $4, $5, $6, $7) ON CONFLICT DO NOTHING 
		    RETURNING id, client_id, redirect_uri, state, scope, code_verifier, code, created;`
	err = s.pool.QueryRow(ctx, sql, data.ClientID, data.RedirectURI, data.CodeChallenge, data.CodeVerifier, data.State, data.Scope, data.Code).Scan(
		&item.ID, &item.ClientID, &item.RedirectURI, &item.State, &item.Scope, &item.CodeVerifier, &item.Code, &item.Created)

	if err == pgx.ErrNoRows {
		log.Println(err)
		return nil, ErrEmptyRows
	}
	if err != nil {
		log.Println(err)
		return nil, ErrInternal
	}
	return item, nil
}

// CheckCode - сервис проверка статуса.
func (s *Service) CheckState(ctx context.Context, code string, state string) (string, error) {

	var tableState string
	sql := `SELECT state FROM applications WHERE code = $1;`
	err := s.pool.QueryRow(ctx, sql, code).Scan(&tableState)

	if tableState != state {
		log.Println("Check Your Code - State is not Matched: access_denied")
		return "0", ErrNotFound
	}
	if err == pgx.ErrNoRows {
		log.Println(err)
		return "0", ErrEmptyRows
	}

	if err != nil {
		log.Println(err)
		return "0", ErrInternal
	}

	return tableState, nil
}

// Token - сервис получения токена доступа.
func (s *Service) Token(ctx context.Context, data *types.Token) (*types.Tokens, error) {

	var (
		err          error
		id           int64
		codeVerifier string
		scope        []string
	)
	sql := `SELECT id, code_verifier, scope 
			FROM applications 
			WHERE code = $1;`
	err = s.pool.QueryRow(ctx, sql, data.Code).Scan(
		&id, &codeVerifier, &scope)

	if err == pgx.ErrNoRows {
		return nil, ErrNotFound
	}

	if err != nil {
		return nil, ErrInternal
	}

	if data.CodeVerifier != codeVerifier {
		log.Println("Code Verifier not Matched!")
		return nil, ErrNotFound
	}

	token, err := bcrypt.GenerateFromPassword([]byte(data.Code), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Can not create a token", err)
		return nil, err
	}

	refreshToken, err := bcrypt.GenerateFromPassword([]byte(token), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Can not create a refresh token", err)
		return nil, err
	}

	item := &types.Tokens{
		TokenType: "Bearer",
		Scope:     scope,
		ExpiresIn: 3600,
	}

	sql1 := `INSERT INTO applications_tokens 
			(application_id, access_token, refresh_token, code_verifier)
		     VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING
			 RETURNING  access_token, refresh_token, created;`
	err = s.pool.QueryRow(ctx, sql1, id, token, refreshToken, codeVerifier).Scan(
		&item.AccessToken, &item.RefreshToken, &item.Created)

	if err == pgx.ErrNoRows {
		log.Println("Please Change Your Data!")
		return nil, ErrTokenNotFound
	}
	if err != nil {
		log.Println(err)
		return nil, ErrInternal
	}
	return item, nil
}

// TokenUpdate - сервис обновления или продления токена доступа.
func (s *Service) TokenUpdate(ctx context.Context, data *types.TokenRefresh) (*types.TokenUpdate, error) {

	var (
		err      error
		appID    int64
		oldToken string
		scope    []string
	)
	sql := `SELECT at.application_id, at.refresh_token, a.scope 
			FROM applications_tokens at
			JOIN applications a ON a.id = at.application_id;`
	err = s.pool.QueryRow(ctx, sql).Scan(&appID, &oldToken, &scope)

	if err == pgx.ErrNoRows {
		return nil, ErrNotFound
	}

	if err != nil {
		return nil, ErrInternal
	}

	comparedStrings := strings.Compare(oldToken, data.RefreshToken)
	if comparedStrings != 0 {
		return nil, ErrInvalidRefreshToken
	}

	newToken, err := bcrypt.GenerateFromPassword([]byte(oldToken), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Can not create a refresh token", err)
		return nil, err
	}

	item := &types.TokenUpdate{
		TokenType: "Bearer",
		Scope:     scope,
		ExpiresIn: 3600,
	}

	sql1 := `UPDATE applications_tokens SET access_token = $1 WHERE application_id = $2
			 RETURNING access_token, created;`
	err = s.pool.QueryRow(ctx, sql1, newToken, appID).Scan(
		&item.AccessToken, &item.Created)

	if err == pgx.ErrNoRows {
		log.Println(err)
		return nil, ErrTokenNotFound
	}
	if err != nil {
		log.Println(err)
		return nil, ErrInternal
	}
	return item, nil
}
