import { httpClient } from "../http"

export const updateReceiptItem = async (id: number, request: UpdateReceiptItemRequest): Promise<SingleReceiptReceiptItem> => {
	const { data } = await httpClient.patch<SingleReceiptReceiptItem>(`/receipt-items/${id}`, request)

	return data;
}
