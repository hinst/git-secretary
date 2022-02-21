import { ActivityReportGroup } from './ActivityReportGroup';

export class WebTask {
    total: number = 0;
    done: number = 0;
    error?: string;
    activityReportGroups?: ActivityReportGroup[];
}