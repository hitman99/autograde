import React from 'react'
import {Button, Form, Grid, Header, Segment, Container} from 'semantic-ui-react'
import {register} from '../utils/api'
import SignupCard from './card'
class SignupForm extends React.Component {

    loadState() {
        let reg = localStorage.getItem("registrationData")
        if (reg) {
            return JSON.parse(reg)
        } else {
            return {
                firstName: '',
                lastName: '',
                githubUsername: '',
                dockerhubUsername: '',
                k8sNamespace: ''
            }
        }
    }

    persistState() {
        localStorage.setItem("registrationData", JSON.stringify(this.state.regData))
    }

    canRegister() {
        const { regData } = this.state;
        return Object.values(regData)
            .filter(field => field.length > 3).length === Object.keys(regData).length && !this.state.registered;
    }

    constructor(props) {
        super(props);
        let regData = this.loadState();
        this.state = {
            regData,
            registered: regData.githubUsername !== '',
            isError: false,
            isLoading: false
        };
    }

    handleInputChange(which, ev) {
        let regData = {...this.state.regData};
        regData[which] = ev.value;
        if (which === 'githubUsername') {
            regData['k8sNamespace'] = `ktu-stud-${ev.value}`;
        }
        this.setState({regData})
    }

    async submit() {
        this.setState({isLoading: true});
        let res = await register(this.state.regData);
        if (res !== 'ok') {
            console.log(res);
            this.setState({isError: true, isLoading: false});
        } else {
            this.persistState();
            this.setState({registered: true, isLoading: false});
        }
    }

    render() {
        const {firstName, lastName, githubUsername, dockerhubUsername} = this.state.regData;
        const { registered, isError, isLoading } = this.state;
        let title;
        let formOrCard;

        if (!registered) {
            title = <span>Registracija į <br/> "Clouds, Containers and Code" lab darbą</span>
            formOrCard =
                <Form size='large'>
                    <Segment>
                        <Form.Input inverted fluid icon='user' iconPosition='left' placeholder='Vardas'
                                    value={firstName} onChange={(e, d)=>{this.handleInputChange('firstName', d)}}
                                    disabled={registered}
                        />
                        <Form.Input fluid icon='user' iconPosition='left' placeholder='Pavardė'
                                    value={lastName} onChange={(e, d)=>{this.handleInputChange('lastName', d)}}
                                    disabled={registered}
                        />
                        <Form.Input fluid icon='github' iconPosition='left' placeholder='GitHub username'
                                    value={githubUsername} onChange={(e, d)=>{this.handleInputChange('githubUsername', d)}}
                                    disabled={registered}
                        />
                        <Form.Input fluid icon='docker' iconPosition='left' placeholder='DockerHub username'
                                    value={dockerhubUsername} onChange={(e, d)=>{this.handleInputChange('dockerhubUsername', d)}}
                                    disabled={registered}
                        />
                        <Button loading={isLoading} color={isError ? 'red' : 'teal'} fluid size='large' onClick={()=>{
                            this.submit()
                        }} disabled={!this.canRegister()}>
                            Registruotis
                        </Button>
                    </Segment>
                </Form>
        } else {
            title = 'Registracija sėkminga';
            formOrCard =
                <Container textAlign='center'>
                    <SignupCard firstName={firstName}
                             lastName={lastName}
                             githubUsername={githubUsername}
                             dockerhubUsername={dockerhubUsername} />
                </Container>
        }
        return (
            <Grid textAlign='center' style={{height: '100vh'}} verticalAlign='middle'>
                <Grid.Column style={{maxWidth: 550}}>
                    <Header as='h2' color='teal' textAlign='center'>
                        {title}
                    </Header>
                    { formOrCard }
                </Grid.Column>
            </Grid>
        )
    }
}

export default SignupForm