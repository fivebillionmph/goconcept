import React, { Component } from "react";
import styled, { css } from "styled-components";
import axios from "axios";

import AddConcept from "./add-concept";
import SingleConcept from "./single-concept";
import { Td, Th, Table } from "../util/styles";

export default class Comp extends Component {
	constructor(props) {
		super(props);
		this.getConcepts = this.getConcepts.bind(this);
		this.changeMode = this.changeMode.bind(this);
		this.changeToConcept = this.changeToConcept.bind(this);
		this.refreshConcepts = this.refreshConcepts.bind(this);
		this.unsetConcept = this.unsetConcept.bind(this);
		this.getComponentTypes = this.getComponentTypes.bind(this);
		this.getRelationshipTypes = this.getRelationshipTypes.bind(this);

		this.state = {
			concept_types: [],
			relationship_types: [],
			concepts: [],
			selected_concept: null,
			concept_type_filter: null,
			page: 1,
			mode: "view"
		};

		this.per_page = 20;
	}

	componentDidMount() {
		this.refreshConcepts();
		this.getComponentTypes()
			.then(() => {
				this.getRelationshipTypes();
			});
	}

	getComponentTypes() {
		let promise = new Promise((resolve, reject) => {
			axios.get("/api/v1/ca/concept/types")
				.then((response) => {
					let state = Object.assign({}, this.state, {
						concept_types: response.data
					});
					this.setState(state, () => {
						resolve();
					});
				}, () => {
					reject();
				});
		});
		return promise;
	}

	getRelationshipTypes() {
		let promise = new Promise((resolve, reject) => {
			axios.get("/api/v1/ca/concept/relationships")
				.then((response) => {
					let state = Object.assign({}, this.state, {
						relationship_types: response.data
					});
					this.setState(state, () => {
						resolve();
					});
				}, () => {
					reject();
				});
		});
		return promise;
	}

	refreshConcepts() {
		this.getConcepts(this.state.concept_type_filter, this.state.page);
	}

	getConcepts(type, page) {
		let count = this.per_page;
		let offset = (page - 1) * count;

		let url = "/api/v1/ca/concept/data";
		if(type) {
			url += "/" + type;
		}
		url += "?count=" + count + "&offset=" + offset;

		axios.get(url)
			.then((response) => {
				let state = Object.assign({}, this.state, {
					concepts: response.data,
					concept_type_filter: type,
					page: page
				});
				this.setState(state);
			});
	}

	changeToConcept(concept) {
		const state = Object.assign({}, this.state, {
			selected_concept: concept,
			mode: "single-concept"
		});
		this.setState(state);
	}

	changeMode(mode) {
		let state = Object.assign({}, this.state, {
			mode: mode
		});
		this.setState(state, () => {
			if(mode == "view") {
				this.refreshConcepts();
			}
		});
	}

	unsetConcept() {
		const state = Object.assign({}, this.state, {
			selected_concept: null
		});
		this.setState(state, () => {
			this.changeMode("view");
		});
	}

	render() {
		return (
			<Wrapper>
				<PanelRow>
					<PanelButton onClick={() => this.changeMode("view")}>
						View Concepts
					</PanelButton>
					<PanelButton onClick={() => this.changeMode("add-concept")}>
						Add concept
					</PanelButton>
				</PanelRow>
				<BodyRow>
					{this.state.mode == "view" &&
						<div>
							<h2>Concept list</h2>
							<Table>
								<thead>
									<tr>
										<th>Type</th>
										<th>Name</th>
									</tr>
								</thead>
								<tbody>
									{this.state.concepts && this.state.concepts.map((concept, ix) => {
										return (
											<tr key={ix} onClick={() => this.changeToConcept(concept)}>
												<Td>{ concept.type }</Td><Td>{ concept.name }</Td>
											</tr>
										);
									})}
								</tbody>
							</Table>
						</div>
					}
					{this.state.mode == "single-concept" &&
						<SingleConcept concept={this.state.selected_concept} types={this.state.concept_types} relationship_types={this.state.relationship_types} unsetConceptFun={this.unsetConcept} />
					}
					{this.state.mode == "add-concept" &&
						<AddConcept refreshFun={this.refreshConcepts} newConceptFun={this.changeToConcept} concept_types={this.state.concept_types} />
					}
				</BodyRow>
			</Wrapper>
		);
	}
}

const Wrapper = styled.div`
	display: grid;
	grid-template-columns: 1fr 5fr;
	grid-gap: 10px;
`;

const PanelRow = styled.div`
	margin: 10px;
`;

const BodyRow = styled.div`
`;

const PanelButton = styled.div`
	cursor: pointer;
	padding: 5px;
	border: 1px black solid;
	margin: 5px;
`;
