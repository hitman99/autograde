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
class App extends React.Component {
    render() {
        return (
            <Router>
                <Switch>
                    <Route path="/">
                        <SignupForm />
                    </Route>
                    <Route path="/state">
                        {/*signup state for admin TODO*/}
                    </Route>
                </Switch>
            </Router>
        )
    }
}

ReactDOM.render(<App/>, document.getElementById('react-root'));