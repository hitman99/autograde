import React from "react";
import PropTypes from "prop-types";
import {Header, Image, Table} from "semantic-ui-react";

import img from '../static/gopher.png'

export default class StudentTable extends React.Component {
    fetchData() {

    }
    render() {
        return(
            <Table basic='very' celled collapsing>
                <Table.Header>
                    <Table.Row>
                        <Table.HeaderCell>Employee</Table.HeaderCell>
                        <Table.HeaderCell>Correct Guesses</Table.HeaderCell>
                    </Table.Row>
                </Table.Header>

                <Table.Body>
                    <Table.Row>
                        <Table.Cell>
                            <Header as='h4' image>
                                <Image src={img} rounded size='mini' />
                                <Header.Content>
                                    Lena
                                    <Header.Subheader>Human Resources</Header.Subheader>
                                </Header.Content>
                            </Header>
                        </Table.Cell>
                        <Table.Cell>22</Table.Cell>
                    </Table.Row>
                    <Table.Row>
                        <Table.Cell>
                            <Header as='h4' image>
                                <Image src={img} rounded size='mini' />
                                <Header.Content>
                                    Matthew
                                    <Header.Subheader>Fabric Design</Header.Subheader>
                                </Header.Content>
                            </Header>
                        </Table.Cell>
                        <Table.Cell>15</Table.Cell>
                    </Table.Row>
                    <Table.Row>
                        <Table.Cell>
                            <Header as='h4' image>
                                <Header.Content>
                                    Lindsay
                                    <Header.Subheader>Entertainment</Header.Subheader>
                                </Header.Content>
                            </Header>
                        </Table.Cell>
                        <Table.Cell>12</Table.Cell>
                    </Table.Row>
                    <Table.Row>
                        <Table.Cell>
                            <Header as='h4' image>
                                <Header.Content>
                                    Mark
                                    <Header.Subheader>Executive</Header.Subheader>
                                </Header.Content>
                            </Header>
                        </Table.Cell>
                        <Table.Cell>11</Table.Cell>
                    </Table.Row>
                </Table.Body>
            </Table>
        )
    }
}

StudentTable.propTypes = {
    students: PropTypes.arrayOf(PropTypes.object).isRequired
}