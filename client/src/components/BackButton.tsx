import { styled } from "@mui/material";
import { FunctionComponent, ReactElement } from "react";
import ChevronLeftIcon from '@mui/icons-material/ChevronLeft';
import { useLocation, useNavigate } from "react-router-dom";

const BackButtonContainer = styled("button")(({ theme }) => ({
	display: "flex",
	alignItems: "center",
	background: "none",
	border: "none",
	cursor: "pointer",
	padding: 0,
	textDecoration: "none",
	fontSize: "18px",
	color: theme.palette.primary.dark,
}));

const BackButton: FunctionComponent = (): ReactElement => {
	const navigate = useNavigate();
  const location = useLocation();

	const back = () => {
    // If there's a previous location in the history, navigate back to it
    if (location.state && location.state.from) {
      navigate(location.state.from);
    } else {
      // If there's no previous location, navigate to a default route
      navigate('/');
    }
	}

	return (
		<BackButtonContainer onClick={back}>
			<ChevronLeftIcon /> Nazad
		</BackButtonContainer>
	)
}

export default BackButton;

