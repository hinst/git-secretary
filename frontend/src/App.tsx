import { RepoHistoryViewer } from './RepoHistoryViewer';
import './App.css';
import './git-stories.css';
import './external/w3.css';
import { BrowserRouter, Route, Routes } from 'react-router-dom';
import { Common } from './Common';

function App() {
  document.title = 'Git Stories';
  return <BrowserRouter>
    <Routes>
      <Route path={Common.baseUrl + '/open-repository'} />
      <Route path={Common.baseUrl} element={<RepoHistoryViewer/>} />
    </Routes>
  </BrowserRouter>;
}

export default App;
