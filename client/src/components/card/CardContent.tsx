import { Box } from "@mui/material";
import { FunctionComponent, ReactElement, ReactNode } from "react";

interface CardContentProps {
	children: ReactNode;
}

const CardContent: FunctionComponent<CardContentProps> = ({ children }): ReactElement => {
	return (
		<Box padding="1rem">{children}</Box>
	)
}

export default CardContent;
