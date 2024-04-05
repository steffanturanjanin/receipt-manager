import { httpClient } from "../http"

export const createReceipt = async (request: CreateReceiptRequest): Promise<{ message: string }> => {
	const { data } = await httpClient.post("/receipts", request);

	return data;
}


export const getReceiptsAggregatedByDate = async (params: GetReceiptsParams): Promise<ReceiptsAggregatedByDate[]> => {
	const { data } = await httpClient.get<ReceiptsAggregatedByDate[]>("/stats/receipts", { params });

	return data;
}
