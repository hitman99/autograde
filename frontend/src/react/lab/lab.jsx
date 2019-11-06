import React, {Fragment} from 'react'
import PropTypes from 'prop-types'
import {Card, Container, Header, Icon, Input, List, Grid, Segment, Divider, Label, Button} from 'semantic-ui-react'

import img from '../static/gopher.png'
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
    render(){
        return(
            <Fragment>
                <Input
                    placeholder={'admin token'}
                    onChange={e => {
                        this.setState({token: e.target.value})
                    }}
                />
                <Button color='teal' content='Set' icon='key' labelPosition='right' onClick={() => this.setToken()} />
            </Fragment>
        )
    }
}

AdminToken.propTypes = {
    onTokenSet: PropTypes.func.isRequired
};

export default class LabScenario extends React.Component {
    constructor(props){
        super(props);
        this.state = {
            adminToken: this.getToken(),
            scenario: {
                studentsKey: "",
                cycle: "",
                duration: "",
                tasksKey: "",
                name: ""
            }
        };
        this.fHandle = setInterval(async ()=>{
            let state = await this.fetchState()
        }, 5 * 1000)
    }
    getToken() {
        return localStorage.getItem("adminToken")
    }

    async createScenario() {
        const { adminToken, scenario } = this.state;
        try {
            let res = await fetch('/lab/scenario', {
                headers: {
                    'Authorization': `Bearer ${btoa(adminToken)}`
                },
                method: 'POST',
                body: JSON.stringify(scenario)
            });
            return res.statusText;
        } catch(err) {
            console.log(err);
            return null;
        }
    }

    async startLab(fromState = false) {
        const { adminToken } = this.state;
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
        } catch(err) {
            console.log(err);
            return null;
        }
    }
    async stopLab() {
        const { adminToken } = this.state;
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
        } catch(err) {
            console.log(err);
            return null;
        }
    }

    async createDeps(redisKey){
        const { adminToken } = this.state;
        try {
            let res = await fetch(`/lab/deps/${redisKey}`, {
                headers: {
                    'Authorization': `Bearer ${btoa(adminToken)}`
                },
                method: 'POST'
            });
            return res.statusText;
        } catch(err) {
            console.log(err);
            return null;
        }
    }

    async deleteDeps(redisKey){
        const { adminToken } = this.state;
        try {
            let res = await fetch(`/lab/deps/${redisKey}`, {
                headers: {
                    'Authorization': `Bearer ${btoa(adminToken)}`
                },
                method: 'DELETE'
            });
            return res.statusText;
        } catch(err) {
            console.log(err);
            return null;
        }
    }

    async fetchState() {
        const { adminToken } = this.state;
        try {
            let res = await fetch('/lab/scenario/state', {
                headers: {
                    'Authorization': `Bearer ${btoa(adminToken)}`
                }
            });
            let state = await res.json();
            return state;
        } catch(err) {
            console.log(err);
            return null;
        }
    }

    render() {
        const {adminToken} = this.state;
        if (!adminToken) {
            return (
                <Segment basic vertical style={{ margin: '1em 1em 1em', padding: '1em 1em' }}>
                    <Grid centered>
                        <AdminToken onTokenSet={(adminToken) => {
                            this.setState({adminToken})
                        }}/>
                    </Grid>
                </Segment>
            );
        }
        return (
            <Segment basic vertical style={{ margin: '1em 1em 1em', padding: '1em 1em' }}>
                <Segment raised>
                    <Label as='div' ribbon>
                        Lab Control
                    </Label>
                    <Grid columns='four'>
                        <Grid.Column>
                        </Grid.Column>
                        <Grid.Column>
                            <Input size='mini' icon='users' iconPosition='left' placeholder='Name'
                                   onChange={e => {
                                       const {scenario} = {...this.state};
                                       scenario.name = e.target.value;
                                       this.setState({scenario})
                                   }}/>
                            <Input size='mini' icon='users' iconPosition='left' placeholder='Lab Cycle'
                                   onChange={e => {
                                       const {scenario} = {...this.state};
                                       scenario.cycle = e.target.value;
                                       this.setState({scenario});
                                   }}/>
                            <Input size='mini' icon='users' iconPosition='left' placeholder='Redis Key'
                                   onChange={e => {
                                       const {scenario} = {...this.state};
                                       scenario.studentsKey = e.target.value;
                                       this.setState({scenario})
                                   }}/>
                        </Grid.Column>
                        <Grid.Column>
                            <Input size='mini' icon='users' iconPosition='left' placeholder='Duration'
                                   onChange={e => {
                                       const {scenario} = {...this.state};
                                       scenario.duration = e.target.value;
                                       this.setState({scenario})
                                   }}/>
                            <Input size='mini' icon='users' iconPosition='left' placeholder='Tasks Key'
                                   onChange={e => {
                                       const {scenario} = {...this.state};
                                       scenario.tasksKey = e.target.value;
                                       this.setState({scenario})
                                   }}/>
                            <Button basic size={'mini'} color='red' content='Create' onClick={() => this.createScenario()} />
                        </Grid.Column>
                        <Grid.Column>
                            <Button.Group size={'mini'}>
                                <Button basic color='green' content='Start' />
                                <Button basic color='green' content='Start From State' />
                                <Button basic color='red' content='Stop' />
                            </Button.Group>
                            <Divider />
                            <Button.Group size='mini'>
                                <Button basic color='orange' content='Provision Deps' />
                                <Button basic color='orange' content='Delete Deps' />
                            </Button.Group>
                        </Grid.Column>
                    </Grid>
                </Segment>
                <Segment raised>
                    <Container textAlign='center'>
                        <Grid columns='two' divided>
                            <Grid.Row>
                                <Grid.Column>
                                    <Header size='large'> Lab Cycle 1 </Header>
                                </Grid.Column>
                                <Grid.Column>
                                    <LabState count={12}/>
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
                        <StudentTable />
                    </Grid>
                </Segment>

            </Segment>
        );
    }
}

LabScenario.propTypes = {};