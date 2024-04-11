import { ChangeEvent, FormEvent, FunctionComponent, ReactElement, useMemo, useState } from "react";
import { useNavigate } from "react-router-dom";
import { useMutation } from "react-query";
import { AxiosError } from "axios";
import { Stack, Typography, TextField, Button } from "@mui/material";
import PageLayout from "../../layouts/PageLayout/PageLayout";
import LoadingButton from "../../../components/LoadingButton";
import { createReceipt } from "../../../api/receipts";
import { FormErrors, FormFieldsTranslator } from "../../../shared/types/errors";
import { transformValidationErrors } from "../../../util/errors";

type CreateReceiptFormErrors = FormErrors<CreateReceiptForm>;

type CreateReceiptForm = {
	url: string;
}

const DEFAULT_CREATE_RECEIPT_FORM: CreateReceiptForm = {
	url: "",
}

const DEFAULT_CREATE_RECEIPT_FORM_ERRORS: CreateReceiptFormErrors = {
	url: "",
};

const createReceiptFormFieldsTranslator: FormFieldsTranslator<CreateReceiptForm> = {
	url: "Url",
}

const CreateReceiptPage: FunctionComponent = (): ReactElement => {
	const navigate = useNavigate();

	const [createReceiptForm, setCreateReceiptForm] = useState<CreateReceiptForm>(
		DEFAULT_CREATE_RECEIPT_FORM
	);
	const [createReceiptFormErrors, setCreateReceiptFormErrors] = useState<CreateReceiptFormErrors>(
		DEFAULT_CREATE_RECEIPT_FORM_ERRORS
	);

	const { isLoading, mutate } = useMutation({
		mutationFn: (request: CreateReceiptRequest) => createReceipt({ url: request.url}),
		onSuccess: () => {
			navigate("/receipts");
		},
		onError: (error: AxiosError) => {
			if (error.response?.status === 422) {
				const validationError = error as AxiosError<ValidationError<CreateReceiptForm>>;
				const { errors } = validationError.response?.data || {};

				if (errors) {
					const errorsTranslation = transformValidationErrors<CreateReceiptForm>(
						errors,
						createReceiptFormFieldsTranslator
					);

					setCreateReceiptFormErrors({
						...createReceiptFormErrors,
						...errorsTranslation,
					});
				}
			}
		}
	})

	const onSubmit = (event: FormEvent<HTMLFormElement>) => {
		event.preventDefault();
		setCreateReceiptFormErrors(DEFAULT_CREATE_RECEIPT_FORM_ERRORS);

		const { url } = createReceiptForm;
		if (url.length) {
			const request = { url };
			mutate(request);
		}
	}

	const onCancel = () => {
		navigate("/receipts", { replace: true });
	}

	const isDisabled = useMemo(() => {
		return !createReceiptForm.url.length;
	}, [createReceiptForm]);

	return (
		<PageLayout title="Ručni unos računa">
			<Stack direction="column" gap="2rem">
				{/* Titles and descriptions */}
				<Stack direction="column" gap="1rem">
					<Typography variant="h6" component="h2">
						Ukoliko imate URL računa, možete uneti račun uz pomoć sledeće forme.
					</Typography>
					<Typography variant="body1" component="p">
						Prosledite URL računa u input polje i započnite obradu.
					</Typography>
				</Stack>
				{/* Form */}
				<Stack component="form" onSubmit={onSubmit} gap="2rem">
					<Stack direction="column">
						<TextField
							id="receiptUrl"
							required
							label="URL računa"
							variant="outlined"
							value={createReceiptForm?.url}
							onChange={(event: ChangeEvent<HTMLInputElement>) =>
								setCreateReceiptForm({...createReceiptForm, url: event.currentTarget.value})
							}
							error={!!createReceiptFormErrors?.url}
							helperText={createReceiptFormErrors?.url}
						/>
					</Stack>
					<Stack direction="row" justifyContent="flex-end">
						<Button variant="text" onClick={onCancel}>Otkaži</Button>
						<LoadingButton type="submit" variant="contained" loading={isLoading} disabled={isDisabled}>
							Unesi račun
						</LoadingButton>
					</Stack>
				</Stack>
			</Stack>
		</PageLayout>
	);
}

export default CreateReceiptPage;
