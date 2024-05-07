import { FunctionComponent, ReactElement } from "react";
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

const ProfileItem = styled(Stack)<StackProps>({
	justifyContent: "space-between",
	flexDirection: "row",
});

const ProfileLogoutItem = styled(ProfileItem)<StackProps>({
	justifyContent: "center",
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

	const { mutate: onSubmit } = useMutation({
		mutationFn: () => logout(),
		onSuccess: () => {
			removeAuth();
			navigate("/auth/login", { replace: true });
		}
	});

	const { firstName, lastName, email, registeredAt, receiptCount } = profile || {};
	const formattedRegisteredAt = registeredAt ? dayjs(registeredAt).format("DD.MM.YYYY.") : "";

	return (
		<PageLayout
			title="Profil"
			showBackdrop={isProfileFetching}
		>
			<Typography variant="h5" component="h2" mb="1rem">
				{`${firstName} ${lastName}`}
			</Typography>
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
					<ProfileLogoutItem component="form" onSubmit={() => onSubmit()}>
						<LogoutButton type="submit" variant="text">
							Odjavi se <LogoutIcon />
						</LogoutButton>
					</ProfileLogoutItem>
				</Card>

			</Stack>
		</PageLayout>
	)
}

export default ProfilePage;
