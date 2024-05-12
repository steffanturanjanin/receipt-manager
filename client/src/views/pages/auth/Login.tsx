import { AxiosError } from 'axios';
import { FunctionComponent, FormEvent, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useMutation } from 'react-query';
import { TextField, Box, Typography, Container } from '@mui/material';
import { SxProps } from '@mui/material/styles';
import { login } from '../../../api/auth';
import LoadingButton from '../../../components/LoadingButton';
import { FormErrors, FormFieldsTranslator } from '../../../shared/types/errors';

interface LoginForm {
	email: string,
	password: string,
}

const LoginFormFieldsTranslation: FormFieldsTranslator<LoginForm> = {
	email: "Email",
	password: "Password",
}

const DefaultLoginForm: LoginForm = {
	email: '',
	password: '',
}

const Copyright: FunctionComponent<{ sx: SxProps }> = ({ sx }) => (
	<Typography variant="body2" color="text.secondary" align="center" sx={sx}>
		{'Copyright © '} Receipt manager {' '}
		{new Date().getFullYear()} {'.'}
	</Typography>
);

const Login: FunctionComponent = () => {
	const [loginForm, setLoginForm] = useState<LoginForm>(DefaultLoginForm);
	const [loginFormErrors, setLoginFormErrors] = useState<FormErrors<LoginForm>>(DefaultLoginForm);

	const navigate = useNavigate();

	const { mutate, isLoading } = useMutation({
		mutationFn: async (request: LoginRequest) => await login(request),
		onSuccess: (response: AuthResponse) => {
			// Save auth token to local storage
			// Redirect to Home
			localStorage.setItem("auth", JSON.stringify(response));
			navigate("/");
		},
		onError: (error: AxiosError<ValidationError<LoginForm>>) => {
			if (error.status = 422) {
				setLoginFormErrors({
					...loginFormErrors,
					...error.response?.data.errors,
					...Object.fromEntries(
						Object.entries(error.response?.data.errors || {}).map(
							([field, value]) => ([field, `${LoginFormFieldsTranslation[field as keyof LoginForm]} ${value}` ])
					)),
				});
			}
		}
	})

  const handleSubmit = (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
		setLoginFormErrors(DefaultLoginForm);
		mutate(loginForm)
  };

  return (
		<Box component="main" width={"100vw"}>
			<Container maxWidth="xs">
				<Box
					sx={{
						marginTop: 8,
						display: 'flex',
						flexDirection: 'column',
						alignItems: 'center',
					}}
				>
					<Typography component="h1" variant="h5">
						Login
					</Typography>
					<Box component="form" onSubmit={handleSubmit} noValidate sx={{ mt: 1 }}>
						<TextField
							margin="normal"
							required
							fullWidth
							type="email"
							id="email"
							label="Email Address"
							name="email"
							autoComplete="email"
							value={loginForm.email}
							onChange={(event) => setLoginForm({ ...loginForm, email: event.target.value })}
							error={!!loginFormErrors.email}
							helperText={loginFormErrors.email}
						/>
						<TextField
							margin="normal"
							required
							fullWidth
							name="password"
							label="Password"
							type="password"
							id="password"
							autoComplete="current-password"
							value={loginForm.password}
							onChange={(event) => setLoginForm({ ...loginForm, password: event.target.value })}
							error={!!loginFormErrors.password}
							helperText={loginFormErrors.password}
						/>
						<LoadingButton
							type="submit"
							variant="contained"
							fullWidth sx={{ mt: 3, mb: 2 }}
							loading={isLoading}
						>
							Login
						</LoadingButton>
					</Box>
				</Box>
				<Copyright sx={{ mt: 8, mb: 4 }} />
			</Container>
		</Box>
  );
}

export default Login;
