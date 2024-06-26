import { Fragment, FunctionComponent, ReactElement } from "react";
import { Divider, Stack, Typography } from "@mui/material";
import CardLinkContent from "../../components/card/CardLinkContent";

interface ExpensesByCategoryItemProps {
	category: ExpensesByCategoryBreakdown;
	divider?: boolean;
}

const ExpensesByCategoryItem: FunctionComponent<ExpensesByCategoryItemProps> = ({ category, divider }) => {
	const { id, name, total, percentage, receiptCount } = category;

	return (
		<Fragment>
			<CardLinkContent to={`/categories/${id}`}>
				<Stack direction="column" gap="0.5rem">
					<Stack direction="row" justifyContent="space-between">
						<Typography>{name}</Typography>
						<Typography>{total}</Typography>
					</Stack>
					<Stack direction="row" justifyContent="space-between">
						<Typography variant="body2" color="grey.600">{receiptCount} račun</Typography>
						<Typography variant="body2" color="grey.600">{percentage}%</Typography>
					</Stack>
				</Stack>
			</CardLinkContent>
			{divider && <Divider />}
		</Fragment>
	)
}

interface ExpensesByCategoryProps {
	categories: ExpensesByCategoryBreakdown[]
}

const ExpensesByCategory: FunctionComponent<ExpensesByCategoryProps> = ({ categories }): ReactElement => {
	return (
		<Fragment>
			{categories.map((category, index) => (
				<ExpensesByCategoryItem
					key={index}
					category={category}
					divider={index !== categories.length -1}
				/>
			))}
		</Fragment>
	);
}

export default ExpensesByCategory;
