import React from 'react'
import PropTypes from 'prop-types'
import {Container, Statistic} from 'semantic-ui-react'

export default class LabState extends React.Component {
    render() {
        return (
            <Container>
                <Statistic color='blue'>
                    <Statistic.Value>{this.props.count}</Statistic.Value>
                    <Statistic.Label>Students</Statistic.Label>
                </Statistic>
            </Container>
        );
    }
}

LabState.propTypes = {
    count: PropTypes.number.isRequired
};