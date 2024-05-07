export const categoryStats: CategoryStats = {
	total: "10878.08",
	categories: [
		{
			category: {
				id: 1,
				color: "#4f3531",
				name: "Hrana",
			},
			total: "3228.21"
		},
		{
			category: {
				id: 2,
				color: "#5f8d5e",
				name: "Piće"
			},
			total: "1328.90"
		},
		{
			category: {
				id: 3,
				color: "#85fa33",
				name: "Kućna hemija"
			},
			total: "6128.90"
		},
		{
			category: {
				id: 4,
				color: "#5dade1",
				name: "Lična higijena"
			},
			total: "2328.17"
		},
		{
			category: {
				id: 5,
				color: "#2a57b2",
				name: "Domaćinstvo"
			},
			total: "3252.42"
		},
	]
}

export const expensesYearlyBreakdown: ExpensesDateBreakdown[] = [
	{
		date: "2023-05",
		total: "2324.23",
	},
	{
		date: "2023-06",
		total: "2323.53",
	},
	{
		date: "2023-07",
		total: "54534.53",
	},
	{
		date: "2023-08",
		total: "4244.23",
	},
	{
		date: "2023-09",
		total: "23543.12",
	},
	{
		date: "2023-10",
		total: "23122.26",
	},
	{
		date: "2023-11",
		total: "2324.23",
	},
	{
		date: "2023-12",
		total: "2323.53",
	},
	{
		date: "2024-01",
		total: "54534.53",
	},
	{
		date: "2024-02",
		total: "4244.23",
	},
	{
		date: "2024-03",
		total: "23543.12",
	},
	{
		date: "2024-4",
		total: "23122.26",
	},
];

export const expensesByCategoryYearlyBreakdown: ExpensesByCategoryBreakdown[] = [
	{
		id: 1,
		name: "Automobil",
		total: "1999.20",
		percentage: 71,
		receiptCount: 1,
	},
	{
		id: 2,
		name: "Piće",
		total: "3328.24",
		percentage: 12,
		receiptCount: 4,
	},
	{
		id: 3,
		name: "Tehnika",
		total: "12612.24",
		percentage: 23,
		receiptCount: 1,
	},
	{
		id: 4,
		name: "Hrana",
		total: "31612.33",
		percentage: 44,
		receiptCount: 15,
	},
];

export const expensesByStoreYearlyBreakdown: ExpensesByStoreBreakdown[] = [
	{
		tin: "1",
		name: "NAFTNA INDUSTRIJA SRBIJE",
		total: "2990.00",
		receiptCount: 5,
	},
	{
		tin: "2",
		name: "METLA DISKONT",
		total: "804.70",
		receiptCount: 3,
	},
]
