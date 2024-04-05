interface Receipt {
	id: number;
	amount: string;
	date: string;
	store: {
		name: string;
	},
	categories: string[],
}

export interface ReceiptByDate {
	date: string;
	total: string;
	receipts: Receipt[];
}

export const receipts: ReceiptByDate[] = [
	{
		date: "2024-03-16",
		total: "804,70",
		receipts: [
			{
				id: 1,
				amount: "804,70",
				date: "2024-03-16 15:00",
				store: {
					name: "METLA DISKONT",
				},
				categories: ["Piće", "Hrana"],
			}
		],
	},
	{
		date: "2024-03-17",
		total: "6204,64",
		receipts: [
			{
				id: 2,
				amount: "804,70",
				date: "2024-03-17 11:25",
				store: {
					name: "IDEA",
				},
				categories: ["Piće", "Hrana"],
			},
			{
				id: 3,
				amount: "5804,70",
				date: "2024-03-17 14:28",
				store: {
					name: "MAXI",
				},
				categories: ["Piće", "Hrana"],
			}
		],
	}
]
