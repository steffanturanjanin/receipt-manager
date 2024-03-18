import { httpClient } from "../http"
import { AuthResponse, RegisterRequest } from "./types"

export const register = async (request: RegisterRequest) => {
	const { data } = await httpClient.post<AuthResponse>('/auth/register', request);

	return data;
}
