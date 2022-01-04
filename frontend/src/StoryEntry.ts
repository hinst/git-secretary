

export class StoryEntry {
    Time: string = new Date().toISOString();
    CommitHash: string = '';
    ParentHash: string = '';
    Description: string = '';
    SourceFilePath: string = '';

    getTime() {
        return new Date(this.Time);
    }
}