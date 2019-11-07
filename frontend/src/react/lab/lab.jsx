import React, {Fragment} from 'react'
import PropTypes from 'prop-types'
import {Button, Container, Divider, Grid, Header, Input, Label, Segment} from 'semantic-ui-react'
import LabState from "./state";
import StudentTable from "./studentTable";

class AdminToken extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      token: ""
    };
  }

  setToken() {
    if (this.state.token !== "") {
      localStorage.setItem("adminToken", this.state.token);
      this.props.onTokenSet(this.state.token)
    }
  }

  getToken() {
    return localStorage.getItem("adminToken")
  }

  render() {
    return (
      <Fragment>
        <Input
          placeholder={'admin token'}
          onChange={e => {
            this.setState({token: e.target.value})
          }}
        />
        <Button color='teal' content='Set' icon='key' labelPosition='right' onClick={() => this.setToken()}/>
      </Fragment>
    )
  }
}

AdminToken.propTypes = {
  onTokenSet: PropTypes.func.isRequired
};

export default class LabScenario extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      adminToken: this.getToken(),
      scenario: this.getScenario(),
      labState: {},
      errorFetching: null
    };
    this.tHandle = null;
    this.fHandle = setInterval(async () => {
      await this.fetchState();

    }, 10 * 1000)

  }

  async componentDidMount() {
    await this.fetchState();
  }

  getToken() {
    return localStorage.getItem("adminToken")
  }

  async createScenario() {
    const {adminToken, scenario} = this.state;
    try {
      let res = await fetch('/lab/scenario', {
        headers: {
          'Authorization': `Bearer ${btoa(adminToken)}`
        },
        method: 'POST',
        body: JSON.stringify(scenario)
      });
      return res.statusText;
    } catch (err) {
      console.log(err);
      return null;
    }
  }

  async startLab(fromState = false) {
    const {adminToken} = this.state;
    try {
      let res = await fetch('/lab/scenario', {
        headers: {
          'Authorization': `Bearer ${btoa(adminToken)}`
        },
        method: 'PATCH',
        body: JSON.stringify({
          action: !fromState ? 'start' : 'startFromState'
        })
      });
      return res.statusText;
    } catch (err) {
      console.log(err);
      return null;
    }
  }

  async stopLab() {
    const {adminToken} = this.state;
    try {
      let res = await fetch('/lab/scenario', {
        headers: {
          'Authorization': `Bearer ${btoa(adminToken)}`
        },
        method: 'PATCH',
        body: JSON.stringify({
          action: 'stop'
        })
      });
      return res.statusText;
    } catch (err) {
      console.log(err);
      return null;
    }
  }

  async createDeps() {
    const {adminToken, scenario} = this.state;
    try {
      let res = await fetch(`/lab/deps/${scenario.studentsKey}`, {
        headers: {
          'Authorization': `Bearer ${btoa(adminToken)}`
        },
        method: 'POST'
      });
      return res.statusText;
    } catch (err) {
      console.log(err);
      return null;
    }
  }

  async deleteDeps() {
    const {adminToken, scenario} = this.state;
    try {
      let res = await fetch(`/lab/deps/${scenario.studentsKey}`, {
        headers: {
          'Authorization': `Bearer ${btoa(adminToken)}`
        },
        method: 'DELETE'
      });
      return res.statusText;
    } catch (err) {
      console.log(err);
      return null;
    }
  }

  async fetchState() {
    const {adminToken} = this.state;
    try {
      let res = await fetch('/lab/scenario/state', {
        headers: {
          'Authorization': `Bearer ${btoa(adminToken)}`
        }
      });
      let state = await res.json();
      this.setState({labState: state, errorFetching: null});
    } catch (err) {
      console.log(err);
      this.setState({labState: {}, errorFetching: null});
    }
  }

  saveScenario() {
    localStorage.setItem("scenario", JSON.stringify(this.state.scenario))
  }

  getScenario() {
    let scenario = localStorage.getItem("scenario");
    if (scenario) {
      return JSON.parse(scenario);
    } else {
      return {
        studentsKey: "",
        cycle: "",
        duration: "",
        tasksKey: "",
        name: ""
      }
    }
  }

  render() {
    const {adminToken, labState} = this.state;
    const {studentsKey, cycle, duration, tasksKey, name} = this.state.scenario;
    if (!adminToken) {
      return (
        <Segment basic vertical style={{margin: '1em 1em 1em', padding: '1em 1em'}}>
          <Grid centered>
            <AdminToken onTokenSet={(adminToken) => {
              this.setState({adminToken})
            }}/>
          </Grid>
        </Segment>
      );
    }
    return (
      <Segment basic vertical style={{margin: '1em 1em 1em', padding: '1em 1em'}}>
        <Segment raised>
          <Label as='div' ribbon>
            Lab Control
          </Label>
          <Grid columns='four'>
            <Grid.Column>
            </Grid.Column>
            <Grid.Column>
              <Input size='mini' icon='lab' iconPosition='left' placeholder='Name'
                     onChange={e => {
                       const {scenario} = {...this.state};
                       scenario.name = e.target.value;
                       this.setState({scenario});
                       this.saveScenario();
                     }}
                     value={name}
              />
              <Input size='mini' icon='numbered list' iconPosition='left' placeholder='Lab Cycle'
                     onChange={e => {
                       const {scenario} = {...this.state};
                       scenario.cycle = e.target.value;
                       this.setState({scenario});
                       this.saveScenario();
                     }}
                     value={cycle}

              />
              <Input size='mini' icon='users' iconPosition='left' placeholder='Redis Key'
                     onChange={e => {
                       const {scenario} = {...this.state};
                       scenario.studentsKey = e.target.value;
                       this.setState({scenario});
                       this.saveScenario();
                     }}
                     value={studentsKey}
              />
            </Grid.Column>
            <Grid.Column>
              <Input size='mini' icon='clock outline' iconPosition='left' placeholder='Duration'
                     onChange={e => {
                       const {scenario} = {...this.state};
                       scenario.duration = e.target.value;
                       this.setState({scenario});
                       this.saveScenario();
                     }}
                     value={duration}
              />
              <Input size='mini' icon='tasks' iconPosition='left' placeholder='Tasks Key'
                     onChange={e => {
                       const {scenario} = {...this.state};
                       scenario.tasksKey = e.target.value;
                       this.setState({scenario});
                       this.saveScenario();
                     }}
                     value={tasksKey}
              />
              <Button size={'mini'} color='red' content='Create' onClick={() => this.createScenario()}/>
            </Grid.Column>
            <Grid.Column>
              <Button.Group size={'mini'}>
                <Button color='green' content='Start' onClick={() => this.startLab()}/>
                <Button color='green' content='Start From State' onClick={() => this.startLab(true)}/>
                <Button color='red' content='Stop' onClick={() => this.stopLab()}/>
              </Button.Group>
              <Divider/>
              <Button.Group size='mini'>
                <Button color='orange' content='Provision Deps' onClick={() => this.createDeps()}/>
                <Button color='orange' content='Delete Deps' onClick={() => this.deleteDeps()}/>
              </Button.Group>
            </Grid.Column>
          </Grid>
        </Segment>
        <Segment raised>
          <Container textAlign='center'>
            <Grid columns='two' divided>
              <Grid.Row>
                <Grid.Column>
                  <Header size='large'> {labState ? labState.name : 'Unknown'} {labState ? labState.cycle : null} </Header>
                </Grid.Column>
                <Grid.Column>
                  <LabState count={labState.assignments ? labState.assignments.length : 0}/>
                </Grid.Column>
              </Grid.Row>
            </Grid>
          </Container>
        </Segment>

        <Segment raised>
          <Label as='div' color='green' ribbon>
            Results
          </Label>
          <Grid centered>
            <StudentTable students={labState.assignments ? labState.assignments : []}/>
          </Grid>
        </Segment>

      </Segment>
    );
  }
}

LabScenario.propTypes = {};