interface StoreSearch {
	id: number;
	name: string;
	location: string;
	city: string;
	address: string;
}

interface StoreListItem {
	id: number;
	name: string;
	total: string;
	receiptCount: number;
}

interface StoreLocation {
	locationId: string;
	locationName: string;
	address: string;
	city: string;
	amount: string;
	receiptCount: number;
}

interface StoreLocations {
	data: StoreLocation[];
	total: string;
	receiptsCount: number;
}

interface StoreExpenses {
	id: number;
	locationName: string;
	date: string;
	amount: string;
}

interface Store {
	id: number;
	name: string;
	tin: string;
	locations: StoreLocations;
	expenses: StoreExpenses[];
}
