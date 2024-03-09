import { Fragment } from 'react';
import { Navigate, RouteObject } from 'react-router-dom';

import Login from "../views/pages/auth/Login";

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
				element: <Fragment>Register</Fragment>,
			}
		]
	}
];

export default auth;
