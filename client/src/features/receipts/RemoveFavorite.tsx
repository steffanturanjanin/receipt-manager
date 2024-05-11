import { FormEvent, FunctionComponent, ReactElement } from "react";
import { useMutation } from "react-query";
import {
	Box,
	Button,
	Dialog,
	DialogActions,
	DialogContent,
	DialogContentText,
	DialogProps,
	DialogTitle
} from "@mui/material";
import { setFavorite } from "../../api/receipts";
import LoadingButton from "../../components/LoadingButton";

const DialogPaperProps: DialogProps["PaperProps"] = {
	sx: {
		width: "600px",
		maxWidth: "100%",
	}
}

interface RemoveFavoriteProps {
	receiptId: number;
	open: boolean;
	onClose: () => void;
	onRemovedFromFavorites: () => void;
}

const RemoveFavorite: FunctionComponent<RemoveFavoriteProps> = ({
	open,
	onClose,
	receiptId,
	onRemovedFromFavorites,
}): ReactElement => {
	const { isLoading, mutate } = useMutation({
		mutationFn: () => setFavorite(receiptId, { isFavorite: false }),
		onSuccess: () => {
			onRemovedFromFavorites();
		}
	});

	const onSubmit = (event: FormEvent<HTMLFormElement>) => {
		event.preventDefault();
		mutate();
	}

	return (
		<Dialog open={open} PaperProps={DialogPaperProps} onClose={onClose}>
			<DialogTitle>Uklanjanje računa iz <strong>omiljenih</strong></DialogTitle>
			<Box component="form" onSubmit={onSubmit}>
				<DialogContent>
					<DialogContentText>
						Da li ste sigurni da želite da uklonite ovaj račun iz <strong>omiljenih?</strong>
					</DialogContentText>
				</DialogContent>
				<DialogActions sx={{ padding: "1rem 1.5rem" }}>
					<Button size="large" onClick={onClose}>Otkaži</Button>
					<LoadingButton type="submit" variant="contained" size="large" loading={isLoading}>
						Potvrdi
					</LoadingButton>
				</DialogActions>
			</Box>
		</Dialog>
	)
}

export default RemoveFavorite;
