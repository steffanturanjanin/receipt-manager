import { httpClient } from "../http";

export const getCategories = async (): Promise<Category[]> => {
	const { data } = await httpClient.get<Category[]>("/categories");

	return data;
}
