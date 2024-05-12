import { FunctionComponent, useMemo } from "react";
import { Stack, Typography } from "@mui/material";
import { Pie } from 'react-chartjs-2';
import {
	Chart as ChartJS,
	Legend,
	Tooltip,
	ArcElement,
	ChartEvent,
	ActiveElement
} from 'chart.js';
import Card from "../../components/card/Card";
import CardContent from "../../components/card/CardContent";
import { useNavigate } from "react-router-dom";

ChartJS.register(
	ArcElement,
	Tooltip,
	Legend
);

interface CategoryStatsProps {
	categoryStats?: CategoryStats
}

const CategoryStats: FunctionComponent<CategoryStatsProps> = ({ categoryStats }) => {
	const navigate = useNavigate();

	const { total, categories } = categoryStats || {};

	const labels = useMemo(
		() => categories?.map((statistic) => statistic.category.name) || [],
		[categories]
	);

	return (
		<Card>
			<CardContent>
				<Stack direction="column">
					<Typography>Potrošeno ovog meseca:</Typography>
					<Typography variant="h4">{total || "0.00"}</Typography>
				</Stack>
				{categories && categories.length > 0 &&
					<Stack justifyContent="center" alignItems="center" sx={{ maxHeight: "350px" }}>
						<Pie
							data={{
								labels: labels,
								datasets: [{
									data: categories.map((statistic) => statistic.total),
									backgroundColor: categories.map((statistic) => statistic.category.color),
								}]
							}}
							options={{
								plugins: {
									legend: {
										display: true,
										position: "bottom",
										align: "start",
										labels: {
											boxWidth: 12,
											boxHeight: 12,
										}
									}
								},
								onClick: (_: ChartEvent, elements: ActiveElement[]) => {
									if (elements.length) {
										const index = elements[0].index;
										const { category } = categories[index];
										navigate(`/categories/${category.id}`);
									}
								}
							}}
						/>
					</Stack>
			}
			</CardContent>
		</Card>
	)
}

export default CategoryStats;
