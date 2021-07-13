package playlists

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log"
	"reflect"
	"strings"

	"github.com/SardorMS/project/pkg/types"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v4/pgxpool"
)

var (
	ErrInternal  = errors.New("internal error")
	ErrEmptyRows = errors.New("app data is empty")
	ErrNotFound  = errors.New("not found")
	ErrName      = errors.New("empty name")
	ErrTrack     = errors.New("track is exist")
)

// Service - структура для сервиса playlists.
type Service struct {
	pool *pgxpool.Pool
}

// NewService - создание сервиса для playlists.
func NewService(pool *pgxpool.Pool) *Service {
	return &Service{pool: pool}
}

// PlaylistCreate - сервис создания плейлистов для текущего пользователя.
func (s *Service) PlaylistCreate(ctx context.Context,
	items *types.UserPlaylist, id int64, userID string) (*types.Playlists, error) {

	data := &types.Playlists{}
	var ownerName string

	sql1 := `SELECT display_name FROM users WHERE id = $1 AND user_id = $2;`
	err := s.pool.QueryRow(ctx, sql1, id, userID).Scan(&ownerName)

	if errors.Is(err, pgx.ErrNoRows) {
		log.Println("Incorrect user_id For The Current User!")
		return nil, ErrNotFound
	}

	if err != nil {
		log.Println(err)
		return nil, ErrNotFound
	}
	changedName := strings.ToLower(ownerName)

	buffer := make([]byte, 8)
	n, err := rand.Read(buffer)
	if n != len(buffer) || err != nil {
		log.Println(err)
		return nil, err
	}
	playlistID := hex.EncodeToString(buffer)

	href := "https://api.spotify.com/v1/users/" +
		changedName + "/playlists/" + playlistID

	uri := "spotify:user:" + changedName + ":playlist:" + playlistID

	sql2 := `INSERT INTO playlists 
		(name, owner_name, descriptions, is_public, images, user_id, playlist_id, href, uri) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) ON CONFLICT DO NOTHING
		RETURNING id, name, owner_name, descriptions, is_public, 
		images, tracks, playlist_id, href, uri, created;`

	err = s.pool.QueryRow(ctx, sql2, items.Name, ownerName, items.Descriptions,
		items.IsPublic, items.Images, id, playlistID, href, uri).Scan(
		&data.ID, &data.Name, &data.OwnerName, &data.Descriptions, &data.IsPublic,
		&data.Images, &data.Tracks, &data.PlaylistID, &data.Endpoint, &data.URI, &data.Created)

	if err == pgx.ErrNoRows {
		log.Println(err)
		return nil, ErrEmptyRows
	}
	if err != nil {
		log.Println(err)
		return nil, ErrNotFound
	}

	return data, nil
}

