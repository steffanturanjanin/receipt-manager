import { FunctionComponent, ReactElement } from "react";
import { Box, BoxProps, styled } from "@mui/material";

const ColorCircle = styled(Box)<BoxProps>({
	width: "0.6em",
	height: "0.6em",
	borderRadius: "100%",
});

interface CategoryCircleProps {
	color: string;
}

const CategoryCircle: FunctionComponent<CategoryCircleProps> = ({ color }): ReactElement => {
	return (
		<ColorCircle component="span" sx={{ backgroundColor: color }}/>
	)
}

export default CategoryCircle
