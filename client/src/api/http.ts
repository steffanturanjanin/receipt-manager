import axios from 'axios';

export const httpClient = axios.create({
	baseURL: import.meta.env.VITE_BASE_API_URL as string,
	headers: { 'Content-Type': 'application/json' },
});
