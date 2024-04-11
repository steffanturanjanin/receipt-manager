import { httpClient } from "../http"
import { profile } from "./mock";

export const register = async (request: RegisterRequest): Promise<AuthResponse> => {
	const { data } = await httpClient.post<AuthResponse>("/auth/register", request);

	return data;
}

export const login = async (request: LoginRequest): Promise<AuthResponse> => {
	const { data } = await httpClient.post<AuthResponse>("/auth/login", request);

	return data;
}

export const logout = async (): Promise<void> => {
	const { data } = await httpClient.post<void>("/auth/logout");

	return data;
}

export const getProfile = async (): Promise<Profile> => {
	// const { data } = await httpClient.get<Profile>("/auth/me");

	// return data;

	return profile;
}
