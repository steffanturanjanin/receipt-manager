import { FunctionComponent, useMemo, useState } from "react";
import { useMutation, useQuery, useQueryClient } from "react-query";
import { getReceipt, setFavorite } from "../../../api/receipts";
import { useParams } from "react-router-dom";
import PageLayout from "../../layouts/PageLayout";
import { Stack, StackProps, Typography, styled } from "@mui/material";
import ReceiptPaymentOverview from "../../../features/receipts/show-receipt/ReceiptPaymentOverview";
import ReceiptItemsList from "../../../features/receipts/show-receipt/ReceiptItemsList";
import ReceiptDetails from "../../../features/receipts/show-receipt/ReceiptDetails";
import ReceiptItemUpdateDialog from "../../../features/receipt-items/ReceiptItemUpdateDialog";
import DeleteReceipt from "../../../features/receipts/DeleteReceipt";
import BackButton from "../../../components/BackButton";
import { DeleteActionButton, FavoriteActionButton } from "../../../features/receipts/ActionButtons";

const ReceiptContainer = styled(Stack)<StackProps>({
	gap: "2rem"
});

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
	const { isLoading: isReceiptLoading, data: receipt } = useQuery({
		queryKey: ["single_receipt", receiptId],
		queryFn: () => getReceipt(receiptId!),
		keepPreviousData: true,
		enabled: !!receiptId,
	});

	// Favorite
	const { isLoading: isSetFavoriteLoading, mutate: setFavoriteMutate } = useMutation({
		mutationFn: (request: SetFavoriteRequest) => setFavorite(receipt!.id, request),
		onSuccess: () => {
			// When favorite status is changed
			// Invalidate query and refetch receipt
			refetchReceipt()
		}
	});

	const onSetFavorite = () => {
		if (receipt) {
			const request: SetFavoriteRequest = { isFavorite: !receipt.isFavorite };
			setFavoriteMutate(request);
		}
	}

	const isLoading = useMemo(
		() => isReceiptLoading || isSetFavoriteLoading,
		[isReceiptLoading, isSetFavoriteLoading]
	);

	const queryClient = useQueryClient();

	const refetchReceipt = () => {
		queryClient.invalidateQueries(["single_receipt", receiptId]);
	}

	const controls = (
		<Stack direction="row" gap="0.5rem">
			<DeleteActionButton onClick={() => setDeleteReceiptForm({ ...deleteReceiptForm, open: true })} />
			<FavoriteActionButton isFavorite={!!receipt?.isFavorite} onClick={onSetFavorite}/>
		</Stack>
	)

	return (
		<PageLayout
			headerPrefix={<BackButton />}
			headerSuffix={controls}
			showBackdrop={isLoading}
		>
			<ReceiptContainer>
				<Typography variant="h4" component="h1">Pregled računa</Typography>
				<ReceiptPaymentOverview {...{
						storeName: receipt?.store.name || "",
						totalPurchaseAmount: receipt?.totalPurchaseAmount || "0.00",
						date: receipt?.date,
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
						counter: receipt.counter,
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
				onSubmitted={refetchReceipt}
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
