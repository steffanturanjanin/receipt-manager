import { FunctionComponent, ReactElement } from "react";
import { Navigate } from "react-router-dom";
import { getAuth } from "../../../util/auth";

const Home: FunctionComponent = (): ReactElement => {
	const { access_token } = getAuth() || {};

	if (access_token) {
		return <Navigate to="/receipts" />
	}

	return <Navigate to="/auth/login" />
}

export default Home;
