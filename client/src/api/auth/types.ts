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

interface Profile {
	id: number;
	firstName: string;
	lastName: string;
	email: string;
	receiptsCount: number;
	registeredAt: string;
}
