import { httpClient } from "../http"
import { AuthResponse, LoginRequest, RegisterRequest } from "./types"

export const register = async (request: RegisterRequest) => {
	const { data } = await httpClient.post<AuthResponse>("/auth/register", request);

	return data;
}

export const login = async (request: LoginRequest) => {
	const { data } = await httpClient.post<AuthResponse>("/auth/login", request);

	return data;
}
