import { Component } from "react";
import { Common } from "./Common";

class Props {
}

interface FileInfo {
    name: string;
    path: string;
    isDirectory: boolean;
}

class State {
    currentDirectory: string = '';
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
        return <div className="w3-bar w3-border-bottom" style={{position: 'sticky', top: 0}}>
            <div className="w3-bar-item">Please pick directory containing a Git repository</div>
        </div>;
    }

    private async loadFileList() {
        const url = Common.apiUrl + '/fileList?directory=' +
            encodeURIComponent(this.state.currentDirectory);
        const response = await fetch(url);
        const files = await response.json();
        this.setState({files: files});
    }
}