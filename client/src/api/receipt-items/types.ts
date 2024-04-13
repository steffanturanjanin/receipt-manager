interface UpdateReceiptItemRequest {
	categoryId: number;
}

interface ReceiptItem {
	id: number;
	name: string;
	date: string;
	store: string;
	amount: string;
	unit: string;
	quantity: number;
}
