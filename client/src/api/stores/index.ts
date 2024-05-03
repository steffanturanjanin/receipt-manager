import { httpClient } from "../http"
import { store, storesList } from "./mocks";

export const getStoresSearch = async ({ searchText }: SearchQuery): Promise<StoreSearch[]> => {
	const { data } = await httpClient.get<StoreSearch[]>("/stores", { params: { searchText }});

	return data;
}

export const getStoresList = async ({ searchText }: Partial<SearchQuery>): Promise<StoreListItem[]> => {
	if (searchText === undefined) {
		return storesList;
	}

	return storesList.filter((store) => store.name.toLowerCase().includes(searchText.toLocaleLowerCase()));
}

export const getStore = async (tin: string): Promise<Store> => {
	// const { data } = await httpClient.get<Store>(`/stores/${tin}`);

	// return data;

	return store;
}
