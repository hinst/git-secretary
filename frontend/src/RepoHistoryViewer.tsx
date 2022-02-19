import { Component, ReactNode } from 'react';
import { Common } from './Common';
import { StoryEntryChangeset } from './StoryEntry';
import FolderOpenIcon from '@mui/icons-material/FolderOpen';
import ErrorIcon from '@mui/icons-material/Error';
import RefreshIcon from '@mui/icons-material/Refresh';
import DoNotDisturbIcon from '@mui/icons-material/DoNotDisturb';
import { replaceAll } from './string';
import { Link, Navigate } from 'react-router-dom';
import { WebTask } from './WebTask';
import { LinearProgress } from '@mui/material';
import { StoriesView } from './StoriesView';

class Props {
    directory?: string;
}

class State {
    stories: StoryEntryChangeset[] = [];
    error?: string;
    taskId?: number;
    isLoading: boolean = false;
    loadingTotal?: number;
    loadingDone?: number;
    goTo?: string;
}

export class RepoHistoryViewer extends Component<Props, State> {
    private static readonly DAY_LIMIT = 100;
    private loadingTaskTimer?: number;

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
                    <RefreshIcon className={ this.state.isLoading ? 'rotating' : undefined }/>
                </button>
                <div className="w3-bar-item" style={{fontSize: 16}}>
                    {this.repositoryName}
                </div>
            </div>
            <div>
                {this.state.error
                    ? this.renderError()
                    : undefined
                }
            </div>
            <div>
                { this.state.isLoading
                    ? this.renderLoading()
                    : undefined
                }
            </div>
            <div>
                {this.state.stories != null
                    ? this.renderStories()
                    : undefined
                }
            </div>
        </div>;
    }

    override componentDidMount() {
        this.receiveLoadClick();
    }

    private async receiveLoadClick() {
        this.setState({ isLoading: true, error: undefined });
        const url = Common.apiUrl + '/stories' +
            '?directory=' + encodeURIComponent(this.props.directory || '') +
            '&lengthLimit=' + encodeURIComponent(RepoHistoryViewer.DAY_LIMIT * 10) +
            '&timeZone=' + encodeURIComponent(Intl.DateTimeFormat().resolvedOptions().timeZone);
        try {
            const response = await fetch(url);
            if (response.ok) {
                const taskId = parseInt(await response.text());
                this.setState({ isLoading: true, taskId: taskId, error: undefined });
                this.loadingTaskTimer = window.setTimeout(() => this.checkStoriesLoaded(), 500);
            } else {
                const errorText = await response.text();
                this.setState({ isLoading: false, taskId: undefined, error: errorText });
            }
        } catch (e) {
            this.setState({ isLoading: false, error: (e as any).message });
        }
    }

    private stopStoriesLoading() {
        if (this.loadingTaskTimer != null)
            window.clearInterval(this.loadingTaskTimer);
        this.loadingTaskTimer = undefined;
        this.setState({ taskId: undefined, isLoading: false});
    }

    private async checkStoriesLoaded() {
        this.loadingTaskTimer = undefined;
        if (!this.state.taskId)
            return this.stopStoriesLoading();
        const url = Common.apiUrl + '/task?id=' + encodeURIComponent(this.state.taskId);
        const response = await fetch(url);
        if (response.ok) {
            const task: WebTask = await response.json();
            if (task.error?.length) {
                this.setState({ error: task.error, stories: [] });
                this.stopStoriesLoading();
            } else if (task.storyEntries) {
                const stories: StoryEntryChangeset[] = task.storyEntries;
                for (let i = 0; i < stories.length; i++)
                    stories[i] = Object.assign(new StoryEntryChangeset(), stories[i]);
                this.setState({ error: undefined, stories: stories });
                this.stopStoriesLoading();
            } else {
                this.setState({ loadingTotal: task.total, loadingDone: task.done });
                this.loadingTaskTimer = window.setTimeout(() => this.checkStoriesLoaded(), 500);
            }
        } else {
            this.setState({ error: await response.text(), stories: [] });
            this.stopStoriesLoading();
        }
    }

    private renderError() {
        return <span>
            <ErrorIcon style={{ verticalAlign: 'middle' }}/> { this.state.error }
        </span>;
    }

    private renderLoading() {
        var progressRatio = this.state.loadingTotal != null && this.state.loadingDone != null
            ? Math.min(1, this.state.loadingDone / Math.max(1, this.state.loadingTotal)) * 100
            : null;
        return <span>
            { progressRatio != null && (this.state.loadingTotal || 0) > 100
                ? <span> Loading entries: {this.state.loadingDone} of {this.state.loadingTotal} </span>
                : undefined
            }
            { progressRatio != null
                ? <LinearProgress variant='determinate' value={progressRatio}/>
                : <LinearProgress variant='indeterminate'/>
            }
        </span>;
    }

    private renderStories(): ReactNode {
        return this.state.stories.length
            ? <StoriesView entries={this.state.stories}/>
            : <div className='w3-panel'>
                <span>
                    <DoNotDisturbIcon style={{ verticalAlign: 'middle' }}/>&nbsp;
                    There are no entries to show. Likely causes are:
                </span>
                <ul style={{marginTop: 4}}>
                    <li>The source repository has zero commits</li>
                    <li>There is an error in the plug-in</li>
                    <li>There is an error in git-secretary</li>
                    <li>An incomplete alpha version of plug-in is used</li>
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