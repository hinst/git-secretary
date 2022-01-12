import { Component } from 'react';
import ArticleIcon from '@mui/icons-material/Article';
import FolderOpenIcon from '@mui/icons-material/FolderOpen';
import { Common } from './Common';

class Props {
}

interface FileInfo {
    name: string;
    path: string;
    isDirectory: boolean;
}

class State {
    currentDirectory: string = 'C:';
    files: FileInfo[] = [];
}

export class DirectoryPicker extends Component<Props, State> {
    constructor(props: Props) {
        super(props);
        this.state = new State();
    }

    componentDidMount() {
        this.loadFileList();
    }

    override render() {
        return <div>
            <div className="w3-bar w3-border-bottom" style={{position: 'sticky', top: 0}}>
                <div className="w3-bar-item">Please pick directory containing a Git repository</div>
            </div>
            <div className="w3-container"> {this.renderFiles()} </div>
        </div>
    }

    private async loadFileList() {
        const url = Common.apiUrl + '/fileList?directory=' +
            encodeURIComponent(this.state.currentDirectory);
        const response = await fetch(url);
        const files = await response.json();
        this.setState({files: files});
    }

    private renderFiles() {
        return this.state.files.map(this.renderFile.bind(this));
    }

    private renderFile(file: FileInfo) {
        const iconStyle = {verticalAlign: 'middle', marginRight: 4};
        const icon = file.isDirectory
            ? <FolderOpenIcon style={iconStyle}/>
            : <ArticleIcon style={iconStyle}/>;
        return <div style={{marginTop: 4}}> {icon} {file.name} </div>
    }
}