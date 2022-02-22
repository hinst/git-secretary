import { Component, CSSProperties, ReactNode } from 'react';
import { ActivityReportGroup } from './ActivityReportGroup';
import ChevronRightIcon from '@mui/icons-material/ChevronRight';
import { ActivityReport } from './ActivityReport';
import ArticleIcon from '@mui/icons-material/Article';

class Props {
    constructor (
        public activityReportGroup: ActivityReportGroup,
        public isExpanded: boolean
    ) {
    }
}

class State {
    constructor(public isExpanded: boolean = false) {
    }
}

export class ActivityReportGroupView extends Component<Props, State> {
    constructor(props: Props) {
        super(props);
        this.state = new State(props.isExpanded);
    }

    private static readonly POINTS_BADGE_STYLE: CSSProperties = {
        backgroundColor: '#3b0b1d',
        borderRadius: 4,
        paddingTop: 2,
        paddingRight: 6,
        paddingBottom: 4,
        paddingLeft: 6,
    };

    render(): ReactNode {
        return <div>
            <button className="w3-btn w3-black" onClick={this.toggleExpanded.bind(this)}>
                <ChevronRightIcon style={{ verticalAlign: 'middle' }}/>
            </button> &nbsp;
            <div style={{ display: 'inline-block', verticalAlign: -2 }}>
                { new Date(this.props.activityReportGroup.time).toLocaleDateString() }.&nbsp;
                <span style={ActivityReportGroupView.POINTS_BADGE_STYLE}>
                    <b>∑</b>&nbsp;
                    <b>{ this.props.activityReportGroup.activity.points}</b>
                </span>
            </div>
            { this.state.isExpanded ? this.renderActivityReport(this.props.activityReportGroup.activity) : undefined }
        </div>;
    }

    private toggleExpanded() {
        this.setState({ isExpanded: !this.state.isExpanded });
    }

    private renderActivityReport(activityReport: ActivityReport) {
        return <div style={{marginLeft: 64}}>
            <ul style={{listStyle: 'none', paddingLeft: 0, marginTop: 4, marginBottom: 4}}>
                <li>{'{'}&nbsp;
                    <ArticleIcon style={{verticalAlign: 'bottom', scale: '0.8'}}/> commits: {activityReport.changesetCount},
                    + insertions: { activityReport.insertionCount },
                    - deletions: { activityReport.deletionCount }
                &nbsp;{'}'}</li>
            </ul>
        </div>;
    }
}