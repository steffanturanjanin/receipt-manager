interface CreateReceiptRequest {
	url: string;
}

interface SetFavoriteRequest {
	isFavorite: boolean;
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
	categories: {
		id: number;
		name: string;
		color: string;
	}[],
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
	isFavorite: boolean;
	receiptItems: SingleReceiptReceiptItem[];
	store: SingleReceiptStore;
	createdAt: string;
}

interface FavoriteReceiptStore {
	id: number;
	name: string;
	location: string;
	address: string;
	city: string;
}

interface FavoriteReceiptCategory {
	id: number;
	name: string;
	color: string;
}

interface FavoriteReceipt {
	id: number;
	amount: string;
	date: string;
	store: FavoriteReceiptStore;
	categories: FavoriteReceiptCategory[];
}
