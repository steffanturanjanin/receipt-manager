import { Paper, Stack, StackProps, styled } from "@mui/material";
import { FunctionComponent, ReactElement, ReactNode } from "react";

const CardItemContainer = styled(Stack)<StackProps>({
	borderRadius: "0.5rem",
	boxShadow: "#959da533 0 8px 24px",
	overflow: "hidden",
});

interface CardItemProps {
	children: ReactNode;
}

const CardItem: FunctionComponent<CardItemProps> = ({ children }): ReactElement => {
	return (
		<CardItemContainer component={Paper}>
			{children}
		</CardItemContainer>
	)
}

export default CardItem;
