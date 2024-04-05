import { AppBar, Box, BoxProps, Stack, Toolbar, Typography, styled } from "@mui/material";
import { Fragment, FunctionComponent, ReactElement, ReactNode } from "react";

interface PageLayoutProps {
	title: string;
	headerPrefix?: ReactNode;
	headerSuffix?: ReactNode;
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

const PageLayout: FunctionComponent<PageLayoutProps> = ({ title, headerPrefix, headerSuffix, children }): ReactElement => {
	return (
		<Fragment>
			<AppBar position="fixed" sx={{ display: "flex", alignItems: "center", backgroundColor: "#fff", color: "black" }}>
				<Toolbar sx={{ justifyContent: "center", alignItems: "center", width: "100%", maxWidth: "600px" }}>
					<Stack direction="row" justifyContent="space-between" alignItems="center" sx={{ width: "100%", gap: "16px" }}>
						<Box component="div" sx={{ minWidth: "64px"}}>{headerPrefix && headerPrefix}</Box>
						<Typography variant="h5" component="h1">{title}</Typography>
						<Box component="div" sx={{ minWidth: "64px" }}>{headerSuffix && headerSuffix}</Box>
					</Stack>
				</Toolbar>
			</AppBar>
			<MainContent component="main">
				{children}
			</MainContent>
		</Fragment>
	);
}

export default PageLayout;
