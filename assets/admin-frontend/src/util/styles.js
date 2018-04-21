import styled, { css } from "styled-components";

export default {
	back_background_color: "#E0F0F0",
	background_color: "#C0C0E0",
	danger_color: "#FF0000"
};

const t_style = css`
	cursor: pointer;
	padding: 5px;
	margin: 5px 0px;
	border: 1px solid black;
`;

const Td = styled.td`
	${t_style}
`;

const Th = styled.th`
	${t_style}
`;

const Table = styled.table`
	border-collapse: collapse;
`;

export {
	Td,
	Th,
	Table
};
