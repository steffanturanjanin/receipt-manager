import { FunctionComponent, ReactElement } from "react";
import { useQuery } from "react-query";
import { Button, ButtonProps, Stack, Typography, styled } from "@mui/material";
import ChevronLeftIcon from '@mui/icons-material/ChevronLeft';
import ChevronRightIcon from '@mui/icons-material/ChevronRight';
import { getReceiptsAggregatedByDate } from "../../../api/receipts";
import { getCategoriesStats } from "../../../api/stats";
import PageLayout from "../../layouts/PageLayout/PageLayout";
import CategoryStats from "../../../features/categories/CategoryStats";
import ReceiptList from "../../../features/receipts/ReceiptList";
import { useMonthRange } from "../../../shared/hooks.ts/useMonthRange";

const EMPTY_STATE_STRING = "Nema unetih troškova za ovaj mesec.";

const MonthButton = styled(Button)<ButtonProps>({
	paddingX: 0,
	textTransform: "none",
})

const ReceiptsPage: FunctionComponent = (): ReactElement => {
	const {
		monthRange,
		prevMonth,
		nextMonth,
		currentMonth,
		calculatePrevMonth,
		calculateNextMonth
	} = useMonthRange();

	const { data: categoriesStats } = useQuery({
		queryKey: ["categories_stats", monthRange?.from, monthRange?.to],
		queryFn: () => getCategoriesStats({ fromDate: monthRange!.from, toDate: monthRange!.to }),
		keepPreviousData: true,
		enabled: !!monthRange,
	});

	const { data: receipts } = useQuery({
		queryKey: ["receipts", monthRange?.from, monthRange?.to],
		queryFn: () => getReceiptsAggregatedByDate({ fromDate: monthRange!.from, toDate: monthRange!.to }),
		keepPreviousData: true,
		enabled: !!monthRange,
	});

	return (
		<PageLayout
			title={currentMonth || ""}
			headerPrefix={prevMonth &&
				<MonthButton onClick={calculatePrevMonth}><ChevronLeftIcon />
					{prevMonth}
				</MonthButton>
			}
			headerSuffix={nextMonth &&
				<MonthButton onClick={calculateNextMonth}>
					{nextMonth} <ChevronRightIcon />
				</MonthButton>
			}
		>
			<CategoryStats
				categoryStats={categoriesStats}
			/>
			<Stack direction="column" gap="2rem" marginTop="3rem">
				<Typography variant="h4" component="h2" fontWeight="bold">Troškovi</Typography>
				{(!receipts|| !receipts.length) ?
					<Typography component="p" variant="body1">{EMPTY_STATE_STRING}</Typography> :
					<ReceiptList receiptsAggregatedByDate={receipts}/>
				}
			</Stack>
		</PageLayout>
	);
}

export default ReceiptsPage;
