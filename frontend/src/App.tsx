import React, { ChangeEvent, Component } from 'react';
import './w3.css';
import './git-stories.css'

class Props {
}

class State {
    filePath: string = '';
}

class App extends Component<Props, State> {
    constructor(props: Props) {
        super(props);
        this.state = new State();
    }

    render() {
        return <div className="w3-panel">
            <input type="text"
                value={this.state.filePath}
                onChange={this.receiveFilePathChange.bind(this)}
            />
            <button className="w3-btn w3-black">WHY</button>
        </div>;
    }

    private receiveFilePathChange(event: ChangeEvent<HTMLInputElement>) {
        const filePath = event.target['value'];
        this.setState({filePath});
    }
}

export default App;
