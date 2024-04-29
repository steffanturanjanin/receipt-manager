import { ChangeEvent, FormEvent, Fragment, FunctionComponent, ReactElement, useState } from "react";
import { useQuery } from "react-query";
import { Box, Divider, Paper, Stack, TextField } from "@mui/material";
import PageLayout from "../../layouts/PageLayout";
import StoreListItem from "../../../features/stores/StoreListItem";
import { getStoresList } from "../../../api/stores";

const StoresPage: FunctionComponent = (): ReactElement => {
	const [searchTerm, setSearchTerm] = useState<string>();
	const [searchInput, setSearchInput] = useState<string>("");

	const onSubmit = (event: FormEvent<HTMLFormElement>) => {
		event.preventDefault();
		setSearchTerm(searchInput);
	}

	const { data: stores } = useQuery({
		queryKey: ["stores_list", searchTerm],
		queryFn: () => getStoresList({ searchText: searchTerm }),
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
				<Stack direction="column">
					<Box component={Paper} sx={{ borderRadius: "0.5rem" }}>
						{stores?.map((store, index) => (
							<Fragment>
								<StoreListItem {...store} />
								{index !== stores.length - 1 && <Divider />}
							</Fragment>
						))}
					</Box>
				</Stack>
			</Stack>
		</PageLayout>
	)
}

export default StoresPage;
