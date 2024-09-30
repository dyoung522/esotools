# Notes

## API Reference

### ESOUI.com (via MMOUI.com)

- Get full AddOn listing: `https://api.mmoui.com/v4/game/ESO/filelist.json` (ModListItem)
- Get AddOn Details:  `https://api.mmoui.com/v4/game/ESO/filedetails/<modIds>.json` (ModAItem) # modIds can be one or more IDs separated by commas

```ts
type ModAddon = {
  path: string;
  addOnVersion: string;
  apiVersion: string;
  requiredDependencies?: string[];
  optionalDependencies?: string[];
};

export type ModListItem = {
  id: number;
  categoryId: number;
  version: string;
  lastUpdate: number;
  title: string;
  author: string;
  fileInfoUri: string;
  downloads: number;
  downloadsMonthly: number;
  favorites: number;
  addons: ModAddon[];
  checksum: string;
};

export type ModItem = ModListItem & {
  cacheExpiry: number;
  description: string;
  changeLog: string;
  downloadUri: string;
  images?: Image[];
  fileName: string;

  /**
   * This exists specifically and only for the modlist backup integration
   */
  installDisabled?: boolean;
};

export type Image = {
  thumbUrl: string;
  imageUrl: string;
  description: string;
};
```

## Example list entry

```JSON
  {
    "id": 13,
    "categoryId": 18,
    "version": "3.8.3.0",
    "lastUpdate": 1724243910000,
    "title": "Ravalox' Quest Tracker",
    "author": "calia1120",
    "fileInfoUri": "https:\/\/www.esoui.com\/downloads\/info13-RavaloxQuestTracker.html",
    "downloads": 1373056,
    "downloadsMonthly": 7977,
    "favorites": 1337,
    "gameVersions": [
      "10.1.0"
    ],
    "checksum": "9320e2f52d18662809b1d71276cb3062",
    "addons": [
      {
        "path": "RavaloxsQuestTracker\/libs\/LibQuestInfo",
        "addOnVersion": "3080300",
        "apiVersion": "101043",
        "library": true
      },
      {
        "path": "RavaloxsQuestTracker",
        "addOnVersion": "3080300",
        "apiVersion": "101043",
        "requiredDependencies": [
          "LibAddonMenu-2.0>=37",
          "LibCustomMenu>=722",
          "LibQuestInfo>=3080300"
        ]
      }
    ]
  },
```

## Example AddOn Details entry

`curl -o test_details.json https://api.mmoui.com/v4/game/ESO/filedetails/13.json`

