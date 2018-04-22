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
		this.getConceptsCount = this.getConceptsCount.bind(this);
		this.refreshConcepts = this.refreshConcepts.bind(this);
		this.changeMode = this.changeMode.bind(this);
		this.changeToConcept = this.changeToConcept.bind(this);
		this.unsetConcept = this.unsetConcept.bind(this);
		this.getComponentTypes = this.getComponentTypes.bind(this);
		this.getRelationshipTypes = this.getRelationshipTypes.bind(this);
		this.changePage = this.changePage.bind(this)
		this.searchFilterOnChange = this.searchFilterOnChange.bind(this);

		this.state = {
			concept_types: [],
			relationship_types: [],
			concepts: [],
			selected_concept: null,
			mode: "view",
			sort: "",
			page: 1,
			total_count: 0,
			search_filter: ""
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
		this.getConcepts();
		this.getConceptsCount();
	}

	getConcepts() {
		let count = this.per_page;
		let offset = (this.state.page - 1) * count;

		let url = "/api/v1/ca/concept/data?q=" + this.state.search_filter;
		url += "&count=" + count + "&offset=" + offset;

		axios.get(url)
			.then((response) => {
				let state = Object.assign({}, this.state, {
					concepts: response.data,
				});
				this.setState(state);
			});
	}

	getConceptsCount() {
		let url = "/api/v1/ca/concept/data-count?q=" + this.state.search_filter;
		axios.get(url)
			.then((response) => {
				const state = Object.assign({}, this.state, {
					total_count: response.data
				});
				this.setState(state);
			});
	}

	changePage(page) {
		const state = Object.assign({}, this.state, {
			page: page
		});
		this.setState(state, () => {
			this.getConcepts();
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

	searchFilterOnChange(event) {
		if(!event.target) return;
		const new_value = event.target.value;
		const state = Object.assign({}, this.state, {
			search_filter: new_value
		});
		this.setState(state, () => {
			this.refreshConcepts();
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
							<input type="text" value={this.state.search_filter} onChange={this.searchFilterOnChange} />
							<Table>
								<thead>
									<tr>
										<Th>Type</Th>
										<Th>Name</Th>
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
							<Pagination count={this.state.total_count} page={this.state.page} per_page={this.per_page} changePageFun={this.changePage} />
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

class Pagination extends Component {
	render() {
		let last_page = Math.floor(this.props.count / this.props.per_page);
		if(this.props.count % this.props.per_page != 0) {
			last_page += 1;
		}
		return (
			<div>
				{this.props.page != 1 &&
					<span><a href="javascript:void(0)" onClick={() => this.props.changePageFun(1)}>1</a>&nbsp;</span>
				}
				{this.props.page - 1 > 1 &&
					<span><a href="javascript:void(0)" onClick={() => this.props.changePageFun(this.props.page - 1)}>{this.props.page - 1}</a>&nbsp;</span>

				}
				<span>({this.props.page})&nbsp;</span>
				{this.props.page + 1 < last_page &&
					<span><a href="javascript:void(0)" onClick={() => this.props.changePageFun(this.props.page + 1)}>{this.props.page + 1}</a>&nbsp;</span>
				}
				{last_page != 1 && last_page != this.props.page &&
					<span><a href="javascript:void(0)" onClick={() => this.props.changePageFun(last_page)}>{last_page}</a>&nbsp;</span>
				}
			</div>
		);
	}
}
