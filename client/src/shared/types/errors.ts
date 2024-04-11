export type FormErrors<Form extends object> = Record<keyof Form, string>;

export type FormFieldsTranslator<Form extends object> = Record<keyof Form, string>;
