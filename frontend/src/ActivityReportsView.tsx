import { Component, ReactNode } from "react";
import { ActivityReportGroup } from "./ActivityReportGroup";

class Props {
    activityReportGroups: ActivityReportGroup[] = [];
}

class State {
}

export class ActivityReportsView extends Component<Props, State> {
    override render(): ReactNode {
        return <div>
            ActivityReportsView
        </div>;
    }
}