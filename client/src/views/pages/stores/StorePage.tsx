import { FunctionComponent, ReactElement } from "react";
import { useQuery } from "react-query";
import { useParams } from "react-router-dom";
import { Divider, Stack, Typography } from "@mui/material";
import { getStore } from "../../../api/stores";
import PageLayout from "../../layouts/PageLayout";
import CardItem from "../../../features/stores/CardItem";
import StoreLocations from "../../../features/stores/StoreLocations";
import StoreExpenses from "../../../features/stores/StoreExpenses";

const StorePage: FunctionComponent = (): ReactElement => {
	const { tin } = useParams();

	const { data: store } = useQuery({
		queryKey: ["store", tin],
		queryFn: () => getStore(tin!),
		enabled: !!tin,
	});

	return (
		<PageLayout title={store?.name || ""}>
			<Stack direction="column" gap="2rem">

				<CardItem>
					<Stack direction="row" justifyContent="space-between" sx={{ padding: "1rem" }}>
						<Typography>Ime:</Typography>
						<Typography fontWeight="bold">{store?.name}</Typography>
					</Stack>
					<Divider />
					<Stack direction="row" justifyContent="space-between" sx={{ padding: "1rem" }}>
						<Typography>PIB:</Typography>
						<Typography fontWeight="bold">{store?.tin}</Typography>
					</Stack>
				</CardItem>

				{store?.locations &&
					<StoreLocations storeLocations={store.locations} />
				}

				{store?.expenses &&
					<StoreExpenses storeExpenses={store.expenses} />
				}

			</Stack>
		</PageLayout>
	)
}

export default StorePage;
