export const receipts: ReceiptsAggregatedByDate[] = [
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
];

export const singleReceipt: SingleReceipt = {
	id: 15,
	user: {
		id: 1,
		firstName: "Stefan",
		lastName: "Turanjanin",
		email: "stefan@turanjanin.net",
	},
	isFavorite: false,
	status: "pending" as ReceiptStatus,
	pfrNumber: "746DUV64-746DUV64-16901",
	counter: "16890/16901ПП",
	totalPurchaseAmount: "2880.82",
	totalTaxAmount: "480.14",
	date: "2022-12-31T15:55:32Z",
	createdAt: "2022-12-31T15:55:32Z",
	qrCode: "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",
	meta: {
			"ЕСИР број": "431/2.0.0.2",
			"Касир": "1011902"
	},
	receiptItems: [
			{
				id: 3,
				name: "VRAT SUVI SLAJS 100G KOSMAJSKA D",
				unit: "KOM",
				quantity: 1,
				singleAmount: "389.99",
				totalAmount: "389.99",
				category: {
						id: 28,
						name: "Slatkiši i grickalice",
						color: "#4C2F27",
				},
					tax: {
						id: 6,
						identifier: "Ђ",
						name: "О-ПДВ",
						rate: 20
					}
			},
			{
				id: 4,
				name: "KOLAC MALA GRCKA TULUMBA 500G DO",
				unit: "KOM",
				quantity: 1,
				singleAmount: "379.99",
				totalAmount: "379.99",
				category: {
					id: 28,
					name: "Slatkiši i grickalice",
					color: "#4A192C",
				},
				tax: {
					id: 7,
					identifier: "Ђ",
					name: "О-ПДВ",
					rate: 20
				}
			},
			{
				id: 5,
				name: "SOK SCHWEPPES B L 1,5L COCA COLA",
				unit: "KOM",
				quantity: 1,
				singleAmount: "113.99",
				totalAmount: "113.99",
				category: {
					id: 31,
					name: "Sokovi i voda",
					color: "#A03472",
				},
				tax: {
					id: 8,
					identifier: "Ђ",
					name: "О-ПДВ",
					rate: 20
				}
			},
			{
				id: 6,
				name: "CIDER SOMERSBY JAB 0.33L NPB",
				unit: "KOM",
				quantity: 1,
				singleAmount: "99.99",
				totalAmount: "99.99",
				category: {
					id: 30,
					name: "Alkoholna pića",
					color: "#293133",
				},
				tax: {
					id: 9,
					identifier: "Ђ",
					name: "О-ПДВ",
					rate: 20
				}
			},
			{
				id: 7,
				name: "PIVO BUDWEISER 0,5L CAN CARLSBER",
				unit: "KOM",
				quantity: 4,
				singleAmount: "107.99",
				totalAmount: "431.96",
				category: {
					id: 30,
					name: "Alkoholna pića",
					color: "#7E7B52",
				},
				tax: {
						id: 10,
						identifier: "Ђ",
						name: "О-ПДВ",
						rate: 20
					}
			},
			{
				id: 8,
				name: "SOK POMORANDZA LIFE 100% 1L NECT",
				unit: "KOM",
				quantity: 1,
				singleAmount: "179.99",
				totalAmount: "179.99",
				category: {
					id: 31,
					name: "Sokovi i voda",
					color: "#909090",
				},
				tax: {
					id: 11,
					identifier: "Ђ",
					name: "О-ПДВ",
					rate: 20
				}
			},
			{
				id: 9,
				name: "VL MARAMICE BABY PLAVE 72/1 NEVE",
				unit: "KOM",
				quantity: 1,
				singleAmount: "109.99",
				totalAmount: "109.99",
				category: {
					id: 38,
					name: "Domaćinstvo",
					color: "#999950",
				},
				tax: {
					id: 12,
					identifier: "Ђ",
					name: "О-ПДВ",
					rate: 20
				}
			},
			{
				id: 10,
				name: "ETIL ALKOHOL 70% 1L VINEX",
				unit: "KOM",
				quantity: 1,
				singleAmount: "199.99",
				totalAmount: "199.99",
				category: {
					id: 30,
					name: "Alkoholna pića",
					color: "#B32821",
				},
				tax: {
					id: 13,
					identifier: "Ђ",
					name: "О-ПДВ",
					rate: 20
				}
			},
			{
				id: 11,
				name: "BAR BOUNTY 57G MARS",
				unit: "KOM",
				quantity: 2,
				singleAmount: "74.99",
				totalAmount: "149.98",
				category: {
					id: 28,
					name: "Slatkiši i grickalice",
					color: "#7FB5B5",
				},
				tax: {
					id: 14,
					identifier: "Ђ",
					name: "О-ПДВ",
					rate: 20
				}
			},
			{
				id: 12,
				name: "PEDIGREE RODEO GOVEDINA 140G",
				unit: "KOM",
				quantity: 1,
				singleAmount: "159.99",
				totalAmount: "159.99",
				category: {
					id: 23,
					name: "Mesne prerađevine",
					color: "#DC9D00",
				},
				tax: {
					id: 15,
					identifier: "Ђ",
					name: "О-ПДВ",
					rate: 20
				}
			},
			{
				id: 13,
				name: "BAR KINDER BUENO WHITE 39G FERRE",
				unit: "KOM",
				quantity: 1,
				singleAmount: "74.99",
				totalAmount: "74.99",
				category: {
					id: 28,
					name: "Slatkiši i grickalice",
					color: "#A52019",
				},
				tax: {
					id: 16,
					identifier: "Ђ",
					name: "О-ПДВ",
					rate: 20
				}
			},
			{
				id: 14,
				name: "ZV GUMA ORBIT EUKALIPTU 14G WRIG",
				unit: "KOM",
				quantity: 1,
				singleAmount: "59.99",
				totalAmount: "59.99",
				category: {
					id: 28,
					name: "Slatkiši i grickalice",
					color: "#5B3A29",
				},
				tax: {
					id: 17,
					identifier: "Ђ",
					name: "О-ПДВ",
					rate: 20
				}
			},
			{
				id: 15,
				name: "MAJONEZ 280ML DOJPAK THOMY",
				unit: "KOM",
				quantity: 1,
				singleAmount: "149.99",
				totalAmount: "149.99",
				category: {
					id: 27,
					name: "Namirnice za pripremu jela",
					color: "#F8F32B",
				},
				tax: {
					id: 18,
					identifier: "Ђ",
					name: "О-ПДВ",
					rate: 20
				}
			},
			{
				id: 16,
				name: "TP KAMILICA 3SL 150L 10/1 NATUR",
				unit: "KOM",
				quantity: 1,
				singleAmount: "379.99",
				totalAmount: "379.99",
				category: {
					id: 32,
					name: "Čaj i kafa",
					color: "#75151E",
				},
				tax: {
					id: 19,
					identifier: "Ђ",
					name: "О-ПДВ",
					rate: 20
				}
			}
	],
	store: {
		id: 12,
		tin: "101670560",
		name: "METLA DISKONT",
		locationId: "1108934",
		locationName: "Diskont br.32",
		address: "Knjazebacka 163",
		city: "Nis"
	}
}

