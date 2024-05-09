import { FunctionComponent } from "react";
import { Stack, Typography, } from "@mui/material";
import { StackProps, styled } from "@mui/system";
import dayjs from "dayjs";
import Card from "../../components/card/Card";
import CardLinkContent from "../../components/card/CardLinkContent";

interface ReceiptListProps {
	receiptsAggregatedByDate: ReceiptsAggregatedByDate[];
}

const ReceiptListItem: FunctionComponent<ReceiptAggregatedByDateItem> = ({ id, amount, date, store, categories }) => {
	const timeFormatted = dayjs(date).format("HH:mm");
	const categoriesList = categories.join(", ");

	return (
		<Card key={id}>
			<CardLinkContent to={`/receipts/${id}`}>
				<Stack direction="column" gap="1rem">
					<Stack direction="row" justifyContent="space-between">
						<Typography component="span" fontWeight="bold">{store.name}</Typography>
						<Typography component="span">{amount}</Typography>
					</Stack>
					<Stack direction="row" justifyContent="space-between" alignItems="center">
						<Typography component="span" variant="body2">{categoriesList}</Typography>
						<Typography component="span">{timeFormatted}</Typography>
					</Stack>
				</Stack>
			</CardLinkContent>
		</Card>
	)
}

const ReceiptListGroup: FunctionComponent<ReceiptsAggregatedByDate> = ({ date, total, receipts }) => {
	const formattedDate = dayjs(date).format("DD.MM.YYYY");

	const ReceiptGroupHeader = styled(Stack)<StackProps>({
		gap: "2rem",
		paddingLeft: "1rem",
		paddingRight: "1rem",
	})

	return (
			<Stack direction="column" gap="0.5rem" >
				<ReceiptGroupHeader direction="row" justifyContent="space-between">
					<Typography component="span" color="text.secondary">{formattedDate}</Typography>
					<Typography component="span">{total}</Typography>
				</ReceiptGroupHeader>
				{receipts.map((receipt, index) => <ReceiptListItem key={index} {...receipt} />)}
			</Stack>
	)
}

const ReceiptList: FunctionComponent<ReceiptListProps> = ({ receiptsAggregatedByDate }) => {
	return (
		receiptsAggregatedByDate.map((aggregatedReceipts, index) =>
			<ReceiptListGroup key={index} {...aggregatedReceipts} />
		)
	)
}

export default ReceiptList;
