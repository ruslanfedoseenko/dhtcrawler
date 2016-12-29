/**
 * Created by devrus on 21.08.16.
 */


function PaginationHelper()
{

}

PaginationHelper.PageIdToPageNumber = function(pageId, paginateConfig) {
    var pageNum = parseInt(pageId)
    if (isNaN(pageNum))
    {
        switch (pageId){
            case "next": {
                pageNum = paginateConfig.page + 2;
                break;
            }
            case 'prev': {
                pageNum = paginateConfig.page;
                break;
            }
            case 'first': {
                pageNum = 1;
                break;
            }
            case 'last': {
                pageNum = paginateConfig.count / paginateConfig.size;
            }

        }
    }
    else {
        pageNum++;
    }
    return pageNum;
}