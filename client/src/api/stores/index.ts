import { httpClient } from "../http"

export const getStores = async ({ searchText }: SearchQuery): Promise<Store[]> => {
	const { data } = await httpClient.get<Store[]>("/stores", { params: { searchText }});

	return data;
}

export const getCompanies = async ({ searchText }: Partial<SearchQuery>): Promise<CompanyListItem[]> => {
	const { data } = await httpClient.get<CompanyListItem[]>("/stores/companies", { params: { searchText } });

	return data;
}

export const getCompany = async (tin: string): Promise<Company> => {
	const { data } = await httpClient.get<Company>(`/stores/companies/${tin}`);

	return data;
}
