import { httpClient } from "../http"
import { receipts } from "./mock/data";

export const createReceipt = async (request: CreateReceiptRequest): Promise<{ message: string }> => {
	const { data } = await httpClient.post("/receipts", request);

	return data;
}


export const getReceipts = async (params: GetReceiptsParams): Promise<PaginatedReceipts> => {
	//const { data } = await httpClient.get("/receipts", { params });

	return receipts;
}
