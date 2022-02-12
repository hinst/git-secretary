export class StoryEntryChangeset {
    time: string = new Date().toISOString();
    commitHash: string = '';
    fileEntries: StoryEntryFileChange[] = [];

    getTime() {
        return new Date(this.time);
    }
}

export class StoryEntryFileChange {
    sourceFilePath: string = '';
    description: string = '';
}
