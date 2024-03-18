
interface ValidationError<T extends object> {
	message: string;
	code: number;
	errors: {
		[key in keyof T]: string;
	}
}
