import { httpClient } from "../http";
import { categoryStats } from "./mocks";

export const getCategories = async (): Promise<Category[]> => {
	const { data } = await httpClient.get<Category[]>("/categories");

	return data;
}

export const getCategoryStats = async (id: number): Promise<SingleCategoryStats> => {
	return categoryStats;
}
