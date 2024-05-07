import { Fragment, FunctionComponent, ReactElement, useMemo } from "react";
import { Divider, Stack, Typography } from "@mui/material";
import Card from "../../components/card/Card";
import CardContent from "../../components/card/CardContent";

interface MostPopularReceiptItemsProps {
	receiptItems: MostPopularReceiptItem[];
}

const MostPopularReceiptItems: FunctionComponent<MostPopularReceiptItemsProps> = ({ receiptItems }): ReactElement => {
	const ReceiptItems = (
		<Card>
			{receiptItems.map(({ total, name, receiptCount }, index) => (
				<Fragment key={index}>
					<CardContent>
						<Stack direction="row" alignItems="center" justifyContent="space-between">
							<Stack direction="column">
								<Typography>{name}</Typography>
								<Typography variant="body2" color="grey.600">{receiptCount} troškova</Typography>
							</Stack>
							<Typography>{total}</Typography>
						</Stack>
					</CardContent>
					{index !== receiptItems.length - 1 && <Divider /> }
				</Fragment>
			))}
		</Card>
	);

	const Empty = <Typography>Nema artikala...</Typography>

	const Content = useMemo(
		() => receiptItems.length ? ReceiptItems : Empty,
		[receiptItems]
	);

	return (
		<Stack direction="column" gap="1rem">
			<Typography variant="h6">Najplaćeniji artikli</Typography>
			{Content}
		</Stack>
	)
}

export default MostPopularReceiptItems;
