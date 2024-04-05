import { useEffect, useState, useMemo } from "react";
import dayjs from "dayjs";

export interface CurrentMonthRange {
	from: dayjs.Dayjs;
	to: dayjs.Dayjs;
}

export interface MonthRange {
	from: string;
	to: string;
}

export interface UseMonthRange {
	monthRange?: MonthRange;
	currentMonthRange?: CurrentMonthRange;
	prevMonth?: string;
	nextMonth?: string;
	currentMonth?: string;
	calculatePrevMonth: () => void;
	calculateNextMonth: () => void;
}

export const useMonthRange = (): UseMonthRange => {
	const [currentMonthRange, setCurrentMonthRange] = useState<CurrentMonthRange>();

	// Initialize `currentMonthRange`
	useEffect(() => {
		const firstDateOfMonth = dayjs().startOf("month");
		const lastDateOfMonth = dayjs().endOf("month");
		setCurrentMonthRange({ from: firstDateOfMonth, to: lastDateOfMonth});
	}, []);

	// Previous month string
	const prevMonth = useMemo(() => {
		if (currentMonthRange) {
			return currentMonthRange.from.subtract(1, "day").format("MMM");
		}
	}, [currentMonthRange]);

	// Next month string
	const nextMonth = useMemo(() => {
		if (!currentMonthRange) return;

		const currentDate = dayjs();
		const nextMonthDate = currentMonthRange.to.add(1, "day");

		if (nextMonthDate.isBefore(currentDate)) {
			return nextMonthDate.format("MMM");
		}
	}, [currentMonthRange]);

	// Current month string
	const currentMonth = useMemo(() => {
		if (currentMonthRange) {
			return currentMonthRange.from.format("MMMM YYYY").toLowerCase();
		}
	}, [currentMonthRange]);

	// Transform CurrentMonthRage that holds dayjs instances to date strings
	const monthRange: MonthRange | undefined = useMemo(() => {
		if (!currentMonthRange) return;

		return {
			from: currentMonthRange.from.format("YYYY-MM-DD"),
			to: currentMonthRange.to.format("YYYY-MM-DD")
		}
	}, [currentMonthRange]);

	// Calculate previous month date range
	const calculatePrevMonth = () => {
		const { from } = currentMonthRange || {};
		if (!from) return;

		const lastDayOfPreviousMonth = from.subtract(1, "day");
		const firstDayOfPreviousMonth = lastDayOfPreviousMonth?.startOf("month");
		setCurrentMonthRange({ from: firstDayOfPreviousMonth, to: lastDayOfPreviousMonth });
	}

	// Calculate next month date range
	const calculateNextMonth = () => {
		const { to } = currentMonthRange || {};
		if (!to) return;

		const firstDayOfNextMonth = to.add(1, "day");
		const lastDayOfNextMonth = firstDayOfNextMonth.endOf("month");
		setCurrentMonthRange({ from: firstDayOfNextMonth, to: lastDayOfNextMonth });
	}

	return {
		monthRange,
		currentMonthRange,
		currentMonth,
		prevMonth,
		nextMonth,
		calculatePrevMonth,
		calculateNextMonth
	}
}
