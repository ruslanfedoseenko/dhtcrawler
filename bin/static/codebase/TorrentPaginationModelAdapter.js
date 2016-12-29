/**
 * Created by devrus on 21.08.16.
 */
webix.DataDriver.torrentsAdapter = webix.extend({
    getRecords: function (data) {
        console.log(data)
        return data.Torrents;
    },
    getInfo: function (data) {
        var info = webix.DataDriver.json.getInfo(data)
        info.size = data.ItemsPerPage * data.PageCount;
        info.from = (data.Page - 1) * data.ItemsPerPage;
        info._page = data.Page;
        info._page_size = data.ItemsPerPage;
        return info;
    }
}, webix.DataDriver.json);

webix.DataDriver.torrentFilesAdapter = webix.extend({
    child: function (data) {

        return data.Children;
    }

}, webix.DataDriver.json);


webix.DataDriver.categoriesAdapter = webix.extend({
    child: function (data) {

        return data.Children;
    },
    getDetails: function(obj){
        obj.value = obj.Name;
        return obj;
    }
}, webix.DataDriver.json);


