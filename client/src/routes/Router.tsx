import { RouteObject, RouterProvider, createBrowserRouter } from "react-router-dom";
import auth from "./auth";
import { FunctionComponent } from "react";

const routes: RouteObject[] = [
	...auth,
];

const Router: FunctionComponent = () =>
	<RouterProvider router={createBrowserRouter(routes)} />

export default Router;
