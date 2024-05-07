import { Fragment, FunctionComponent, ReactElement } from "react";
import { Divider, Stack, Typography } from "@mui/material";
import CardLinkContent from "../../components/card/CardLinkContent";

interface ExpenseByStoreItemProps {
	expenseByStore: ExpensesByStoreBreakdown;
	divider?: boolean;
}

const ExpenseByStoreItem: FunctionComponent<ExpenseByStoreItemProps> = ({ expenseByStore, divider }) => {
	const { tin, name, total, receiptCount } = expenseByStore;

	return (
		<Fragment>
			<CardLinkContent to={`/stores/companies/${tin}`}>
				<Stack direction="row" justifyContent="space-between" gap="0.5rem">
					<Stack direction="row" alignItems="center">
						<Typography>{name}</Typography>
					</Stack>
					<Stack direction="column" alignItems="flex-end">
						<Typography>{total}</Typography>
						<Typography variant="caption">{receiptCount} račun</Typography>
					</Stack>
			</Stack>
			</CardLinkContent>
			{divider && <Divider />}
		</Fragment>
	)
}

interface ExpensesByStoreProps {
	stores: ExpensesByStoreBreakdown[],
}

const ExpensesByStore: FunctionComponent<ExpensesByStoreProps> = ({ stores }): ReactElement => {
	return (
		<Fragment>
			{stores.map((store, index) => (
				<ExpenseByStoreItem
					key={index}
					expenseByStore={store}
					divider={index !== stores.length -1}
				/>
			))}
		</Fragment>
	)
}

export default ExpensesByStore;
