import { FunctionComponent, ReactElement } from "react";
import { Link } from "react-router-dom";
import { Link as MuiLink, Stack, Typography } from "@mui/material";
import MainContent from "../layouts/MainContent";
import NotFound from "../../components/errors/NotFound";

const NotFoundPage: FunctionComponent = (): ReactElement => {
	return (
		<MainContent>
			<Stack direction="column" gap="2rem" alignItems="center" justifyContent="center">
				<Typography variant="h4" textAlign="center">
					Strana koju ste tražili nije pronađena
				</Typography>
				<NotFound />
				<MuiLink component={Link} to={"/"}>Nazad na početnu stranu</MuiLink>
			</Stack>
		</MainContent>
	)
}

export default NotFoundPage;
