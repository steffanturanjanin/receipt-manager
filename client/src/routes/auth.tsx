import { Navigate, RouteObject } from 'react-router-dom';
import Login from "../views/pages/auth/Login";
import Register from "../views/pages/auth/Register";

const auth: RouteObject[] = [
	{
		path: '/auth',
		children: [
			{
				path: '',
				element: <Navigate to="login" />
			},
			{
				index: true,
				path: 'login',
				element: <Login />
			},
			{
				path: 'register',
				element: <Register />,
			}
		]
	}
];

export default auth;
