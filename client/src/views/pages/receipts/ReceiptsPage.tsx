import { Fragment, FunctionComponent, ReactElement, useMemo } from "react";
import { useQuery } from "react-query";
import { getReceipts } from "../../../api/receipts";
import { Card, CardContent, Stack, Typography } from "@mui/material";
import PageLayout from "../../layouts/PageLayout/PageLayout";

const ReceiptCard: FunctionComponent<{ receipt: Receipt}> = ({ receipt }): ReactElement => {
	const { store } = receipt;
	const categoryNames = [...new Set(receipt.receiptItems.map(receiptItem => receiptItem.category.name))];

	return (
		<Card key={receipt.id}>
			<CardContent>
				<Stack>
					<Typography component="p">{receipt.store.name}</Typography>
					<Typography variant="body2">{categoryNames.join(", ")}</Typography>
					<Stack direction="row" justifyContent="space-between">
						<Typography>{receipt.totalPurchaseAmount}</Typography>
						<Typography>{receipt.date}</Typography>
					</Stack>
				</Stack>
			</CardContent>
		</Card>
	)
}

const ReceiptsPage: FunctionComponent = (): ReactElement => {

	const { isLoading, data: receipts } = useQuery({
		queryKey: "receipts",
		queryFn: () => getReceipts({ fromDate: "", toDate: ""}),
	})

	return (
		<PageLayout title="Receipts">
			{receipts?.data.map(receipt => <ReceiptCard receipt={receipt} />)}
		</PageLayout>
	);
}

export default ReceiptsPage;
