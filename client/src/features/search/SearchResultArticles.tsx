import { Fragment, FunctionComponent, ReactElement } from "react";
import dayjs from "dayjs";
import { CircularProgress, Divider, Stack, Typography } from "@mui/material";
import Card from "../../components/card/Card";
import CardLinkContent from "../../components/card/CardLinkContent";

interface SearchResultArticleProps {
	receiptItem: ReceiptItem;
	divider?: boolean;
}

const SearchResultArticle: FunctionComponent<SearchResultArticleProps> = ({ receiptItem, divider }): ReactElement => {
	const { receiptId, name, date, store, amount, unit, quantity } = receiptItem;
	const formattedDate = dayjs(date).format("DD.MM.YYYY.");

	return (
		<Fragment>
			<CardLinkContent to={`/receipts/${receiptId}`}>
				<Stack direction="row" justifyContent="space-between">
					<Stack direction="column" justifyContent="space-between" gap="0.25rem">
						<Typography>{name}</Typography>
						<Typography variant="body2" color="grey.700">{`${formattedDate} - ${store}`}</Typography>
					</Stack>
					<Stack direction="column" alignItems="flex-end" gap="0.25rem">
						<Typography>{amount}</Typography>
						<Typography
							component={Stack}
							direction="row"
							alignItems="center"
							variant="body2"
							color="grey.700"
							gap="0.25rem"
						>
							x <Typography>{`${quantity} ${unit}`}</Typography>
						</Typography>
					</Stack>
				</Stack>
			</CardLinkContent>
			{divider && <Divider />}
		</Fragment>
	)
}

interface SearchResultArticlesProps {
	receiptItems: ReceiptItem[];
	isLoading?: boolean;
}

const SearchResultArticles: FunctionComponent<SearchResultArticlesProps> = ({
	receiptItems,
	isLoading,
}): ReactElement => {
	if (isLoading) {
		return (
			<Stack alignItems="center">
				<CircularProgress color="primary" />
			</Stack>
		)
	}

	if (!receiptItems.length && !isLoading) {
		return <Typography>Nema artikala koji odgovaraju pretraženom terminu.</Typography>
	}

	return (
		<Card>
			{receiptItems.map((receiptItem, index) => (
				<SearchResultArticle
					key={receiptItem.id}
					receiptItem={receiptItem}
					divider={index !== receiptItems.length - 1}
				/>
			))}
		</Card>
	)
}

export default SearchResultArticles;
