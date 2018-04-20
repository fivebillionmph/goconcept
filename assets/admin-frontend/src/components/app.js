import React, { Component } from "react";
import styled from "styled-components";
import axios from "axios";

import UserLogin from "./user-login";
import AdminPanel from "./admin-panel";

export default class Comp extends Component {
	constructor(props) {
		super(props);
		this.checkUser = this.checkUser.bind(this);

		this.state = {
			user: null
		};

		this.checkUser()
	}

	checkUser() {
		axios.get("/api/v1/user/info")
			.then((response) => {
				this.setState({
					user: response.data
				});
			}, (err) => {
				this.setState({
					user: null
				});
			});
	}

	render() {
		return (
			<div>
				{this.state.user === null &&
					<UserLogin checkUserFun={this.checkUser} />
				}
				{this.state.user !== null &&
					<AdminPanel />
				}
			</div>
		);
	}
}
