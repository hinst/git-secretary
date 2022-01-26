import { StoryEntry } from "./StoryEntry";

export class WebTask {
    total: number = 0;
    done: number = 0;
    error?: string;
    storyEntries?: StoryEntry[];
}