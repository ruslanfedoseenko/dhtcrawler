
import getMuiTheme from 'material-ui/styles/getMuiTheme';
import {white, grey200} from 'material-ui/styles/colors';

const theme = getMuiTheme({
    appBarRightContainerStyle:{
        marginTop: 0
    },
    searchButtonStyle:{
        marginTop: 5,
    },
    searchBoxHintStyle: {
        color: grey200
    },
    searchBoxUnderLineStyle: {
        borderColor: white,
        color: white,
    },
    searchBoxInputStyle: {
        color: white,

    },
    searchBoxStyle: {
        marginTop: 5,
    },
    palette:{
        primary1Color:'#78909C'
    }
});

export default theme;

