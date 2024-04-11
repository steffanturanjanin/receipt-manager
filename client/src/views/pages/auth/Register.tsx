import { AxiosError } from 'axios';
import { FunctionComponent, FormEvent, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useMutation } from 'react-query';
import { TextField, Box, Typography, Container } from '@mui/material';
import { SxProps } from '@mui/material/styles';
import { register } from '../../../api/auth';
import { AuthResponse, RegisterRequest } from '../../../api/auth/types';
import LoadingButton from '../../../components/LoadingButton';
import { FormFieldsTranslator, FormErrors } from '../../../shared/types/errors';

interface RegisterForm {
	firstName: string,
	lastName: string,
	email: string,
	password: string,
}

const RegisterFormFieldsTranslation: FormFieldsTranslator<RegisterForm> = {
	firstName: "First name",
	lastName: "Last name",
	email: "Email",
	password: "Password",
}

const DefaultRegisterForm: RegisterForm = {
	firstName: '',
	lastName: '',
	email: '',
	password: '',
}

const Copyright: FunctionComponent<{ sx: SxProps }> = ({ sx }) => (
	<Typography variant="body2" color="text.secondary" align="center" sx={sx}>
		{'Copyright © '} Receipt manager {' '}
		{new Date().getFullYear()} {'.'}
	</Typography>
);

const Register: FunctionComponent = () => {
	const [registerForm, setRegisterForm] = useState<RegisterForm>(DefaultRegisterForm);
	const [registerFormErrors, setRegisterFormErrors] = useState<FormErrors<RegisterForm>>(DefaultRegisterForm);

	const navigate = useNavigate();

	const { mutate, isLoading } = useMutation({
		mutationFn: async (request: RegisterRequest) => await register(request),
		onSuccess: (response: AuthResponse) => {
			// Save auth token to local storage
			// Redirect to Home
			localStorage.setItem("auth", JSON.stringify(response));
			navigate("/");
		},
		onError: (error: AxiosError<ValidationError<RegisterForm>>) => {
			if (error.status = 422) {
				setRegisterFormErrors({
					...registerFormErrors,
					...error.response?.data.errors,
					...Object.fromEntries(
						Object.entries(error.response?.data.errors || {}).map(
							([field, value]) => ([field, `${RegisterFormFieldsTranslation[field as keyof RegisterForm]} ${value}` ])
					)),
				});
			}
		}
	})

  const handleSubmit = (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
		setRegisterFormErrors(DefaultRegisterForm);
		mutate(registerForm)
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
						Register
					</Typography>
					<Box component="form" onSubmit={handleSubmit} noValidate sx={{ mt: 1 }}>
					<TextField
							margin="normal"
							required
							fullWidth
							id="firstName"
							label="First name"
							name="firstName"
							autoFocus
							value={registerForm.firstName}
							onChange={(event) => setRegisterForm({ ...registerForm, firstName: event.target.value })}
							error={!!registerFormErrors.firstName}
							helperText={registerFormErrors.firstName}
						/>
						<TextField
							margin="normal"
							required
							fullWidth
							id="last_name"
							name="lastName"
							label="Last name"
							value={registerForm.lastName}
							onChange={(event) => setRegisterForm({ ...registerForm, lastName: event.target.value })}
							error={!!registerFormErrors.lastName}
							helperText={registerFormErrors.lastName}
						/>
						<TextField
							margin="normal"
							required
							fullWidth
							type="email"
							id="email"
							label="Email Address"
							name="email"
							autoComplete="email"
							value={registerForm.email}
							onChange={(event) => setRegisterForm({ ...registerForm, email: event.target.value })}
							error={!!registerFormErrors.email}
							helperText={registerFormErrors.email}
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
							value={registerForm.password}
							onChange={(event) => setRegisterForm({ ...registerForm, password: event.target.value })}
							error={!!registerFormErrors.password}
							helperText={registerFormErrors.password}
						/>
						<LoadingButton
							type="submit"
							variant="contained"
							fullWidth sx={{ mt: 3, mb: 2 }}
							loading={isLoading}
						>
							Register
						</LoadingButton>
					</Box>
				</Box>
				<Copyright sx={{ mt: 8, mb: 4 }} />
			</Container>
		</Box>
  );
}

export default Register;
