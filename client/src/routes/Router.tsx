import { FunctionComponent } from "react";
import { RouteObject, RouterProvider, createBrowserRouter } from "react-router-dom";
import auth from "./auth";
import Home from "../views/pages/home/Home";

const routes: RouteObject[] = [
	...auth,
	{
		path: "/",
		element: <Home />,
	}
];

const Router: FunctionComponent = () =>
	<RouterProvider router={createBrowserRouter(routes)} />

export default Router;
