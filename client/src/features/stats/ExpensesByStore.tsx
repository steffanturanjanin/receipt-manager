import { Fragment, FunctionComponent, ReactElement } from "react";
import { useQuery } from "react-query";
import { Box, Divider, Stack, Typography } from "@mui/material";
import { getExpensesByStoreBreakdownStats } from "../../api/stats";

interface ExpenseByStoreItemProps {
	expenseByStore: ExpensesByStoreBreakdown;
	divider?: boolean;
}

const ExpenseBySoreItem: FunctionComponent<ExpenseByStoreItemProps> = ({ expenseByStore, divider }) => {
	const { name, total, receiptCount } = expenseByStore;

	return (
		<Fragment>
			<Stack direction="row" justifyContent="space-between" sx={{ padding: "1rem" }} gap="0.5rem">
				<Stack direction="row" alignItems="center">
					<Typography>{name}</Typography>
				</Stack>
				<Stack direction="column" alignItems="flex-end">
					<Typography>{total}</Typography>
					<Typography variant="caption">{receiptCount} račun</Typography>
				</Stack>
			</Stack>
			{divider && <Divider />}
		</Fragment>
	)
}

const ExpensesByStore: FunctionComponent = (): ReactElement => {
	const { data: expensesByStore } = useQuery({
		queryKey: ["expenses_by_store_yearly_breakdown"],
		queryFn: () => getExpensesByStoreBreakdownStats(),
	});

	return (
		<Box>
			{expensesByStore?.map((expenseByStore, index) => (
				<ExpenseBySoreItem
					key={expenseByStore.id}
					expenseByStore={expenseByStore}
					divider={index !== expensesByStore.length -1}
				/>
			))}
		</Box>
	)
}

export default ExpensesByStore;
