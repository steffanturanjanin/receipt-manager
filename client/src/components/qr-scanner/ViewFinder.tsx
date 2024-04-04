import { FunctionComponent, ReactElement, ReactNode } from "react";
import { Box, BoxProps, Button, ButtonProps, Stack, StackProps, Typography, styled } from "@mui/material";
import CloseIcon from '@mui/icons-material/Close';

// Components
const QrScannerContainer = styled(Box)<BoxProps>({
	width:"100%",
	height:"100%",
	padding: "1.5rem",
	position: "relative",
});

const QrScannerContent = styled(Box)<BoxProps>({
	position: "relative",
	width: "100%",
	height:  "100%",
	alignItems: "center",
	justifyContent: "center",
	zIndex: 10,
	color: "#fff",
});

const QrScannerCloseButton = styled(Button)<ButtonProps>({
	zIndex: 10,
	padding: 0,
	minWidth: 0,
	color: "#fff",
});

const QrScannerCameraContainer = styled(Stack)<StackProps>({
	width: "85%",
	position: "absolute",
	top: "50%",
	left: "50%",
	transform: "translate(-50%, -50%)",
})

const QrScannerCamera = styled(Box)<BoxProps>({
	position: "relative",
	width: "100%",
	minWidth: "200px",
	maxWidth: "400px",
	aspectRatio: 1,
});

const QrScannerCameraOuter = styled(Box)<BoxProps>({
	position: "absolute",
	inset: 0,
	zIndex: 30,
	background:
		`linear-gradient(to right, white 4px, transparent 4px) 0 0,
		linear-gradient(to right, white 4px, transparent 4px) 0 100%,
		linear-gradient(to left, white 4px, transparent 4px) 100% 0,
		linear-gradient(to left, white 4px, transparent 4px) 100% 100%,
		linear-gradient(to bottom, white 4px, transparent 4px) 0 0,
		linear-gradient(to bottom, white 4px, transparent 4px) 100% 0,
		linear-gradient(to top, white 4px, transparent 4px) 0 100%,
		linear-gradient(to top, white 4px, transparent 4px) 100% 100%`,
	backgroundRepeat: "no-repeat",
	backgroundSize: "1.5rem 1.5rem",
	padding: "1.5rem",
});

const QrScannerCameraInner = styled(Box)<BoxProps>({
	position: "absolute",
	display: "flex",
	alignItems: "center",
	inset: "1.5rem",
	borderRadius: "1rem",
	background: "none",
	boxShadow: "0 0 0 1600px rgba(0, 0, 0, 0.65)",
})

interface ViewFinderProps {
	onClose: () => void;
	notification?: ReactNode;
}

const ViewFinder: FunctionComponent<ViewFinderProps> = ({ onClose, notification }): ReactElement => {
	return (
		<QrScannerContainer>
			<QrScannerContent>
				<Stack direction="row" justifyContent="end">
					<QrScannerCloseButton onClick={onClose}>
						<CloseIcon fontSize="large" sx={{fontSize: "48px"}}/>
					</QrScannerCloseButton>
				</Stack>
				<Stack direction="row" alignItems="top" justifyContent="center">
					<Typography variant="h4" component="p" zIndex={1} textAlign="center" marginY="1.5rem">
						Skeniraj QR kod sa računa
					</Typography>
				</Stack>
				<QrScannerCameraContainer alignItems="center" justifyContent="center">
					<QrScannerCamera>
						<QrScannerCameraOuter />
						<QrScannerCameraInner>
							{notification}
						</QrScannerCameraInner>
						</QrScannerCamera>
				</QrScannerCameraContainer>
			</QrScannerContent>
		</QrScannerContainer>
	)
}

export default ViewFinder;
