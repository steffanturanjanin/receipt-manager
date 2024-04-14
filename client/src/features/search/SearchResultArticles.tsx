import { CircularProgress, Divider, Paper, Stack, StackProps, Typography, styled } from "@mui/material";
import dayjs from "dayjs";
import { Fragment, FunctionComponent, ReactElement } from "react";

interface SearchResultArticleProps {
	receiptItem: ReceiptItem;
	divider?: boolean;
}

const SearchResultArticle: FunctionComponent<SearchResultArticleProps> = ({ receiptItem, divider }): ReactElement => {
	const { name, date, store, amount, unit, quantity } = receiptItem;
	const formattedDate = dayjs(date).format("DD.MM.YYYY.");

	return (
		<Fragment>
			<Stack direction="row" justifyContent="space-between" padding="1rem">
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
			{divider && <Divider />}
		</Fragment>
	)
}

interface SearchResultArticlesProps {
	receiptItems: ReceiptItem[];
	isLoading?: boolean;
}

const SearchResultContainer = styled(Stack)<StackProps>({
	borderRadius: "0.75rem",
	boxShadow: "#959da533 0 0.5rem 1.5rem",
})

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
		<SearchResultContainer direction="column" component={Paper}>
			{receiptItems.map((receiptItem, index) => (
				<SearchResultArticle
					key={receiptItem.id}
					receiptItem={receiptItem}
					divider={index !== receiptItems.length - 1}
				/>
			))}
		</SearchResultContainer>
	)
}

export default SearchResultArticles;
