import { AppBar, Box, BoxProps, Stack, Toolbar, Typography, styled, Backdrop as MuiBackdrop, CircularProgress } from "@mui/material";
import { Fragment, FunctionComponent, ReactElement, ReactNode } from "react";

interface PageLayoutProps {
	title: string;
	headerPrefix?: ReactNode;
	headerSuffix?: ReactNode;
	showBackdrop?: boolean;
	children: ReactNode;
}

const MainContent = styled(Box)<BoxProps>({
	width: "100%",
	maxWidth: "600px",
	marginTop: "64px",
	padding: "1.5rem",
	marginLeft: "auto",
	marginRight: "auto",
})

interface BackdropProps {
	open: boolean;
}

const Backdrop: FunctionComponent<BackdropProps> = ({ open }) => (
	<MuiBackdrop
		sx={{ color: '#fff', zIndex: (theme) => theme.zIndex.drawer + 1 }}
		open={open}
	>
		<CircularProgress color="inherit" />
	</MuiBackdrop>
)

const PageLayout: FunctionComponent<PageLayoutProps> = ({ title, showBackdrop, headerPrefix, headerSuffix, children }): ReactElement => {
	return (
		<Fragment>
			<AppBar position="fixed" sx={{ display: "flex", alignItems: "center", backgroundColor: "#fff", color: "black" }}>
				<Toolbar sx={{ justifyContent: "center", alignItems: "center", width: "100%", maxWidth: "600px" }}>
					<Stack direction="row" justifyContent="space-between" alignItems="center" sx={{ width: "100%", gap: "16px" }}>
						<Box component="div" sx={{ minWidth: headerPrefix ? "0" : "64px"}}>{headerPrefix && headerPrefix}</Box>
						<Typography variant="h5" component="h1">{title}</Typography>
						<Box component="div" sx={{ minWidth: headerSuffix ? "0" : "64px" }}>{headerSuffix && headerSuffix}</Box>
					</Stack>
				</Toolbar>
			</AppBar>
			<MainContent component="main">
				{children}
			</MainContent>
			<Backdrop open={!!showBackdrop} />
		</Fragment>
	);
}

export default PageLayout;
