import { FunctionComponent, ReactElement } from "react";
import { Box, Divider, Stack, Typography } from "@mui/material";
import { ReceiptCard, ReceiptCardContent } from "./components";
import dayjs from "dayjs";

const StoreDetails: FunctionComponent<SingleReceiptStore> = ({ name, locationName, address, city }): ReactElement => {
	return (
		<ReceiptCard>
			<ReceiptCardContent component={Stack} gap="1rem" direction="column">
				<Stack direction="column">
					<Typography variant="h5" fontWeight="bold">{name}</Typography>
					<Typography variant="body1">{locationName}</Typography>
				</Stack>
				<Stack direction="column">
					<Typography variant="h5">{address}</Typography>
					<Typography variant="body1">{city}</Typography>
				</Stack>
			</ReceiptCardContent>
		</ReceiptCard>
	)
}

interface ReceiptMetaDetailsProps {
	totalTaxAmount: SingleReceipt["totalTaxAmount"];
	pfrNumber: SingleReceipt["pfrNumber"];
	counter: SingleReceipt["counter"];
}

const ReceiptMetaDetails: FunctionComponent<ReceiptMetaDetailsProps> = ({ totalTaxAmount, pfrNumber, counter }): ReactElement => {
	return (
		<ReceiptCard sx={{ marginTop: "1.5rem"}}>
			<ReceiptCardContent direction="column">
				<Stack direction="column">
					<Stack direction="row" gap="0.5rem" alignItems="center" justifyContent="space-between">
						<Typography variant="h6">Vrednost poreza:</Typography>
						<Typography variant="body1">{totalTaxAmount}</Typography>
					</Stack>
					<Divider
						orientation="horizontal"
						sx={{ marginY: "1rem"}}
					/>
					<Stack direction="row" gap="0.5rem" alignItems="center" justifyContent="space-between">
						<Typography variant="h6">PFR broj:</Typography>
						<Typography variant="body1">{pfrNumber}</Typography>
					</Stack>
					<Divider
						orientation="horizontal"
						sx={{ marginY: "1rem"}}
					 />
					<Stack direction="row" gap="0.5rem" alignItems="center" justifyContent="space-between">
						<Typography variant="h6">Brojač:</Typography>
						<Typography variant="body1">{counter}</Typography>
					</Stack>
				</Stack>
			</ReceiptCardContent>
		</ReceiptCard>
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
		<ReceiptCard sx={{ marginTop: "1.5rem" }}>
			<ReceiptCardContent direction="row" alignItems="center" justifyContent="space-between">
				<Typography>Dodao:</Typography>
				<Stack direction="column" alignItems="flex-end">
					<Typography variant="body1">{email}</Typography>
					<Typography variant="body1">{formattedCreatedAt}</Typography>
				</Stack>
			</ReceiptCardContent>
		</ReceiptCard>
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
