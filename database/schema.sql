

--               Таблица приложений.
-- Данная таблица для хранения данных исключительно для каждого пользователя,
-- который будет использовать в дальнейшем ваше приложение.
CREATE TABLE IF NOT EXISTS applications 
(
    id             BIGSERIAL PRIMARY  KEY,
    client_id      TEXT      NOT NULL,
    redirect_uri   TEXT      NOT NULL,
    code_challenge TEXT      NOT NULL UNIQUE,
    code_verifier  TEXT      NOT NULL UNIQUE,
    state          TEXT      NOT NULL UNIQUE,
    scope          TEXT[]    NOT NULL DEFAULT '{}',
    code           TEXT      NOT NULL UNIQUE,
    created        TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

--                 Таблица токенов 
-- Для каждого пользователя на очновании кода приложения, генерируется токен доступа.

CREATE TABLE IF NOT EXISTS applications_tokens 
(
    application_id BIGINT    NOT NULL REFERENCES applications,
    access_token   TEXT      NOT NULL,
    refresh_token  TEXT      NOT NULL, 
    code_verifier  TEXT      NOT NULL UNIQUE,
    expires_in     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP + INTERVAL '1 hour',
    created        TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP 
);


--                Таблица пользователей.
-- Таблица для манипуляций с пользователями.

CREATE TABLE IF NOT EXISTS users
(
    id             BIGSERIAL PRIMARY  KEY,
    country        TEXT      NOT NULL,     
	display_name   TEXT      NOT NULL UNIQUE,  
	email          TEXT      NOT NULL UNIQUE,
	href           TEXT      NOT NULL,    
	id_from_pass   TEXT      NOT NULL,
    user_id        TEXT      NOT NULL UNIQUE,        
	images         JSONB     NOT NULL DEFAULT '[]',       
	product        TEXT      NOT NULL,
	birthdate      TEXT      NOT NULL,   
	uri            TEXT      NOT NULL,
	created        TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);


-- Таблица токенов для пользователя.

CREATE TABLE IF NOT EXISTS users_tokens
(
    user_id      BIGINT    NOT NULL REFERENCES users,
    access_token TEXT      NOT NULL UNIQUE,
    expires_in   TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP + INTERVAL '1 hour',
    created      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);


--                  Таблица плейлистов 
-- для зарегистрированного пользователя.

CREATE TABLE IF NOT EXISTS playlists
(
    id           BIGSERIAL PRIMARY  KEY,
    user_id      BIGINT    NOT NULL REFERENCES users,
    playlist_id  TEXT      NOT NULL UNIQUE,
    name         TEXT      NOT NULL,
    owner_name   TEXT      NOT NULL,
    descriptions TEXT,
    is_public    BOOLEAN   NOT NULL DEFAULT TRUE,
    images       JSONB     NOT NULL DEFAULT '[]',
    tracks       JSONB     DEFAULT '[]',
    href         TEXT      NOT NULL,
    uri          TEXT      NOT NULL,
    created      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);



--                    Таблица треков 

CREATE TABLE IF NOT EXISTS tracks 
(
    id           BIGSERIAL PRIMARY  KEY,
    track_id     TEXT      NOT NULL UNIQUE,
    name         TEXT      NOT NULL,
    artist_name  TEXT      NOT NULL,
    album_name   TEXT      NOT NULL,
    track_number BIGINT    NOT NULL DEFAULT 0,
    genres       TEXT[]    NOT NULL DEFAULT '{}',
    duration     BIGINT    NOT NULL DEFAULT 0,
    href         TEXT      NOT NULL,
    uri          TEXT      NOT NULL,
    created      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);



-- DROP   TABLE  applications;
-- DROP   TABLE  applications_tokens;
-- DROP   TABLE  users;
-- DROP   TABLE  users_tokens;
-- DROP   TABLE  playlists;
-- DROP   TABLE  tracks;

-- DELETE FROM  applications;
-- DELETE FROM  applications_tokens;
-- DELETE FROM  users;
-- DELETE FROM  users_tokens;
-- DELETE FROM  playlists;
-- DELETE FROM  tracks;


-- ALTER  SEQUENCE applications_id_seq RESTART WITH 1;
-- ALTER  SEQUENCE users_id_seq RESTART WITH 1;
-- ALTER  SEQUENCE playlists_id_seq RESTART WITH 1;
-- ALTER  SEQUENCE tracks_id_seq RESTART WITH 1;

























------------------------------------------------------------
/* create table objects
(
    id BIGSERIAL PRIMARY KEY,
    word TEXT
    --image JSONB  DEFAULT '[]'
);


insert into objects (id, word) VALUES (1, null);
update objects set word = 'err' where id = 1;
insert into objects (id, image) VALUES (1, '[]');
insert into objects (id, image) VALUES (1, '[{"height": null, "url": null, "width": null}]');

insert into objects (id, image) VALUES ('[{"height": 10, "url": "https://../", "width": 20}]');

select image from objects
where id = 1; 



--- Добавляет новый объект в слайс JSONB
UPDATE objects 
SET image = image || '{"url": "https://spotify.com", "width": 12, "height": 10 }' 
WHERE id = 2;


-- Обновляет конкретный элемент (по индексу) в слайсе.
update objects set image = 
jsonb_set(jsonb_set(jsonb_set(image, '{0, url}', '"https://spotify.com"'), '{0, width}', '200'), '{0, height}', '200')
where id = 2;

update playlists set tracks = '[]' where id = 1;


delete from objects;
drop table objects;



-- Добавление по элементу
UPDATE objects SET
image = jsonb_set( image, '{1}', array_to_json(
	ARRAY(
		SELECT DISTINCT( UNNEST( ARRAY(
			SELECT json_array_elements_text( COALESCE( image::json->'1', '[]' ) )
		) || {'url: httpl', 'height: 120'} ) )
	)
)::jsonb )
WHERE id = 2
RETURNING *;


-- Удаление по элементу оставляет [] скобки
UPDATE objects SET
image = jsonb_set( image, '{1}', array_to_json(
	array_remove( ARRAY(
		SELECT json_array_elements_text( COALESCE( image::json->'1', '[]' ) )
	), 'url' )
)::jsonb )
WHERE id = 1
RETURNING *;


 */