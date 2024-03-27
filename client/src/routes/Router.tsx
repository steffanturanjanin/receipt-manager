import { FunctionComponent } from "react";
import { RouteObject, RouterProvider, createBrowserRouter } from "react-router-dom";
import auth from "./auth";
import Home from "../views/pages/home/Home";
import AppLayout from "../views/layouts/AppLayout";

const routes: RouteObject[] = [
	...auth,
	{
		path: "/",
		element: <AppLayout><Home /></AppLayout>,
	}
];

const Router: FunctionComponent = () =>
	<RouterProvider router={createBrowserRouter(routes)} />

export default Router;
