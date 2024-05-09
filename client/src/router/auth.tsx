import { Navigate, Outlet, RouteObject } from 'react-router-dom';
import Login from "../views/pages/auth/Login";
import Register from "../views/pages/auth/Register";
import GuestGuard from './guards/GuestGuard';

const auth: RouteObject[] = [
	{
		path: "/auth",
		element: <GuestGuard element={<Outlet />} />,
		children: [
			{
				path: "",
				element: <Navigate to="login" replace />
			},
			{
				index: true,
				path: "login",
				element: <Login />
			},
			{
				path: "register",
				element: <Register />,
			}
		]
	}
];

export default auth;
