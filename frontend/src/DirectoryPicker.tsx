import { Component, ReactNode } from "react";

class Props {
}

class State {
    currentDirectory?: string = undefined;
}

export class DirectoryPicker extends Component<Props, State> {
    override render(): ReactNode {
        return <div className="w3-panel">
            <div className="w3-bar-item">Choose a directory containing a Git repository</div>
        </div>;
    }
}