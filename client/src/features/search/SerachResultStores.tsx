import { Divider, Paper, Stack, StackProps, Typography, styled } from "@mui/material";
import { Fragment, FunctionComponent, ReactElement } from "react";

interface SearchResultStoreProps {
	store: Store;
	divider?: boolean;
}

const SearchResultStore: FunctionComponent<SearchResultStoreProps> = ({ store, divider }): ReactElement => {
	const { name, location, city, address } = store;
	return (
		<Fragment>
			<Stack direction="column" sx={{ padding: "1rem" }}>
				<Typography fontWeight="bold">{name}</Typography>
				<Typography>{location}</Typography>
				<Typography variant="body2" color="grey.700">{`${address}, ${city}`}</Typography>
			</Stack>
			{divider && <Divider />}
		</Fragment>
	)
}

interface SearchResultStoresProps {
	stores: Store[];
}

const SearchResultContainer = styled(Stack)<StackProps>({
	borderRadius: "0.75rem",
	boxShadow: "#959da533 0 0.5rem 1.5rem",
})

const SearchResultStores: FunctionComponent<SearchResultStoresProps> = ({ stores }): ReactElement => {
	if (!stores.length) {
		return <Typography>Nema prodavnica koje odgovaraju pretraženom terminu.</Typography>
	}

	return (
		<SearchResultContainer direction="column" component={Paper}>
			{stores.map((store, index) => (
				<SearchResultStore
					key={store.id}
					store={store}
					divider={index !== stores.length - 1}
				/>
			))}
		</SearchResultContainer>
	)
}

export default SearchResultStores;
