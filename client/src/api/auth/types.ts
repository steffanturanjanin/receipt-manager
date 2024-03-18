interface RegisterRequest {
	firstName: string;
	lastName: string;
	email: string;
	password: string;
}

interface LoginRequest {
	email: string;
	password: string;
}

interface AuthResponse {
	access_token: string;
}

export type {
	AuthResponse,
	LoginRequest,
	RegisterRequest,
}
