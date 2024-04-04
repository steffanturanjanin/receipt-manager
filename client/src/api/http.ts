import axios from 'axios';

export const httpClient = axios.create({
	baseURL: import.meta.env.VITE_BASE_API_URL as string,
	headers: { 'Content-Type': 'application/json' },
});

httpClient.interceptors.request.use((config) => {
	const requestConfig = { ...config };

	const auth = localStorage.getItem("auth");
	const { access_token: accessToken } = JSON.parse(auth || "{}");

	if (accessToken) {
		requestConfig.headers.Authorization = `Bearer ${accessToken}`;
	}

	return requestConfig;
});
