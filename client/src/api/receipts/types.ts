interface CreateReceiptRequest {
	url: string;
}

enum ReceiptStatus {
	PENDING = "pending",
	PROCESSED = "processed",
}

interface Store {
	id: number;
	tin: string;
	name: string;
	locationId: string;
	locationName: string;
	address: string;
	city: string;
}

interface Category {
	id: number;
	name: string;
}

interface Tax {
	id: number;
	identifier: string;
	name: string;
	rate: number;
}

interface ReceiptItem {
	id: number;
	receiptId: number;
	name: string;
	unit: string;
	quantity: number;
	singleAmount: number;
	totalAmount: number;
	category: Category;
	tax: Tax;
}

interface Receipt {
	id: number;
	userId: number;
	//status: ReceiptStatus;
	status: string;
	pfrNumber?: string;
	counter?: string;
	totalPurchaseAmount: number;
	totalTaxAmount: number;
	//date: Date;
	date: string;
	qrCode: string;
	meta?: object;
	receiptItems: ReceiptItem[];
	store: Store;
}


type PaginatedReceipts = Paginated<Receipt> & { meta: Paginated<Receipt>["meta"] & { total: number } }


interface GetReceiptsParams {
	fromDate: string;
	toDate: string;
}
