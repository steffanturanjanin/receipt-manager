import { Fragment, FunctionComponent, ReactElement } from "react";
import { Link, LinkProps } from "react-router-dom";
import dayjs from "dayjs";
import { Divider, Stack, StackProps, Typography, styled } from "@mui/material";
import Card from "../../components/card/Card";

interface CompanyExpensesProps {
	companyExpenses: CompanyExpense[];
}

const CompanyExpense = styled(Stack)<StackProps & LinkProps>(({ theme }) => ({
	display: "flex",
	flexDirection: "row",
	justifyContent: "space-between",
	alignItems: "center",
	padding: "1rem",
	color: "inherit",
	textDecoration: "none",
	"&:hover": {
		backgroundColor: theme.palette.grey[50],
	}
}));

interface CompanyExpenseItemProps {
	expense: CompanyExpense;
	divider?: boolean
}

const CompanyExpenseItem: FunctionComponent<CompanyExpenseItemProps> = ({ expense, divider}): ReactElement => {
	const { id, date, locationName, amount } = expense;
	const formattedDate = dayjs(date).format("MM.DD.YYYY. HH:MM");

	return (
		<Fragment>
			<CompanyExpense component={Link} to={`/receipts/${id}`}>
				<Stack direction="column" gap="0.25rem">
					<Typography>{locationName}</Typography>
					<Typography color="grey.700" variant="body2">{formattedDate}</Typography>
				</Stack>
				<Typography>{amount}</Typography>
			</CompanyExpense>
			{divider && <Divider />}
		</Fragment>
	)
}

const CompanyExpenses: FunctionComponent<CompanyExpensesProps> = ({ companyExpenses }): ReactElement => {
	return (
		<Stack direction="column" gap="1rem">
			<Typography variant="h6" component="h2">Troškovi</Typography>
			<Card>
				{companyExpenses.map((expense, index) => (
					<CompanyExpenseItem
						key={index}
						expense={expense}
						divider={index !== companyExpenses.length - 1}
					/>
				))}
			</Card>
		</Stack>
	)
}

export default CompanyExpenses;
