import { FunctionComponent, ReactElement } from "react";
import { useQuery } from "react-query";
import {  useParams } from "react-router-dom";
import { Divider, Stack, Typography } from "@mui/material";
import { getCompany } from "../../../api/stores";
import PageLayout from "../../layouts/PageLayout";
import Card from "../../../components/card/Card";
import CompanyLocations from "../../../features/stores/CompanyLocations";
import CompanyExpenses from "../../../features/stores/CompanyExpenses";
import BackButton from "../../../components/BackButton";

const CompanyPage: FunctionComponent = (): ReactElement => {
	const { tin } = useParams();

	const { isLoading: isCompanyLoading, data: company } = useQuery({
		queryKey: ["company", tin],
		queryFn: () => getCompany(tin!),
		enabled: !!tin,
	});

	return (
		<PageLayout title={company?.name || ""} showBackdrop={isCompanyLoading} headerPrefix={<BackButton />}>
			{company && (
				<Stack direction="column" gap="2rem">
					<Card>
						<Stack direction="row" justifyContent="space-between" padding="1rem">
							<Typography>Ime:</Typography>
							<Typography fontWeight="bold">{company?.name}</Typography>
						</Stack>
						<Divider />
						<Stack direction="row" justifyContent="space-between" padding="1rem">
							<Typography>PIB:</Typography>
							<Typography fontWeight="bold">{company?.tin}</Typography>
						</Stack>
					</Card>
					<CompanyLocations
						companyLocations={company.locations}
					/>
					<CompanyExpenses
						companyExpenses={company.expenses}
					 />
				</Stack>
			)}
		</PageLayout>
	)
}

export default CompanyPage;
