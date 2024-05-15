import { FormEvent, FunctionComponent, ReactElement } from "react";
import {
	Box,
	Button,
	Dialog,
	DialogActions,
	DialogContent,
	DialogContentText,
	DialogProps,
	DialogTitle,
} from "@mui/material";
import LoadingButton from "../../components/LoadingButton";
import { useMutation } from "react-query";
import { deleteReceipt } from "../../api/receipts";
import dayjs from "dayjs";
import { useNavigate } from "react-router-dom";

const DialogPaperProps: DialogProps["PaperProps"] = {
	sx: {
		width: "600px",
		maxWidth: "100%",
	}
}

interface DeleteReceiptProps {
	open: boolean;
	onClose: () => void;
	receipt?: SingleReceipt;
}

const DeleteReceipt: FunctionComponent<DeleteReceiptProps> = ({ open, onClose, receipt }): ReactElement => {
	const navigate = useNavigate();

	const { isLoading, mutate } = useMutation({
		mutationFn: (id: number) => deleteReceipt(id),
		onSuccess: () => {
			if (receipt?.date) {
				const month = dayjs(receipt.date).format("YYYY-MM");
				navigate(`/receipts?${month}`);
			}
		}
	});

	const onSubmit = (event: FormEvent<HTMLFormElement>) => {
		event.preventDefault();
		if (receipt?.id) {
			mutate(receipt.id);
		}
	}

	return (
		<Dialog open={open} PaperProps={DialogPaperProps} onClose={onClose}>
			<DialogTitle>Račun {receipt?.pfrNumber}</DialogTitle>
			<Box component="form" onSubmit={onSubmit}>
				<DialogContent>
					<DialogContentText >
						Da li ste sigurni da želite da obrišite ovaj račun? <br />
						Ova akcija je nepovratna i ne može biti poništena.
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
	);
}

export default DeleteReceipt;
