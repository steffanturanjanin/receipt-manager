interface Store {
	tin: string;
	name: string;
	location: string;
	city: string;
	address: string;
}

interface CompanyListItem {
	tin: string;
	name: string;
	total: string;
	receiptCount: number;
}

interface CompanyLocation {
	locationId: string;
	locationName: string;
	address: string;
	city: string;
	amount: string;
	receiptCount: number;
}

interface CompanyLocations {
	data: CompanyLocation[];
	total: string;
	receiptsCount: number;
}

interface CompanyExpense {
	id: number;
	locationName: string;
	date: string;
	amount: string;
}

interface Company {
	tin: string;
	name: string;
	locations: CompanyLocations;
	expenses: CompanyExpense[];
}
