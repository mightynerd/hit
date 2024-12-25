I don't know how to go

# API

## Auth

The client should redirect the user to `/login`. A `redirect_to` query parameter should be added in order for the user to be redirected back to the client. When this happens, a `token` query parameter will be added to the `redirect_to` URL. This JWT token must be used as a bearer token for all requests below.

Once the token expires, the authentication flow needs to be re-triggered.

## Playlists

### Create

Creates a playlist from a Spotify playlist and returns its id.

```
POST /playlists
{
    "name": "{name}",
    "from": {
        "source": "spotify_playlist",
        "id": "{spotify_playlist_id}"
    }
}
```

### List

```
GET /playlists
```

## Tracks

### List

```
GET /playlists/:playlist_id/tracks
```

### Update

All body properties are optional.

```
PATCH /playlists/:playlist_id/tracks/:track_id
{
  "title": "New Title",
  "artist": "New Artist",
  "year": 2000
}
```

### Delete

```
DELETE /playlists/:playlist_id/tracks/:track_id
```

## Games

### Create

Creates a game and returns its id.

```
POST /games
{
  "playlist_id": "{playlist_id}"
}
```

### Advance

Advances the game by one track returning the currently playing track.

```
POST /games/{gameId}/advance
```
