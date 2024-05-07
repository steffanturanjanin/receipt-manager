import { FunctionComponent, ReactElement, useMemo } from "react";
import { useParams } from "react-router-dom";
import { useQuery } from "react-query";
import { Stack, Typography } from "@mui/material";
import { getCategoryStats } from "../../../api/categories";
import PageLayout from "../../layouts/PageLayout";
import Card from "../../../components/card/Card";
import CardContent from "../../../components/card/CardContent";
import BackButton from "../../../components/BackButton";
import MostPopularReceiptItems from "../../../features/categories/MostPopularReceiptItems";
import MostPopularStores from "../../../features/categories/MostPopularStores";

const CategoryStatsPage: FunctionComponent = (): ReactElement => {
	const { id } = useParams();

	const categoryId = useMemo(() => {
		if (id === undefined) return undefined;
		return parseInt(id)
	}, [id]);

	const { isLoading, data: categoryStats } = useQuery({
		queryKey: ["category", categoryId],
		queryFn: () => getCategoryStats(categoryId!),
		enabled: !!categoryId,
	})

	return (
		<PageLayout
			title={categoryStats?.name || ""}
			showBackdrop={isLoading}
			headerPrefix={<BackButton />}
		>
			<Stack direction="column" gap="2rem">
				<Card>
					<CardContent>
						<Stack gap="0.5rem">
							<Typography>Potrošeno u proteklih 12 meseci:</Typography>
							<Typography variant="h5">{categoryStats?.total}</Typography>
						</Stack>
					</CardContent>
				</Card>

				<MostPopularReceiptItems
					receiptItems={categoryStats?.mostPopularReceiptItems || []}
				/>

				<MostPopularStores
					stores={categoryStats?.mostPopularStores || []}
				/>
			</Stack>
		</PageLayout>
	)
}

export default CategoryStatsPage;
