//* TORRENT SEARCH FLOW
/* Pre : fetch title and metadata from notion in currently watching (still need to think reg. the filter)
-> //? Movie
-> search for movie name and year (specific pattern "Name Year Quality")
   -start search
   -stop search after x seconds?
   -get results and sort them by seed/peers ratio? wont work always as sometimes ratio < 1 
        -SORTING
            -order by no of seeds 
      -add torrent 
      -check sequentional optionm, first last parts option
      -set the download dir to plex dir D:\\Media\\Movies
      -Add to current downloading torrents obj with key: tmdb id and value: name of the torrent
      -add torrent to download
      -after x time or upon download init, replace value with the hash of the torrent in the active torrent obj 
*/
//! Duplicate check
/* Prevent duplicate titles getting downloaded
-> maintain a js object with title as key and status as value for all the current torrents present or the moment when the servie is triggered for that title
-> before proceeding with the event flow, check for exist in obj
*/

//! Valid Session check
/* -> if response is forbiddem - invalid then refresh
 */

// ==================================
// imports
const { qBittorent } = require("./qbt")
require("dotenv").config()
