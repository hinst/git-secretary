import { RepoHistoryViewer } from './RepoHistoryViewer';
import './App.css';
import './git-stories.css';
import './external/w3.css';

function App() {
  document.title = 'Git Stories';
  return <RepoHistoryViewer/>;
}

export default App;
