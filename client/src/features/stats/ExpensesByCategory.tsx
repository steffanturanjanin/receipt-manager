import { Fragment, FunctionComponent, ReactElement } from "react";
import { Box, Divider, Stack, Typography } from "@mui/material";

interface ExpensesByCategoryItemProps {
	category: ExpensesByCategoryBreakdown;
	divider?: boolean;
}

const ExpensesByCategoryItem: FunctionComponent<ExpensesByCategoryItemProps> = ({ category, divider }) => {
	const { name, total, percentage, receiptCount } = category;

	return (
		<Fragment>
			<Stack direction="column" sx={{ padding: "1rem" }} gap="0.5rem">
				<Stack direction="row" justifyContent="space-between">
					<Typography>{name}</Typography>
					<Typography>{total}</Typography>
				</Stack>
				<Stack direction="row" justifyContent="space-between">
					<Typography variant="caption">{receiptCount} račun</Typography>
					<Typography variant="caption">{percentage}%</Typography>
				</Stack>
			</Stack>
			{divider && <Divider />}
		</Fragment>
	)
}

interface ExpensesByCategoryProps {
	categories: ExpensesByCategoryBreakdown[]
}

const ExpensesByCategory: FunctionComponent<ExpensesByCategoryProps> = ({ categories }): ReactElement => {
	return (
		<Box>
			{categories.map((category, index) => (
				<ExpensesByCategoryItem
					key={category.id}
					category={category}
					divider={index !== categories.length -1}
				/>
			))}
		</Box>
	);
}

export default ExpensesByCategory;
