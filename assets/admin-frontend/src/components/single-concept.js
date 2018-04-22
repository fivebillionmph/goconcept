import React, { Component } from "react";
import styled, { css } from "styled-components";
import axios from "axios";

import { Table, Th, Td } from "../util/styles";

export default class Comp extends Component {
	constructor(props) {
		super(props);
		this.toggleShowNewData = this.toggleShowNewData.bind(this);
		this.toggleShowNewRelationship = this.toggleShowNewRelationship.bind(this);
		this.submitData = this.submitData.bind(this);
		this.submitRelationship = this.submitRelationship.bind(this);
		this.deleteDataConfirm = this.deleteDataConfirm.bind(this);
		this.deleteData = this.deleteData.bind(this);
		this.deleteConceptConfirm = this.deleteConceptConfirm.bind(this);
		this.deleteConcept = this.deleteConcept.bind(this);
		this.deleteRelationshipConfirm = this.deleteRelationshipConfirm.bind(this);
		this.deleteRelationship = this.deleteRelationship.bind(this);
		this.closeModals = this.closeModals.bind(this);
		this.changeNewRelationshipType = this.changeNewRelationshipType.bind(this);
		this.relOtherNameChange = this.relOtherNameChange.bind(this);
		this.relOtherNameOnChange = this.relOtherNameOnChange.bind(this);
		this.relOtherNameAutocomplete = this.relOtherNameAutocomplete.bind(this);
		this.relOtherNameAutocompleteShow = this.relOtherNameAutocompleteShow.bind(this);

		const addable_data = this.getAddableData(this.props.types, this.props.concept);
		const addable_relationships = [];
		for(let i = 0; i < this.props.relationship_types.length; i++) {
			let this_rt = this.props.relationship_types[i];
			let reversible = false;
			let forward = true;
			let other_type = null;
			let display_string = "";
			let other_string = "";
			if(this_rt.type1 != this.props.concept.type && this_rt.type2 != this.props.concept.type) {
				continue;
			}
			if(this_rt.type1 == this_rt.type2) {
				reversible = true;
			}
			if(this_rt.type1 == this.props.concept.type) {
				forward = true;
				other_type = this_rt.type2;
				other_string = this_rt.string2;
				display_string = this_rt.string1;
			} else {
				forward = false;
				other_type = this_rt.type1;
				other_string = this_rt.string1;
				display_string = this_rt.string2;
			}

			addable_relationships.push({
				other_type,
				reversible,
				forward,
				display_string,
				other_string
			});
		}
		this.state = {
			show_new_data: false,
			show_new_relationship: false,
			addable_data: addable_data,
			addable_relationships: addable_relationships,
			concept: this.props.concept,
			selected_addable_relationship: 0,
			rel_other_name: "",
			rel_other_name_autocomplete: false,
			rel_other_concepts: []
		};

		this.refreshConcept();
	}

	getAddableData(types, concept) {
		let data_names = [];
		let type = null;
		for(let i = 0; i < types.length; i++) {
			if(types[i].type == concept.type) {
				type = types[i];
				break;
			}
		}
		if(type == null) {
			return data_names;
		}
		let existing_keys = new Set();
		if(concept.data) {
			for(let i = 0; i < concept.data.length; i++) {
				existing_keys.add(concept.data[i].key)
			}
		}
		for(let i = 0; i < type.concept_data.length; i++) {
			if(type.concept_data[i].single && existing_keys.has(type.concept_data[i].type)) {
				continue;
			}
			data_names.push(type.concept_data[i].type);
		}
		return data_names;
	}

	submitData(event) {
		event.preventDefault();
		event.stopPropagation();
		const key = event.target.elements["key"].value;
		const data = event.target.elements["data"].value;

		const req_data = {
			type: this.state.concept.type,
			name: this.state.concept.name,
			data_key: key,
			data_value: data
		};

		axios.post("/api/v1/ca/concept/add/data", req_data)
			.then((response) => {
				this.refreshConcept();
			});
	}

