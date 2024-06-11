# CICK Playlister Usage

The CICK Playlister is a tool to help create playlists following a CICK radio show. It fetches track information from online streaming platforms to automate some of the work of logging which tracks were played.

It currently supports the following streaming platforms:
- Spotify

> [!TIP]
> Additional streaming platforms may be supported at a later date depending on demand, enthusiasm, and financial constraints.

## How to Use

The tool is activated by a bookmarklet. A bookmarklet is a special type of web browser bookmark that doesn't take you to a new web page. Look for a bookmark called `CICK Playlister` and click it. If you don't see the bookmark ensure that the bookmarks bar is visible via web browser settings. If the bookmark is still not visible contact someone who knows more about the tool.

> [!NOTE]
> The tool will only work properly if you are viewing the smithersradio.com Create Program Playlist page.

Clicking the bookmark will bring up a window with a single text input field. Paste the URL of a playlist, album, or track into this field and hit `return` or click the `Fill` button next to it.

The tool will fetch track information and fill the playlist's `ARTIST`, `TITLE`, `ALBUM`, and `NEW` fields.

> [!IMPORTANT]  
> The tool cannot fill the CAN (Canadian Content) field - this is your responsibility.

## Limitations

The tool will not duplicate track information. If you attempt to add the same track twice no change will occur.

The tool will only fill empty rows in the playlist table. Once all rows are filled you must add rows to the table for it to make further changes.

## Issues

If you encounter any issues with the tool please log an issue [here](https://github.com/captaincoordinates/cick-playlister/issues/). You will need a GitHub account to create an issue. If you do not want to create a GitHub account please report the problem to someone who knows more about the tool.
