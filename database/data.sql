---
---
--- Раздел исключительно для наглядности. 
--- Если собираетесь вручную добавлять треки, то для плейлиста обязательно 
--- нужно скопировать и вводить данные точно такие же как и те которые были введены, а имеено
--- столбцы name, artist_name и uri.
---
--- Но лучше воспользоваться запросами в файле test2.http.


INSERT INTO tracks 
(name, artist_name, album_name, track_number, genres, duration, href, uri)
VALUES 
('American_Idiot.mp3', 'Green_Day', 'American_Idiot', 1, '{punk, punk-rock, alternative, rock, pop-rock}', 173, 'https://accounts.spotify.com/v1/tracks/american_idiot.mp3', 'spotify:track:american_idiot'),
('Boulevard_of_Broken_Dreams.mp3', 'Green_Day', 'American_Idiot', 4, '{punk, punk-rock, alternative, rock, pop-rock}', 259, 'https://accounts.spotify.com/v1/tracks/boulevard_of_broken_dreams.mp3', 'spotify:track:boulevard_of_broken_dreams'),
('Wake_Me_When_September_Ends.mp3', 'Green_Day', 'American_Idiot', 11, '{punk, punk-rock, alternative, rock, pop-rock}', 282, 'https://accounts.spotify.com/v1/tracks/wake_me_when_september_ends.mp3', 'spotify:track:wake_me_when_september_ends')
('Holiday.mp3', 'Green_Day', 'American_Idiot', 3, '{punk, punk-rock, alternative, rock, pop-rock}', 232, 'https://accounts.spotify.com/v1/tracks/holiday.mp3', 'spotify:track:holiday')
('Animal_I_Have_Become.mp3', 'Three_Days_Grace', 'One_X', 12, '{alternative, rock, grunge}', 286, 'https://accounts.spotify.com/v1/tracks/animal_i_have_become.mp3', 'spotify:track:animal_i_have_become')
ON CONFLICT DO NOTHING;


DROP TABLE tracks;
DELETE FROM  tracks;
ALTER  SEQUENCE tracks_id_seq RESTART WITH 1;


