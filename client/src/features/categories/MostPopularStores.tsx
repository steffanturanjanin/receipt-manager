import { Fragment, FunctionComponent, ReactElement, useMemo } from "react";
import { Divider, Stack, Typography } from "@mui/material";
import Card from "../../components/card/Card";
import CardLinkContent from "../../components/card/CardLinkContent";

interface MostPopularStoresProps {
	stores: MostPopularStore[];
}

const MostPopularStores: FunctionComponent<MostPopularStoresProps> = ({ stores }): ReactElement => {
	const Stores = (
		<Card>
			{stores.map((store, index) => (
				<Fragment key={index}>
					<CardLinkContent to={`/stores/companies/${store.tin}`}>
						<Stack direction="row" justifyContent="space-between" alignItems="center">
							<Stack direction="column">
								<Typography>{store.name}</Typography>
								<Typography>{store.location}</Typography>
								<Typography variant="body2" color="grey.600">
									{`${store.address} - ${store.city}`}
								</Typography>
							</Stack>
							<Stack direction="column" alignItems="end">
								<Typography>{store.total}</Typography>
								<Typography variant="body2" color="grey.600">
									{store.receiptCount} računa
								</Typography>
								<Typography variant="body2" color="grey.600">
									{store.percent}%
								</Typography>
							</Stack>
						</Stack>
					</CardLinkContent>
					{index !== stores.length - 1 && <Divider />}
				</Fragment>
			))}
		</Card>
	);

	const Content = useMemo(
		() => stores.length ? Stores : <Typography>Nema prodavnica...</Typography>,
		[stores]
	);

	return (
		<Stack direction="column" gap="1rem">
			<Typography variant="h6">Potrošnja po prodavnicama</Typography>
			{Content}
		</Stack>
	)
}

export default MostPopularStores;
