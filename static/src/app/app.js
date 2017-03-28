/**
 * Created by Ruslan on 07-Feb-17.
 */

import React from 'react';
import ReactDOM from 'react-dom';
import AppBar from 'material-ui/AppBar';
import AutoComplete from 'material-ui/AutoComplete';
import IconButton from 'material-ui/IconButton';
import FontIcon from 'material-ui/FontIcon';

import TorrentStatsCountUp from './TorrentStatsCountUp.js';
import {Request} from 'superagent'
import injectTapEventPlugin from 'react-tap-event-plugin';
import TorrentsList from './TorrentsList.js'
import MuiThemeProvider from 'material-ui/styles/MuiThemeProvider';
import {white} from 'material-ui/styles/colors';
import theme from './theme.js'
import {Page} from 'react-layout-components'
injectTapEventPlugin();

class App extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            autoCompleteTerms: [],
            drawerOpen: false,
            menuIconName: 'menu',
            autoCompleteOpen: false
        };
    }

    onSearch(input) {
        new Request('get', '/search/suggest/' + this.getSearchWord(input)).end(this.setTerms.bind(this));
        if (this.props.searchInputTimeId !== -1) {
            clearTimeout(this.props.searchInputTimeId);
        }
        this.props.searchInputTimeId = setTimeout(this.doSearch.bind(this), 1200, input)
    }

    doSearch(what) {
        this.TorrentList.search(what);
        clearTimeout(this.props.searchInputTimeId);
        this.props.searchInputTimeId = -1;
    }

    UpdateDimensions(){
        this.setState({width: window.innerWidth, height: window.innerHeight});
    }

    onCategoryChanged(_, cat) {
        let parts = cat.toString().split('.');
        this.TorrentList.setGroup(parts[parts.length - 1]);
    }
    componentDidMount(){
        window.addEventListener('resize', this.UpdateDimensions.bind(this));
    }
    render() {

        return <MuiThemeProvider muiTheme={theme}>
            <Page column>

                    <AppBar title="Torrent Search Engine"
                            iconStyleRight={theme.appBarRightContainerStyle}
                            iconElementRight={<div>
                                <TorrentStatsCountUp/>
                                <AutoComplete hintText="Search..." inputStyle={theme.searchBoxInputStyle}
                                              style={theme.searchBoxStyle}
                                              dataSource={this.state.autoCompleteTerms}
                                              underlineStyle={theme.searchBoxUnderLineStyle}
                                              underlineFocusStyle={theme.searchBoxUnderLineStyle}
                                              open={this.state.autoCompleteOpen}
                                              hintStyle={theme.searchBoxHintStyle}
                                              onUpdateInput={this.onSearch.bind(this)}/>
                                <IconButton onTouchTap={this.onSearch.bind(this)} style={theme.searchButtonStyle}>
                                    <FontIcon
                                        className="material-icons" color={white}>search</FontIcon>
                                </IconButton>
                            </div>}>

                    </AppBar>

                    <TorrentsList ref={(tList) => this.TorrentList = tList} location={{}} />

            </Page>
        </MuiThemeProvider>
    }

    setTerms(err, res) {
        if (err === null) {
            this.setState({
                autoCompleteTerms: res.body.data,
                autoCompleteOpen: true,
            });
        }
    }

    getSearchWord(input) {
        let words = input.split(' ').filter(w => w.length !== 0);
        if (this.props.previousSearchWords.length === 0) {

            return words[words.length - 1];
        }


    }
}
App.defaultProps = {
    searchInputTimeId: 0,
    previousSearchWords: []
};
ReactDOM.render(
    <App width={window.innerWidth}
         height={window.innerHeight}/>,
    document.getElementById('app')
);


// WEBPACK FOOTER //
// ./src/app/app.jsx