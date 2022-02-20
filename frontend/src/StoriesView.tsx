import { Component, ReactNode } from 'react';
import { StoryEntryChangeset, StoryEntryFileChange } from './StoryEntry';
import WarningIcon from '@mui/icons-material/Warning';
import lodash from 'lodash';
import { getStartOfDay } from './dateTime';

class Props {
    entries: StoryEntryChangeset[] = [];
}

class State {
}

export class StoriesView extends Component<Props, State> {
    private static readonly DAY_LIMIT = 100;

    override render(): ReactNode {
        const storyDays: StoryEntryChangeset[][] = Object.values(
            lodash.groupBy(this.props.entries, (story: StoryEntryChangeset) => getStartOfDay(story.getTime()))
        );
        const isDayLimitExceeded = storyDays.length > StoriesView.DAY_LIMIT;
        const totalStoryDays = storyDays.length;
        if (isDayLimitExceeded) {
            storyDays.splice(StoriesView.DAY_LIMIT);
        }
        return <div>
            { storyDays.map(storyDay => this.renderStoryDay(storyDay)) }
            { isDayLimitExceeded
                ? <div>
                    <WarningIcon style={{ verticalAlign: 'middle' }} />&nbsp;
                    A limited number of days is displayed:
                    {StoriesView.DAY_LIMIT} of {totalStoryDays}
                </div>
                : undefined
            }
        </div>;
    }

    private renderStoryDay(entries: StoryEntryChangeset[]): ReactNode {
        const key = getStartOfDay(entries[0].getTime()).toUTCString();
        const dayTitle = getStartOfDay(entries[0].getTime()).toLocaleDateString();
        return <div className='w3-panel w3-leftbar' style={{paddingLeft: 0}} key={key}>
            <div style={{ marginLeft: '8px' }}>
                {dayTitle}
            </div>
            <ul>
                { entries.map(entry => this.renderStoryChangeset(entry)) }
            </ul>
        </div>;
    }

    private renderStoryChangeset(changeset: StoryEntryChangeset): ReactNode {
        const key = changeset.commitHash;
        return <li key={key}>
            <ul>
                { (changeset.fileEntries || []).map((fileEntry, index) => this.renderFileEntry(index, fileEntry)) }
            </ul>
        </li>;
    }

    private renderFileEntry(index: number, fileEntry: StoryEntryFileChange): ReactNode {
        return <li key={index}>
            {fileEntry.description}
        </li>;
    }
}