import { Divider, Stack, Typography } from "@mui/material";
import { Fragment, FunctionComponent, ReactElement } from "react";
import CardItem from "./CardItem";

interface StoreLocationsProps {
	storeLocations: StoreLocations;
}

const StoreLocations: FunctionComponent<StoreLocationsProps> = ({ storeLocations }): ReactElement => {
	const { data, total, receiptsCount } = storeLocations;

	return (
		<Stack direction="column" gap="1rem">
			<Typography variant="h6" component="h2">Lokacije</Typography>
			<CardItem>
				{data.map((location, index) => (
					<Fragment key={index}>
						<Stack direction="column" gap="0.25rem" padding="1rem">
							<Stack direction="row" justifyContent="space-between">
								<Typography>{location.locationName}</Typography>
								<Typography>{location.amount}</Typography>
							</Stack>
							<Stack direction="row" justifyContent="space-between">
								<Typography color="grey.700">{location.address}</Typography>
								<Typography color="grey.700" variant="body2">{location.receiptCount} račun</Typography>
							</Stack>
						</Stack>
						<Divider />
					</Fragment>
				))}
				<Stack direction="row" justifyContent="space-between" alignItems="center" padding="1rem">
					<Typography>Ukupno:</Typography>
					<Stack direction="column">
						<Typography>{total}</Typography>
						<Typography color="grey.700" variant="body2">{receiptsCount} račun</Typography>
					</Stack>
				</Stack>
			</CardItem>
		</Stack>
	);
}

export default StoreLocations;