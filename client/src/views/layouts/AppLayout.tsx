import { FunctionComponent, ReactElement, useState } from "react";
import { NavLink, NavLinkProps } from "react-router-dom";
import {
	BottomNavigation,
	BottomNavigationAction,
	BottomNavigationActionProps,
	Box,
	Paper,
	Stack,
	StackProps,
	SxProps,
	styled
} from "@mui/material";
import ReceiptIcon from "@mui/icons-material/Receipt";
import CameraIcon from "@mui/icons-material/Camera";
import "./app-layout.scss";
import QrScannerDialog from "../../features/qr-scanner/QrScannerDialog";

interface AppLayoutProps {
	children: ReactElement;
}

const AppLayoutContainer = styled(Stack)<StackProps>({
	width: "100vw",
	height: "100vh",
});

const AppLayoutContent = styled(Stack)<StackProps>({
	padding: '1.5rem',
	marginBottom: '5.5rem',
	overflow: 'auto',
});

const BottomNavigationContainerStyle: SxProps = {
	position: 'fixed',
	bottom: 0,
	left: 0,
	right: 0,
	paddingY: '1rem',
}

const NavigationAction = styled(BottomNavigationAction)<BottomNavigationActionProps>({
	".MuiBottomNavigationAction-label": { fontSize: "1rem" },
})

const NavigationActionLink = styled(BottomNavigationAction)<BottomNavigationActionProps & NavLinkProps>({
	".MuiBottomNavigationAction-label": { fontSize: "1rem" },
})


const AppLayout: FunctionComponent<AppLayoutProps> = ({ children }): ReactElement => {
	const [receiptScannerOpened, setReceiptScannerOpened] = useState<boolean>(false);

	return (
		<AppLayoutContainer>
			<AppLayoutContent component="main">
				{children}
			</AppLayoutContent>

			<Box component={Paper} sx={BottomNavigationContainerStyle} elevation={3}>
				<BottomNavigation showLabels>
					<NavigationActionLink
						component={NavLink}
						to="/"
						label="Receipts"
						icon={<ReceiptIcon fontSize="large" />}
					/>
					<NavigationAction
						label="Scan"
						icon={<CameraIcon fontSize="large" />}
						onClick={() => setReceiptScannerOpened(true)}
					/>
				</BottomNavigation>
			</Box>

			<QrScannerDialog open={receiptScannerOpened} onClose={() => setReceiptScannerOpened(false)} />
		</AppLayoutContainer>
	)
}

export default AppLayout;
