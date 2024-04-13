import { Fragment, FunctionComponent, ReactElement } from "react";
import { Box, Divider, Stack, Typography } from "@mui/material";

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

interface ExpensesByStoreProps {
	stores: ExpensesByStoreBreakdown[],
}

const ExpensesByStore: FunctionComponent<ExpensesByStoreProps> = ({ stores }): ReactElement => {
	return (
		<Box>
			{stores.map((store, index) => (
				<ExpenseBySoreItem
					key={store.id}
					expenseByStore={store}
					divider={index !== stores.length -1}
				/>
			))}
		</Box>
	)
}

export default ExpensesByStore;
