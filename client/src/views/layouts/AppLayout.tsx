import { FunctionComponent, ReactElement } from "react";
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

const NavigationAction = styled(BottomNavigationAction)<BottomNavigationActionProps & NavLinkProps>({
	".MuiBottomNavigationAction-label": { fontSize: "1rem" },
})

const AppLayout: FunctionComponent<AppLayoutProps> = ({ children }): ReactElement => {
	return (
		<AppLayoutContainer>
			<AppLayoutContent component="main">
				{children}
			</AppLayoutContent>

			<Box component={Paper} sx={BottomNavigationContainerStyle} elevation={3}>
				<BottomNavigation showLabels>
					<NavigationAction
						component={NavLink}
						to="/"
						label="Receipts"
						icon={<ReceiptIcon fontSize="large" />}
					/>
					<NavigationAction
						component={NavLink}
						to="/scan"
						label="Scan"
						icon={<CameraIcon fontSize="large" />}
					/>
				</BottomNavigation>
			</Box>
		</AppLayoutContainer>
	)
}

export default AppLayout;
