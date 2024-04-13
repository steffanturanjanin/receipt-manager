import { httpClient } from "../http"
import { receiptItems } from "./mock";

export const updateReceiptItem = async (id: number, request: UpdateReceiptItemRequest): Promise<SingleReceiptReceiptItem> => {
	const { data } = await httpClient.patch<SingleReceiptReceiptItem>(`/receipt-items/${id}`, request)

	return data;
}

export const getReceiptItems = async ({ searchText }: SearchQuery): Promise<ReceiptItem[]> => {
	// const { data } = await httpClient.get<ReceiptItem[]>("/receipt-items");

	// return data;

	const search = searchText.toLowerCase();

	const filteredReceiptItems = receiptItems.filter((receiptItem) => {
		const name = receiptItem.name.toLocaleLowerCase();
		const store = receiptItem.store.toLowerCase();

		return name.includes(search) || store.includes(search);
	})

	return filteredReceiptItems;
}
