import { Stack, Typography } from "@mui/material";
import { FunctionComponent, ReactElement } from "react";

const StoreListItem: FunctionComponent<StoreListItem> = ({ id, name, total, receiptCount }): ReactElement => {
	return (
		<Stack direction="row" justifyContent="space-between" sx={{ padding: "1rem" }} gap="0.5rem">
			<Stack direction="row" alignItems="center">
				<Typography>{name}</Typography>
			</Stack>
			<Stack direction="column" alignItems="flex-end">
				<Typography>{total}</Typography>
				<Typography variant="caption">{receiptCount} račun</Typography>
			</Stack>
		</Stack>
	)
}

export default StoreListItem;
