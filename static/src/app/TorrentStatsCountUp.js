/**
 * Created by Ruslan on 08-Feb-17.
 */


import React from 'react'
import CountUp from 'react-countup'
import {Request} from 'superagent'

export default class TorrentStatsCountUp extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            torrentsPrevCount: 0,
            torrentsCount: 0,
            filesPrevCount: 0,
            filesCount: 0
        }
    }

    responseHandler(err, res) {
        if (err === null){
            let newTorrentCount = this.state.torrentsCount;
            if (this.state.torrentsCount !== res.TorrentCount) {
                newTorrentCount = res.body.TorrentCount;
            }
            let newFilesCount = this.state.filesCount
            if (this.state.filesCount !== res.FileCount) {
                newFilesCount = res.body.FileCount;
            }
            this.setState({
                torrentsPrevCount: this.state.torrentsCount,
                torrentsCount: newTorrentCount,
                filesPrevCount: this.state.filesCount,
                filesCount: newFilesCount
            })
        }

    }

    tick() {
        let handler = this.responseHandler.bind(this);
        new Request('get', '/torrents/count/').end(handler);

    }

    componentDidMount() {
        this.tick();
        setInterval(this.tick.bind(this), 5000);
    }

    render() {
        return <span className="CountUp">
            Torrents <CountUp start={this.state.torrentsPrevCount} separator="," useGrouping={true} end={this.state.torrentsCount}
                              duration={2}/> Files <CountUp start={this.state.filesPrevCount}
                                                            end={this.state.filesCount} useGrouping={true}  separator="," duration={2}/>&nbsp;
            </span>
    }
}
