import { FunctionComponent, ReactElement, ReactNode } from "react";
import { Navigate } from "react-router-dom";
import { getAuth } from "../../util/auth";

interface GuestGuardProps {
	element: ReactNode;
}

const GuestGuard: FunctionComponent<GuestGuardProps> = ({ element }): ReactElement => {
	const { access_token } = getAuth() || {};

	if (access_token) {
		return <Navigate to="/" replace />
	}

	return <>{element}</>
}

export default GuestGuard;
