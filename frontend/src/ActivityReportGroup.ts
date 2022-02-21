import { ActivityReport } from "./ActivityReport";

export class ActivityReportGroup {
    time: string = new Date().toISOString();
    period: number = 0; // nanoseconds
    activity: ActivityReport = new ActivityReport()
    authors: { [authorName: string]: ActivityReport } = {};
}