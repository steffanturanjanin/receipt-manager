import { FunctionComponent, ReactElement, ReactNode } from "react";
import { QrReader, QrReaderProps } from "react-qr-reader";
import ViewFinder from "./ViewFinder";

type QrScannerProps = Pick<QrReaderProps, "onResult">  & {
	onScanStop: () => void;
	notification?: ReactNode;
}

const containerStyle = {
	margin: 'auto',
	height: '100%',
	width: '100%',
	maxWidth: "820px",
}

const videoContainerStyle = {
	paddingTop: 0,
	position: "relative",
	height: "100%",
	overflow: 'hidden',
}

const videoStyle = {
	width: '100%',
	//objectFit: "cover",
}

const QrScanner: FunctionComponent<QrScannerProps> = ({
	onResult,
	onScanStop,
	notification
}): ReactElement => {
	return (
		<QrReader
			constraints={{ facingMode: "environment" }}
			scanDelay={500}
			onResult={onResult}
			containerStyle={containerStyle}
			videoContainerStyle={videoContainerStyle}
			videoStyle={videoStyle}
			ViewFinder={() => <ViewFinder onClose={onScanStop} notification={notification} />}
		/>
	)
}

export default QrScanner;
