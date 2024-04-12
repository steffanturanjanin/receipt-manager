import { Fragment, FunctionComponent, ReactElement } from "react";
import { useQuery } from "react-query";
import { Box, Divider, Stack, Typography } from "@mui/material";
import { getExpensesByCategoryBreakdownStats } from "../../api/stats";


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

const ExpensesByCategory: FunctionComponent = (): ReactElement => {
	const { data: categoriesStats } = useQuery({
		queryKey: ["expenses_by_category_stats_yearly_breakdown"],
		queryFn: () => getExpensesByCategoryBreakdownStats(),
	});

	return (
		<Box>
			{categoriesStats?.map((category, index) => (
				<ExpensesByCategoryItem
					key={category.id}
					category={category}
					divider={index !== categoriesStats.length -1}
				/>
			))}
		</Box>
	);
}

export default ExpensesByCategory;
