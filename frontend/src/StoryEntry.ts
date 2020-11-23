

export class StoryEntry {
    Time: string;
    CommitHash: string;
    ParentHash: string;
    Description: string;
    getTime() {
        return new Date(this.Time);
    }
}