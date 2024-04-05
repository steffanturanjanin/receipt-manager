interface CategoriesStatsParams {
	fromDate: string;
	toDate: string;
}

interface CategoryStats {
	category: {
		id: number;
		name: string;
		color: string;
	};
	total: string;
}
