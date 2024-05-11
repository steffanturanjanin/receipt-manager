import { Button, ButtonProps, styled } from "@mui/material";
import { FunctionComponent, ReactElement } from "react";
import DeleteIcon from '@mui/icons-material/Delete';
import StarEmptyIcon from "@mui/icons-material/StarBorderOutlined";
import StarFullIcon from "@mui/icons-material/Star";

const ActionButtonStyled = styled(Button)<ButtonProps>(({ theme }) => ({
	minWidth: "auto",
	border: `1px solid ${theme.palette.grey[400]}`,
	"&:hover": {
		borderColor: theme.palette.grey[600],
	}
}));

export const DeleteActionButton: FunctionComponent<ButtonProps> = ({ onClick }): ReactElement => {
	return (
		<ActionButtonStyled onClick={onClick} sx={{ color: "error.light" }}>
			<DeleteIcon />
		</ActionButtonStyled>
	)
}

interface FavoriteActionButtonProps extends ButtonProps {
	isFavorite: boolean;
}

export const FavoriteActionButton: FunctionComponent<FavoriteActionButtonProps> = ({ isFavorite, onClick }): ReactElement => {
	return (
		<ActionButtonStyled onClick={onClick} sx={{ color: "warning.light" }}>
			{ isFavorite ? <StarFullIcon /> : <StarEmptyIcon /> }
		</ActionButtonStyled>
	)
}

