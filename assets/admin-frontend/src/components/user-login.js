import React, { Component } from "react";
import styled from "styled-components";
import axios from "axios";

import colorShade from "../util/color-shade.js";
import styles from "../util/styles.js";

export default class Comp extends Component {
	constructor(props) {
		super(props);
		this.handleSubmit = this.handleSubmit.bind(this);
		this.state = {
			login_error: ""
		};
	}

	handleSubmit(event) {
		event.preventDefault();
		event.stopPropagation();
		const email = event.target.elements["email"].value;
		const password = event.target.elements["password"].value;
		axios.post("/api/v1/user/login", {
			"email": email,
			"password": password
		})
			.then((response) => {
				this.props.checkUserFun();
			}, (err) => {
				this.setState({
					login_error: "Email/password combination does not exist"
				});
			});
	}

	render() {
		return (
			<Wrapper>
				<LoginWrapper>
					<form onSubmit={this.handleSubmit}>
						<FormGrid>
							<InputLabels>Email:</InputLabels>
							<div><input type="text" name="email"/></div>
							<InputLabels>Password:</InputLabels>
							<div><input type="password" name="password" /></div>
							<div><Button type="submit" /></div>
							<div></div>
						</FormGrid>
						<LoginError>{this.state.login_error}</LoginError>
					</form>
				</LoginWrapper>
			</Wrapper>
		);
	}
}

const Wrapper = styled.div`
	background-color: ${styles.back_background_color};
	width: 100%;
	text-align: center;
	min-height: 100vh;
`;

const LoginWrapper = styled.div`
	background-color: ${styles.background_color};
	border: solid ${colorShade(styles.background_color, -0.5)} 1px;
	margin-top: 100px;
	display: inline-block;
	min-width: 200px;
	color: white;
	padding: 10px;
`;

const FormGrid = styled.div`
	display: grid;
	grid-template-columns: 1fr 2fr;
	grid-gap: 10px;
`;

const InputLabels = styled.div`
	text-align: right;
`;

const LoginError = styled.div`
	color: ${styles.danger_color};
`;

const Button = styled.input`
	border: none;
	color: white;
	text-align: center;
	background-color: ${colorShade(styles.background_color, 0.1)};
	cursor: pointer;
`;
