import { Box, Paper, Stack, Typography } from "@mui/material";
import { FunctionComponent, ReactElement, ReactNode } from "react";

interface StatItemContentProps {
	children: ReactNode;
}

export const StatItemContent: FunctionComponent<StatItemContentProps> = ({ children }): ReactElement => {
	return (
		<Box sx={{ padding: "1rem" }}>
			{children}
		</Box>
	)
}

interface StatItemProps {
	title: ReactNode;
	children: ReactNode;
}

const StatItem: FunctionComponent<StatItemProps> = ({ title, children }): ReactElement => {
	return (
		<Stack direction="column" gap="1rem">
			<Typography variant="h5" component="h2">{title}</Typography>
			<Paper>{children}</Paper>
		</Stack>
	)
}

export default StatItem;
