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
  user_id UUID NOT NULL references users(id) ON DELETE CASCADE
);
CREATE TABLE tracks (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  playlist_id UUID NOT NULL references playlists(id) ON DELETE CASCADE,
  title TEXT NOT NULL,
  artist TEXT NOT NULL,
  year INT NOT NULL,
  spotify_uri TEXT NOT NULL
);
CREATE TABLE games (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  user_id UUID NOT NULL references users(id) ON DELETE CASCADE,
  playlist_id UUID NOT NULL references playlists(id) ON DELETE CASCADE
);
CREATE TABLE game_tracks (
  game_id UUID NOT NULL references games(id) ON DELETE CASCADE,
  track_id UUID NOT NULL references tracks(id) ON DELETE CASCADE,
  PRIMARY KEY(game_id, track_id)
);