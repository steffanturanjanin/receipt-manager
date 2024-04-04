import { FunctionComponent } from "react";
import { RouteObject, RouterProvider, createBrowserRouter } from "react-router-dom";
import auth from "./auth";
import Home from "../views/pages/home/Home";
import AppLayout from "../views/layouts/AppLayout/AppLayout";
import AuthGuard from "./guards/AuthGuard";
import ReceiptsPage from "../views/pages/receipts/ReceiptsPage";

const routes: RouteObject[] = [
	...auth,
	{
		path: "/",
		element: <AuthGuard element={<AppLayout />} />,
		children: [
			{
				index: true,
				path: "/receipts",
				element: <ReceiptsPage />
			}
		]
	}
];

const Router: FunctionComponent = () => <RouterProvider router={createBrowserRouter(routes)} />

export default Router;
