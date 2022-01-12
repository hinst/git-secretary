import { RepoHistoryViewer } from './RepoHistoryViewer';
import './App.css';
import './external/w3.css';
import { BrowserRouter, Route, Routes } from 'react-router-dom';
import { Common } from './Common';
import { DirectoryPicker } from './DirectoryPicker';

function App() {
    document.title = 'Git Stories';
    const directoryPicker = <DirectoryPicker setDirectory={(directory) => console.log(directory)}/>;
    return <div style={{margin: 4}}>
        <BrowserRouter>
            <div className="w3-bar w3-dark-grey" style={{marginBottom: 4, position: 'sticky', top: 0}}>
                <a href={Common.baseUrl + '/'} className="w3-bar-item w3-black w3-btn">GIT-STORIES</a>
            </div>
            <Routes>
                <Route path={Common.baseUrl + '/open-repository'} element={directoryPicker} />
                <Route path={Common.baseUrl} element={<RepoHistoryViewer/>} />
            </Routes>
        </BrowserRouter>
    </div>;
}

export default App;
