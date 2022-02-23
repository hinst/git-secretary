import { Component, CSSProperties, ReactNode } from 'react';
import { ActivityReportGroup } from './ActivityReportGroup';
import ChevronRightIcon from '@mui/icons-material/ChevronRight';
import ExpandMoreIcon from '@mui/icons-material/ExpandMore';
import { ActivityReport } from './ActivityReport';
import ArticleIcon from '@mui/icons-material/Article';

class Props {
    constructor (
        public activityReportGroup: ActivityReportGroup,
        public isExpanded: boolean,
    ) {
    }
}

class State {
    constructor(
        public isExpanded: boolean = false,
        public isAuthorsExpanded: boolean = false,
    ) {
    }
}

export class ActivityReportGroupView extends Component<Props, State> {
    private static readonly POINTS_BADGE_STYLE: CSSProperties = {
        backgroundColor: '#3b0b1d',
        borderRadius: 4,
        paddingTop: 2,
        paddingRight: 6,
        paddingBottom: 4,
        paddingLeft: 6,
    };

    private static readonly AUTHORS_BUTTON_STYLE: CSSProperties = {
        paddingTop: 2,
        paddingBottom: 2,
        paddingLeft: 9,
        paddingRight: 9,
        height: 30
    };

    constructor(props: Props) {
        super(props);
        this.state = new State(props.isExpanded, props.isExpanded);
    }

    render(): ReactNode {
        return <div>
            <button onClick={this.toggleExpanded.bind(this)}
                className="w3-btn w3-black"
            >
                { this.state.isExpanded
                    ? <ExpandMoreIcon style={{ verticalAlign: 'middle' }}/>
                    : <ChevronRightIcon style={{ verticalAlign: 'middle' }}/>
                }
            </button> &nbsp;
            <div style={{ display: 'inline-block', verticalAlign: -2 }}>
                { new Date(this.props.activityReportGroup.time).toLocaleDateString() }.&nbsp;
                <span style={ActivityReportGroupView.POINTS_BADGE_STYLE}>
                    <b>∑</b>&nbsp;
                    <b>{ this.props.activityReportGroup.activity.points}</b>
                </span>
            </div>
            { this.state.isExpanded
                ? this.renderExpandedInfo()
                : undefined }
        </div>;
    }

    private toggleExpanded() {
        this.setState({ isExpanded: !this.state.isExpanded });
    }

    private toggleAuthorsExpanded() {
        this.setState({ isAuthorsExpanded: !this.state.isAuthorsExpanded });
    }

    private renderExpandedInfo() {
        const authorNames = Object.keys(this.props.activityReportGroup.authors).sort();
        return <div style={{marginLeft: 64}}>
            {this.renderActivityReport(this.props.activityReportGroup.activity)}
            <div style={{marginTop: 6}}>
                <button onClick={this.toggleAuthorsExpanded.bind(this)}
                    className="w3-btn w3-black"
                    style={ActivityReportGroupView.AUTHORS_BUTTON_STYLE}
                >
                    { this.state.isAuthorsExpanded
                        ? <ExpandMoreIcon style={{ verticalAlign: 'middle' }}/>
                        : <ChevronRightIcon style={{ verticalAlign: 'middle' }}/>
                    }
                </button> &nbsp;
                { authorNames.length ? <span>Authors: [{authorNames.join(', ')}] </span> : undefined}
            </div>
        </div>;
    }

    private renderActivityReport(activityReport: ActivityReport) {
        return <span>
            <ul style={{listStyle: 'none', paddingLeft: 0, marginTop: 4, marginBottom: 4}}>
                <li>{'{'}&nbsp;
                    <ArticleIcon style={{verticalAlign: 'bottom', scale: '0.8'}}/>
                    commits: {activityReport.changesetCount},
                    + insertions: { activityReport.insertionCount },
                    - deletions: { activityReport.deletionCount }
                &nbsp;{'}'}</li>
            </ul>
        </span>;
    }
}