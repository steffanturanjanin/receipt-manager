import { httpClient } from "../http"
import { categories } from "./mocks"

export const getCategories = async (): Promise<Category[]> => {
	const { data } = await httpClient.get<Category[]>("/categories");

	return data;
}
