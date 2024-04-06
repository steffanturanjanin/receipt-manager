import { FunctionComponent, ReactElement } from "react";
import { Stack, Typography } from "@mui/material";
import { ReceiptCard, ReceiptCardContent } from "./components";
import dayjs from "dayjs";

interface ReceiptPaymentOverviewProps {
	totalPurchaseAmount: string;
	date: string;
	storeName: string;
}

const ReceiptPaymentOverview: FunctionComponent<ReceiptPaymentOverviewProps> = ({
	totalPurchaseAmount,
	date,
	storeName
}): ReactElement => {
	const formattedDate = dayjs(date).format("DD.MM.YYYY HH:mm:ss");

	return (
		<ReceiptCard>
			<ReceiptCardContent direction="column" gap="1rem">
				<Stack direction="column">
					<Typography variant="h4">{totalPurchaseAmount || "0.00" } RSD</Typography>
					<Typography variant="body1">{formattedDate}</Typography>
				</Stack>
				<Typography variant="h5">{storeName}</Typography>
			</ReceiptCardContent>
		</ReceiptCard>
	)
}

export default ReceiptPaymentOverview;
