import { FunctionComponent, ReactElement } from "react";
import { useQuery } from "react-query";
import { useParams } from "react-router-dom";
import { Divider, Stack, Typography } from "@mui/material";
import { getCompany } from "../../../api/stores";
import PageLayout from "../../layouts/PageLayout";
import CardItem from "../../../features/stores/CardItem";
import CompanyLocations from "../../../features/stores/CompanyLocations";
import CompanyExpenses from "../../../features/stores/CompanyExpenses";

const CompanyPage: FunctionComponent = (): ReactElement => {
	const { tin } = useParams();

	const { data: company } = useQuery({
		queryKey: ["company", tin],
		queryFn: () => getCompany(tin!),
		enabled: !!tin,
	});

	return (
		<PageLayout title={company?.name || ""}>
			<Stack direction="column" gap="2rem">

				<CardItem>
					<Stack direction="row" justifyContent="space-between" padding="1rem">
						<Typography>Ime:</Typography>
						<Typography fontWeight="bold">{company?.name}</Typography>
					</Stack>
					<Divider />
					<Stack direction="row" justifyContent="space-between" padding="1rem">
						<Typography>PIB:</Typography>
						<Typography fontWeight="bold">{company?.tin}</Typography>
					</Stack>
				</CardItem>

				{company?.locations &&
					<CompanyLocations companyLocations={company.locations} />
				}

				{company?.expenses &&
					<CompanyExpenses companyExpenses={company.expenses} />
				}

			</Stack>
		</PageLayout>
	)
}

export default CompanyPage;
