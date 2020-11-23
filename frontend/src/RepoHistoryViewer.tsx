import React, { ChangeEvent, Component } from 'react';
import { Common } from './Common';
import { localStorageAppPrefix } from './localStorage';
import { StoryEntry } from './StoryEntry';

class Props {
}

class State {
    repoDirectory: string;
    stories: StoryEntry[];
}

export class RepoHistoryViewer extends Component<Props, State> {
    constructor(props: Props) {
        super(props);
        const state = new State();
        state.repoDirectory = localStorage.getItem(localStorageAppPrefix + '.repoDirectory');
        if (null == state.repoDirectory)
            state.repoDirectory = '';
        this.state = state;
    }

    render() {
        return <div className="w3-panel">
            <input type="text"
                value={this.state.repoDirectory}
                onChange={this.receiveFilePathChange.bind(this)}
            />
            <button className="w3-btn w3-black" onClick={this.receiveLoadClick.bind(this)}>
                LOAD
            </button>
            {this.state.stories != null
                ? this.renderStories()
                : null
            }
        </div>;
    }

    componentDidUpdate() {
        localStorage.setItem(localStorageAppPrefix + '.repoDirectory', this.state.repoDirectory);
    }

    private receiveFilePathChange(event: ChangeEvent<HTMLInputElement>) {
        const repoDirectory = event.target['value'];
        this.setState({repoDirectory: repoDirectory});
    }

    private async receiveLoadClick() {
        const url = Common.apiUrl + '/stories?' +
            'directory=' + encodeURIComponent(this.state.repoDirectory) + '&' +
            'lengthLimit=20';
        const response = await fetch(url);
        const stories: StoryEntry[] = await response.json();
        this.setState({stories});
    }

    private renderStoryEntry(entry: StoryEntry) {
        return <li>
            {entry.Description}
        </li>;
    }

    private renderStories() {
        return <ul>
            { this.state.stories.map(story => this.renderStoryEntry(story)) }
        </ul>;
    }
}