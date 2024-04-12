import { FunctionComponent, ReactElement } from "react";
import { Stack } from "@mui/material";
import PageLayout from "../../layouts/PageLayout";
import ExpensesChart from "../../../features/stats/ExpensesChart";
import StatItem, { StatItemContent } from "../../../features/stats/StatItem";
import ExpensesByCategory from "../../../features/stats/ExpensesByCategory";
import ExpensesByStore from "../../../features/stats/ExpensesByStore";

const StatsPage: FunctionComponent = (): ReactElement => {
	return (
		<PageLayout title="Statistika">
			<Stack direction="column" gap="2rem">
				<StatItem title="Potrošnja u proteklih 12 meseci">
					<StatItemContent>
						<ExpensesChart />
					</StatItemContent>
				</StatItem>

				<StatItem title="Potrošnja po kategoriji">
					<ExpensesByCategory />
				</StatItem>

				<StatItem title="Top 10 prodavnica">
					<ExpensesByStore />
				</StatItem>
			</Stack>
		</PageLayout>
	)
}

export default StatsPage;
