
interface ValidationError<T extends object> {
	message: string;
	code: number;
	errors: {
		[key in keyof T]: string;
	}
}

interface BadRequestError {
	error: string;
	code: number;
}

interface PaginationMeta {
	page: number;
	prevPage: number | null;
	nextPage: number | null;
	perPage: number;
	totalPages: number;
	totalEntries: number;
}

interface Paginated<T> {
	data: T[];
	meta: {
		pagination: PaginationMeta;
	}
}
