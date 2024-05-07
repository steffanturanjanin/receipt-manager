import { Paper, Stack, StackProps, SxProps, styled } from "@mui/material";
import { FunctionComponent, ReactElement, ReactNode } from "react";

const CardContainer = styled(Stack)<StackProps>({
	borderRadius: "0.5rem",
	boxShadow: "#959da533 0 8px 24px",
	overflow: "hidden",
});

interface CardProps {
	children: ReactNode;
	sx?: SxProps;
}

const Card: FunctionComponent<CardProps> = ({ children, sx }): ReactElement => {
	return (
		<CardContainer component={Paper} sx={sx}>
			{children}
		</CardContainer>
	)
}

export default Card;
