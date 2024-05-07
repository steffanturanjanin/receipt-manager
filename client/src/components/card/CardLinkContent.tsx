import { styled } from "@mui/system";
import { FunctionComponent, ReactElement, ReactNode } from "react";
import { Link, LinkProps } from "react-router-dom";
import CardContent from "./CardContent";


const CardLinkContentContainer = styled(Link)<LinkProps>(({ theme }) => ({
	color: "inherit",
	textDecoration: "none",
	"&:hover": {
		backgroundColor: theme.palette.grey[50],
		boxShadow: "#959da533 0 8px 24px",
	}
}));

interface CardLinkContentProps extends LinkProps {
	children: ReactNode;
}

const CardLinkContent: FunctionComponent<CardLinkContentProps> = ({ children, ...rest }): ReactElement => {
	return (
		<CardLinkContentContainer {...rest}>
			<CardContent>
				{children}
			</CardContent>
		</CardLinkContentContainer>
	)
}

export default CardLinkContent;
