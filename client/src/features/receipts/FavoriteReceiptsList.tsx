import { FunctionComponent, MouseEvent, ReactElement } from "react";
import dayjs from "dayjs";
import { ButtonProps, Stack, Typography, styled } from "@mui/material";
import StarFullIcon from "@mui/icons-material/Star";
import Card from "../../components/card/Card";
import CardLinkContent from "../../components/card/CardLinkContent";
import CategoryCircle from "../categories/CategoryCircle";

const FavoriteButton = styled("button")<ButtonProps>(({ theme }) => ({
	maxWidth: "2rem",
	maxHeight: "2rem",
	color: theme.palette.warning.light,
	border: "none",
	background: "none",
	padding: 0,
	cursor: "pointer",
}));

interface FavoriteReceiptsListItemProps {
	receipt: FavoriteReceipt;
	onClick: () => void;
}

const FavoriteReceiptsListItem: FunctionComponent<FavoriteReceiptsListItemProps> = ({ receipt, onClick }): ReactElement => {
	const { store, categories, amount, date } = receipt;

	const formattedDate = dayjs(date).format("DD.MM.YYYY.");
	const formattedTime = dayjs(date).format("HH:mm");

	const onRemoveFavoriteClick = (event: MouseEvent<HTMLButtonElement>) => {
		event.stopPropagation();
		event.preventDefault();
		onClick();
	}

	return (
		<Card>
			<CardLinkContent to={`/receipts/${receipt.id}`}>
				<Stack direction="row" justifyContent="space-between" gap="1.5rem">
					<Stack direction="column" gap="1rem" justifyContent="space-between">
						<Stack direction="column">
							<Typography fontWeight="bold">{store.name}</Typography>
							<Typography>{store.location}</Typography>
							<Typography>{`${store.address} - ${store.city}`}</Typography>
						</Stack>
						<Stack direction="row" gap="0.5rem" flexWrap="wrap">
							{categories.map((category, index) =>
								<Stack
									key={index}
									component={Typography}
									direction="row"
									justifyContent="center"
									alignItems="center"
									gap="0.25rem"
									variant="caption"
									color="grey.600"
								>
									<CategoryCircle color={category.color} /> {category.name}
								</Stack>
							)}
						</Stack>
					</Stack>
					<Stack direction="column" alignItems="flex-end" justifyContent="space-between" gap="1rem">
						<FavoriteButton onClick={onRemoveFavoriteClick}>
							<StarFullIcon sx={{ fontSize: "2rem" }}/>
						</FavoriteButton>
						<Stack direction="column" alignItems="flex-end" gap="0.5rem">
							<Typography fontWeight="bold">{amount}</Typography>
							<Stack alignItems="flex-end">
								<Typography variant="body2" color="grey.600">{formattedDate}</Typography>
								<Typography variant="body2" color="grey.600">{formattedTime}</Typography>
							</Stack>
						</Stack>
					</Stack>
				</Stack>
			</CardLinkContent>
		</Card>
	)
}

interface FavoriteReceiptsListProps {
	receipts: FavoriteReceipt[];
	onItemClicked: (receiptItem: number) => void;
}

const FavoriteReceiptsList: FunctionComponent<FavoriteReceiptsListProps> = ({ receipts, onItemClicked }): ReactElement => {
	const empty = (
		<Stack direction="column" alignItems="center" justifyContent="center" marginTop="50%">
			<Typography variant="h4">Nema omiljenih računa.</Typography>
			<Typography>Dodajte omiljene račune i oni će biti prikazani na ovoj strani.</Typography>
		</Stack>
	);

	const list = (
		receipts.map((receipt, index) => (
			<FavoriteReceiptsListItem
				key={index}
				receipt={receipt}
				onClick={() => onItemClicked(receipt.id)}
			/>
		))
	)

	const content = receipts.length ? list : empty;

	return (
		<Stack direction="column" gap="1rem">
			{content}
		</Stack>
	)
}

export default FavoriteReceiptsList;