	submitRelationship(event) {
		event.preventDefault();
		event.stopPropagation();
		const addable_relationship = this.state.addable_relationships[this.state.selected_addable_relationship];
		let type1 = this.state.concept.type;
		let type2 = addable_relationship.other_type;
		let name1 = this.state.concept.name;
		let name2 = event.target.elements["other_name"].value;
		let string1 = addable_relationship.display_string;
		let string2 = addable_relationship.other_string;

		if(!addable_relationship.forward) {
			[type1, type2] = [type2, type1];
			[name1, name2] = [name2, name1];
			[string1, string2] = [string2, string1];
		}

		if(addable_relationship.reversible && event.target.elements["reverse"].checked) {
			[name1, name2] = [name2, name1];
		}

		const req_data = {
			type1,
			type2,
			name1,
			name2,
			string1,
			string2
		};
		axios.post("/api/v1/ca/concept/add/rel", req_data)
			.then(() => {
				this.refreshConcept();
			});
	}

	refreshConcept() {
		axios.get("/api/v1/ca/concept/data/" + this.state.concept.type + "/" + this.state.concept.name)
			.then((response) => {
				const concept = response.data;
				const addable_data = this.getAddableData(this.props.types, concept);
				const state = Object.assign({}, this.state, {
					addable_data: addable_data,
					concept: concept
				});
				this.setState(state);
			});
	}

	toggleShowNewData() {
		const state = Object.assign({}, this.state, {
			show_new_data: !this.state.show_new_data
		});
		this.setState(state);
	}

	toggleShowNewRelationship() {
		const state = Object.assign({}, this.state, {
			show_new_relationship: !this.state.show_new_relationship
		});
		this.setState(state);
	}

	deleteDataConfirm(data) {
		const state = Object.assign({}, this.state, {
			modal_delete_data: data
		});
		this.setState(state);
	}

	deleteData(data) {
		const req_data = {
			type: this.state.concept.type,
			name: this.state.concept.name,
			data_key: data.key,
			data_value: data.value
		};
		axios.post("/api/v1/ca/concept/delete/data", req_data)
			.then((response) => {
				this.closeModals()
					.then(() => {
						this.refreshConcept();
					});
			});
	}

	deleteConceptConfirm() {
		const state = Object.assign({}, this.state, {
			modal_delete_concept: true
		});
		this.setState(state);
	}

	deleteConcept() {
		const req_data = {
			type: this.state.concept.type,
			name: this.state.concept.name
		};

		axios.post("/api/v1/ca/concept/delete", req_data)
			.then(() => {
				this.closeModals();
				this.props.unsetConceptFun();
			});
	}

	deleteRelationshipConfirm(rel) {
		const state = Object.assign({}, this.state, {
			modal_delete_relationship: rel
		});
		this.setState(state);
	}

	deleteRelationship(rel) {
		let type1 = this.state.concept.type;
		let type2 = rel.item.type;
		let name1 = this.state.concept.name;
		let name2 = rel.item.name;
		let string1 = rel.string1;
		let string2 = rel.string2;
		if(rel.reverse) {
			[type1, type2] = [type2, type1];
			[name1, name2] = [name2, name1];
		}
		const req_data = {
			type1,
			type2,
			name1,
			name2,
			string1,
			string2
		};

		axios.post("/api/v1/ca/concept/delete/rel", req_data)
			.then(() => {
				this.closeModals()
					.then(() => {
						this.refreshConcept();
					});
			});
	}

	closeModals() {
		let promise = new Promise((resolve, reject) => {
			const state = Object.assign({}, this.state, {
				modal_delete_data: null,
				modal_delete_concept: null,
				modal_delete_relationship: null
			});

			this.setState(state, () => {
				resolve();
			});
		});

		return promise;
	}

	changeNewRelationshipType(event) {
		const state = Object.assign({}, this.state, {
			selected_addable_relationship: event.target.value
		});
		this.setState(state, () => {
			this.relOtherNameChange("")
				.then(() => {
					this.relOtherNameAutocomplete("");
				});
		});
	}

	relOtherNameAutocompleteShow(change) {
		const state = Object.assign({}, this.state, {
			rel_other_name_autocomplete: change
		});
		this.setState(state);
	}

