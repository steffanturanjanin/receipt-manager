import { FunctionComponent } from "react";
import { Card, CardContent, Stack, Typography } from "@mui/material";
import { Pie } from 'react-chartjs-2';
import {
	Chart as ChartJS,
	Legend,
	Tooltip,
	ArcElement
} from 'chart.js';

ChartJS.register(
	ArcElement,
	Tooltip,
	Legend
);

interface CategoryStatsProps {
	categoryStats?: CategoryStats
}

const CategoryStats: FunctionComponent<CategoryStatsProps> = ({ categoryStats }) => {
	const { total, categories } = categoryStats || {};

	return (
		<Card>
			<CardContent component={Stack} direction="column" gap="1rem">
				<Stack direction="column">
					<Typography variant="body1" component="p">Potrošeno ovog meseca:</Typography>
					<Typography variant="h4" component="p">{total || "0.00"}</Typography>
				</Stack>
				{categories && categories.length > 0 &&
					<Stack justifyContent="center" alignItems="center" sx={{ maxHeight: "350px"}}>
						<Pie
							data={{
								labels: categories.map(statistic => statistic.category.name),
								datasets: [{
									data: categories.map(statistic => statistic.total),
									backgroundColor: categories.map(statistic => statistic.category.color),
								}]
							}}
							options={{
								plugins: {
									legend: {
										display: true,
										position: "bottom",
										align: "center",
										labels: {
											boxWidth: 20,
											boxHeight: 20,
										}
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
