import React, { Component } from "react";
import styled from "styled-components";
import Concepts from "./concepts";

export default class Comp extends Component {
	constructor(props) {
		super(props);
	}

	render() {
		return (
			<Wrapper>
				<Concepts />
			</Wrapper>
		);
	}
}

const Wrapper = styled.div`
	padding: 10px;
`;
