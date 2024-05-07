import { FunctionComponent, ReactElement } from "react";
import { Box, Divider, Stack, Typography } from "@mui/material";
import dayjs from "dayjs";
import Card from "../../../components/card/Card";
import CardContent from "../../../components/card/CardContent";
import CardLinkContent from "../../../components/card/CardLinkContent";

const StoreDetails: FunctionComponent<SingleReceiptStore> = ({ tin, name, locationName, address, city }): ReactElement => {
	return (
		<Card>
			<CardLinkContent to={`/stores/companies/${tin}`}>
				<Stack direction="column" gap="1rem">
					<Stack direction="column">
						<Typography variant="h5" fontWeight="bold">{name}</Typography>
						<Typography variant="body1">{locationName}</Typography>
					</Stack>
					<Stack direction="column">
						<Typography variant="h5">{address}</Typography>
						<Typography variant="body1">{city}</Typography>
				</Stack>
				</Stack>
			</CardLinkContent>
		</Card>
	)
}

interface ReceiptMetaDetailsProps {
	totalTaxAmount: SingleReceipt["totalTaxAmount"];
	pfrNumber: SingleReceipt["pfrNumber"];
	counter: SingleReceipt["counter"];
}

const ReceiptMetaDetails: FunctionComponent<ReceiptMetaDetailsProps> = ({ totalTaxAmount, pfrNumber, counter }): ReactElement => {
	return (
		<Card sx={{ mt: "1.5rem"}}>
			<Stack direction="column">
				<CardContent>
					<Stack direction="row" gap="0.5rem" alignItems="center" justifyContent="space-between">
						<Typography variant="h6">Vrednost poreza:</Typography>
						<Typography variant="body1">{totalTaxAmount}</Typography>
					</Stack>
				</CardContent>
				<Divider
					orientation="horizontal"
				/>
				<CardContent>
					<Stack direction="row" gap="0.5rem" alignItems="center" justifyContent="space-between">
						<Typography variant="h6">PFR broj:</Typography>
						<Typography variant="body1">{pfrNumber}</Typography>
					</Stack>
				</CardContent>
				<Divider
					orientation="horizontal"
				/>
				<CardContent>
					<Stack direction="row" gap="0.5rem" alignItems="center" justifyContent="space-between">
						<Typography variant="h6">Brojač:</Typography>
						<Typography variant="body1">{counter}</Typography>
					</Stack>
				</CardContent>
			</Stack>
		</Card>
	)
}

interface ReceiptUserDetailsProps {
	user: SingleReceiptUser;
	createdAt: SingleReceipt["createdAt"];
}

const ReceiptUserDetails: FunctionComponent<ReceiptUserDetailsProps> = ({ user, createdAt }): ReactElement => {
	const { email } = user;
	const formattedCreatedAt = dayjs(createdAt).format("DD.MM.YYYY HH:mm");

	return (
		<Card sx={{ marginTop: "1.5rem" }}>
			<CardContent>
				<Stack direction="row" alignItems="center" justifyContent="space-between">
					<Typography>Dodao:</Typography>
					<Stack direction="column" alignItems="flex-end">
						<Typography variant="body1">{email}</Typography>
						<Typography variant="body1">{formattedCreatedAt}</Typography>
					</Stack>
				</Stack>
			</CardContent>
		</Card>
	)
}

interface ReceiptDetailsProps {
	store: SingleReceiptStore;
	user: SingleReceiptUser;
	totalTaxAmount: SingleReceipt["totalTaxAmount"];
	pfrNumber: SingleReceipt["pfrNumber"];
	counter: SingleReceipt["counter"];
	createdAt: SingleReceipt["createdAt"];
}

const ReceiptDetails: FunctionComponent<ReceiptDetailsProps> = ({
	store,
	user,
	totalTaxAmount,
	pfrNumber,
	counter,
	createdAt
}): ReactElement => {
	return (
		<Box component="section">
			<Typography variant="h4" component="h2" marginY="2rem">Detalji sa računa</Typography>
			<StoreDetails {...store} />
			<ReceiptMetaDetails {...({ totalTaxAmount, pfrNumber, counter})} />
			<ReceiptUserDetails {...({ user, createdAt})} />
		</Box>
	);
}

export default ReceiptDetails
