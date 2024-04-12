interface CategoriesStatsParams {
	fromDate: string;
	toDate: string;
}

interface CategoryStat {
	id: number;
	color: string;
	name: string;
}

interface CategoryStatItem {
	category: CategoryStat;
	total: string;
}

interface CategoryStats {
	total: string;
	categories: CategoryStatItem[];
}

interface ExpensesDateBreakdown {
	date: string;
	total: string;
}

interface ExpensesByCategoryBreakdown {
	id: number;
	name: string;
	total: string;
	percentage: number;
	receiptCount: number;
}

interface ExpensesByStoreBreakdown {
	id: number;
	name: string;
	total: string;
	receiptCount: number;
}
