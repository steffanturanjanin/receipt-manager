import { FormEvent, FunctionComponent, ReactElement, useState, MouseEvent, ChangeEvent } from "react";
import PageLayout from "../../layouts/PageLayout";
import { Stack, TextField, ToggleButton, ToggleButtonGroup, Typography } from "@mui/material";
import { useQuery } from "react-query";
import { getStores } from "../../../api/stores";
import SearchResultStores from "../../../features/search/SerachResultStores";
import { getReceiptItems } from "../../../api/receipt-items";
import SearchResultArticles from "../../../features/search/SearchResultArticles";

type SearchCriteria = "articles" | "stores";

const SearchPage: FunctionComponent = (): ReactElement => {
	const [searchTerm, setSearchTerm] = useState<string>();
	const [searchInput, setSearchInput] = useState<string>("");
	const [searchCriteria, setSearchCriteria] = useState<SearchCriteria>();

	const onCriteriaChange = (_: MouseEvent<HTMLElement>, criteria: SearchCriteria) => {
		setSearchCriteria(criteria)
	}

	const onSubmit = (event: FormEvent<HTMLFormElement>) => {
		event.preventDefault();
		setSearchTerm(searchInput);

		if (!searchCriteria) {
			setSearchCriteria("articles");
		}
	}

	const { data: stores } = useQuery({
		queryKey: ["stores", searchCriteria, searchTerm],
		queryFn: () => getStores({ searchText: searchTerm! }),
		enabled:
			(searchCriteria === "stores" || searchCriteria === undefined) &&
			!!searchTerm?.length
	});

	const { data: receiptItems } = useQuery({
		queryKey: ["receipt_items", searchCriteria, searchTerm],
		queryFn: () => getReceiptItems({ searchText: searchTerm! }),
		enabled:
			(searchCriteria === "articles" || searchCriteria === undefined) &&
			!!searchTerm?.length
	});

	return (
		<PageLayout title="Pretraga">
			<Stack direction="column" gap="2rem">
				<Stack component="form" onSubmit={onSubmit}>
					<TextField
						fullWidth
						size="small"
						variant="outlined"
						placeholder="Pretraži artikle i predovnice"
						value={searchInput}
						onChange={(event: ChangeEvent<HTMLInputElement>) => setSearchInput(event.target.value)}
						sx={{ backgroundColor: "white"}} />
				</Stack>

				<Stack direction="column" gap="1rem">
					{searchTerm === undefined && (
						<Typography variant="h5">Unesite termin u polje da započnete pretragu...</Typography>
					)}

					{(searchCriteria && searchTerm) &&
						<ToggleButtonGroup
								color="primary"
								value={searchCriteria}
								exclusive
								onChange={onCriteriaChange}
								aria-label="Platform"
								fullWidth
						>
								<ToggleButton value="articles" sx={{ flex: 1, padding: "0.5rem" }}>Artikli</ToggleButton>
								<ToggleButton value="stores" sx={{ flex: 1, padding: "0.5rem" }}>Prodavnice</ToggleButton>
						</ToggleButtonGroup>
					}

					{searchCriteria === "articles" && <SearchResultArticles receiptItems={receiptItems || []} /> }
					{searchCriteria === "stores" && <SearchResultStores stores={stores || []} />}
				</Stack>
			</Stack>
		</PageLayout>
	)
}

export default SearchPage;
