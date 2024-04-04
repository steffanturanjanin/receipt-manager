import { AppBar, Box, BoxProps, Stack, Toolbar, Typography, styled } from "@mui/material";
import { Fragment, FunctionComponent, ReactElement, ReactNode } from "react";

interface PageLayoutProps {
	title: string;
	headerPrefix?: ReactElement;
	headerSuffix?: ReactElement;
	children: ReactNode;
}

const MainContent = styled(Box)<BoxProps>({
	width: "100%",
	maxWidth: "1248px",
	marginTop: "64px",
	padding: "1.5rem",
	marginLeft: "auto",
	marginRight: "auto",
})

const PageLayout: FunctionComponent<PageLayoutProps> = ({ title, headerPrefix, headerSuffix, children }): ReactElement => {
	return (
		<Fragment>
			<AppBar position="fixed" color="transparent">
				<Toolbar sx={{ justifyContent: "center", alignItems: "center" }}>
					<Stack justifyContent="space-between">
						<Box component="div">{headerPrefix && headerPrefix}</Box>
						<Typography variant="h4" component="h1">{title}</Typography>
						<Box component="div">{headerSuffix && headerSuffix}</Box>
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
