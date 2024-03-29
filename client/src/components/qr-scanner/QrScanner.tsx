import { FunctionComponent, ReactElement } from "react";
import { QrReader, QrReaderProps } from "react-qr-reader";
import ViewFinder from "./ViewFinder";

type QrScannerProps = Pick<QrReaderProps, "onResult">  & {
	onScanStop: () => void;
}

const containerStyle = {
	margin: 'auto',
	height: '100%',
	width: '100%',
	maxWidth: "600px",
}

const videoContainerStyle = {
	paddingTop: 0,
	position: "relative",
	height: "100%",
	overflow: 'hidden',
}

const videoStyle = {
	width: '100%',
	objectFit: "cover",
}

const QrScanner: FunctionComponent<QrScannerProps> = ({
	onResult,
	onScanStop,
}): ReactElement => {
	return (
		<QrReader
			constraints={{ facingMode: "environment" }}
			scanDelay={500}
			onResult={onResult}
			containerStyle={containerStyle}
			videoContainerStyle={videoContainerStyle}
			videoStyle={videoStyle}
			ViewFinder={() => <ViewFinder onClose={onScanStop} />}
		/>
	)
}

export default QrScanner;
