import { FunctionComponent } from "react";
import { RouteObject, RouterProvider, createBrowserRouter } from "react-router-dom";
import auth from "./auth";
import AuthGuard from "./guards/AuthGuard";
import Home from "../views/pages/home/Home";
import AppLayout from "../views/layouts/AppLayout/AppLayout";
import ReceiptsPage from "../views/pages/receipts/ReceiptsPage";
import ShowReceiptPage from "../views/pages/receipts/ShowReceiptPage";

const routes: RouteObject[] = [
	...auth,
	{
		path: "",
		element: <AuthGuard element={<AppLayout />} />,
		children: [
			{
				index: true,
				path: "",
				element: <Home />
			},
			{
				path: "/receipts",
				element: <ReceiptsPage />
			},
			{
				path: "/receipts/:id",
				element: <ShowReceiptPage />
			}
		]
	}
];

const Router: FunctionComponent = () => <RouterProvider router={createBrowserRouter(routes)} />

export default Router;
