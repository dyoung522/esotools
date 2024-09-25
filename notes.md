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
