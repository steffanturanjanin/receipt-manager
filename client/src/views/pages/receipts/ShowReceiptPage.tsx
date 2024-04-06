import { FunctionComponent } from "react";
import { useQuery } from "react-query";
import { getReceipt } from "../../../api/receipts";
import { Link, LinkProps, useParams } from "react-router-dom";
import PageLayout from "../../layouts/PageLayout/PageLayout";
import { Stack, StackProps, styled } from "@mui/material";
import ChevronLeftIcon from '@mui/icons-material/ChevronLeft';
import ReceiptPaymentOverview from "../../../features/receipts/show-receipt/ReceiptPaymentOverview";
import ReceiptItemsList from "../../../features/receipts/show-receipt/ReceiptItemsList";
import ReceiptDetails from "../../../features/receipts/show-receipt/ReceiptDetails";


const ShowReceiptPage: FunctionComponent = () => {
	const { id: receiptId } = useParams();

	const { data: receipt } = useQuery({
		queryKey: ["single_receipt", receiptId],
		queryFn: () => getReceipt(receiptId!),
		keepPreviousData: true,
		enabled: !!receiptId,
	});

	const ReceiptContainer = styled(Stack)<StackProps>({
		gap: "2rem"
	});

	const BackButton = styled(Stack)<StackProps & LinkProps>(({ theme }) => ({
		textDecoration: "none",
		fontSize: "22px",
		color: theme.palette.primary.dark,
	}))

	return (
		<PageLayout
			title="Pregled računa"
			headerPrefix={
				<BackButton component={Link} to="/" direction="row" alignItems="center">
					<ChevronLeftIcon /> Nazad
				</BackButton>
			}>
			<ReceiptContainer>
				<ReceiptPaymentOverview {...{
						storeName: receipt?.store.name || "",
						totalPurchaseAmount: receipt?.totalPurchaseAmount || "0.00",
						date: receipt?.date || "",
					}}
				/>
				<ReceiptItemsList
					receiptItems={receipt?.receiptItems || [] }
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
		</PageLayout>
	)
}

export default ShowReceiptPage;
