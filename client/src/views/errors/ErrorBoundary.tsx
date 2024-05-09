import { FunctionComponent, ReactElement } from "react";
import { useRouteError } from "react-router-dom";
import { AxiosError } from "axios";
import NotFoundPage from "./NotFoundPage";

const ErrorBoundary: FunctionComponent = (): ReactElement => {
	const error = useRouteError() as AxiosError;

	if (error.status === 404) {
		return <NotFoundPage />
	}

	return <>Something went wrong!</>
}

export default ErrorBoundary;
