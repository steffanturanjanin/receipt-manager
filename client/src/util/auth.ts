interface Auth {
	access_token: string;
}

const AUTH_KEY = "auth";

export const getAuth = (): Auth | null  => {
	const authData = localStorage.getItem(AUTH_KEY);
	if (!authData) {
		return null;
	}

	return JSON.parse(authData) as Auth;
}

export const removeAuth = (): void => {
	localStorage.removeItem(AUTH_KEY);
}
