import { FunctionComponent, ReactElement } from "react";
import { Stack, Typography } from "@mui/material";
import dayjs from "dayjs";
import Card from "../../../components/card/Card";
import CardContent from "../../../components/card/CardContent";

interface ReceiptPaymentOverviewProps {
	totalPurchaseAmount: string;
	date?: string;
	storeName: string;
}

const ReceiptPaymentOverview: FunctionComponent<ReceiptPaymentOverviewProps> = ({
	totalPurchaseAmount,
	date,
	storeName
}): ReactElement => {
	const formattedDate = date ? dayjs(date).format("DD.MM.YYYY HH:mm:ss") : null;

	return (
		<Card>
			<CardContent>
				<Stack direction="column" gap="1rem">
					<Stack direction="column">
						<Typography variant="h4">{totalPurchaseAmount || "0.00" } RSD</Typography>
						<Typography variant="body1" color="grey.600">{formattedDate}</Typography>
					</Stack>
					<Typography variant="h6">{storeName}</Typography>
				</Stack>
			</CardContent>
		</Card>
	)
}

export default ReceiptPaymentOverview;
