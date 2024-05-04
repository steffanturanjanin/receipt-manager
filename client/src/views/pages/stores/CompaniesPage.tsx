import { ChangeEvent, FormEvent, FunctionComponent, ReactElement, useMemo, useState } from "react";
import { useQuery } from "react-query";
import { Stack, TextField, Typography } from "@mui/material";
import { getCompanies } from "../../../api/stores";
import PageLayout from "../../layouts/PageLayout";
import CompanyList from "../../../features/stores/CompanyList";

const CompaniesPage: FunctionComponent = (): ReactElement => {
	const [searchTerm, setSearchTerm] = useState<string>();
	const [searchInput, setSearchInput] = useState<string>("");

	const onSubmit = (event: FormEvent<HTMLFormElement>) => {
		event.preventDefault();
		setSearchTerm(searchInput);
	}

	const { isLoading, data: companies } = useQuery({
		queryKey: ["company_list", searchTerm],
		queryFn: () => getCompanies({ searchText: searchTerm }),
	});

	const empty = useMemo(
		() => !isLoading && companies?.length === 0,
		[companies, isLoading]
	);

	return (
		<PageLayout title="Prodavnice" showBackdrop={isLoading}>
			<Stack direction="column" gap="2rem">
				<Stack component="form" onSubmit={onSubmit}>
					<TextField
						fullWidth
						size="small"
						variant="outlined"
						placeholder="Pretraži artikle i prodavnice"
						value={searchInput}
						onChange={(event: ChangeEvent<HTMLInputElement>) => setSearchInput(event.target.value)}
						sx={{ backgroundColor: "white" }} />
				</Stack>
				{empty ?
					<Typography variant="h4">Nema rezultata</Typography> :
					<CompanyList companies={companies || []} />
				}
			</Stack>
		</PageLayout>
	)
}

export default CompaniesPage;
