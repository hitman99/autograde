import React from "react";
import PropTypes from "prop-types";
import {Header, Icon, Image, Table} from "semantic-ui-react";

import img from '../static/gopher.png'

export default class StudentTable extends React.Component {
  constructor(props) {
    super(props);
  }

  render() {
    const {students} = this.props;
    let tasks = students.length ? students[0].tasks : [];
    return (
      <Table basic='very' celled collapsing>
        <Table.Header>
          <Table.Row>
            <Table.HeaderCell>Student</Table.HeaderCell>
            <Table.HeaderCell>Score</Table.HeaderCell>
            {
              tasks.map((t, i) => {
                const {taskDefinition: td} = t;
                return (
                  <Table.HeaderCell key={i}>Task #{i+1}</Table.HeaderCell>
                )
              })
            }
          </Table.Row>
        </Table.Header>

        <Table.Body>

          {
            students.map(s => {
              return (
                <Table.Row key={s.student.k8sNamespace} textAlign='center'>
                  <Table.Cell>
                    <Header as='h4' image>
                      <Image src={img} rounded size='mini'/>
                      <Header.Content>
                        {s.student.firstName} {s.student.lastName}
                      </Header.Content>
                    </Header>
                  </Table.Cell>
                  <Table.HeaderCell>
                    {
                      s.tasks.reduce((score, t) => {
                        return score + (t.completed ? t.taskDefinition.score : 0);
                      }, 0)
                    }
                  </Table.HeaderCell>
                  {
                    s.tasks.map((t, i) => {
                      const {taskDefinition: td} = t;
                      return (
                        <Table.Cell key={i}>{t.completed ? <Icon name='check'/> : <Icon name='close'/>}</Table.Cell>
                      )
                    })
                  }
                </Table.Row>
              )
            })
          }

        </Table.Body>
      </Table>
    )
  }
}

StudentTable.propTypes = {
  students: PropTypes.arrayOf(PropTypes.object).isRequired
}