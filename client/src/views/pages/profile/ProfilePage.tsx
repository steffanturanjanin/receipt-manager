import { FunctionComponent, FormEvent, ReactElement, useMemo } from "react";
import { useNavigate } from "react-router-dom";
import { useMutation, useQuery } from "react-query";
import dayjs from "dayjs";
import {
	Divider,
	Stack,
	StackProps,
	Typography,
	alpha,
	styled
} from "@mui/material";
import LogoutIcon from '@mui/icons-material/Logout';
import { getProfile, logout } from "../../../api/auth";
import PageLayout from "../../layouts/PageLayout";
import LoadingButton, { LoadingButtonProps } from "../../../components/LoadingButton";
import { removeAuth } from "../../../util/auth";
import Card from "../../../components/card/Card";
import CardContent from "../../../components/card/CardContent";
import CardLinkContent from "../../../components/card/CardLinkContent";
import StarFullIcon from "@mui/icons-material/Star";

const ProfileItem = styled(Stack)<StackProps>({
	justifyContent: "space-between",
	flexDirection: "row",
});

const LogoutButton = styled(LoadingButton)<LoadingButtonProps>(({ theme }) => ({
	color: theme.palette.error.main,
	width: "100%",
	gap: "0.5rem",
	"&:hover": {
		backgroundColor: alpha(theme.palette.error.light, 0.1),
	}
}));

const ProfilePage: FunctionComponent = (): ReactElement => {
	const navigate = useNavigate();

	const { data: profile, isFetching: isProfileFetching } = useQuery({
		queryKey: ["profile"],
		queryFn: () => getProfile(),
	});

	const { mutate: onLogout } = useMutation({
		mutationFn: () => logout(),
		onSuccess: () => {
			removeAuth();
			navigate("/auth/login", { replace: true });
		}
	});

	const onSubmit = (event: FormEvent<HTMLFormElement>) => {
		event.preventDefault();
		onLogout();
	}

	const {
		firstName,
		lastName,
		email,
		registeredAt,
		receiptCount
	} = profile || {};

	const formattedRegisteredAt = useMemo(
		() => registeredAt ? dayjs(registeredAt).format("DD.MM.YYYY.") : "",
		[registeredAt]
	);
	const username = useMemo(
		() => (firstName && lastName) ? `${firstName} ${lastName}` : "",
		[firstName, lastName]
	);


	return (
		<PageLayout
			title="Profil"
			showBackdrop={isProfileFetching}
		>
			<Typography variant="h5" component="h2" mb="1rem">{username}</Typography>
			<Stack direction="column" gap="2rem">
				<Card>
					<CardContent>
						<ProfileItem>
							<Typography>E-mail:</Typography>
							<Typography>{email}</Typography>
						</ProfileItem>
					</CardContent>
					<Divider />
					<CardContent>
						<ProfileItem>
							<Typography>Datum pridruživanja:</Typography>
							<Typography>{formattedRegisteredAt}</Typography>
						</ProfileItem>
					</CardContent>
					<Divider />
					<CardContent>
						<ProfileItem>
							<Typography>Skeniranih računa:</Typography>
							<Typography>{receiptCount}</Typography>
						</ProfileItem>
					</CardContent>
				</Card>

				<Card>
					<CardLinkContent to="/receipts/favorites">
						<Stack component={Typography} direction="row" alignItems="center" justifyContent="center" gap="0.25rem">
							Omiljeni računi <StarFullIcon color="warning" />
						</Stack>
					</CardLinkContent>
				</Card>

				<Card>
					<ProfileItem>
						<Stack component="form" onSubmit={onSubmit} width="100%">
							<LogoutButton type="submit" variant="text">
								Odjavi se <LogoutIcon />
							</LogoutButton>
						</Stack>
					</ProfileItem>
				</Card>

			</Stack>
		</PageLayout>
	)
}

export default ProfilePage;