	relOtherNameAutocomplete(value) {
		const type = this.state.addable_relationships[this.state.selected_addable_relationship].other_type;
		axios.get("/api/v1/ca/concept/data/" + type + "?q=" + value + "&count=100")
			.then((response) => {
				const state = Object.assign({}, this.state, {
					rel_other_concepts: response.data
				});
				this.setState(state);
			});
	}

	relOtherNameOnChange(event) {
		const value = event.target.value;
		const state = Object.assign({}, this.state, {
			rel_other_name: value
		});
		this.setState(state, () => {
			this.relOtherNameAutocomplete(value);
		});
	}

	relOtherNameChange(name) {
		const promise = new Promise((resolve, reject) => {
			const state = Object.assign({}, this.state, {
				rel_other_name: name
			});
			this.setState(state, () => {
				resolve();
			});
		});

		return promise;
	}

	render() {
		return (
			<Wrapper>
				<h2>Single concept</h2>
				<Table>
					<tbody>
						<tr>
							<Td>Type</Td>
							<Td>{ this.state.concept.type }</Td>
						</tr>
						<tr>
							<Td>Name</Td>
							<Td>{ this.state.concept.name }</Td>
						</tr>
					</tbody>
				</Table>
				<DeleteTypeButton onClick={this.deleteConceptConfirm}>Delete</DeleteTypeButton>

				<h4>Data</h4>
				{this.state.addable_data.length > 0 &&
					<div>
						<AddDataButton onClick={this.toggleShowNewData}>Add data</AddDataButton>
						<div>
							{this.state.show_new_data &&
								<form onSubmit={this.submitData}>
									<NewDataFormWrapper>
										<div>Field</div>
										<div>
											<select name="key">
												{
													this.state.addable_data.map((data_name, ix) => {
														return (
															<option key={ix}>
																{ data_name }
															</option>
														);
													})
												}
											</select>
										</div>
										<div>
											Data
										</div>
										<div>
											<textarea name="data"></textarea>
										</div>
										<input type="submit" value="Submit" />
									</NewDataFormWrapper>
								</form>
							}
						</div>
					</div>
				}
				<Table>
					<tbody>
						{this.state.concept.data && this.state.concept.data.map((data, ix) => {
							return (
								<tr key={ix}>
									<Td>{data.key}</Td>
									<Td>{data.value}</Td>
									<Td><a href="javascript:void(0)" onClick={() => this.deleteDataConfirm(data)}>Delete</a></Td>
								</tr>
							);
						})}
					</tbody>
				</Table>

				<h4>Relationships</h4>
				{this.state.addable_relationships.length > 0 &&
					<div>
						<AddDataButton onClick={this.toggleShowNewRelationship}>Add relationship</AddDataButton>
						<div>
							{this.state.show_new_relationship &&
								<form onSubmit={this.submitRelationship}>
									<NewDataFormWrapper>
										<div>Relationship</div>
										<div>
											<select name="relationship" onChange={this.changeNewRelationshipType}>
												{this.state.addable_relationships.map((rt, ix) => {
													return (
														<option key={ix} value={ix}>
															{ rt.display_string }
														</option>
													)
												})}
											</select>
										</div>
										{this.state.addable_relationships[this.state.selected_addable_relationship].reversible &&
											<div>
												Reverse
											</div>
										}
										{this.state.addable_relationships[this.state.selected_addable_relationship].reversible &&
											<div>
												<input type="checkbox" name="reverse" />
											</div>
										}
										<div>Other type</div>
										<div>{ this.state.addable_relationships[this.state.selected_addable_relationship].other_type }</div>
										<div>Other type name</div>
										<div>
											<input type="text" value={this.state.rel_other_name} onChange={this.relOtherNameOnChange} onFocus={() => this.relOtherNameAutocompleteShow(true)} onBlur={() => this.relOtherNameAutocompleteShow(false)} name="other_name" />
											{this.state.rel_other_name_autocomplete &&
												<Table>
													<tbody>
														{this.state.rel_other_concepts && this.state.rel_other_concepts.map((oc, ix) => {
															return (
																<tr onClick={() => this.relOtherNameChange(oc.name)} key={ix}><Td>{oc.name}</Td></tr>
															)
														})}
													</tbody>
												</Table>
											}
										</div>
										<div>
											<input type="submit" value="Submit" />
										</div>
									</NewDataFormWrapper>
								</form>
							}
						</div>
					</div>
				}
				<Table>
					<tbody>
						{this.state.concept.relationships && this.state.concept.relationships.map((rel, ix) => {
							return (
								<tr key={ix}>
									<Td>{rel.reverse ? rel.string2 : rel.string1}</Td>
									<Td>{rel.item.type}</Td>
									<Td>{rel.item.name}</Td>
									<Td><a href="javascript:void(0)" onClick={() => this.deleteRelationshipConfirm(rel)}>Delete</a></Td>
								</tr>
							);
						})}
					</tbody>
				</Table>

				{this.state.modal_delete_data &&
					<DeleteDataModal data={this.state.modal_delete_data} confirmFun={() => this.deleteData(this.state.modal_delete_data)} cancelFun={this.closeModals}>
					</DeleteDataModal>
				}
				{this.state.modal_delete_concept &&
					<DeleteConceptModal confirmFun={this.deleteConcept} cancelFun={this.closeModals}>
					</DeleteConceptModal>
				}
				{this.state.modal_delete_relationship &&
					<DeleteRelationshipModal confirmFun={() => this.deleteRelationship(this.state.modal_delete_relationship)} cancelFun={this.closeModals}>
					</DeleteRelationshipModal>
				}
			</Wrapper>
		);
	}
}

