import { RepoHistoryViewer } from './RepoHistoryViewer';
import './App.css';
import './git-stories.css';
import './external/w3.css';
import { BrowserRouter, Route, Routes } from 'react-router-dom';

function App() {
  document.title = 'Git Stories';
  return (<BrowserRouter>
    <Routes>
      <Route path="/" element={<RepoHistoryViewer/>}>
      </Route>
    </Routes>
  </BrowserRouter>);
}

export default App;
