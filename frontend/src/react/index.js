import React from 'react';
import ReactDOM from 'react-dom';
import SignupForm from "./signup/signup";
import {
    BrowserRouter as Router,
    Switch,
    Route,
    Link
} from "react-router-dom";
import 'semantic-ui-css/semantic.min.css'
import LabScenario from "./lab/lab";
class App extends React.Component {
    render() {
        return (
            <Router>
                <Switch>
                    <Route path="/control">
                        <LabScenario />
                    </Route>
                    <Route path="/">
                        <SignupForm />
                    </Route>
                </Switch>
            </Router>
        )
    }
}

ReactDOM.render(<App/>, document.getElementById('react-root'));