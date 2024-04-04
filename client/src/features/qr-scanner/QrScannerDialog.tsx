import { FunctionComponent, ReactElement, ReactNode, forwardRef, useEffect, useState } from "react";
import { useMutation } from "react-query";
import { OnResultFunction } from "react-qr-reader";
import { AxiosError } from "axios";
import { Box, BoxProps, Button, Dialog, Slide, Stack, StackProps, Typography, styled } from "@mui/material";
import CheckCircleOutlineIcon from '@mui/icons-material/CheckCircleOutline';
import HighlightOffIcon from '@mui/icons-material/HighlightOff';
import CircularProgress from '@mui/material/CircularProgress';
import { TransitionProps } from "@mui/material/transitions";
import QrScanner from "../../components/qr-scanner";
import { createReceipt } from "../../api/receipts";

const BottomTopTransition = forwardRef(
	function Transition(
		props: TransitionProps & {
			children: React.ReactElement;
		},
		ref: React.Ref<unknown>,
	) {
		return <Slide direction="up" ref={ref} {...props} />;
	}
);

const QrScannerDialogContainer = styled(Box)<BoxProps>({
	position: "fixed",
	inset: "0",
	width: "100%",
	height: "100%",
	backgroundColor: "#000",
	overflow: "hidden",
})

interface QrScannerDialogProps {
	open: boolean;
	onClose: () => void;
}

const QrScanNotification = styled(Stack)<StackProps>({
	display: "flex",
	flexDirection: "column",
	alignItems: "center",
	textAlign: "center",
	zIndex: 30,
	gap: "0.5rem",
	padding: "1rem",
});

interface ReceiptScanNotificationProps {
	message: ReactNode;
	action?: () => void;
}

const ReceiptScanned: FunctionComponent<ReceiptScanNotificationProps> = ({ message, action }) => (
	<QrScanNotification color="success.light">
		<CheckCircleOutlineIcon sx={{ width: "60px", height: "60px" }} />
		<Typography variant="body1" component="p" color="success.light">
			{message}
		</Typography>
		<Button onClick={action} variant="outlined" size="large" sx={{ color: "success.light", borderColor: "success.light" }}>
			Skeniraj ponovo
		</Button>
	</QrScanNotification>
);

const ReceiptScanError: FunctionComponent<ReceiptScanNotificationProps> = ({ message, action }) => (
	<QrScanNotification color="error.light">
		<HighlightOffIcon sx={{ width: "60px", height: "60px" }} />
		<Typography variant="body1" component="p" color="error.light">
			{message}
		</Typography>
		<Button onClick={action} variant="outlined" size="large" sx={{ color: "error.light", borderColor: "error.light", ":hover": {borderColor: "error.dark"} }}>
			Pokušaj ponovo
		</Button>
	</QrScanNotification>
);

const LoadingSpinner: FunctionComponent<ReceiptScanNotificationProps> = ({ message }) => (
	<QrScanNotification>
		<CircularProgress sx={{ color: "#fff"}} size="60px" />
		<Typography variant="h6" component="p" color="#fff">
			{message}
		</Typography>
	</QrScanNotification>
)

const QrScannerDialog: FunctionComponent<QrScannerDialogProps> = ({ open, onClose }): ReactElement => {
	const [isScanned, setIsScanned] = useState<boolean>(false);
	const [notification, setNotification] = useState<ReactNode | null>(null);

	// When open dialog, make sure the state is restarted to default values
	useEffect(() => {
		if (open) {
			setIsScanned(false);
			setNotification(null);
		}
	}, [open]);

	const { isLoading, mutate: onQrCodeScanned } = useMutation({
		mutationFn: (request: CreateReceiptRequest) => createReceipt(request),
		onSuccess: ({ message }) => {
			setNotification(
				<ReceiptScanned
					message={message}
					action={() => { setIsScanned(false); setNotification(null) }}
				/>
			);
		},
		onError: (error: AxiosError) => {
			if (error.response?.status === 400) {
				const badRequestError = error as AxiosError<BadRequestError>;
				const { error: errorMessage } = badRequestError.response?.data || {};

				setNotification(
					<ReceiptScanError
						message={errorMessage}
						action={() => { setIsScanned(false); setNotification(null) }}
					/>
				);
			}

			if (error.response?.status === 422) {
				const validationError = error as AxiosError<ValidationError<CreateReceiptRequest>>;
				const { errors } = validationError.response?.data || {};
				const { url: urlErrorMessage } = errors || {};

				if (urlErrorMessage) {
					setNotification(
						<ReceiptScanError
							message={urlErrorMessage}
							action={() => { setIsScanned(false); setNotification(null) }}
						/>
					);
				}
			}

			const { error: errorMessage } = (error as AxiosError<BadRequestError>).response?.data || {};
			if (errorMessage) {
				setNotification(
					<ReceiptScanError
						message={errorMessage}
						action={() => { setIsScanned(false); setNotification(null) }}
					/>
				);
			}
		}
	});

	// Set notification when `isLoading` is true
	useEffect(() => {
		if (isLoading) {
			setNotification(<LoadingSpinner message="Račun je skeniran. Obrada je u toku..." />);
		}
	}, [isLoading]);

	const onResult: OnResultFunction = (result) => {
		if (result && !isScanned) {
			setIsScanned(prevIsScanned => {
				if (!prevIsScanned) {
					onQrCodeScanned({ url: result.getText() });
				}
				return true;
			});
		}
	}

	return (
		<Dialog fullScreen open={open} TransitionComponent={BottomTopTransition}>
			<QrScannerDialogContainer>
				{open &&
					<QrScanner
						notification={notification}
						onResult={onResult}
						onScanStop={onClose}
					/>
				}
			</QrScannerDialogContainer>
		</Dialog>
	);
}

export default QrScannerDialog;
