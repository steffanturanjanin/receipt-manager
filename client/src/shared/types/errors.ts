type FormErrors<FormType extends object> = {
	[key in keyof FormType]: string;
}

type FormFieldsTranslation<FormType extends object> = Record<keyof FormType, string>;
