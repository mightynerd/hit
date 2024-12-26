import axios, { AxiosError } from 'axios';
import type { Playlist, Track } from './types';
import { BASE_URL } from './consts';
import { redirectToLogin } from './auth';

const client = axios.create({
	baseURL: BASE_URL
});

client.interceptors.request.use((config) => {
	const token = localStorage.getItem('token');

	if (token) {
		config.headers.Authorization = `Bearer ${token}`;
	}

	return config;
});

client.interceptors.response.use(
	(resp) => resp,
	(err) => {
		if (err instanceof AxiosError && err.response?.status === 401) {
			localStorage.removeItem('token');
			redirectToLogin();
		}

		throw err;
	}
);

export default client;

export const getPlaylists = async (): Promise<Playlist[]> => {
	const resp = await client.get<Playlist[]>('/playlists');
	return resp.data;
};

export const deletePlaylist = async (playlistId: string) => {
	await client.delete(`/playlists/${playlistId}`);
};

export const createPlaylist = async (name: string, spotifyId: string) => {
	await client.post('/playlists', {
		name,
		from: {
			source: 'spotify_playlist',
			id: spotifyId
		}
	});
};

export const getTracks = async (playlistId: string, page = 0, size = 20): Promise<Track[]> => {
	const resp = await client.get<Track[]>(
		`/playlists/${playlistId}/tracks?page=${page}&size=${size}`
	);
	return resp.data;
};

export const deleteTrack = async (playlistId: string, trackId: string) => {
	await client.delete(`/playlists/${playlistId}/tracks/${trackId}`);
};

export const createGame = async (playlistId: string): Promise<{ id: string }> => {
	const resp = await client.post<{ id: string }>('/games', { playlist_id: playlistId });
	return resp.data;
};

export const advanceGame = async (gameId: string): Promise<Track> => {
	const resp = await client.post<Track>(`/games/${gameId}/advance`);
	return resp.data;
};
