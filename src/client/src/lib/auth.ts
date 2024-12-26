import { LOGIN_URL } from './consts';

export const redirectToLogin = () => {
	window.location.href = `${LOGIN_URL}?redirect_to=${window.location.href}`;
};

export const ensureToken = () => {
	const url = new URL(window.location.href);
	let token = url.searchParams.get('token');

	if (token) {
		localStorage.setItem('token', token);
		url.searchParams.delete('token');
		window.history.replaceState({}, url.pathname);
	} else {
		token = localStorage.getItem('token');

		if (!token) {
			redirectToLogin();
		}
	}
};
