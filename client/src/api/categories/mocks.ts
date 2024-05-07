export const categories: Category[] = [
	{
		id: 1,
		name: "Hrana",
		color: " #FF7514",
	},
	{
		id: 2,
		name: "Prevoz",
		color: "#343E40",
	},
	{
		id: 3,
		name: "Tehnologija",
		color: "#7FB5B5",
	}
];

export const categoryStats: SingleCategoryStats = {
	id: 1,
	name: "Hrana",
	total: "7554.89",
	mostPopularReceiptItems: [
		{
			name: "Hleb",
			receiptCount: 2,
			total: "2559,90",
		},
		{
			name: "Coca-cola",
			receiptCount: 8,
			total: "1778,28",
		},
		{
			name: "Meso juneće 100%",
			receiptCount: 8,
			total: "1778,28",
		}
	],
	mostPopularStores: [
		{
			tin: "1",
			name: "METLA DISKONTI",
			location: "Diskont br. 32",
			address: "Knjaževačka 223",
			city: "Niš - Pantelej",
			receiptCount: 4,
			total: "2555.25",
			percent: 33,
		},
		{
			tin: "2",
			name: "IDEA D.O.O.",
			location: "Diskont br. 54",
			address: "Matejevački put 25",
			city: "Niš - Pantelej",
			receiptCount: 7,
			total: "5212.25",
			percent: 28,
		},
		{
			tin: "3",
			name: "DIS D.O.O.",
			location: "Diskont br. 12",
			address: "Majakovskog 32",
			city: "Niš - Duvanište",
			receiptCount: 7,
			total: "12500.25",
			percent: 38,
		}
	]
}
