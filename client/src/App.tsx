import { FunctionComponent } from 'react';
import { QueryClient, QueryClientProvider } from 'react-query';
import { createTheme } from '@mui/material';
import { ThemeProvider } from '@emotion/react';
import CssBaseline from '@mui/material/CssBaseline';
import Router from './router/Router';
import './App.scss'

const queryClient = new QueryClient({
	defaultOptions: {
		queries: {
			refetchOnWindowFocus: false,
		}
	}
})

const theme = createTheme();

const App: FunctionComponent = () => (
	<ThemeProvider theme={theme}>
		<CssBaseline />
		<QueryClientProvider client={queryClient}>
			<Router />
		</QueryClientProvider>
	</ThemeProvider>
);

export default App;
