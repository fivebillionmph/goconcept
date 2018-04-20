import React, { Component } from "react";
import styled from "styled-components";
import axios from "axios";

export default class Comp extends Component {
	constructor(props) {
		super(props);
		this.submit = this.submit.bind(this);
	}

	submit(event) {
		event.preventDefault();
		event.stopPropagation();
		const type = event.target.elements["type"].value;
		const name = event.target.elements["name"].value;

		const data = {
			type: type,
			name: name
		};
		axios.post("/api/v1/ca/concept/add", data)
			.then((response) => {
				const new_concept = response.data;
				this.props.refreshFun();
				this.props.newConceptFun(new_concept);
			});

	}

	render() {
		return (
			<Wrapper>
				<form onSubmit={this.submit}>
					<select name="type">
						{
							this.props.concept_types.map((ct, ix) => {
								return (
									<option key={ix}>
										{ct.type}
									</option>
								);
							})
						}
					</select>
					<input type="text" name="name" />
					<input type="submit" value="Submit" />
				</form>
			</Wrapper>
		);
	}
}

const Wrapper = styled.div`
	
`;
