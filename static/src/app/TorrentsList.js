/**
 * Created by Ruslan on 09-Feb-17.
 */

import React from 'react';
import {List, ListItem, Chip} from 'material-ui';
import {Request} from 'superagent'
import Pagination from 'material-ui-pagination';
import Loader from 'react-loader-advanced';
import Paper from 'material-ui/Paper'
import CircularProgress from 'material-ui/CircularProgress';
import Scroll  from 'react-scroll';
import Formatter from "./Formaters";
import {Box, Center} from 'react-layout-components'

let scroll = Scroll.animateScroll;
export default class TorrentsList extends React.Component {


    constructor(props) {
        super(props);

        this.state = {
            isLoading: false,
            Page: 1,
            PageCount: 0,
            ItemsCount: 0,
            ItemsPerPage: 30,
            Torrents: [],
        }


    }

    search(term) {
        this.props.isSearch = term;
        this.fetchTorrents(1);
    }

    setGroup(group) {
        if (group !== 112) {
            this.props.groupId = group

        } else {
            this.props.groupId = 0;
        }
        this.props.isSearch = false;
        this.fetchTorrents(1);
    }

    fetchTorrents(pageNumber) {
        this.setState({
            isLoading: true
        });
        new Request('get', this.buildUrl(pageNumber)).end(this.onGotTorrents.bind(this));
    }

    componentDidMount() {
        this.fetchTorrents(this.state.Page);
    }

    componentDidUpdate() {
        scroll.scrollToTop();
    }

    render() {

        let state = this.state;
        return <Box column>
            <Loader show={state.isLoading} backgroundStyle={{backgroundColor: 'rgba(255,255,255, 0.498039)'}}
                    message={<CircularProgress size={120} thickness={7}/>}>

                    <List style={{overflow: 'auto', height: window.innerHeight - 64 - 55}} >
                        {
                            this.state.Torrents.map(function (torrent) {
                                let torrentSize = 0;
                                for (let i = 0; i < torrent.FilesTree.length; i++) {
                                    torrentSize += torrent.FilesTree[i].Size;
                                }
                                return <ListItem
                                    style={{cursor: 'normal'}}
                                    key={torrent.InfoHash}
                                    primaryText={torrent.Name}
                                    secondaryTextLines={2}
                                    secondaryText={

                                        <div>
                                            <Chip style={{display: 'inline-block', margin: 2}}>
                                                Size: {Formatter.formatBytes(torrentSize, 2)}
                                            </Chip>
                                            <Chip style={{display: 'inline-block', margin: 2}}>
                                                <a href={'magnet:?xt=urn:btih:' + torrent.Infohash + '&tr=udp://tracker.coppersurfer.tk:6969/announce&tr=udp://open.demonii.com:1337/announce&tr=udp://tracker.openbittorrent.com:80&tr=http://tracker.opentrackr.org:1337/announce&tr=http://explodie.org:6969/announce' }>
                                                    <i className="fa fa-download"/>
                                                </a>
                                            </Chip>
                                            {(torrent.Tags || []).map(function (tag) {
                                                return <Chip
                                                    style={{display: 'inline-block', margin: 2}}>{tag.Tag}</Chip>
                                            })}

                                        </div>

                                    }>

                                </ListItem>
                            })
                        }
                    </List>

            </Loader>
            <Center>
                <Paper>
                    <Pagination total={ this.state.PageCount }
                                current={ this.state.Page }
                                display={ 10 }

                                onChange={this.setPageInternal.bind(this)}/>
                </Paper>
            </Center>
        </Box>
    }

    calculateSize(filesTree) {
        let size = 0;
        for (let i = 0; i < filesTree.length; i++) {
            size += filesTree[i].Size;
            if (filesTree[i].Children) {
                size += this.calculateSize(filesTree[i].Children);
            }
        }
        return size;
    }

    calculateSeeds(obj) {
        if (obj.TrackersInfo) {
            let result = 0;
            for (let i = 0; i < obj.TrackersInfo.length; i++) {
                result += obj.TrackersInfo[i].Seeds;
            }
            return result;
        }
        return '-';
    }

    calculateLeechers(obj) {
        if (obj.TrackersInfo) {
            let result = 0;
            for (let i = 0; i < obj.TrackersInfo.length; i++) {
                result += obj.TrackersInfo[i].Leaches;
            }
            return result;
        }
        return '-';
    }

    setPageInternal(pageNumber) {
        this.fetchTorrents(pageNumber);
    }

    buildUrl(pageNumber) {
        let url = '/torrents';
        if (this.props.isSearch) {
            url += '/search/' + this.props.isSearch;

        } else {
            if (this.props.groupId > 0) {
                url += '/group/' + this.props.groupId;

            }
        }

        if (pageNumber > 1) {
            url += '/page/' + pageNumber;
        }

        return url;
    }

    onGotTorrents(err, res) {
        if (err === null) {
            res.body.isLoading = false;
            this.setState(res.body);
        }
    }
}

TorrentsList.defaultProps = {
    isSearch: false,
    groupId: -1
};