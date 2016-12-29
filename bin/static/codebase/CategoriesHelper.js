/**
 * Created by devrus on 23.08.16.
 */



function CategoriesHelper(){

}

/**
 * @return {string}
 */
CategoriesHelper.GetCategoryTemplate = function(name, withText){
    switch (name) {
        case "Audio": {
            return '<i style="margin-right: 6px;" class="fa fa-music" aria-hidden="true"></i>' +(withText ? name : '');
        }
        case "Video": {
            return '<i style="margin-right: 6px;" class="fa fa-film" aria-hidden="true"></i>' + (withText ? name : '');
        }
        case "TV": {
            return '<i style="margin-right: 6px;" class="fa fa-film" aria-hidden="true"></i>' + (withText ? 'TV' : '');
        }

        case "Text": {
            return '<i style="margin-right: 6px;" class="fa fa-book" aria-hidden="true"></i>' + (withText ? name : '');
        }
        case "Book": {
            return '<i style="margin-right: 6px;" class="fa fa-book" aria-hidden="true"></i>' + (withText ? 'Book' : '');
        }
        case "Magazine": {
            return '<i style="margin-right: 6px;" class="fa fa-book" aria-hidden="true"></i>' + (withText ? 'Magazine' : '');
        }
        case "Soft": {
            return '<i style="margin-right: 6px;" class="fa fa-desktop" aria-hidden="true"></i>' + (withText ? name : '');
        }
        case "Game": {
            return '<i style="margin-right: 6px;" class="fa fa-desktop" aria-hidden="true"></i>' + (withText ? 'Games': '');
        }
        case "Picture": {
            return '<i style="margin-right: 6px;" class="fa fa-picture-o" aria-hidden="true"></i>' + (withText ? name : '');
        }
        case "Other": {
            return '<i style="margin-right: 6px;" class="fa fa-file" aria-hidden="true"></i>' + (withText ? name : '');

        }
        default: {
            return name;
        }
    }
};

