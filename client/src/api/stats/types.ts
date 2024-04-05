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
