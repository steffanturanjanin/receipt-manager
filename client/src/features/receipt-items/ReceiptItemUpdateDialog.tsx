import { Box, Button, Dialog, DialogActions, DialogContent, DialogContentText, DialogProps, DialogTitle, FormControl, InputLabel, MenuItem, Select, SelectChangeEvent, Stack } from "@mui/material";
import { FormEvent, FunctionComponent, ReactElement } from "react";
import LoadingButton from "../../components/LoadingButton";
import { useMutation, useQuery } from "react-query";
import { getCategories } from "../../api/categories";
import { updateReceiptItem as apiUpdateReceiptItem } from "../../api/receipt-items";

const DialogPaperProps: DialogProps["PaperProps"] = {
	sx: {
		width: "600px",
		maxWidth: "100%",
	}
}

interface UpdateReceiptItemProps {
	id: number;
	request: UpdateReceiptItemRequest;
}

interface ReceiptItemUpdateDialogProps {
	open: boolean;
	onClose: () => void;
	receiptItem?: SingleReceiptReceiptItem;
	updateReceiptItem: (value: SingleReceiptReceiptItem) => void;
	onSubmitted: () => void;
}

const ReceiptItemUpdateDialog: FunctionComponent<ReceiptItemUpdateDialogProps> = ({
	open,
	onClose,
	receiptItem,
	updateReceiptItem,
	onSubmitted,
}): ReactElement => {
	// Get categories
	const { data: categories } = useQuery({
		queryKey: ["categories"],
		queryFn: () => getCategories(),
		keepPreviousData: true,
	});

	// Perform mutation
	const { isLoading, mutate } = useMutation({
		mutationFn: ({ id, request }: UpdateReceiptItemProps) => apiUpdateReceiptItem(id, request),
		onSuccess: () => {
			onSubmitted();
			onClose();
		}
	});

	const onSubmit = (event: FormEvent<HTMLFormElement>) => {
		event.preventDefault();
		const { id } = receiptItem || {};
		const { id: categoryId } =  receiptItem?.category || {};

		if (id && categoryId) {
			mutate({ id, request: { categoryId } })
		}
	}

	const onCategoryChange = (event: SelectChangeEvent<number>) => {
		if (receiptItem) {
			const updateCategory = {...receiptItem?.category, id: event.target.value } as SingleReceiptCategory;
			updateReceiptItem({ ...receiptItem, category: updateCategory });
		}
	}

	return (
		<Dialog open={open} PaperProps={DialogPaperProps} onClose={onClose}>
			<DialogTitle>{receiptItem?.name}</DialogTitle>
			<Box component="form" onSubmit={onSubmit}>
				<DialogContent>
					<Stack direction="column" gap="1rem">
						<DialogContentText>
							Izaberite novu kategoriju iz padajućeg menija
						</DialogContentText>
						<FormControl fullWidth>
							<InputLabel id="categoryId">Kategorija</InputLabel>
							<Select
								labelId="categoryIdLabelId"
								id="categoryId"
								value={receiptItem?.category?.id || ""}
								label="Kategorija"
								onChange={onCategoryChange}
							>
								{categories?.map(category =>
									<MenuItem key={category.id} value={category.id}>{category.name}</MenuItem>
								)}
							</Select>
						</FormControl>
					</Stack>
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

export default ReceiptItemUpdateDialog;
