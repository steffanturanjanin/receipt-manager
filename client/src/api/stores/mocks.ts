export const stores: Store[] = [
	{
		tin: "1",
		name: "METLA DISKONT",
		location: "Diskont br.32",
		city: "Ниш - Пантелеј",
		address: "КЊАЖЕВАЧКА 205",
	},
	{
		tin: "2",
		name: "NAFTNA INDUSTRIJA SRBIJE",
		location: "БС Ниш Исток",
		city: "Ниш - Пантелеј",
		address: "МАТЕЈЕВАЧКИ ПУТ ББ",
	},
	{
		tin: "3",
		name: "MERCATOR-S",
		location: "Roda Megamarket 345",
		city: "Ниш - Пантелеј",
		address: "ВИЗАНТИЈСКИ БУЛЕВАР 1",
	},
	{
		tin: "4",
		name: "PRIVREDNO DRUŠTVO ZAPLANJKA KOMERC DOO NIŠ",
		location: "Продавница бр.9",
		city: "Ниш - Пантелеј",
		address: "ЗОРАНА РАДОСАВЉЕВИЋА ЧУПЕ 1",
	}
];

export const companiesList: CompanyListItem[] = [
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
	{
		tin: "3",
		name: "MERCATOR-S",
		total: "3213.50",
		receiptCount: 4,
	},
	{
		tin: "4",
		name: "PRIVREDNO DRUŠTVO ZAPLANJKA KOMERC DOO NIŠ",
		total: "6501.20",
		receiptCount: 5,
	}
];

export const company: Company = {
	tin: "12345",
	name: "MERCATOR - S D.O.O",
	locations: {
		data: [
			{
				locationId: "123",
				locationName: "Roda Megamarker 345",
				address: "ВИЗАНТИЈСКИ БУЛЕВАР 1",
				city: "Ниш",
				amount: "2880.82",
				receiptCount: 4,
			},
			{
				locationId: "1524",
				locationName: "Roda Merkator 1",
				address: "Bulevar Nikole Tesle 23",
				city: "Ниш",
				amount: "4150.28",
				receiptCount: 7,
			}
		],
		total: "7540.36",
		receiptsCount: 11,
	},
	expenses: [
		{
			id: 1,
			locationName: "Roda Megamarker 345",
			date: "30-12-2023 15:55:00",
			amount: "322.90",
		},
		{
			id: 2,
			locationName: "Roda Megamarker 345",
			date: "30-12-2023 11:55:00",
			amount: "322.90",
		},
		{
			id: 3,
			locationName: "Roda Megamarker 345",
			date: "30-12-2023 12:27:00",
			amount: "322.90",
		},
		{
			id: 4,
			locationName: "Roda Merkator",
			date: "11-12-2023 09:55:00",
			amount: "5822.90",
		},
		{
			id: 5,
			locationName: "Roda Merkator",
			date: "11-12-2023 21:22:00",
			amount: "1530.50",
		}
	]
}
