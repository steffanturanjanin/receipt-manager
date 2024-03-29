import { FunctionComponent, ReactElement } from "react";
import { Box, BoxProps, Button, ButtonProps, Stack, StackProps, styled } from "@mui/material";
import CloseIcon from '@mui/icons-material/Close';

// Components
const QrScannerContainer = styled(Box)<BoxProps>({
	width:"100%",
	height:"100%",
	padding: "1.5rem"
});

const QrScannerContent = styled(Stack)<StackProps>({
	position: "relative",
	width: "100%",
	height:  "100%",
	display: "flex",
	alignItems: "center",
	justifyContent: "center",
	zIndex: "9999",
	color: "#fff",
});

const QrScannerCloseButton = styled(Button)<ButtonProps>({
	position: "absolute",
	right: 0,
	top: 0,
	zIndex: 1,
	padding: 0,
	minWidth: 0,
	color: "#fff",
});

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
	zIndex: 2,
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
	inset: "1.5rem",
	borderRadius: "1rem",
	background: "none",
	boxShadow: "0 0 0 1600px rgba(0, 0, 0, 0.65)",
})

interface ViewFinderProps {
	onClose: () => void;
}

const ViewFinder: FunctionComponent<ViewFinderProps> = ({ onClose }): ReactElement => {
	return (
		<QrScannerContainer>
			<QrScannerContent>
				<QrScannerCloseButton onClick={onClose}>
					<CloseIcon fontSize="large" />
				</QrScannerCloseButton>
				<QrScannerCamera>
					<QrScannerCameraOuter />
					<QrScannerCameraInner />
				</QrScannerCamera>
			</QrScannerContent>
		</QrScannerContainer>
	)
}

export default ViewFinder;
