import { httpClient } from "../http"

export const getCategoriesStats = async ({ fromDate, toDate }: CategoriesStatsParams): Promise<CategoryStats> => {
	const { data } = await httpClient.get<CategoryStats>("/stats/categories", {
		params: { fromDate, toDate }
	});

	return data;
}
