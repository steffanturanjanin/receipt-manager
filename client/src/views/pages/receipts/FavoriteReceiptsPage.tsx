import { FunctionComponent, ReactElement, useState } from "react";
import { useQuery, useQueryClient } from "react-query";
import { getFavoriteReceipts } from "../../../api/receipts";
import PageLayout from "../../layouts/PageLayout";
import RemoveFavorite from "../../../features/receipts/RemoveFavorite";
import FavoriteReceiptsList from "../../../features/receipts/FavoriteReceiptsList";

interface RemoveFavoriteDialog {
	open: boolean;
	receiptId?: number;
}

const FavoriteReceiptsPage: FunctionComponent = (): ReactElement => {
	const [removeFavoriteDialog, setRemoveFavoriteDialog] = useState<RemoveFavoriteDialog>({
		open: false
	});

	const { isLoading: isFavoriteReceiptsLoading, data: favoriteReceipts } = useQuery({
		queryKey: ["favorite_receipts"],
		queryFn: () => getFavoriteReceipts(),
		keepPreviousData: true,
	});

	const queryClient = useQueryClient();

	const refetchFavorites = () => {
		queryClient.invalidateQueries(["favorite_receipts"]);
	}

	return (
		<PageLayout
			title="Omiljeni računi"
			showBackdrop={isFavoriteReceiptsLoading}
		>
			<FavoriteReceiptsList
				receipts={favoriteReceipts || []}
				onItemClicked={(receiptId) => setRemoveFavoriteDialog({ open: true, receiptId })}
			/>
			{removeFavoriteDialog.receiptId &&
				<RemoveFavorite
					receiptId={removeFavoriteDialog.receiptId}
					open={removeFavoriteDialog.open}
					onClose={() => setRemoveFavoriteDialog({ ...removeFavoriteDialog, open: false })}
					onRemovedFromFavorites={refetchFavorites}
				/>
			}
		</PageLayout>
	)
}

export default FavoriteReceiptsPage;
