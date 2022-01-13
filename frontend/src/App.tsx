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
    goToHome: boolean = false;
}

export class App extends Component<Props, State> {
    constructor(props: Props) {
        super(props);
        const state = new State();
        state.directory = localStorage.getItem(localStorageAppPrefix + '.directory') || undefined;
        this.state = state;
        document.title = 'Git Stories';
    }

    override render() {
        if (this.state.goToHome)
            setTimeout(() => this.setState({ goToHome: false }));
        return <div style={{margin: 4}}>
            <BrowserRouter>
                <div className="w3-bar w3-dark-grey" style={{marginBottom: 4, position: 'sticky', top: 0}}>
                    <a href={Common.baseUrl + '/'} className="w3-bar-item w3-black w3-btn">GIT-STORIES</a>
                </div>
                { this.state.goToHome ? <Navigate to={Common.baseUrl + '/'} /> : undefined }
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
        this.setState({ directory, goToHome: true });
        localStorage.setItem(localStorageAppPrefix + '.directory', directory);
    }

    private renderRepoHistoryViewer() {
        return <RepoHistoryViewer directory={this.state.directory}/>
    }
}

export default App;
