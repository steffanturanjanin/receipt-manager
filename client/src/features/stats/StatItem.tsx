import { Stack, Typography } from "@mui/material";
import { FunctionComponent, ReactElement, ReactNode } from "react";
import Card from "../../components/card/Card";

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
				<Card>{children}</Card>
			}
		</Stack>
	)
}

export default StatItem;
