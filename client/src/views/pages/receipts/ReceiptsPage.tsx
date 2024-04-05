import { FunctionComponent, ReactElement, useEffect, useMemo, useState } from "react";
import { useQuery } from "react-query";
import { Link } from "react-router-dom";
import { Button, Card, CardContent, Stack, Typography } from "@mui/material";
import ChevronLeftIcon from '@mui/icons-material/ChevronLeft';
import ChevronRightIcon from '@mui/icons-material/ChevronRight';
import { Pie } from 'react-chartjs-2';
import {
	Chart as ChartJS,
	Legend,
	Tooltip,
	ArcElement
} from 'chart.js';
import "dayjs/locale/sr";
import { getReceipts } from "../../../api/receipts";
import PageLayout from "../../layouts/PageLayout/PageLayout";
import { getCategoriesStats } from "../../../api/stats";
import dayjs from "dayjs";

ChartJS.register(
	ArcElement,
	Tooltip,
	Legend
);

dayjs.locale("sr");

interface CurrentMonthRange {
	from: dayjs.Dayjs;
	to: dayjs.Dayjs;
}

const ReceiptsPage: FunctionComponent = (): ReactElement => {
	const [currentMonthRange, setCurrentMonthRange] = useState<CurrentMonthRange>();

	useEffect(() => {
		const firstDateOfMonth = dayjs().startOf("month");
		const lastDateOfMonth = dayjs().endOf("month");
		setCurrentMonthRange({ from: firstDateOfMonth, to: lastDateOfMonth});
	}, []);

	const prevMonth = useMemo(() => {
		if (currentMonthRange) {
			return currentMonthRange.from.subtract(1, "day").format("MMM");
		}
	}, [currentMonthRange]);

	const nextMonth = useMemo(() => {
		if (!currentMonthRange) return;
		const currentDate = dayjs();
		const nextMonthDate = currentMonthRange.to.add(1, "day");

		if (nextMonthDate.isBefore(currentDate)) {
			return nextMonthDate.format("MMM");
		}
	}, [currentMonthRange]);

	const currentMonth = useMemo(() => {
		if (currentMonthRange) {
			return currentMonthRange.from.format("MMMM YYYY").toLowerCase();
		}
	}, [currentMonthRange]);

	const calculatePrevMonth = () => {
		const { from } = currentMonthRange || {};
		if (from) {
			const lastDayOfPreviousMonth = from.subtract(1, "day");
			const firstDayOfPreviousMonth = lastDayOfPreviousMonth?.startOf("month");
			setCurrentMonthRange({ from: firstDayOfPreviousMonth, to: lastDayOfPreviousMonth });
		}
	}

	const calculateNextMonth = () => {
		const { to } = currentMonthRange || {};
		if (!to) return;

		const firstDayOfNextMonth = to.add(1, "day");
		const lastDayOfNextMonth = firstDayOfNextMonth.endOf("month");
		setCurrentMonthRange({ from: firstDayOfNextMonth, to: lastDayOfNextMonth });

	}

	const { data: categoriesStats } = useQuery({
		queryKey: "categories_stats",
		queryFn: () => getCategoriesStats({ fromDate: "2023-01-01", toDate: "2024-04-01"}),
		keepPreviousData: true,
	});

	const { data: receipts } = useQuery({
		queryKey: "receipts",
		queryFn: () => getReceipts({ fromDate: "", toDate: ""}),
		keepPreviousData: true,
	});

	return (
		<PageLayout
			title={currentMonth || ""}
			headerPrefix={prevMonth && <Button onClick={calculatePrevMonth} sx={{ paddingX: 0, textTransform: "none"}}><ChevronLeftIcon />{prevMonth}</Button>}
			headerSuffix={nextMonth && <Button onClick={calculateNextMonth} sx={{ paddingX: 0, textTransform: "none"}}>{nextMonth} <ChevronRightIcon /></Button>}
		>
			{categoriesStats &&
				<Card>
					<CardContent>
						<Stack justifyContent="center" alignItems="center">
							<Pie
								data={{
									labels: categoriesStats.map(stat => stat.category.name),
									datasets: [{
										data: categoriesStats.map(stat => stat.total),
										backgroundColor: categoriesStats.map(stat => stat.category.color),
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
					</CardContent>
				</Card>
			}

			<Stack direction="column" gap="2rem" marginTop="3rem">
				<Typography variant="h4" component="h2" fontWeight="bold">Troškovi</Typography>
				{receipts?.map(receipt => (
					<Stack direction="column" gap="0.5rem" key={receipt.date}>
						<Stack direction="row" justifyContent="space-between" gap="2rem" paddingX="1rem">
							<Typography component="span" color="text.secondary">
								{dayjs(receipt.date).format("DD.MM.YYYY")}
							</Typography>
							<Typography component="span">{receipt.total}</Typography>
						</Stack>

						{receipt.receipts.map(r => (
							<Card component={Link} to="/" sx={{ textDecoration: "none"}} key={r.id}>
								<CardContent component={Stack} direction="column" gap="1rem">
									<Stack direction="row" justifyContent="space-between">
										<Typography component="span" fontWeight="bold">{r.store.name}</Typography>
										<Typography component="span">{r.amount}</Typography>
									</Stack>
									<Stack direction="row" justifyContent="space-between" alignItems="center">
										<Typography component="span" variant="body2">{r.categories.join(", ")}</Typography>
										<Typography component="span">{dayjs(r.date).format("HH:mm")}</Typography>
									</Stack>
								</CardContent>
							</Card>
						))}
					</Stack>
				))}
		</Stack>

		</PageLayout>
	);
}

export default ReceiptsPage;
