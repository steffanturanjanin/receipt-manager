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
	empty?: ReactNode;
}

const StatItem: FunctionComponent<StatItemProps> = ({ title, children, empty }): ReactElement => {
	return (
		<Stack direction="column" gap="1rem">
			<Typography variant="h5" component="h2">{title}</Typography>
			{ empty ? empty :
				<Paper>{children}</Paper>
			}
		</Stack>
	)
}

export default StatItem;
