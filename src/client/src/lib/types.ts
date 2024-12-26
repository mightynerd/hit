export type Playlist = {
	id: string;
	created_at: string;
	name: string;
	user_id: string;
};

export type Track = {
	id: string;
	created_at: string;
	playlist_id: string;
	title: string;
	artist: string;
	year: number;
	spotify_uri: string;
};
