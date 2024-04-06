interface CreateReceiptRequest {
	url: string;
}

enum ReceiptStatus {
	PENDING = "pending",
	PROCESSED = "processed",
}

interface ReceiptAggregatedByDateItem {
	id: number;
	amount: string;
	date: string;
	store: {
		name: string;
	},
	categories: string[],
}

interface ReceiptsAggregatedByDate {
	date: string;
	total: string;
	receipts: ReceiptAggregatedByDateItem[];
}

interface GetReceiptsParams {
	fromDate: string;
	toDate: string;
}

// Single Receipt interfaces

interface SingleReceiptCategory {
	id: number;
	name: string;
	color: string;
}

interface SingleReceiptReceiptItemTax {
	id: number;
	identifier: string;
	name: string;
	rate: number;
}

interface SingleReceiptUser {
	id: number;
	firstName: string;
	lastName: string;
	email: string;
}

interface SingleReceiptStore {
	id: number;
	tin: string;
	name: string;
	locationId: string;
	locationName: string;
	address: string;
	city: string;
}

interface SingleReceiptReceiptItem {
	id: number;
	name: string;
	unit: string;
	quantity: number;
	singleAmount: string;
	totalAmount: string;
	category: SingleReceiptCategory | null;
	tax: SingleReceiptReceiptItemTax;
}

interface SingleReceipt {
	id: number;
	user: SingleReceiptUser;
	status: ReceiptStatus;
	pfrNumber: string;
	counter: string;
	totalPurchaseAmount: string;
	totalTaxAmount: string;
	date: string;
	meta: object;
	qrCode: string;
	receiptItems: SingleReceiptReceiptItem[];
	store: SingleReceiptStore;
	createdAt: string;
}
