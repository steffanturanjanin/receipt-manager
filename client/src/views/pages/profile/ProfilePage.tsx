import { FunctionComponent, ReactElement } from "react";
import { useNavigate } from "react-router-dom";
import { useMutation, useQuery } from "react-query";
import dayjs from "dayjs";
import { alpha, styled } from "@mui/system";
import { Divider, Paper, Stack, StackProps, Typography } from "@mui/material";
import { getProfile, logout } from "../../../api/auth";
import PageLayout from "../../layouts/PageLayout";
import LoadingButton, { LoadingButtonProps } from "../../../components/LoadingButton";
import { removeAuth } from "../../../util/auth";

const ProfileItem = styled(Stack)<StackProps>({
	padding: "1rem",
	justifyContent: "space-between",
	flexDirection: "row",
});

const ProfileLogoutItem = styled(ProfileItem)<StackProps>({
	justifyContent: "center",
});

const LogoutButton = styled(LoadingButton)<LoadingButtonProps>(({ theme }) => ({
	color: theme.palette.error.main,
	"&:hover": {
		backgroundColor: alpha(theme.palette.error.light, 0.1),
	}
}));

const ProfilePage: FunctionComponent = (): ReactElement => {
	const navigate = useNavigate();

	const { data: profile } = useQuery({
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

	const { firstName, lastName, email, registeredAt, receiptsCount } = profile || {};
	const formattedRegisteredAt = registeredAt ? dayjs(registeredAt).format("DD.MM.YYYY.") : "";

	return (
		<PageLayout title="Profil">
			<Typography variant="h5" component="h2" mb="1rem">
				{`${firstName} ${lastName}`}
			</Typography>
			<Stack direction="column" gap="2rem">
				<Paper>
					<ProfileItem>
						<Typography>E-mail:</Typography>
						<Typography>{email}</Typography>
					</ProfileItem>
					<Divider />
					<ProfileItem>
						<Typography>Datum pridruživanja:</Typography>
						<Typography>{formattedRegisteredAt}</Typography>
					</ProfileItem>
					<Divider />
					<ProfileItem>
						<Typography>Skeniranih računa:</Typography>
						<Typography>{receiptsCount}</Typography>
					</ProfileItem>
				</Paper>
				<Paper>
					<ProfileLogoutItem component="form" onSubmit={() => onSubmit()}>
						<LogoutButton type="submit" variant="text">Odjavi se</LogoutButton>
					</ProfileLogoutItem>
				</Paper>
			</Stack>
		</PageLayout>
	)
}

export default ProfilePage;
