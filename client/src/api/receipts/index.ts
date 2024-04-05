import { httpClient } from "../http"
import { ReceiptByDate, receipts } from "./mock";

export const createReceipt = async (request: CreateReceiptRequest): Promise<{ message: string }> => {
	const { data } = await httpClient.post("/receipts", request);

	return data;
}


export const getReceipts = async (params: GetReceiptsParams): Promise<ReceiptByDate[]> => {
	//const { data } = await httpClient.get("/receipts", { params });

	return receipts;
}
