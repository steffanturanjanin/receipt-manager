import { httpClient } from "../http"
import { storesList } from "./mocks";

export const getStoresSearch = async ({ searchText }: SearchQuery): Promise<StoreSearch[]> => {
	const { data } = await httpClient.get<StoreSearch[]>("/stores", { params: { searchText }});

	return data;
}

export const getStoresList = async ({ searchText }: Partial<SearchQuery>): Promise<ExpensesByStoreBreakdown[]> => {
	if (searchText === undefined) {
		return storesList;
	}

	return storesList.filter((store) => store.name.toLowerCase().includes(searchText.toLocaleLowerCase()));
}
