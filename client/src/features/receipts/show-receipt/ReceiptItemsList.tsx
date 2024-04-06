import { Fragment, FunctionComponent, ReactElement } from "react";
import { Box, BoxProps, Divider, Stack, Typography, styled } from "@mui/material";
import { ReceiptCard, ReceiptCardContent } from "./components";

const ColorCircle = styled(Box)<BoxProps>({
	width: "0.6em",
	height: "0.6em",
	borderRadius: "100%",
});

interface CategoryCircleProps {
	color: string;
}

const CategoryCircle: FunctionComponent<CategoryCircleProps> = ({ color }) =>
	<ColorCircle component="span" sx={{ backgroundColor: color }}/>


interface ReceiptItemProps {
	receiptItem: SingleReceiptReceiptItem;
	divider: boolean;
}

const ReceiptItem: FunctionComponent<ReceiptItemProps> = ({ divider, receiptItem }): ReactElement => {
	const { name, totalAmount, category, quantity, unit, singleAmount } = receiptItem;
	const breakdownPerUnit = `${quantity} ${unit} x ${singleAmount}`;

	return (
		<Fragment>
			<Stack direction="column" gap="0.5rem">
				<Stack direction="row" gap="0.5rem" justifyContent="space-between">
					<Typography variant="h6">{name}</Typography>
					<Typography variant="h6">{totalAmount}</Typography>
				</Stack>
				<Stack direction="row" gap="0.5rem" justifyContent="space-between">
					<Stack component={Typography} variant="body2" direction="row" alignItems="center" gap="0.5rem">
						<CategoryCircle color={category.color} />
						{category.name}
					</Stack>
					<Typography variant="body2">{breakdownPerUnit}</Typography>
				</Stack>
			</Stack>
			{divider && <Divider orientation="horizontal" sx={{ marginY: "1rem" }}/> }
		</Fragment>
	)
}

interface ReceiptItemsListProps {
	receiptItems: SingleReceiptReceiptItem[];
}

const ReceiptItemsList: FunctionComponent<ReceiptItemsListProps> = ({ receiptItems }): ReactElement => {
	return (
		<Box component="section">
			<Typography variant="h4" component="h2" marginY="2rem">Stavke sa računa</Typography>
			<ReceiptCard>
				<ReceiptCardContent>
					{receiptItems.map((receiptItem, index) => (
						<ReceiptItem
							key={receiptItem.id}
							receiptItem={receiptItem}
							divider={index !== receiptItems.length - 1}
						/>
					))}
				</ReceiptCardContent>
			</ReceiptCard>
		</Box>
	)
}

export default ReceiptItemsList;
