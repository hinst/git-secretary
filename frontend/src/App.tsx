import { RepoHistoryViewer } from './RepoHistoryViewer';
import './App.css';
import './external/w3.css';
import { BrowserRouter, Route, Routes } from 'react-router-dom';
import { Common } from './Common';
import { DirectoryPicker } from './DirectoryPicker';
import { Component } from 'react';

class Props {
}

class State {
    directory?: string;
}

export class App extends Component<Props, State> {
    constructor(props: Props) {
        super(props);
        this.state = new State();
        document.title = 'Git Stories';
    }

    override render() {
        return <div style={{margin: 4}}>
            <BrowserRouter>
                <div className="w3-bar w3-dark-grey" style={{marginBottom: 4, position: 'sticky', top: 0}}>
                    <a href={Common.baseUrl + '/'} className="w3-bar-item w3-black w3-btn">GIT-STORIES</a>
                </div>
                <Routes>
                    <Route path={Common.baseUrl + '/open-repository'} element={this.renderDirectoryPicker()} />
                    <Route path={Common.baseUrl} element={<RepoHistoryViewer/>} />
                </Routes>
            </BrowserRouter>
        </div>;
    }

    private renderDirectoryPicker() {
        return <DirectoryPicker setDirectory={(directory) => console.log(directory)}/>;
    }
}

export default App;
