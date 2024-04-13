import { httpClient } from "../http"

export const getCategoriesStats = async ({ fromDate, toDate }: CategoriesStatsParams): Promise<CategoryStats> => {
	const { data } = await httpClient.get<CategoryStats>("/stats/categories", {
		params: { fromDate, toDate }
	});

	return data;
}

export const getExpensesDateBreakdownStats = async (): Promise<ExpensesDateBreakdown[]> => {
	const { data } = await httpClient.get<ExpensesDateBreakdown[]>("/stats/expenses/breakdown");

	return data;
}

export const getExpensesByCategoryBreakdownStats = async (): Promise<ExpensesByCategoryBreakdown[]> => {
	const { data } = await httpClient.get<ExpensesByCategoryBreakdown[]>("/stats/categories/breakdown");

	return data;
}

export const getExpensesByStoreBreakdownStats = async (): Promise<ExpensesByStoreBreakdown[]> => {
	const { data } = await httpClient.get<ExpensesByStoreBreakdown[]>("/stats/stores/breakdown");

	return data;
}
