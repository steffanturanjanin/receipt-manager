import { CircularProgress, Divider, Stack, Typography } from "@mui/material";
import { Fragment, FunctionComponent, ReactElement } from "react";
import Card from "../../components/card/Card";
import CardLinkContent from "../../components/card/CardLinkContent";

interface SearchResultStoreProps {
	store: Store;
	divider?: boolean;
}

const SearchResultStore: FunctionComponent<SearchResultStoreProps> = ({ store, divider }): ReactElement => {
	const { tin, name, location, city, address } = store;
	return (
		<Fragment>
			<CardLinkContent to={`/stores/companies/${tin}`}>
				<Typography fontWeight="bold">{name}</Typography>
				<Typography>{location}</Typography>
				<Typography variant="body2" color="grey.700">{`${address}, ${city}`}</Typography>
			</CardLinkContent>
			{divider && <Divider />}
		</Fragment>
	)
}

interface SearchResultStoresProps {
	stores: Store[];
	isLoading?: boolean;
}

const SearchResultStores: FunctionComponent<SearchResultStoresProps> = ({ stores, isLoading }): ReactElement => {
	if (isLoading) {
		return (
			<Stack alignItems="center">
				<CircularProgress color="primary"/>
			</Stack>
		)
	}

	if (!stores.length && !isLoading) {
		return <Typography>Nema prodavnica koje odgovaraju pretraženom terminu.</Typography>
	}

	return (
		<Card>
			{stores.map((store, index) => (
				<SearchResultStore
					key={index}
					store={store}
					divider={index !== stores.length - 1}
				/>
			))}
		</Card>
	)
}

export default SearchResultStores;
