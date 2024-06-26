import { Fragment, FunctionComponent, ReactElement } from "react";
import { Box, BoxProps, ButtonProps, Divider, Stack, Typography, styled, useTheme } from "@mui/material";
import Card from "../../../components/card/Card";
import CategoryCircle from "../../categories/CategoryCircle";

interface ReceiptItemProps {
	receiptItem: SingleReceiptReceiptItem;
	divider: boolean;
	onClick: (receiptItem: SingleReceiptReceiptItem) => void;
}

const ReceiptItemContainer = styled(Box)<BoxProps & ButtonProps>(({ theme }) => ({
	cursor: "pointer",
	padding: "1rem",
	"&:hover": {
		backgroundColor: theme.palette.grey[50],
	}
}));

const ReceiptItem: FunctionComponent<ReceiptItemProps> = ({ divider, receiptItem, onClick }): ReactElement => {
	const { name, totalAmount, category, quantity, unit, singleAmount } = receiptItem;
	const breakdownPerUnit = `${quantity} ${unit} x ${singleAmount}`;

	const theme = useTheme();
	const UNCATEGORIZED_COLOR = theme.palette.grey['500'];
	const UNCATEGORIZED_NAME = "Nekategorisano";

	return (
		<Fragment>
			<ReceiptItemContainer onClick={() => onClick(receiptItem)}>
				<Stack direction="column" gap="0.5rem">
					<Stack direction="row" gap="0.5rem" justifyContent="space-between">
						<Typography fontWeight="bold">{name}</Typography>
						<Typography fontWeight="bold">{totalAmount}</Typography>
					</Stack>
					<Stack direction="row" gap="0.5rem" justifyContent="space-between">
						<Stack component={Typography} variant="body2" direction="row" alignItems="center" gap="0.5rem">
							<CategoryCircle color={category?.color || UNCATEGORIZED_COLOR} />
							{category?.name || UNCATEGORIZED_NAME}
						</Stack>
						<Typography variant="body2" color="grey.600">{breakdownPerUnit}</Typography>
					</Stack>
				</Stack>
			</ReceiptItemContainer>
			{divider && <Divider orientation="horizontal" /> }
		</Fragment>
	)
}

interface ReceiptItemsListProps {
	receiptItems: SingleReceiptReceiptItem[];
	onClick: (receiptItem: SingleReceiptReceiptItem) => void;
}

const ReceiptItemsList: FunctionComponent<ReceiptItemsListProps> = ({ receiptItems, onClick }): ReactElement => {
	return (
		<Box component="section">
			<Typography variant="h5" component="h2" marginY="1rem">Stavke sa računa</Typography>
			<Card>
				{receiptItems.map((receiptItem, index) => (
					<ReceiptItem
						key={receiptItem.id}
						receiptItem={receiptItem}
						divider={index !== receiptItems.length - 1}
						onClick={onClick}
					/>
				))}
			</Card>
		</Box>
	)
}

export default ReceiptItemsList;
