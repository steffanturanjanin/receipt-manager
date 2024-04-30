import { Divider, Stack, Typography } from "@mui/material";
import { Fragment, FunctionComponent, ReactElement } from "react";
import CardItem from "./CardItem";

interface StoreExpensesProps {
	storeExpenses: StoreExpenses[];
}

const StoreExpenses: FunctionComponent<StoreExpensesProps> = ({ storeExpenses }): ReactElement => {
	return (
		<Stack direction="column" gap="1rem">
			<Typography variant="h6" component="h2">Troškovi</Typography>
			<CardItem>
				{storeExpenses.map((expense, index) => (
					<Fragment key={index}>
						<Stack direction="row" justifyContent="space-between" alignItems="center" padding="1rem">
							<Stack direction="column" gap="0.25rem">
								<Typography>{expense.locationName}</Typography>
								<Typography color="grey.700" variant="body2">{expense.date}</Typography>
							</Stack>
							<Typography>{expense.amount}</Typography>
						</Stack>
						{index !== storeExpenses.length - 1 && <Divider />}
					</Fragment>
				))}
			</CardItem>
		</Stack>
	)
}

export default StoreExpenses;
