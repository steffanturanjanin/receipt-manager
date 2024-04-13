import { FunctionComponent, ReactElement, useMemo } from "react";
import { useQuery } from "react-query";
import { Stack } from "@mui/material";
import PageLayout from "../../layouts/PageLayout";
import ExpensesChart from "../../../features/stats/ExpensesChart";
import StatItem, { StatItemContent } from "../../../features/stats/StatItem";
import ExpensesByCategory from "../../../features/stats/ExpensesByCategory";
import ExpensesByStore from "../../../features/stats/ExpensesByStore";
import {
	getExpensesByCategoryBreakdownStats,
	getExpensesByStoreBreakdownStats,
	getExpensesDateBreakdownStats
 } from "../../../api/stats";

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

	return (
		<PageLayout title="Statistika" showBackdrop={isLoading}>
			<Stack direction="column" gap="2rem">
				<StatItem title="Potrošnja u proteklih 12 meseci">
					<StatItemContent>
						<ExpensesChart expenses={expenses || []}/>
					</StatItemContent>
				</StatItem>

				<StatItem
					title="Potrošnja po kategoriji"
					empty={!categories?.length && "Nema potrošnji po kategorijama..."}
				>
					<ExpensesByCategory categories={categories || []}/>
				</StatItem>

				<StatItem
					title="Top 10 prodavnica"
					empty={!stores?.length && "Nema potrošnji po prodavnicama..."}
				>
					<ExpensesByStore stores={stores || []}/>
				</StatItem>
			</Stack>
		</PageLayout>
	)
}

export default StatsPage;
