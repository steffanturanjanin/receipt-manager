interface RegisterRequest {
	firstName: string,
	lastName: string,
	password: string,
	email: string;
}

interface AuthResponse {
	access_token: string,
}

export type {
	RegisterRequest,
	AuthResponse,
}
