import ky from 'ky';
import { PUBLIC_API_BASE_URL } from '$env/static/public';
import { browser } from '$app/environment';

let baseURL = PUBLIC_API_BASE_URL;

if (browser) {
	if (baseURL == undefined) {
		baseURL = location.origin;
	}
}

console.info('API_BASE_URL:', baseURL);

const lsKey = 'API-KEY';
const getApiKey = function (): string {
	let apiKey = localStorage.getItem(lsKey);
	if (apiKey == null) {
		apiKey = prompt('Please enter the api key');
		if (apiKey == null) {
			throw new Error('API Key is required');
		}
	}

	localStorage.setItem(lsKey, apiKey);
	return apiKey;
};

export const api = ky.create({
	prefixUrl: baseURL,
	hooks: {
		beforeRequest: [
			(request) => {
				request.headers.set('X-API-Key', getApiKey());
			}
		]
	}
});