export const favoriteReceipts: FavoriteReceipt[] = [
	{
		id: 1,
		amount: "2455.99",
		date: "2024-04-01 16:12:54",
		store: {
			id: 1,
			name: "METLA DISKONTI",
			location: "Diskont br. 32",
			address: "Knjazevacka 209",
			city: "Nis - Pantelej",
		},
		categories: [
			{
				id: 1,
				name: "Hrana",
				color: "#D95030"
			},
			{
				id: 2,
				name: "Piće",
				color: "#6C6874",
			},
		]
	},
	{
		id: 2,
		amount: "5455.70",
		date: "2024-3-23 12:02:11",
		store: {
			id: 1,
			name: "IDEA D.O.O.",
			location: "Idea Nis 118",
			address: "Matejevacki put 33",
			city: "Nis - Pantelej",
		},
		categories: [
			{
				id: 1,
				name: "Hrana",
				color: "#D95030"
			},
			{
				id: 2,
				name: "Mleko i jaja ",
				color: "#1B5583",
			},
			{
				id: 3,
				name: "Kućna hemija",
				color: "#84C3BE",
			}
		]
	},
	{
		id: 3,
		amount: "4210.56",
		date: "2024-3-11 09:31:20",
		store: {
			id: 1,
			name: "LIDL",
			location: "LIDL Nis - Duvaniste",
			address: "Vizantijski bulevar 28",
			city: "Nis - Pantelej",
		},
		categories: [
			{
				id: 1,
				name: "Hrana",
				color: "#D95030"
			},
			{
				id: 2,
				name: "Piće",
				color: "#1B5583",
			},
			{
				id: 3,
				name: "Domaćinstvo",
				color: "#D36E70",
			}
		]
	},
	{
		id: 1,
		amount: "2455.99",
		date: "2024-04-01 16:12:54",
		store: {
			id: 1,
			name: "METLA DISKONTI",
			location: "Diskont br. 32",
			address: "Knjazevacka 209",
			city: "Nis - Pantelej",
		},
		categories: [
			{
				id: 1,
				name: "Hrana",
				color: "#D95030"
			},
			{
				id: 2,
				name: "Piće",
				color: "#6C6874",
			},
		]
	},
	{
		id: 2,
		amount: "5455.70",
		date: "2024-3-23 12:02:11",
		store: {
			id: 1,
			name: "IDEA D.O.O.",
			location: "Idea Nis 118",
			address: "Matejevacki put 33",
			city: "Nis - Pantelej",
		},
		categories: [
			{
				id: 1,
				name: "Hrana",
				color: "#D95030"
			},
			{
				id: 2,
				name: "Mleko i jaja ",
				color: "#1B5583",
			},
			{
				id: 3,
				name: "Kućna hemija",
				color: "#84C3BE",
			}
		]
	},
	{
		id: 3,
		amount: "4210.56",
		date: "2024-3-11 09:31:20",
		store: {
			id: 1,
			name: "LIDL",
			location: "LIDL Nis - Duvaniste",
			address: "Vizantijski bulevar 28",
			city: "Nis - Pantelej",
		},
		categories: [
			{
				id: 1,
				name: "Hrana",
				color: "#D95030"
			},
			{
				id: 2,
				name: "Piće",
				color: "#1B5583",
			},
			{
				id: 3,
				name: "Domaćinstvo",
				color: "#D36E70",
			}
		]
	},
]
