import { Component, ReactNode } from "react";
import { ActivityReportGroup } from "./ActivityReportGroup";
import { ActivityReportGroupView } from "./ActivityReportGroupView";

class Props {
    activityReportGroups: ActivityReportGroup[] = [];
}

class State {
}

export class ActivityReportGroupsView extends Component<Props, State> {
    override render(): ReactNode {
        return <div>
            { this.props.activityReportGroups.map(
                (item, index) => this.renderGroup(item, index)
            ) }
        </div>;
    }

    private renderGroup(item: ActivityReportGroup, index: number) {
        return <div key={item.time} className="w3-panel" style={{paddingLeft: 0, marginTop: 4, marginBottom: 4}}>
            <ActivityReportGroupView activityReportGroup={item} isExpanded={index === 0} />
        </div>;
    }
}