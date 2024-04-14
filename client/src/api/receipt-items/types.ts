interface UpdateReceiptItemRequest {
	categoryId: number;
}

interface ReceiptItem {
	id: number;
	receiptId: number;
	name: string;
	date: string;
	store: string;
	amount: string;
	unit: string;
	quantity: number;
}
