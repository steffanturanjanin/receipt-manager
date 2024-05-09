import { httpClient } from "../http";

export const getCategories = async (): Promise<Category[]> => {
	const { data } = await httpClient.get<Category[]>("/categories");

	return data;
}

export const getCategoryStats = async (id: number): Promise<SingleCategoryStats> => {
	const { data } = await httpClient.get<SingleCategoryStats>(`/stats/categories/${id}`);

	return data;
}
