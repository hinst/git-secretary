import lodash from 'lodash';
import { Component, CSSProperties } from 'react';
import ArticleIcon from '@mui/icons-material/Article';
import FolderOpenIcon from '@mui/icons-material/FolderOpen';
import ArrowBackIcon from '@mui/icons-material/ArrowBack';
import { Common } from './Common';
import { LinearProgress } from '@material-ui/core';
import { replaceAll, splitAll } from './string';
import { sleep } from './sleep';

class Props {
    setDirectory: (directory: string) => void = () => {};
}

interface FileInfo {
    name: string;
    path: string;
    isDirectory: boolean;
}

class State {
    isLoading: boolean = false;
    directory: string = '';
    files: FileInfo[] = [];
}

const fileIconStyle: CSSProperties = { verticalAlign: 'middle', marginRight: 4 };
const fileItemStyle: CSSProperties = { paddingTop: 2, paddingBottom: 2, paddingLeft: 4 };

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
            <div className="w3-bar w3-dark-grey" style={{position: 'sticky', top: 0}}>
                { this.state.directory === ''
                    ? <div className="w3-bar-item">Please pick your Git repository</div>
                    : [
                        <button
                            key="OkButton"
                            className="w3-bar-item w3-btn w3-black"
                            onClick={() => this.clickOk()}
                        >
                            OK
                        </button>,
                        <div key="DirectoryString" className="w3-bar-item">
                            { replaceAll('\\', '/', this.state.directory) }
                        </div>
                    ]
                }
            </div>
            <LinearProgress style={{visibility: this.state.isLoading ? 'visible' : 'hidden' }} />
            <div>
                { this.state.directory.length ? this.renderParentDirectory() : undefined }
                {this.renderFiles()}
            </div>
        </div>
    }

    private async loadFileList() {
        this.setState({isLoading: true});
        const url = Common.apiUrl + '/fileList?directory=' +
            encodeURIComponent(this.state.directory);
        const response = await fetch(url);
        if (response.ok) {
            const files = await response.json();
            this.setState({files: files, isLoading: false});
        } else {
            this.setState({files: [], isLoading: false})
            alert(response.statusText + '\n' + await response.text());
        }
    }

    private renderFiles() {
        const files = lodash.sortBy(this.state.files,
            file => file.isDirectory ? 0 : 1,
            file => file.name);
        return files.map(file => this.renderFile(file));
    }

    private renderFile(file: FileInfo) {
        const icon = file.isDirectory
            ? <FolderOpenIcon style={fileIconStyle}/>
            : <ArticleIcon style={fileIconStyle}/>;
        const className = file.isDirectory ? 'GitStories_FilePicker_ClickableItem' : undefined;
        const onClick = file.isDirectory ? () => this.clickFile(file) : () => {};
        return <div
            key={file.path}
            onClick={onClick}
            className={className}
            style={fileItemStyle}
        >
            {icon} {file.name}
        </div>;
    }

    private renderParentDirectory() {
        return <div
            onClick={() => this.goToParent()}
            className="GitStories_FilePicker_ClickableItem"
            style={fileItemStyle}
        >
            <ArrowBackIcon style={fileIconStyle}/> <i>go to parent directory</i>
        </div>
    }

    private goToParent() {
        const parts = splitAll(this.state.directory, ['/', '\\']);
        parts.pop();
        this.setState({ directory: parts.join('/') });
        setTimeout(() => this.loadFileList());
    }

    private async clickFile(file: FileInfo) {
        this.setState({directory: file.path});
        await sleep();
        await this.loadFileList();
        window.scrollTo(0, 0);
    }

    private clickOk() {
        const isGitDirectory = this.state.files.some(file =>
            file.name.toLowerCase() === '.git' && file.isDirectory);
        if (isGitDirectory)
            this.props.setDirectory(this.state.directory);
        else
            alert('This folder does not look like a Git repository because it does not contain .git');
    }
}