CREATE TABLE users (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  spotify_id TEXT UNIQUE,
  name TEXT NOT NULL,
  token TEXT
);
CREATE TABLE playlists (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  name TEXT NOT NULL,
  user_id UUID references users(id) NOT NULL
);
CREATE TABLE tracks (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  playlist_id UUID references playlists(id) NOT NULL,
  title TEXT NOT NULL,
  artist TEXT NOT NULL,
  year INT NOT NULL,
  spotify_uri TEXT NOT NULL
);
CREATE TABLE games (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  user_id UUID references users(id) NOT NULL,
  playlist_id UUID references playlists(id) NOT NULL
);
CREATE TABLE game_tracks (
  game_id UUID references games(id) NOT NULL,
  track_id UUID references tracks(id) NOT NULL,
  PRIMARY KEY(game_id, track_id)
);