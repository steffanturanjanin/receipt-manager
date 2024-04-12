import { httpClient } from "../http"
import { expensesByCategoryYearlyBreakdown, expensesByStoreYearlyBreakdown, expensesYearlyBreakdown } from "./mock";

export const getCategoriesStats = async ({ fromDate, toDate }: CategoriesStatsParams): Promise<CategoryStats> => {
	const { data } = await httpClient.get<CategoryStats>("/stats/categories", {
		params: { fromDate, toDate }
	});

	return data;
}

export const getExpensesDateBreakdownStats = async (): Promise<ExpensesDateBreakdown[]> => {
	// const { data } = await httpClient.get<ExpensesYearlyBreakdown[]>("/stats/expenses/yearly-breakdown");

	// return data;

	return expensesYearlyBreakdown;
}

export const getExpensesByCategoryBreakdownStats = async (): Promise<ExpensesByCategoryBreakdown[]> => {
	// const { data } = await httpClient.get<CategoryYearlyBreakdown[]>("/stats/categories/yearly-breakdown");

	// return data;

	return expensesByCategoryYearlyBreakdown;
}

export const getExpensesByStoreBreakdownStats = async (): Promise<ExpensesByStoreBreakdown[]> => {
	// const { data } = await httpClient.get<ExpensesByCategoryYearlyBreakdown[]>("/stats/stores/yearly-breakdown");

	// return data;

	return expensesByStoreYearlyBreakdown;
}