```json
[
  {
    "id": 13,
    "categoryId": 18,
    "version": "3.8.3.0",
    "lastUpdate": 1724243910000,
    "checksum": "9320e2f52d18662809b1d71276cb3062",
    "fileName": "Ravalox'QuestTracker_SA_240821_3.8.3.0.zip",
    "downloadUri": "https:\/\/cdn.esoui.com\/downloads\/getfile.php?id=13&d=1724243910&minion",
    "pendingUpdate": 0,
    "title": "Ravalox' Quest Tracker",
    "author": "calia1120",
    "description": "[SIZE=\"4\"][COLOR=#e69138][B]History and Credits:[\/B][\/COLOR][\/SIZE]\r\n\r\n[SIZE=\"2\"] Former Add-on \/ UI Design: Wykkyd Gaming, aka [URL=\"https:\/\/x.com\/WykkydGaming\"]Wykkyd[\/URL] (up to v2.3.3.1), [URL=\"https:\/\/www.esoui.com\/forums\/member.php?userid=1034\"]Joviex[\/URL] (tweak patch)[\/SIZE]\r\n[SIZE=\"2\"] Current Add-on \/ Codebase: Exodus Code Group, aka Ravalox and Balkoth[\/SIZE]\r\n\r\n[SIZE=\"2\"][COLOR=#e5e5e5][B]First author: [\/B][\/COLOR][URL=\"https:\/\/x.com\/Ravalox\"][COLOR=#FF2222]Ravalox Darkshire[\/COLOR][\/URL][\/SIZE] \r\n[SIZE=\"2\"][COLOR=#e5e5e5][B]Second author: [\/COLOR][\/B][URL=\"https:\/\/www.esoui.com\/downloads\/author-1653.html\"][COLOR=#00ffff]Calia1120[\/COLOR][\/URL][\/SIZE]\r\n[SIZE=\"2\"][COLOR=#e5e5e5][B]Inactive maintenance team:[\/COLOR][\/B] [URL=\"https:\/\/www.esoui.com\/downloads\/author-1653.html\"][COLOR=#00ffff]Calia1120[\/COLOR][\/URL],   [COLOR=#00ffff]Demiknight[\/COLOR], [COLOR=#00ffff]Lakashi[\/COLOR][\/SIZE]\r\n[SIZE=\"2\"][COLOR=#e5e5e5][B]Temporary maintenance staff: [\/COLOR][\/B][URL=\"https:\/\/www.esoui.com\/downloads\/author-52353.html\"][COLOR=\"PaleGreen\"]Calamath[\/COLOR][\/URL][\/SIZE]\r\n\r\n*****************************************************************************\r\n[SIZE=\"3\"][COLOR=\"Red\"]The add-on name changed as of version 3.8.3. Please refer to the change log.[\/COLOR][\/SIZE]\r\n*****************************************************************************\r\n\r\n[size=4][color=#e69138][b]Dependencies:[\/b][\/color][\/size]\r\n[SIZE=\"3\"][COLOR=\"Red\"]As of version 3.8.2 you will need to download and install the following libraries to use this addon:\r\n[\/COLOR][\/SIZE][SIZE=\"2\"][URL=\"https:\/\/www.esoui.com\/downloads\/info7-LibAddonMenu.html\"]LibAddonMenu[\/URL]  V2.0r37 or later\r\n[URL=\"https:\/\/www.esoui.com\/downloads\/info1146-LibCustomMenu.html\"]LubCustomMenu[\/URL] V7.2.2 or later\r\n[\/SIZE]\r\n[b][color=#B60000]Ravalox' Quest Tracker[\/color][\/b] provides [b]A more intuitive QUEST TRACKER[\/b] than the in-game provided tracker interface.  (see screenshots) The interface can be left on the screen or hidden via settings or \/slash commands.\r\n\r\n[size=4][color=#e69138][b]Features:[\/b][\/color][\/size]\r\n[list][*] Use the ZOS Tracker key-bind to cycle through quests (default \"T\")\r\n[*] Can switch between the game's tracker and Ravalox' Tracker in the settings menu\r\n[*] Option to automatically hide Ravalox' tracker when in combat\r\n[*] Manual or automatic window re-sizing (by width, height or both)\r\n[*] Backdrop color and texture  can be disabled or changed\r\n[*] Background and drag bar can be \"faded\" to become any level of transparent\r\n[*] Window can be locked or unlocked via \/slash command, settings menu, or lock icon\r\n[*] Can optionally display the number of open quests in each Category or Zone\r\n[*] Can optionally display and sort by the level of the quest\r\n[*] Quest Tooltips can be enabled or disabled in the settings menu\r\n[*] Optionally show Quest Tooltips when a quest is accepted\r\n[*] Option to display a chat alert when accepting a quest.  With or without quest details\r\n[*] Save settings either account-wide or by character\r\n[\/list]\r\n\r\n[size=4][color=#e69138][b]How to use:[\/b][\/color][\/size]\r\n[list][*] Lock icon will allow for the UI to be dragged around the screen (by dragging the bar at the top of the UI) when unlocked.  (Lock icon can be hidden via the settings menu)\r\n[*] If the tooltips option is selected, a tooltip window will appear when hovering over a quest.\r\n[*] Right click on a quest to display quest management options. (Share , Abandon, or show on map)\r\n[*] Click on the Category to expand or collapse the list.\r\n[*] Additional options for automatically collapsing Categories, Quests or both\r\n[*] Options for showing all quests at all times or show all quests in the currently selected category at login\r\n[\/list]\r\n\r\n[size=4][color=#e69138][b]Slash Commands:[\/b][\/color][\/size]\r\n[list]\r\n[*][color=#00bbbb][b]\/qt [\/b][\/color]or[color=#00bbbb][b]\/qt help[\/b][\/color] Display the QuestTracker \/slash command help\r\n[*][color=#00bbbb][b]\/qth[\/b][\/color] Toggle the Questtracker UI on or off on the fly.\r\n[*][color=#00bbbb][b]\/qtl[\/b][\/color] Toggle the QuestTracker UI lock on or off on the fly.\r\n[*][color=#00bbbb][b]\/qta[\/b][\/color] Toggle show\/hide Categories and Quest entries on the fly.\r\n[*][color=#00bbbb][b]\/qts[\/b][\/color] Show and Expand all Quests only in the Selected Category.\r\n[*][color=#00bbbb][b]\/questracker[\/b][\/color] or [color=#00bbbb][b]\/qt settings[\/b][\/color] ... Directly access the settings menu.\r\n[\/list]",
    "changeLog": "v3.8.3.0 ~ Calamath\r\n\t* Changed add-on name to RavaloxsQuestTracker.\r\n\t* As a painful side effect, the add-on settings will initialize to default. However, the save data file with the old add-on name remains and should be reusable if renamed with ESO not running.\r\n\t* Updated to API 101043\r\n\r\n\r\nv3.8.2.15 ~ Calamath\r\n\t* Updated to Gold Road API (101042)\r\n\r\nv3.8.2.14 ~ Calamath\r\n\t* Since the 'Automatic Quest Tracking' setting was added to the vanilla UI in Update38, the same setting items in the add-on now work together. \r\n\t* Updated to Necrom API (101038)\r\n\r\nv3.8.2.13 ~ Calamath\r\n\t* Updated to Firesong API (101036)\r\n\r\nv3.8.2.12 ~ Calamath\r\n\t* Added failsafe to avoid rare UI errors caused by unusual quest update events.\r\n\t* Updated to Lost Depths API (101035)\r\n\r\nv3.8.2.11 ~ Calamath\r\n\t* Updated to High Isle API (101034)\r\n\r\nv3.8.2.10 ~ Calamath\r\n\t* Updated to Ascending Tide API (101033)\r\n\r\nv3.8.2.9 ~ Calamath\r\n\t* Addressed an issue where if the player changed the default quest tracker display settings in the system menu before opening the settings panel for this add-on, the settings would not be reflected in the add-on's save data.\r\n\r\nv3.8.2.8 ~ Calamath\r\n\t* Updated to Deadlands API (101032)\r\n\r\nv3.8.2.7 ~ Calamath\r\n\t* Removed a temporary fix implemented in the previous version 3.8.2.6.\r\n\t* Updated to Waking Flame API (101031)\r\n\r\nv3.8.2.6 ~ Calamath\r\n\t* Temporarily fixed an issue where a UI error message was displayed immediately after the player logged in.\r\n\r\nv3.8.2.5 ~ Calamath\r\n\t* Updated to Blackwood API (100035)\r\n\r\nv3.8.2.4 ~ Calamath\r\n\t* Updated to Flames of Ambition API (100034)\r\n\r\nv3.8.2.3 ~ Calamath\r\n\t* Fixed an issue that caused a UI error when showing or hiding the default quest tracker.\r\n\t* Removed \"Arial Narrow\" font choice from UI font settings.\r\n\t* Fixed an issue where the focused quest would automatically switch each time you accepted a new quest.\r\n\r\nv3.8.2.2 ~ Calamath\r\n\t* Updated to Markarth API (100033)\r\n\t* Removed references and uses of LibStub completely.\r\n\t* Embedded LibQuestInfo library now has its own manifest file.\r\n\t* Added 'AddOnVersion' directive to manifest file.\r\n\t* Fixed an issue where the quest tracker was not updated correctly.\r\n\r\nv3.8.2.0\r\n\t* Updated to Harrowstorm API (10030)\r\n\t* Removed references and use of LibStub\r\n\t* LibAddonMenu & LibCustomMenu removed from code and changed to DependsOn components. You will need to install these libraries separately from this addon.\r\n\r\nv3.8.1.0\r\n* Updated to Elsweyr API (100027)\r\n\r\nv3.8.0.0\r\n* Updated to Wrathstone API (10026)\r\n* LibCustomMenu updated to 6.6.3\r\n\r\nv3.7.0.0\r\n* Updated to API 10025\r\n* Updated LAM to 2.0 r26\r\n* Fixed error message when addon initialized. QUEST_TRACKER interface appears to be FOCUSED_QUEST_TRACKER now.\r\n* Fixed error message when quest changed outside of addon. QUEST_TRACKER interface appears to be FOCUSED_QUEST_TRACKER now.\r\n\r\n\r\nv3.6.0.0 ~ Calia\r\n\t* Updated to API 20 (Horns Of The Reach)\r\n\t* Updated LAM to 2.0 r24\r\n\r\nv3.5.0.0 ~Calia\r\n*Updated to API 19 (Morrowind)\r\n\r\nv3.4.4.2\r\n*Updated to API 18\r\n\r\nv3.4.4.2\r\n*Updated to API 17\r\n\r\nv3.4.4.1\r\n* Update to API 16\r\n\r\nv3.4.4.0\r\n* Temporary fix for white square icons\r\n\r\nv3.4.3.2 - June 04 2016\r\n* Updated to API version 15\r\nv3.4.3.1 - March 07 2016\r\n* Updated to API version 14 (1.9)\r\n\r\nV3.4.3.0 - February 12 2016\r\n* Corrected a bug where under certain circumstances the ZOS QT will be visible when first using \r\n   the addon or after creating a new character\r\n* Inserted  workaround for the last of the cutoff text issues.  This problem may be ZOS related, \r\n   a permanent fix will be made after getting more info.\r\n\r\nV3.4.2.0 - February 09 2016\r\n\r\n* Fixed crash when no quests are active (no quests in list)\r\n* Fixed various formatting problems\r\n\r\n[color #00ff00][b]Known Issues:[\/b][\/color] v3.4.2.0:  \r\n\t\t  [b]Issue with text being cut off (under certain conditions)[\/b]\r\n---\r\n\r\n\r\nv3.4.1.0 - February 07 2016\r\n\r\n* Bug fix:  Hide DFLT QuestTracker was restoring itself on player load.\r\n---\r\n\r\n\r\nv3.4.0.46 - February 07 2016\r\n\r\n* Added QuestTracker window Auto Size options\r\n* Added Account-Wide\/Character Settings option\r\n* Additional \/slash commands and help menu updated\r\n* Re-organized Settings Menu\r\n* Option to hide QT window when unit is in combat\r\n* Option to sort quests (within Categories) by level\r\n* Additional tweaks to ensure text is not chopped off in the right margin\r\n---\r\n\r\n\r\nv3.3.0.4 - January 23 2016\r\n\r\n* Added dragbar color option\r\n* Settings menu look\/feel adjustments\r\n* Eliminated left indent\r\n* Added \/slash command for lock\/unlock ui.  (\/qtl)\r\n* Modified some slash commands to be toggles eliminating redundancies\r\n* Added UI lock option in settings menu\r\n* Added Display choices for lock Icon\r\n* Added Option to force all quests in Active category open on login.\r\n\r\n[color #dd0000]Known Issues: v3.3.0.4:[\/color]  \tNone.\r\n\r\n[color #00dd00]More still to come!![\/color]\r\n\r\n--------\r\n\r\n\r\n3.1.1.0\r\n* Added a feature to expand all nodes in the list.  This feature is enabled in the settings menu, and can be controlled during play via slash command.\r\n* Fixed override quest colors where under certain conditions the colors would not display correctly.\r\n* Added option to chat alert to include quest details\r\n* Adjusted items in the settings menu for better \"flow\"\r\n* Changed Saved Variables to account-wide\r\n\r\nKnown Issues: v3.1.1.0\r\n\r\n\tNone.\r\n\r\n3.0.0.2\r\nComplete re-write of the addon.  The re-written addon has the look and feel of the old version but the mechanics work much smoother and the code is optimized.\r\nMigrated Addon to Ravalox' addon catalog.\r\n\r\n2.3.5.2\r\nUpdated to support ESO 2.2.4\r\n\r\n2.3.5.1\r\nUpdated to correct \"T\" key quest cycle behaviour.  (with thanks to Circonian for his assistance!)\r\nCompleted quest steps (for quests with mutliple steps showing) will now be greyed out)\r\n\r\nKnown Issues:\r\n                         * QT UI flashes on the screen during questgiver conversations (and assorted other times).  \r\n                         * Under certain circumstances, Crafting quests (certification and writs) cannot be selected \r\n                            by mouse.  They can still be cycled by \"T\" key or Journal menu.\r\n\r\nI have a partial fix in the works for the flashing issue.  Circonian had written a modification that corrects the flashing but the QT UI scrambles when the Journal menu is used to change quests or a new quest is added.  I've put in additional code that corrects this but causes an additional update issue. (therefore this fix is a WIP)\r\n\r\nSince the Cycling of the Crafting Writs is a sporadic issue (it's affecting only two of my characters on my client, so ... odd.) I will work towards a solution, but it will not be as high a priority.  :-)\r\n                      \r\n\r\n2.3.4.0\r\nCorrected code to work with ZOS changes to internal quest tracker.  (They are preparing for Gamepad and Keyboard integration causing changes to ZOS Quest source code.)\r\n\r\n2.3.3.3\r\nThis update is to support ESOTU 2.1 (API 1.7)\r\n**NOTE:  This addon is currently broken.  Additional code needs to change to make it compatible with ZOS' quest tracker API changes.\r\n\r\n2.3.3.2\r\nThis update (2.3.3.2) only updates supporting documentation to reflect the transition from Wykkyd to Ravalox and Balkoth under the name of Exodus Code Group.  The LUA addon\/library code has not been changed.",
    "downloads": 1373067,
    "downloadsMonthly": 7890,
    "favorites": 1337,
    "images": [
      {
        "thumbUrl": "https:\/\/cdn-eso.mmoui.com\/preview\/tiny\/pvw4095.png",
        "imageUrl": "https:\/\/cdn-eso.mmoui.com\/preview\/pvw4095.png",
        "description": "Background is optional."
      },
      {
        "thumbUrl": "https:\/\/cdn-eso.mmoui.com\/preview\/tiny\/pvw4094.png",
        "imageUrl": "https:\/\/cdn-eso.mmoui.com\/preview\/pvw4094.png",
        "description": "Quest Tracker UI can be resized.  Scroll bar appears if UI is too small to show entire list."
      },
      {
        "thumbUrl": "https:\/\/cdn-eso.mmoui.com\/preview\/tiny\/pvw4093.png",
        "imageUrl": "https:\/\/cdn-eso.mmoui.com\/preview\/pvw4093.png",
        "description": "Quest Tracker Interface with Tooltips opn the right.  (Tooltips can be moved to left side in the set"
      }
    ]
  }
]
```
