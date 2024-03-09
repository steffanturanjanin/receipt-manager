import { FunctionComponent, FormEvent } from 'react';
import { Button, TextField, Link, Box, Typography, Container } from '@mui/material';
import { SxProps } from '@mui/material/styles';

const Copyright: FunctionComponent<{ sx: SxProps }> = ({ sx }) => (
	<Typography variant="body2" color="text.secondary" align="center" sx={sx}>
		{'Copyright © '}
		<Link color="inherit" href="/">Receipt manager</Link>
		{' '} {new Date().getFullYear()} {'.'}
	</Typography>
);

const Login: FunctionComponent = () => {
  const handleSubmit = (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    const data = new FormData(event.currentTarget);
    console.log({
      email: data.get('email'),
      password: data.get('password'),
    });
  };

  return (
		<Container component="main" maxWidth="xs">
			<Box
				sx={{
					marginTop: 8,
					display: 'flex',
					flexDirection: 'column',
					alignItems: 'center',
				}}
			>
				<Typography component="h1" variant="h5">
					Log in
				</Typography>
				<Box component="form" onSubmit={handleSubmit} noValidate sx={{ mt: 1 }}>
					<TextField
						margin="normal"
						required
						fullWidth
						id="email"
						label="Email Address"
						name="email"
						autoComplete="email"
						autoFocus
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
					/>
					<Button
						type="submit"
						fullWidth
						variant="contained"
						sx={{ mt: 3, mb: 2 }}
					>
						Log In
					</Button>
				</Box>
			</Box>
			<Copyright sx={{ mt: 8, mb: 4 }} />
		</Container>
  );
}

export default Login;
