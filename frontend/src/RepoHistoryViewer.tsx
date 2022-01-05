import React, { ChangeEvent, Component } from 'react';
import { Common } from './Common';
import { localStorageAppPrefix } from './localStorage';
import { StoryEntry } from './StoryEntry';
import lodash from 'lodash';
import { getStartOfDay } from './dateTime';
import { LinearProgress } from '@material-ui/core';

class Props {
}

class State {
    repoDirectory?: string;
    stories: StoryEntry[] = [];
    isLoading: boolean = false;
}

export class RepoHistoryViewer extends Component<Props, State> {
    constructor(props: Props) {
        super(props);
        const state = new State();
        state.repoDirectory = localStorage.getItem(localStorageAppPrefix + '.repoDirectory') || undefined;
        if (null == state.repoDirectory)
            state.repoDirectory = '';
        this.state = state;
    }

    render() {
        return <div className="w3-panel">
            <input type="text" className="w3-input"
                value={this.state.repoDirectory}
                onChange={this.receiveFilePathChange.bind(this)}
            />
            { !this.state.isLoading
                ? <button className="w3-btn w3-black" onClick={this.receiveLoadClick.bind(this)}>
                    LOAD
                </button>
                : <div style={{ marginTop: '4px' }}>
                    <LinearProgress/>
                </div>
            }
            {this.state.stories != null
                ? this.renderStories()
                : null
            }
        </div>;
    }

    componentDidUpdate() {
        if (this.state.repoDirectory)
            localStorage.setItem(localStorageAppPrefix + '.repoDirectory', this.state.repoDirectory);
    }

    private receiveFilePathChange(event: ChangeEvent<HTMLInputElement>) {
        const repoDirectory = event.target['value'];
        this.setState({repoDirectory: repoDirectory});
    }

    private async receiveLoadClick() {
        this.setState({isLoading: true});
        try {
            const url = Common.apiUrl + '/stories?' +
                'directory=' + encodeURIComponent(this.state.repoDirectory || '') + '&' +
                'lengthLimit=20';
            const response = await fetch(url);
            const stories: StoryEntry[] = await response.json();
            for (let i = 0; i < stories.length; i++)
                stories[i] = Object.assign(new StoryEntry(), stories[i]);
            this.setState({stories});
        } finally {
            this.setState({isLoading: false});
        }
    }

    private renderStories() {
        const storyDays: StoryEntry[][] = Object.values(
            lodash.groupBy(this.state.stories, (story: StoryEntry) => getStartOfDay(story.getTime()))
        );
        return <div>
            { storyDays.map(storyDay => this.renderStoryDay(storyDay))  }
        </div>;
    }

    private renderStoryEntry(entry: StoryEntry) {
        const key = entry.CommitHash + ' ' + entry.SourceFilePath;
        return <li key={key}>
            {entry.Description}
        </li>;
    }

    private renderStoryDay(entries: StoryEntry[]) {
        const key = getStartOfDay(entries[0].getTime()).toUTCString();
        const dayTitle = getStartOfDay(entries[0].getTime()).toLocaleDateString();
        return <div className='w3-panel w3-leftbar' style={{paddingLeft: 0}} key={key}>
            <div style={{ marginLeft: '8px' }}>
                {dayTitle}
            </div>
            <ul>
                { entries.map(entry => this.renderStoryEntry(entry)) }
            </ul>
        </div>;
    }
}