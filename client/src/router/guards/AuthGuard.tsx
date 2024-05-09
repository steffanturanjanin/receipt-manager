import { FunctionComponent, ReactElement, ReactNode } from "react";
import { getAuth } from "../../util/auth";
import { Navigate } from "react-router-dom";

interface AuthGuardProps {
	element: ReactNode;
}

const AuthGuard: FunctionComponent<AuthGuardProps> = ({ element }): ReactElement => {
	const { access_token } = getAuth() || {};

	if (!access_token) {
		return <Navigate to="/auth/login" replace />
	}

	return <>{element}</>
}

export default AuthGuard;
