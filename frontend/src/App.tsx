import './App.css';
import './external/w3.css';
import { RepoHistoryViewer } from './RepoHistoryViewer';
import { BrowserRouter, Route, Routes, Navigate } from 'react-router-dom';
import { Common } from './Common';
import { DirectoryPicker } from './DirectoryPicker';
import { Component } from 'react';
import { localStorageAppPrefix } from './localStorage';

class Props {
}

class State {
    directory?: string;
    goTo?: string;
}

export class App extends Component<Props, State> {
    constructor(props: Props) {
        super(props);
        const state = new State();
        state.directory = localStorage.getItem(localStorageAppPrefix + '.directory') || undefined;
        if (!state.directory)
            state.goTo = Common.baseUrl + '/open-repository';
        this.state = state;
        document.title = 'Git Stories';
    }

    override render() {
        if (this.state.goTo)
            setTimeout(() => this.setState({ goTo: undefined }));
        return <div style={{margin: 4}}>
            <BrowserRouter>
                <div className="w3-bar w3-dark-grey" style={{marginBottom: 4, position: 'sticky', top: 0}}>
                    <a href={Common.baseUrl + '/'} className="w3-bar-item w3-black w3-btn">GIT-STORIES</a>
                </div>
                { this.state.goTo ? <Navigate to={this.state.goTo} /> : undefined }
                <Routes>
                    <Route path={Common.baseUrl + '/open-repository'} element={this.renderDirectoryPicker()} />
                    <Route path={Common.baseUrl} element={this.renderRepoHistoryViewer()} />
                </Routes>
            </BrowserRouter>
        </div>;
    }

    private renderDirectoryPicker() {
        return <DirectoryPicker
            setDirectory={ directory => this.setDirectory(directory) }
        />;
    }

    private setDirectory(directory: string) {
        this.setState({ directory, goTo: Common.baseUrl + '/' });
        localStorage.setItem(localStorageAppPrefix + '.directory', directory);
    }

    private renderRepoHistoryViewer() {
        return <RepoHistoryViewer directory={this.state.directory}/>
    }
}

export default App;
