import { ChangeEvent, FormEvent, FunctionComponent, ReactElement, useState } from "react";
import { useQuery } from "react-query";
import { Stack, TextField } from "@mui/material";
import PageLayout from "../../layouts/PageLayout";
import { getCompanies } from "../../../api/stores";
import CompanyList from "../../../features/stores/CompanyList";

const CompaniesPage: FunctionComponent = (): ReactElement => {
	const [searchTerm, setSearchTerm] = useState<string>();
	const [searchInput, setSearchInput] = useState<string>("");

	const onSubmit = (event: FormEvent<HTMLFormElement>) => {
		event.preventDefault();
		setSearchTerm(searchInput);
	}

	const { data: companies } = useQuery({
		queryKey: ["company_list", searchTerm],
		queryFn: () => getCompanies({ searchText: searchTerm }),
	})

	return (
		<PageLayout title="Prodavnice">
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
				<CompanyList companies={companies || []} />
			</Stack>
		</PageLayout>
	)
}

export default CompaniesPage;
