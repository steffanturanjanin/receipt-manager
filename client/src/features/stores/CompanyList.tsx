import { Fragment, FunctionComponent, ReactElement } from "react";
import { Link, LinkProps } from "react-router-dom";
import { Divider, Stack, StackProps, Typography, styled } from "@mui/material";
import CardItem from "./CardItem";

const CompanyListItemLink = styled(Stack)<StackProps & LinkProps>(({ theme }) => ({
	display: "flex",
	flexDirection: "row",
	justifyContent: "space-between",
	padding: "1rem",
	gap: "0.5rem",
	color: "inherit",
	textDecoration: "none",
	"&:hover": {
		backgroundColor: theme.palette.grey[50],
	}
}));

interface CompanyListItemProps {
	company: CompanyListItem;
	divider?: boolean;
}

const CompanyListItem: FunctionComponent<CompanyListItemProps> = ({ company, divider }): ReactElement => {
	const { tin, name, total, receiptCount } = company;

	return (
		<Fragment>
			<CompanyListItemLink component={Link} to={`/stores/companies/${tin}`}>
				<Stack direction="row" alignItems="center">
					<Typography>{name}</Typography>
				</Stack>
				<Stack direction="column" alignItems="flex-end">
					<Typography>{total}</Typography>
					<Typography variant="caption">{receiptCount} račun</Typography>
				</Stack>
			</CompanyListItemLink>
			{divider && <Divider />}
		</Fragment>
	)
}

interface CompanyListProps {
	companies: CompanyListItem[],
}

const CompanyList: FunctionComponent<CompanyListProps> = ({ companies }): ReactElement => {
	return (
		<CardItem>
			{companies.map((company, index) => (
				<CompanyListItem
					key={index}
					company={company}
					divider={index !== companies.length - 1}
				/>
			))}
		</CardItem>
	)
}

export default CompanyList;
