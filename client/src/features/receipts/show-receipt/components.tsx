import { Card, CardProps, CardContentProps, Stack, StackProps, styled } from "@mui/material";

export const ReceiptCard = styled(Card)<CardProps>({
	borderRadius: "0.5rem",
});

export const ReceiptCardContent = styled(Stack)<CardContentProps>({
	padding: "1rem",
	"&:last-child": {
		paddingBottom: "1rem"
	},
});

export const ReceiptContainer = styled(Stack)<StackProps>({
	gap: "2rem",
});
