import { Button, ButtonProps, CircularProgress } from '@mui/material';
import { FunctionComponent, ReactElement } from 'react';

export interface LoadingButtonProps extends ButtonProps {
	loaderPosition?: 'left' | 'right';
	loaderSize?: string | number;
	loading?: boolean;
	ariaLabel?: string;
}

const LoadingButton: FunctionComponent<LoadingButtonProps> = ({
	loaderPosition = 'right',
	loaderSize = '1rem',
	loading,
	ariaLabel,
	children,
	disabled,
	...rest
}): ReactElement => {
	const loader = (
		<CircularProgress
			className={`loading-button-spinner--${loaderPosition}`}
			sx={{
				marginLeft: loaderPosition === 'right' ? '0.625rem' : 0,
				marginRight: loaderPosition === 'left' ? '0.625rem' : 0,
			}}
			aria-busy={!!loading}
			aria-label={ariaLabel}
			role="progressbar"
			size={loaderSize}
			color="inherit"
		/>
	);

	return (
		<Button {...rest} disabled={loading || disabled}>
			{loading && loaderPosition === 'left' && loader}
			{children}
			{loading && loaderPosition === 'right' && loader}
		</Button>
	);
};

export default LoadingButton;
