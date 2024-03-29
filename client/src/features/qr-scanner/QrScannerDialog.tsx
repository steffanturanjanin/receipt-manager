import { FunctionComponent, ReactElement, forwardRef, useEffect, useState } from "react";
import { Box, BoxProps, Dialog, Slide, styled } from "@mui/material";
import { TransitionProps } from "@mui/material/transitions";
import QrScanner from "../../components/qr-scanner";

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

const QrScannerDialog: FunctionComponent<QrScannerDialogProps> = ({ open, onClose }): ReactElement => {
	const [isScanning, setIsScanning] = useState<boolean>(true);
	const [scanResult, setScanResult] = useState<string>();

	const onScannerClose = () => {
		setIsScanning(false);
		onClose();
	}

	useEffect(() => {
		open && setIsScanning(true);
	}, [open]);

	return (
		<Dialog fullScreen open={open} TransitionComponent={BottomTopTransition}>
			<QrScannerDialogContainer>
				{isScanning &&
					<QrScanner
						onResult={(result, error) => {
							if (result) {
								console.log(result.getText())
								setScanResult(result.getText())
							}
						}}
						onScanStop={() => onScannerClose()}
					/>
				}
			</QrScannerDialogContainer>
		</Dialog>
	);
}

export default QrScannerDialog;