const Wrapper = styled.div`

`;

const button_css = css`
	color: white;
	cursor: pointer;
	width: 100px;
	text-align: center;
`;

const DeleteTypeButton = styled.span`
	background-color: red;
	${button_css}
`;

const AddDataButton = styled.div`
	background-color: blue;
	${button_css}
`;

//const ElementsWrapper = styled.div`
//	display: grid;
//	grid-template-columns: auto 1fr;
//	grid-gap: 5px;
//`;
//
//const ElementTitle = styled.div`
//	text-align: right;
//`;
//
//const ElementElement = styled.div`
//
//`;

const NewDataFormWrapper = styled.div`
	display: grid;
	grid-gap: 5px;
	grid-template-columns: auto 1fr;
`;

/* #################################################################################################### */

class DeleteDataModal extends Component {
	render() {
		return (
			<ModalWrapper onClick={this.props.cancelFun}>
				<ModalForm>
					Are you sure you want to delete this data?<br/>
					{ this.props.data.key }: { this.props.data.value }<br/>
					<a href="javascript:void(0)" onClick={this.props.confirmFun}>Delete</a>&nbsp;
					<a href="javascript:void(0)" onClick={this.props.cancelFun}>Cancel</a>
				</ModalForm>
			</ModalWrapper>
		);
	}
}

const ModalWrapper = styled.div`
	position: absolute;
	top: 0;
	left: 0;
	height: 100%;
	width: 100%;
	background-color: rgba(0, 0, 0, 0.5);
	z-index: 1;
`;

const ModalForm = styled.div`
	position: relative;
	top: 20%;
	left: 10%;
	width: 80%;
	background-color: white;
	padding: 10px;
`;

/* #################################################################################################### */

class DeleteConceptModal extends Component {
	render() {
		return (
			<ModalWrapper onClick={this.props.cancelFun}>
				<ModalForm>
					Are you sure you want to delete this concept?<br/>
					<a href="javascript:void(0)" onClick={this.props.confirmFun}>Delete</a>&nbsp;
					<a href="javascript:void(0)" onClick={this.props.cancelFun}>Cancel</a>
				</ModalForm>
			</ModalWrapper>
		);
	}
}

/* #################################################################################################### */

class DeleteRelationshipModal extends Component {
	render() {
		return (
			<ModalWrapper onClick={this.props.cancelFun}>
				<ModalForm>
					Are you sure you want to delete this relationship?<br/>
					<a href="javascript:void(0)" onClick={this.props.confirmFun}>Delete</a>&nbsp;
					<a href="javascript:void(0)" onClick={this.props.cancelFun}>Cancel</a>
				</ModalForm>
			</ModalWrapper>
		);
	}
}
