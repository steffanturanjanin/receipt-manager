import { httpClient } from "../http"

export const updateReceiptItem = async (id: number, request: UpdateReceiptItemRequest): Promise<SingleReceiptReceiptItem> => {
	const { data } = await httpClient.patch<SingleReceiptReceiptItem>(`/receipt-items/${id}`, request)

	return data;
}

export const getReceiptItems = async ({ searchText }: SearchQuery): Promise<ReceiptItem[]> => {
	const { data } = await httpClient.get<ReceiptItem[]>("/receipt-items", {params: { searchText }});

	return data;
}
