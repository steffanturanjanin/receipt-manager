import { FormErrors, FormFieldsTranslator } from "../shared/types/errors";

export const transformValidationErrors = <Form extends object>(
	validationErrors: ValidationError<Form>["errors"],
	formErrorsTranslator: FormFieldsTranslator<Form>
): FormErrors<Form> => {
	const translatedErrorsEntries = Object
		.entries(validationErrors)
		.map(([field, error]) => ([field, `${formErrorsTranslator[field as keyof Form]} ${error}`]));

	return Object.fromEntries(translatedErrorsEntries)
}
