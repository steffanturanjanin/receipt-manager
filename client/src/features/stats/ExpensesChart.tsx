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
import {Stack} from "@mui/material"
import dayjs from "dayjs";
import { useQuery } from "react-query";
import { getExpensesDateBreakdownStats } from "../../api/stats";

ChartJS.register(
  CategoryScale,
  LinearScale,
  BarElement,
  Title,
  Tooltip,
  Legend
);

const ExpensesChart: FunctionComponent = (): ReactElement => {
	const { data: expenses } = useQuery({
		queryKey: ["expenses_yearly_breakdown"],
		queryFn: () => getExpensesDateBreakdownStats(),
	});

	const labels = useMemo(() => expenses?.map(expense => dayjs(expense.date).format("MMM")) || [], [expenses]);
	const data = useMemo(() => expenses?.map(expense => expense.total) || [], [expenses])

	return (
		<Stack>
			<Bar
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
					plugins: {
						legend: {
							display: false,
						}
					}
				}}
			/>
		</Stack>
	)
}

export default ExpensesChart;
