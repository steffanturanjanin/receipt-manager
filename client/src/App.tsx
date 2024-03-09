import { FunctionComponent } from 'react';
import Router from './routes/Router';

import { createTheme } from '@mui/material';
import { ThemeProvider } from '@emotion/react';
import CssBaseline from '@mui/material/CssBaseline';
import './App.css'

const theme = createTheme();

const App: FunctionComponent = () => (
	<ThemeProvider theme={theme}>
		<CssBaseline />
		<Router />
	</ThemeProvider>
);

export default App;
