import { Fragment, FunctionComponent } from "react";
import { Link } from "react-router-dom";
import { Card, CardContent, Stack, Typography, } from "@mui/material";
import { StackProps, styled } from "@mui/system";
import dayjs from "dayjs";

interface ReceiptListProps {
	receiptsAggregatedByDate: ReceiptsAggregatedByDate[];
}

const ReceiptListItem: FunctionComponent<ReceiptAggregatedByDateItem> = ({ id, amount, date, store, categories }) => {
	const timeFormatted = dayjs(date).format("HH:mm");
	const categoriesList = categories.join(", ");

	return (
		<Card component={Link} to="/" sx={{ textDecoration: "none"}} key={id}>
			<CardContent component={Stack} direction="column" gap="1rem">
				<Stack direction="row" justifyContent="space-between">
					<Typography component="span" fontWeight="bold">{store.name}</Typography>
					<Typography component="span">{amount}</Typography>
				</Stack>
				<Stack direction="row" justifyContent="space-between" alignItems="center">
					<Typography component="span" variant="body2">{categoriesList}</Typography>
					<Typography component="span">{timeFormatted}</Typography>
				</Stack>
			</CardContent>
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
		<Fragment>
			{receipts.map(receipt => (
				<Stack direction="column" gap="0.5rem" key={receipt.date}>
					<ReceiptGroupHeader direction="row" justifyContent="space-between">
						<Typography component="span" color="text.secondary">{formattedDate}</Typography>
						<Typography component="span">{total}</Typography>
					</ReceiptGroupHeader>
					{receipts.map(receipt => <ReceiptListItem key={receipt.id} {...receipt} />)}
				</Stack>
			))}
		</Fragment>
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
