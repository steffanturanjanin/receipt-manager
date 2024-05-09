import { FunctionComponent, ReactElement, ReactNode } from "react";
import { Box, BoxProps, styled } from "@mui/material";

interface MainContentProps {
	children: ReactNode;
}

const MainContentContainer = styled(Box)<BoxProps>({
	width: "100%",
	maxWidth: "600px",
	marginTop: "64px",
	padding: "1.5rem",
	marginLeft: "auto",
	marginRight: "auto",
})

const MainContent: FunctionComponent<MainContentProps> = ({ children }): ReactElement => {
	return (
		<MainContentContainer component="main">
			{children}
		</MainContentContainer>
	)
}

export default MainContent;
