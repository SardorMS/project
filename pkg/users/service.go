package users

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log"
	"reflect"

	"github.com/SardorMS/project/pkg/types"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInternal        = errors.New("internal error")
	ErrEmptyRows       = errors.New("app data is empty")
	ErrNotFound        = errors.New("not found")
	ErrInvalidPassword = errors.New("invalid password")
)

// Service - структура для сервиса users.
type Service struct {
	pool *pgxpool.Pool
}

// NewService - создание сервиса users.
func NewService(pool *pgxpool.Pool) *Service {
	return &Service{pool: pool}
}

// IDByToken - сервис для middlware (выдача id по его токену с помощью контекста).
func (s *Service) IDByToken(ctx context.Context, token string) (int64, error) {

	var id int64
	sql := `SELECT user_id FROM users_tokens WHERE access_token = $1;`
	err := s.pool.QueryRow(ctx, sql, token).Scan(&id)

	if err == pgx.ErrNoRows {
		return 0, nil
	}

	if err != nil {
		return 0, ErrInternal
	}
	return id, nil
}

// Register - сервис процедуры регистрации пользователя.
func (s *Service) Register(ctx context.Context, reg *types.Registration, country string) (
	*types.UserPrivate, error) {

	var err error
	item := &types.UserPrivate{}

	idFromPass, err := bcrypt.GenerateFromPassword([]byte(reg.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		return nil, ErrNotFound
	}

	buffer := make([]byte, 8)
	n, err := rand.Read(buffer)
	if n != len(buffer) || err != nil {
		return nil, ErrInternal
	}
	userID := hex.EncodeToString(buffer)

	href := "https://api.spotify.com/v1/users/" + userID
	uri := "spotify:" + "user:" + userID

	sql := `INSERT INTO users 
	(country, display_name, email, href, id_from_pass, user_id, product, birthdate, uri) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) ON CONFLICT DO NOTHING 
	RETURNING id, country, display_name, email, href, 
			  user_id, images, product, birthdate, uri, created;`
	err = s.pool.QueryRow(ctx, sql, country, reg.DisplayName, reg.Email, href,
		idFromPass, userID, reg.Product, reg.Birthdate, uri).Scan(
		&item.ID, &item.Country, &item.DisplayName, &item.Email, &item.Endpoint,
		&item.UserID, &item.Images, &item.Product, &item.Birthdate, &item.URI, &item.Created)

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

// Login - сервис для идентификации пользователя.
func (s *Service) Login(ctx context.Context, token string, email string, pass string) (
	item *types.Login, err error) {

	var (
		id          int64
		idFromPass  string
		userID      string
		country     string
		displayName string
		product     string
	)

	sql1 := `SELECT id, country, display_name, id_from_pass, user_id, product 
			 FROM users WHERE email = $1;`
	err = s.pool.QueryRow(ctx, sql1, email).Scan(
		&id, &country, &displayName, &idFromPass, &userID, &product)

	if err == pgx.ErrNoRows {
		return nil, ErrEmptyRows
	}

	if err != nil {
		return nil, ErrInternal
	}

	err = bcrypt.CompareHashAndPassword([]byte(idFromPass), []byte(pass))
	if err != nil {
		return nil, ErrInvalidPassword
	}

	sql2 := `INSERT INTO users_tokens (user_id, access_token) VALUES($1, $2)
			 ON CONFLICT DO NOTHING;`
	_, err = s.pool.Exec(ctx, sql2, id, token)

	if err == pgx.ErrNoRows {
		log.Println(err)
		return nil, ErrEmptyRows
	}

	if err != nil {
		log.Println(err)
		return nil, ErrInternal
	}

	item = &types.Login{
		ID:          id,
		Country:     country,
		DisplayName: displayName,
		Email:       email,
		UserID:      userID,
		Product:     product,
	}

	return item, nil
}

// GetCurrentUser - сервис для получения информации о текущем пользователе.
func (s *Service) GetCurrentUser(ctx context.Context, id int64) (*types.UserPrivate, error) {
	item := &types.UserPrivate{}

	sql := `SELECT id, country, display_name, email, href, 
			user_id, images, product, birthdate, uri, created  
			FROM users WHERE id = $1`

	err := s.pool.QueryRow(ctx, sql, id).Scan(
		&item.ID, &item.Country, &item.DisplayName, &item.Email, &item.Endpoint,
		&item.UserID, &item.Images, &item.Product, &item.Birthdate, &item.URI, &item.Created)

	if errors.Is(err, pgx.ErrNoRows) {
		log.Println("No Rows")
		return nil, ErrNotFound
	}

	if err != nil {
		log.Println(err)
		return nil, ErrInternal
	}
	return item, nil
}

// GetUserByID - сервис для получения информации о пользователе по его user_id.
func (s *Service) GetUserByID(ctx context.Context, userID string) (*types.AnyUser, error) {
	item := &types.AnyUser{}

	href := "https://api.spotify.com/v1/users/" + userID

	sql := `SELECT display_name, href, user_id, images, uri 
			FROM users WHERE href = $1;`

	err := s.pool.QueryRow(ctx, sql, href).Scan(
		&item.DisplayName, &item.Endpoint, &item.UserID, &item.Images, &item.URI)

	if errors.Is(err, pgx.ErrNoRows) {
		log.Println("Incorrect user_id For The Current User!")
		return nil, ErrNotFound
	}

	if err != nil {
		log.Println(err)
		return nil, ErrInternal
	}
	return item, nil
}

// ImageUpload - сервис загрузки фотографии на профиль для текущего пользователя.
func (s *Service) ImageUpload(ctx context.Context, item *types.ImageUpdate, id int64) (
	*types.ImageUpdate, string, error) {

	data1 := &types.ImageChecker{
		Images: []types.Image{},
	}
	data2 := &types.ImageChecker{
		Images: []types.Image{},
	}

	var message string

	sql1 := `SELECT images FROM users WHERE id = $1;`
	err := s.pool.QueryRow(ctx, sql1, id).Scan(&data1.Images)
	if err == pgx.ErrNoRows {
		log.Println(err)
		return nil, "0", ErrEmptyRows
	}

	if reflect.DeepEqual(data1, data2) {

		sql := `UPDATE users SET images = images || $2 WHERE id = $1 
	RETURNING id, country, display_name, href, user_id, images, uri, product, created;`
		err = s.pool.QueryRow(ctx, sql, id, item.Images).Scan(
			&item.ID, &item.Country, &item.DisplayName, &item.Endpoint,
			&item.UserID, &item.Images, &item.URI, &item.Product, &item.Created)

		if err == pgx.ErrNoRows {
			log.Println(err)
			return nil, "0", ErrEmptyRows
		}
		if err != nil {
			log.Println(err)
			return nil, "0", ErrInternal
		}
	} else {
		message = "1"
	}

	return item, message, nil
}

// ImageRemove - сервис удаления фотографии с профиля для текущего пользователя.
func (s *Service) ImageRemove(ctx context.Context, id int64) (
	*types.ImageUpdate, error) {

	item := &types.ImageUpdate{}

	sql := `UPDATE users SET images = '[]' WHERE id = $1 
	RETURNING id, country, display_name, href, user_id, images, uri, product, created;`
	err := s.pool.QueryRow(ctx, sql, id).Scan(
		&item.ID, &item.Country, &item.DisplayName, &item.Endpoint,
		&item.UserID, &item.Images, &item.URI, &item.Product, &item.Created)

	if err != nil {
		log.Println(err)
		return nil, ErrInternal
	}

	return item, nil
}
