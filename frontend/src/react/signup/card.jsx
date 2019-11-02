import React from 'react'
import PropTypes from 'prop-types'
import {Card, Icon, Image} from 'semantic-ui-react'

import img from '../static/gopher.png'

export default class SignupCard extends React.Component {
    render() {
        const {firstName, lastName, dockerhubUsername, githubUsername} = this.props;
        return (
            <Card.Group centered>
                <Card raised>
                    <Image src={img} wrapped ui={false}/>
                    <Card.Content>
                        <Card.Header textAlign={'left'}>{firstName} {lastName}</Card.Header>
                    </Card.Content>
                    <Card.Content textAlign={'left'} extra>
                        <Icon name='github'/>
                        {githubUsername}
                        <span> </span>
                        <Icon name='docker'/>
                        {dockerhubUsername}
                    </Card.Content>
                </Card>
            </Card.Group>
        );
    }
}

SignupCard.propTypes = {
    firstName: PropTypes.string.isRequired,
    lastName: PropTypes.string.isRequired,
    githubUsername: PropTypes.string.isRequired,
    dockerhubUsername: PropTypes.string.isRequired
};