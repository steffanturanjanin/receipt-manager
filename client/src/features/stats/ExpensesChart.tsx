import { FunctionComponent, ReactElement, useMemo } from "react";
import { Bar } from "react-chartjs-2";
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  BarElement,
  Title,
  Tooltip,
  Legend,
} from 'chart.js';
import { Stack, StackProps, styled } from "@mui/material"
import dayjs from "dayjs";

ChartJS.register(
  CategoryScale,
  LinearScale,
  BarElement,
  Title,
  Tooltip,
  Legend
);

const ExpensesChartContainer = styled(Stack)<StackProps>(({ theme }) => ({
	width: "100%",
	height: "250px",
	[theme.breakpoints.up("sm")]: {
		height: "320px",
	}
}));

interface ExpensesChartProps {
	expenses: ExpensesDateBreakdown[]
}

const ExpensesChart: FunctionComponent<ExpensesChartProps> = ({ expenses }): ReactElement => {
	const labels = useMemo(
		() => expenses.map(expense => dayjs(expense.date).format("MMM")) || [],
		[expenses]
	);

	const data = useMemo(
		() => expenses.map(expense => expense.total) || [],
		[expenses]
	);

	return (
		<ExpensesChartContainer direction="column" justifyContent="center" alignItems="center">
			<Bar
				style={{ height: "100%" }}
				data={{
					labels: labels,
					datasets: [
						{
							label: "Potrošnja",
							data: data,
							backgroundColor: "#4caf50",
						}
					],
				}}
				options={{
					responsive: true,
					maintainAspectRatio: false,
					plugins: {
						legend: {
							display: false,
						}
					}
				}}
			/>
		</ExpensesChartContainer>
	)
}

export default ExpensesChart;
