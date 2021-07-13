package tracks

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log"

	"github.com/SardorMS/project/pkg/types"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v4/pgxpool"
)

var (
	ErrInternal  = errors.New("internal error")
	ErrEmptyRows = errors.New("app data is empty")
	ErrNotFound  = errors.New("not found")
)

// Service - структура сервиса для tracks.
type Service struct {
	pool *pgxpool.Pool
}

// NewService - создания сервиса tracks.
func NewService(pool *pgxpool.Pool) *Service {
	return &Service{pool: pool}
}

// TracksCreate - сервис для добавления одного или нескольких треков.
func (s *Service) TracksCreate(ctx context.Context, data []*types.Track) error {

	for _, v := range data {

		buffer := make([]byte, 8)
		n, err := rand.Read(buffer)
		if n != len(buffer) || err != nil {
			log.Println(err)
			return err
		}
		trackID := hex.EncodeToString(buffer)

		href := "https://api.spotify.com/v1/tracks/" + trackID
		uri := "spotify:track:" + trackID

		sql := `INSERT INTO tracks 
		(name, artist_name, album_name, track_number, genres, duration, track_id, href, uri) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) ON CONFLICT DO NOTHING`

		_, err = s.pool.Exec(ctx, sql, v.Name, v.ArtistName,
			v.AlbumName, v.TrackNumber, v.Genres, v.Duration, trackID, href, uri)

		if errors.Is(err, pgx.ErrNoRows) {
			log.Println(err)
			return ErrEmptyRows
		}
		if err != nil {
			log.Println(err)
			return ErrInternal
		}
	}
	return nil
}

// GetAllTracks - сервис для вывода всех треков.
func (s *Service) GetAllTracks(ctx context.Context) (
	[]*types.Track, error) {

	items := make([]*types.Track, 0)

	sql2 := `SELECT id, track_id, name, artist_name, album_name, 
			track_number, genres, duration, href, uri, created 
			FROM tracks`
	rows, err := s.pool.Query(ctx, sql2)
	if err != nil {
		log.Println(err)
		return nil, ErrNotFound
	}
	defer rows.Close()

	for rows.Next() {
		item := &types.Track{}
		err = rows.Scan(
			&item.ID, &item.TrackID, &item.Name, &item.ArtistName,
			&item.AlbumName, &item.TrackNumber, &item.Genres,
			&item.Duration, &item.Endpoint, &item.URI, &item.Created)

		if err != nil {
			log.Println(err)
			return nil, err
		}
		items = append(items, item)
	}

	err = rows.Err()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return items, nil
}

// TracksRemove - сервис для удаления всех треков.
func (s *Service) TracksRemove(ctx context.Context) error {

	sql1 := `DELETE FROM tracks`
	_, err := s.pool.Exec(ctx, sql1)

	if errors.Is(err, pgx.ErrNoRows) {
		log.Println(err)
		return ErrEmptyRows
	}
	if err != nil {
		log.Println(err)
		return ErrInternal
	}

	sql2 := `ALTER SEQUENCE tracks_id_seq RESTART WITH 1;`
	_, err = s.pool.Exec(ctx, sql2)

	if errors.Is(err, pgx.ErrNoRows) {
		log.Println(err)
		return ErrEmptyRows
	}
	if err != nil {
		log.Println(err)
		return ErrInternal
	}

	return nil
}
