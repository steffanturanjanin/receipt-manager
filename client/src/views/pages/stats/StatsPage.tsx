import { FunctionComponent, ReactElement, useMemo } from "react";
import { useQuery } from "react-query";
import { Stack } from "@mui/material";
import PageLayout from "../../layouts/PageLayout";
import ExpensesChart from "../../../features/stats/ExpensesChart";
import StatItem from "../../../features/stats/StatItem";
import ExpensesByCategory from "../../../features/stats/ExpensesByCategory";
import ExpensesByStore from "../../../features/stats/ExpensesByStore";
import {
	getExpensesByCategoryBreakdownStats,
	getExpensesByStoreBreakdownStats,
	getExpensesDateBreakdownStats
 } from "../../../api/stats";
import CardContent from "../../../components/card/CardContent";
import { Link } from "react-router-dom";

const StatsPage: FunctionComponent = (): ReactElement => {
	const { data: expenses, isFetching: isExpensesFetching } = useQuery({
		queryKey: ["expenses_yearly_breakdown"],
		queryFn: () => getExpensesDateBreakdownStats(),
	});

	const { data: categories, isFetching: isCategoriesFetching } = useQuery({
		queryKey: ["expenses_by_category_stats_yearly_breakdown"],
		queryFn: () => getExpensesByCategoryBreakdownStats(),
	});

	const { data: stores, isFetching: isStoresFetching } = useQuery({
		queryKey: ["expenses_by_store_yearly_breakdown"],
		queryFn: () => getExpensesByStoreBreakdownStats(),
	});

	const isLoading = useMemo(
		() => isExpensesFetching || isCategoriesFetching || isStoresFetching,
		[isExpensesFetching, isCategoriesFetching, isStoresFetching]
	);

	const expensesByCategoryEmpty = useMemo(
		() => (!categories || !categories.length) && "Nema potrošnji po kategorijama...",
		[categories]
	);

	const expensesByStoreEmpty = useMemo(
		() => (!stores || !stores.length) && "Nema potrošnji po prodavnicama...",
		[stores]
	)

	return (
		<PageLayout title="Statistika" showBackdrop={isLoading}>
			<Stack direction="column" gap="2rem">
				<StatItem title="Potrošnja u proteklih 12 meseci">
					<CardContent>
						<ExpensesChart expenses={expenses || []}/>
					</CardContent>
				</StatItem>

				<StatItem
					title="Potrošnja po kategoriji"
					empty={expensesByCategoryEmpty}
				>
					<ExpensesByCategory categories={categories || []}/>
				</StatItem>

				<StatItem
					title={<Link to="/stores">Top 10 prodavnica</Link>}
					empty={expensesByStoreEmpty}
				>
					<ExpensesByStore stores={stores || []}/>
				</StatItem>
			</Stack>
		</PageLayout>
	)
}

export default StatsPage;