// GetAllPlaylists - сервис для получения информации всех плейлистов текущего пользователя.
func (s *Service) GetAllPlaylists(ctx context.Context, id int64, userID string) (
	[]*types.Playlists, error) {

	items := make([]*types.Playlists, 0)

	sql2 := `SELECT p.id, p.name, p.owner_name, p.descriptions, p.is_public, 
			p.images, p.tracks, p.playlist_id, p.href, p.uri, p.created 
			FROM playlists p JOIN users u
			ON u.user_id = $1
			WHERE p.user_id = $2;`
	rows, err := s.pool.Query(ctx, sql2, userID, id)
	if err != nil {
		log.Println(err)
		return nil, ErrNotFound
	}
	defer rows.Close()

	for rows.Next() {
		item := &types.Playlists{}
		err = rows.Scan(
			&item.ID, &item.Name, &item.OwnerName, &item.Descriptions, &item.IsPublic,
			&item.Images, &item.Tracks, &item.PlaylistID, &item.Endpoint, &item.URI, &item.Created)

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

// GetPlaylist - сервис для получения информации о плейлисте по его playlist_id.
func (s *Service) GetPlaylist(ctx context.Context, id int64, playlistID string) (
	*types.Playlists, error) {

	result := &types.Playlists{}

	sql1 := `SELECT id, name, owner_name, descriptions, is_public, 
			images, tracks, playlist_id, href, uri, created 
			FROM playlists
			WHERE user_id = $1 AND playlist_id = $2;`

	err := s.pool.QueryRow(ctx, sql1, id, playlistID).Scan(
		&result.ID, &result.Name, &result.OwnerName, &result.Descriptions,
		&result.IsPublic, &result.Images, &result.Tracks, &result.PlaylistID,
		&result.Endpoint, &result.URI, &result.Created)

	if err == pgx.ErrNoRows {
		log.Println(err)
		return nil, ErrEmptyRows
	}
	if err != nil {
		log.Println("Please Check Playlist_ID First, It's Wrong!")
		return nil, ErrNotFound
	}

	return result, nil
}

// RemoveAllPlaylists - сервис для удаления всех плйлистов текущего пользователя.
func (s *Service) RemoveAllPlaylists(ctx context.Context, id int64, userID string) error {

	sql := `DELETE FROM playlists p
			USING users u
			WHERE p.user_id = $1 AND u.user_id = $2;`
	_, err := s.pool.Exec(ctx, sql, id, userID)

	if err != nil {
		log.Println("Something Went Wrong and Playlists not DELETED!")
		return nil
	}
	return nil
}

// PlaylistUploadImage - сервис для загрузки изображения на плейлист, конкретного пользователя.
func (s *Service) PlaylistUploadImage(ctx context.Context, item *types.Image,
	id int64, playlistID string) (*types.Playlists, string, error) {

	data1 := &types.ImageChecker{
		Images: []types.Image{},
	}
	data2 := &types.ImageChecker{
		Images: []types.Image{},
	}

	result := &types.Playlists{}
	var message string

	sql1 := `SELECT images FROM playlists WHERE user_id = $1 AND playlist_id = $2;`
	err := s.pool.QueryRow(ctx, sql1, id, playlistID).Scan(&data1.Images)

	if err == pgx.ErrNoRows {
		log.Println(err)
		return nil, "0", ErrEmptyRows
	}

	if reflect.DeepEqual(data1, data2) {

		result = &types.Playlists{
			Images: []types.Image{
				{
					URL:    item.URL,
					Width:  item.Width,
					Height: item.Height,
				},
			},
		}
		sql := `UPDATE playlists SET images = images || $3
			WHERE user_id = $1 AND playlist_id = $2
			RETURNING id, name, owner_name, descriptions, is_public, images, 
			tracks, playlist_id, href, uri, created;`

		err = s.pool.QueryRow(ctx, sql, id, playlistID, result.Images).Scan(
			&result.ID, &result.Name, &result.OwnerName, &result.Descriptions,
			&result.IsPublic, &result.Images, &result.Tracks, &result.PlaylistID,
			&result.Endpoint, &result.URI, &result.Created)

		if err == pgx.ErrNoRows {
			log.Println(err)
			return nil, "0", ErrEmptyRows
		}

		if err != nil {
			log.Println(err)
			return nil, "0", ErrNotFound
		}

	} else {
		message = "1"
	}
	return result, message, nil
}

// PlaylistImageRemove - сервис для удаления изображения с плейлиста, для конкретного пользователя.
func (s *Service) PlaylistImageRemove(ctx context.Context, id int64, playlistID string) (
	*types.Playlists, error) {

	result := &types.Playlists{}

	sql := `UPDATE playlists SET images = '[]' WHERE user_id = $1 AND playlist_id = $2 
			RETURNING id, name, owner_name, descriptions, is_public, 
			images, tracks, playlist_id, href, uri, created;`
	err := s.pool.QueryRow(ctx, sql, id, playlistID).Scan(
		&result.ID, &result.Name, &result.OwnerName, &result.Descriptions,
		&result.IsPublic, &result.Images, &result.Tracks, &result.PlaylistID,
		&result.Endpoint, &result.URI, &result.Created)

	if err != nil {
		log.Println(err)
		return nil, ErrInternal
	}

	return result, nil
}

// PlaylistChange - сервис для внесения изменений в детали плейлиста конкретного пользователя.
func (s *Service) PlaylistChange(ctx context.Context, item *types.UserPlaylist,
	id int64, playlistID string) (*types.Playlists, error) {

	var err error
	result := &types.Playlists{}

	if item.Name != "" {

		if strings.TrimSpace(item.Name) == "" {
			log.Println("String Is Empty, You Can Not Set White Spaces!")
			return nil, ErrName
		}

		sql := `UPDATE playlists SET name = $3 
				WHERE user_id = $1 AND playlist_id = $2 
	            RETURNING id, name, owner_name, descriptions, is_public, images, 
	            tracks, playlist_id, href, uri, created;`
		err = s.pool.QueryRow(ctx, sql, id, playlistID, item.Name).Scan(
			&result.ID, &result.Name, &result.OwnerName, &result.Descriptions,
			&result.IsPublic, &result.Images, &result.Tracks, &result.PlaylistID,
			&result.Endpoint, &result.URI, &result.Created)

		if err == pgx.ErrNoRows {
			log.Println(err)
			return nil, ErrEmptyRows
		}
		if err != nil {
			log.Println("Something Went Wrong With Name, Please Check Playlist_ID First!")
			return nil, ErrNotFound
		}
	}

	if item.Descriptions != "" {
		sql := `UPDATE playlists SET descriptions = $3 
				WHERE user_id = $1 AND playlist_id = $2
				RETURNING id, name, owner_name, descriptions, is_public, images, 
				tracks, playlist_id, href, uri, created;`
		err = s.pool.QueryRow(ctx, sql, id, playlistID, item.Descriptions).Scan(
			&result.ID, &result.Name, &result.OwnerName, &result.Descriptions,
			&result.IsPublic, &result.Images, &result.Tracks, &result.PlaylistID,
			&result.Endpoint, &result.URI, &result.Created)

		if err == pgx.ErrNoRows {
			log.Println(err)
			return nil, ErrEmptyRows
		}

		if err != nil {
			log.Println("Something Went Wrong With Descriptions, Please Check Playlist_ID First!")
			return nil, ErrNotFound
		}
	}

	sql := `UPDATE playlists SET is_public = $3 
		WHERE user_id = $1 AND playlist_id = $2
		RETURNING id, name, owner_name, descriptions, is_public, images, 
		tracks, playlist_id, href, uri, created;`
	err = s.pool.QueryRow(ctx, sql, id, playlistID, item.IsPublic).Scan(
		&result.ID, &result.Name, &result.OwnerName, &result.Descriptions,
		&result.IsPublic, &result.Images, &result.Tracks, &result.PlaylistID,
		&result.Endpoint, &result.URI, &result.Created)

	if err == pgx.ErrNoRows {
		log.Println(err)
		return nil, ErrEmptyRows
	}

	if err != nil {
		log.Println("Something Went Wrong With Status, Please Check Playlist_ID First!", err)
		return nil, ErrNotFound
	}

	return result, nil
}

// TrackToPlaylist - сервис для добавления трека в конкретный плейлист для определённого пользователя.
func (s *Service) TrackToPlaylist(ctx context.Context,
	items *types.TrackInfo, id int64, playlistID string) (*types.Playlists, error) {

	data := &types.Track{}

	sql1 := `SELECT id, name, artist_name, album_name, track_number, 
		genres, duration, track_id, href, uri, created 
		FROM tracks
		WHERE name = $1 AND artist_name = $2 AND uri = $3;`

	err := s.pool.QueryRow(ctx, sql1, items.Name, items.ArtistName, items.URI).Scan(
		&data.ID, &data.Name, &data.ArtistName, &data.AlbumName, &data.TrackNumber,
		&data.Genres, &data.Duration, &data.TrackID, &data.Endpoint, &data.URI, &data.Created)

	if err == pgx.ErrNoRows {
		log.Println(err)
		return nil, ErrEmptyRows
	}
	if err != nil {
		log.Println("Please Input Data in Correct Case!, and Check Playlist_ID First!")
		return nil, ErrNotFound
	}

	data2 := &types.Playlists{}

	sql2 := `SELECT tracks FROM playlists
		     WHERE user_id = $1 AND playlist_id = $2;`

	err = s.pool.QueryRow(ctx, sql2, id, playlistID).Scan(&data2.Tracks)

	if err == pgx.ErrNoRows {
		log.Println(err)
		return nil, ErrEmptyRows
	}
	if err != nil {
		log.Println(err)
		return nil, ErrNotFound
	}

	var uri types.URI
	for _, value := range data2.Tracks {
		uri = value.URI
		if data.URI != uri {
			continue
		} else {
			return nil, ErrTrack
		}
	}

	result := &types.Playlists{
		Tracks: []types.Track{
			{
				ID:          data.ID,
				TrackID:     data.TrackID,
				Name:        data.Name,
				ArtistName:  data.ArtistName,
				AlbumName:   data.AlbumName,
				TrackNumber: data.TrackNumber,
				Genres:      data.Genres,
				Duration:    data.Duration,
				Endpoint:    data.Endpoint,
				URI:         data.URI,
				Created:     data.Created,
			},
		},
	}

	sql3 := `UPDATE playlists SET tracks = tracks || $3 
				WHERE user_id = $1 AND playlist_id = $2
				RETURNING id, name, owner_name, descriptions, 
				is_public, images, tracks, playlist_id, href, uri, created;`
	err = s.pool.QueryRow(ctx, sql3, id, playlistID, result.Tracks).Scan(
		&result.ID, &result.Name, &result.OwnerName, &result.Descriptions,
		&result.IsPublic, &result.Images, &result.Tracks, &result.PlaylistID,
		&result.Endpoint, &result.URI, &result.Created)

	if err != nil {
		log.Println("No Updated Rows:", err)
		return nil, ErrInternal
	}

	return result, nil
}

// RemoveTrackFromPlaylist - сервис для удаления трека из конкретного плейлиста,
// определённого пользователя.
func (s *Service) RemoveTrackFromPlaylist(ctx context.Context,
	items *types.TrackInfo, id int64, playlistID string) (*types.Playlists, error) {

	result := &types.Playlists{}

	sql1 := `
	   	WITH t AS (  
			SELECT jsonb_agg( (tracks ->> ( idx-1 )::int)::jsonb ) AS js_new  
		    FROM playlists   
		    CROSS JOIN jsonb_array_elements(tracks)   
		    WITH ORDINALITY arr(j,idx)   
		    WHERE j->> 'uri' != $1
		) 
       UPDATE playlists    
	   SET tracks = js_new
       FROM t
	   WHERE user_id = $2 AND playlist_id = $3;`

	_, err := s.pool.Exec(ctx, sql1, items.URI, id, playlistID)

	if err == pgx.ErrNoRows {
		log.Println(err)
		return nil, ErrEmptyRows
	}

	if err != nil {
		log.Println("No Updated Rows:", err)
		return nil, ErrInternal
	}

	sql2 := `UPDATE playlists SET tracks = CASE
			 WHEN tracks IS NULL THEN '[]'::jsonb 
			 ELSE tracks
			 END
		     WHERE user_id = $1 AND playlist_id = $2
			 RETURNING id, name, owner_name, descriptions, is_public, images, 
			 tracks, playlist_id, href, uri, created;`

	err = s.pool.QueryRow(ctx, sql2, id, playlistID).Scan(
		&result.ID, &result.Name, &result.OwnerName, &result.Descriptions,
		&result.IsPublic, &result.Images, &result.Tracks, &result.PlaylistID,
		&result.Endpoint, &result.URI, &result.Created)

	if err == pgx.ErrNoRows {
		log.Println(err)
		return nil, ErrEmptyRows
	}

	if err != nil {
		log.Println("No Updated Rows:", err)
		return nil, ErrInternal
	}

	return result, nil
}

// RemoveAllTracksFromPlaylist - сервис для удаления всех стреков из определённого плейлиста
// для конкретного пользователя.
func (s *Service) RemoveAllTracksFromPlaylist(ctx context.Context, id int64,
	playlistID string) (*types.Playlists, error) {

	result := &types.Playlists{}

	sql2 := `UPDATE playlists SET tracks = '[]'::jsonb  
			WHERE user_id = $1 AND playlist_id = $2
			RETURNING id, name, owner_name, descriptions, is_public, 
			images, tracks, playlist_id, href, uri, created;`

	err := s.pool.QueryRow(ctx, sql2, id, playlistID).Scan(
		&result.ID, &result.Name, &result.OwnerName, &result.Descriptions,
		&result.IsPublic, &result.Images, &result.Tracks, &result.PlaylistID,
		&result.Endpoint, &result.URI, &result.Created)

	if err == pgx.ErrNoRows {
		log.Println(err)
		return nil, ErrEmptyRows
	}

	if err != nil {
		log.Println("No Updated Rows:", err)
		return nil, ErrInternal
	}

	return result, nil
}
