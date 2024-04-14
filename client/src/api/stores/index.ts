import { httpClient } from "../http"

export const getStores = async ({ searchText }: SearchQuery): Promise<Store[]> => {
	const { data } = await httpClient.get<Store[]>("/stores", { params: { searchText }});

	return data;
}
