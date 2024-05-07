interface Category {
	id: number;
	name: string;
	color: string;
}

interface MostPopularReceiptItem {
	name: string;
	receiptCount: number;
	total: string;
}

interface MostPopularStore {
	tin: string;
	name: string;
	location: string;
	address: string;
	city: string;
	receiptCount: number;
	total: string;
	percent: number;
}

interface SingleCategoryStats {
	id: number;
	name: string;
	total: string;
	mostPopularReceiptItems: MostPopularReceiptItem[],
	mostPopularStores: MostPopularStore[],
}
