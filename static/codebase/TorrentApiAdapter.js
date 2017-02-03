/**
 * Created by devrus on 21.08.16.
 */
webix.proxy.torrentApi = webix.extend({
    $proxy: true,
    params: {
        page: 1,
        searchTerm: '',
        groupId: 0,
    },
    _buildSource: function () {
        if (this.params.searchTerm) {
            return this.source + 'search/' + this.params.searchTerm + '/page/' + this.params.page;
        }
        if (this.params.groupId > 0)
        {
            return this.source + 'group/' + this.params.groupId +'/page/' + this.params.page;
        }
        return this.source + 'page/' + this.params.page;
    },
    load: function (view, callback, options) {


        if (options) {
            this.params.page = Math.floor(options.start / 30) + 1;
        }

        webix.ajax(this._buildSource(), callback, view);

    }
}, webix.proxy);
