/**
 * Created by Ruslan on 08-Feb-17.
 */

import React from 'react';
import {Request} from 'superagent'
import {List, ListItem, makeSelectable, Drawer} from 'material-ui';

const SelectableList = makeSelectable(List);
const listItemStyle = {
    backgroundColor: 'transparent'
};

export default class SideBarMenu extends React.Component {

    constructor(props) {
        super(props);
        this.state = {
            opened: false,
            items: [],
            selectedIndex: '',
        }
    }

    onToggleState() {
        this.setState({
            opened: !this.state.opened
        })
    }

    onCategoriesReceived(err, res) {
        if (err === null) {
            let categories = res.body;
            categories.unshift({
                Id: 112,
                icon: "globe",
                value: "All",
                $count: 0
            });
            categories.push({
                Id: 212,
                icon: "bar-chart",
                value: "Stats"
            });
            categories.push({
                Id: 312,
                icon: "info",
                value: "About"
            });
            this.setState({

                opened: this.state.opened,
                items: categories
            });
        }
    }

    componentDidMount() {
        let handler = this.onCategoriesReceived.bind(this);
        new Request('get', '/categories/').end(handler)

    }

    onChangeListInternal(i, a) {
        this.setState({
            opened: !this.state.opened,
            selectedIndex: a
        });
        this.props.onChangeList(i, a);
    }

    render() {
        return <Drawer containerStyle={ {'top': '75px'} } open={this.state.opened}>
            <SelectableList value={this.state.selectedIndex} onChange={this.onChangeListInternal.bind(this)}>
                {
                    this.state.items.map((item) => {
                        if (item.data) {
                            return <ListItem leftIcon={<i className={'fa fa-' + item.icon}/>}
                                             primaryText={item.value}
                                             key={item.Id}
                                             value={item.Id}
                                             style={listItemStyle}
                                             nestedItems={
                                                 item.data.map((nestedItem) => {
                                                     return <ListItem
                                                         value={item.Id + '.' + nestedItem.Id}
                                                         key={nestedItem.Id}
                                                         style={listItemStyle}
                                                         leftIcon={<i className={'fa fa-' + nestedItem.icon}/>}
                                                         primaryText={nestedItem.value}/>
                                                 })
                                             }/>
                        } else {
                            return <ListItem leftIcon={<i className={'fa fa-' + item.icon}/>}
                                             key={item.Id}
                                             value={item.Id}
                                             style={listItemStyle}
                                             primaryText={item.value}/>
                        }


                    })

                }
            </SelectableList>
        </Drawer>

    }
}

SideBarMenu.propTypes = {
    location: React.PropTypes.object.isRequired,
    onChangeList: React.PropTypes.func.isRequired,
};
