export type Playlist = {
	id: string;
	createdAt: string;
	name: string;
	userId: string;
	status: 'active' | 'importing' | 'failed';
};

export type Track = {
	id: string;
	createdAt: string;
	playlistId: string;
	title: string;
	artist: string;
	year: number;
	spotifyURI: string;
};

export type Paginated<D> = {
	data: D[];
	total: number;
};
