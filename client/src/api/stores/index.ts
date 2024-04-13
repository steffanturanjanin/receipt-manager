import { httpClient } from "../http"
import { stores } from "./mocks"

export const getStores = async ({ searchText }: SearchQuery): Promise<Store[]> => {
	// const { data } = await httpClient.get<Store[]>("/stores");

	// return data;

	const search = searchText.toLowerCase();

	const filteredStores = stores.filter((store) => {
		const name = store.name.toLowerCase();
		const location = store.location.toLowerCase();
		const address = store.address.toLowerCase();
		const city = store.city.toLowerCase();

		return name.includes(search) || location.includes(search) || address.includes(search) || city.includes(search);
	});

	return filteredStores;
}
