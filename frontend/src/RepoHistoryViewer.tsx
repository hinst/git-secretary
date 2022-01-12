import React, { ChangeEvent, Component } from 'react';
import { Common } from './Common';
import { localStorageAppPrefix } from './localStorage';
import { StoryEntry } from './StoryEntry';
import lodash from 'lodash';
import { getStartOfDay } from './dateTime';
import { LinearProgress } from '@material-ui/core';
import FolderOpenIcon from '@mui/icons-material/FolderOpen';
import { replaceAll } from './string';
import { Margin } from '@mui/icons-material';
import { Link } from 'react-router-dom';

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

    override render() {
        return <div>
            <div style={{padding: '0px'}}>
                <div className="w3-bar">
                    <Link to={Common.baseUrl + '/open-repository'} className="w3-bar-item w3-btn w3-black">
                        <FolderOpenIcon/>
                    </Link>
                    <div className="w3-bar-item" style={{fontSize: '17px'}}>
                        {this.repositoryName}
                    </div>
                </div>
                <div>
                    {this.state.stories != null
                        ? this.renderStories()
                        : null
                    }
                </div>
            </div>
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

    private get repositoryName(): string {
        const path = replaceAll('\\', '/', this.state.repoDirectory || '');
        const parts = path.split('/');
        const lastPart = parts.length ? parts[parts.length - 1] : '';
        return lastPart;
    }
}