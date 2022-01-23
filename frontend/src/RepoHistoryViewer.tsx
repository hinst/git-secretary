import { Component } from 'react';
import { Common } from './Common';
import { StoryEntry } from './StoryEntry';
import lodash from 'lodash';
import { getStartOfDay } from './dateTime';
import FolderOpenIcon from '@mui/icons-material/FolderOpen';
import RefreshIcon from '@mui/icons-material/Refresh';
import { replaceAll } from './string';
import { Link, Navigate } from 'react-router-dom';

class Props {
    directory?: string;
}

class State {
    stories: StoryEntry[] = [];
    isLoading: boolean = false;
    goTo?: string;
}

export class RepoHistoryViewer extends Component<Props, State> {
    constructor(props: Props) {
        super(props);
        const state = new State();
        if (!props.directory)
            state.goTo = '/open-repository';
        this.state = state;
    }

    override render() {
        if (this.state.goTo)
            setTimeout(() => this.setState({ goTo: undefined }));
        return <div>
            { this.state.goTo ? <Navigate to={this.state.goTo} /> : undefined }
            <div className="w3-bar w3-dark-grey" style={{marginBottom: 4, position: 'sticky', top: 0}}>
                <Link
                    to={'/open-repository'}
                    title="Open Git repository"
                    className="w3-bar-item w3-btn w3-black"
                >
                    <FolderOpenIcon/>
                </Link>
                <button
                    onClick={() => this.receiveLoadClick()}
                    className="w3-btn w3-black w3-bar-item"
                    style={{marginLeft: 4}}
                >
                    <RefreshIcon className={ this.state.isLoading ? "rotating" : undefined }/>
                </button>
                <div className="w3-bar-item" style={{fontSize: 16}}>
                    {this.repositoryName}
                </div>
            </div>
            <div>
                {this.state.stories != null
                    ? this.renderStories()
                    : null
                }
            </div>
        </div>;
    }

    override componentDidMount() {
        this.receiveLoadClick();
    }

    private async receiveLoadClick() {
        this.setState({isLoading: true});
        try {
            const url = Common.apiUrl + '/stories?' +
                'directory=' + encodeURIComponent(this.props.directory || '');
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
        const path = replaceAll('\\', '/', this.props.directory || '');
        const parts = path.split('/');
        const lastPart = parts.length ? parts[parts.length - 1] : '';
        return lastPart;
    }
}