import { FunctionComponent, useState } from "react";
import { useQuery, useQueryClient } from "react-query";
import { getReceipt } from "../../../api/receipts";
import { Link, LinkProps, useParams } from "react-router-dom";
import PageLayout from "../../layouts/PageLayout/PageLayout";
import { Button, ButtonProps, Stack, StackProps, styled } from "@mui/material";
import ChevronLeftIcon from '@mui/icons-material/ChevronLeft';
import ReceiptPaymentOverview from "../../../features/receipts/show-receipt/ReceiptPaymentOverview";
import ReceiptItemsList from "../../../features/receipts/show-receipt/ReceiptItemsList";
import ReceiptDetails from "../../../features/receipts/show-receipt/ReceiptDetails";
import ReceiptItemUpdateDialog from "../../../features/receipt-items/ReceiptItemUpdateDialog";
import DeleteIcon from '@mui/icons-material/Delete';
import DeleteReceipt from "../../../features/receipts/DeleteReceipt";

const ReceiptContainer = styled(Stack)<StackProps>({
	gap: "2rem"
});

const BackButton = styled(Stack)<StackProps & LinkProps>(({ theme }) => ({
	textDecoration: "none",
	fontSize: "22px",
	color: theme.palette.primary.dark,
}));

const RemoveReceiptButton = styled(Button)<ButtonProps>(({theme}) => ({
	minWidth: "auto",
	color: theme.palette.error.light,
	border: `1px solid ${theme.palette.grey[400]}`,
	"&:hover": {
		borderColor: theme.palette.grey[600],
	}
}));

interface UpdateReceiptItemForm {
	open: boolean;
	receiptItem?: SingleReceiptReceiptItem;
}

const ShowReceiptPage: FunctionComponent = () => {
	// Extract `receiptId` from path
	const { id: receiptId } = useParams();

	// Update Receipt Item Form state
	const [updateReceiptItemForm, setUpdateReceiptItemForm] = useState<UpdateReceiptItemForm>({
		open: false,
	});
	const [deleteReceiptForm, setDeleteReceiptForm] = useState<{open: boolean}>({
		open: false
	});

	// Fetch receipt
	const { data: receipt } = useQuery({
		queryKey: ["single_receipt", receiptId],
		queryFn: () => getReceipt(receiptId!),
		keepPreviousData: true,
		enabled: !!receiptId,
	});

	const queryClient = useQueryClient();

	const refetch = () => {
		queryClient.invalidateQueries(["single_receipt", receiptId]);
	}

	const back = (
		<BackButton component={Link} to="/" direction="row" alignItems="center">
			<ChevronLeftIcon /> Nazad
		</BackButton>
	);

	const controls = (
		<RemoveReceiptButton onClick={() => setDeleteReceiptForm({ ...deleteReceiptForm, open: true })}>
			<DeleteIcon />
		</RemoveReceiptButton>
	)

	return (
		<PageLayout
			title="Pregled računa"
			headerPrefix={back}
			headerSuffix={controls}
		>
			<ReceiptContainer>
				<ReceiptPaymentOverview {...{
						storeName: receipt?.store.name || "",
						totalPurchaseAmount: receipt?.totalPurchaseAmount || "0.00",
						date: receipt?.date || "",
					}}
				/>
				<ReceiptItemsList
					receiptItems={receipt?.receiptItems || [] }
					onClick={(receiptItem: SingleReceiptReceiptItem) =>
						setUpdateReceiptItemForm({ open: true, receiptItem})
					}
				/>
				{receipt &&
					<ReceiptDetails {...{
						store: receipt.store,
						user: receipt.user,
						totalTaxAmount: receipt.totalTaxAmount,
						pfrNumber: receipt.pfrNumber,
						counter: receipt.pfrNumber,
						createdAt: receipt.createdAt,
					}} />
				}
			</ReceiptContainer>

			<ReceiptItemUpdateDialog
				open={updateReceiptItemForm.open}
				onClose={() => setUpdateReceiptItemForm({ ...updateReceiptItemForm, open: false })}
				receiptItem={updateReceiptItemForm.receiptItem}
				updateReceiptItem={(value: SingleReceiptReceiptItem) =>
					setUpdateReceiptItemForm({...updateReceiptItemForm, receiptItem: value})
				}
				onSubmitted={refetch}
			/>

			<DeleteReceipt
				open={deleteReceiptForm.open}
				onClose={() => setDeleteReceiptForm({ ...deleteReceiptForm, open: false })}
				receipt={receipt}
			/>

		</PageLayout>
	)
}

export default ShowReceiptPage;
