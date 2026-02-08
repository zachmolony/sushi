package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx        context.Context
	db         *Database
	thumbDir   string
	fileServer *LocalFileServer
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	dataHome := os.Getenv("XDG_DATA_HOME")
	if dataHome == "" {
		home, _ := os.UserHomeDir()
		dataHome = filepath.Join(home, ".local", "share")
	}
	a.thumbDir = filepath.Join(dataHome, "sushi", "thumbnails")
	os.MkdirAll(a.thumbDir, 0755)

	// Open database
	db, err := NewDatabase()
	if err != nil {
		fmt.Printf("FATAL: could not open database: %v\n", err)
		return
	}
	a.db = db

	// Start the local file server on its own port
	a.fileServer = StartLocalFileServer()

	// Scan all existing watch folders on startup
	go ScanAllFolders(a.db)
}

// shutdown is called when the app is closing
func (a *App) shutdown(ctx context.Context) {
	if a.db != nil {
		a.db.Close()
	}
}

// GetFileServerURL returns the base URL of the local file server.
func (a *App) GetFileServerURL() string {
	if a.fileServer == nil || a.fileServer.Port == 0 {
		return ""
	}
	return fmt.Sprintf("http://127.0.0.1:%d", a.fileServer.Port)
}

// --- Watch Folder Methods ---

// AddWatchFolder opens a folder picker, registers the folder, scans it, and returns the updated asset list.
func (a *App) AddWatchFolder() ([]Asset, error) {
	dir, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select a folder to watch for 3D assets",
	})
	if err != nil {
		return nil, err
	}
	if dir == "" {
		return a.GetAssets()
	}

	folder, err := a.db.AddWatchFolder(dir)
	if err != nil {
		return nil, fmt.Errorf("add folder: %w", err)
	}

	count, err := ScanFolder(a.db, *folder)
	if err != nil {
		fmt.Printf("scan error: %v\n", err)
	} else {
		fmt.Printf("scanned %s: found %d assets\n", dir, count)
	}

	return a.GetAssets()
}

// RemoveWatchFolder removes a watch folder and all its assets from the DB.
func (a *App) RemoveWatchFolder(id int64) error {
	if err := a.db.DeleteAssetsByFolder(id); err != nil {
		return err
	}
	return a.db.RemoveWatchFolder(id)
}

// GetWatchFolders returns all registered watch folders.
func (a *App) GetWatchFolders() ([]WatchFolder, error) {
	return a.db.ListWatchFolders()
}

// RescanFolder re-scans a specific watch folder.
func (a *App) RescanFolder(id int64) ([]Asset, error) {
	folder, err := a.db.GetWatchFolder(id)
	if err != nil {
		return nil, err
	}
	ScanFolder(a.db, *folder)
	return a.GetAssets()
}

// --- Asset Methods ---

// GetAssets returns all indexed assets.
func (a *App) GetAssets() ([]Asset, error) {
	assets, err := a.db.ListAssets()
	if err != nil {
		return nil, err
	}
	if assets == nil {
		assets = []Asset{}
	}
	return assets, nil
}

// GetAssetsByTag returns assets that have a specific tag.
func (a *App) GetAssetsByTag(tagName string) ([]Asset, error) {
	assets, err := a.db.GetAssetsByTag(tagName)
	if err != nil {
		return nil, err
	}
	if assets == nil {
		assets = []Asset{}
	}
	return assets, nil
}

// GetUntaggedAssets returns assets with no tags.
func (a *App) GetUntaggedAssets() ([]Asset, error) {
	assets, err := a.db.GetUntaggedAssets()
	if err != nil {
		return nil, err
	}
	if assets == nil {
		assets = []Asset{}
	}
	return assets, nil
}

// GetAssetIDsByTags returns IDs of assets that have any of the given tags.
func (a *App) GetAssetIDsByTags(tagNames []string) ([]int64, error) {
	return a.db.GetAssetIDsByTags(tagNames)
}

// GetFavoritedAssets returns all favorited assets.
func (a *App) GetFavoritedAssets() ([]Asset, error) {
	assets, err := a.db.GetFavoritedAssets()
	if err != nil {
		return nil, err
	}
	if assets == nil {
		assets = []Asset{}
	}
	return assets, nil
}

// ToggleFavorite toggles the favorite status of an asset.
func (a *App) ToggleFavorite(assetID int64) (bool, error) {
	return a.db.ToggleFavorite(assetID)
}

// BulkSetFavorite sets favorite status for multiple assets.
func (a *App) BulkSetFavorite(assetIDs []int64, favorited bool) error {
	return a.db.BulkSetFavorite(assetIDs, favorited)
}

// DeleteAsset deletes a single asset from the database (not from disk).
func (a *App) DeleteAsset(assetID int64) error {
	return a.db.DeleteAssetByID(assetID)
}

// BulkDeleteAssets deletes multiple assets from the database. Returns the count removed.
func (a *App) BulkDeleteAssets(assetIDs []int64) (int, error) {
	return a.db.DeleteAssetsByIDs(assetIDs)
}

// MarkAssetUsed marks an asset as recently used.
func (a *App) MarkAssetUsed(assetID int64) error {
	return a.db.SetAssetUsed(assetID)
}

// GetRecentlyUsedAssets returns recently used assets.
func (a *App) GetRecentlyUsedAssets() ([]Asset, error) {
	assets, err := a.db.GetRecentlyUsedAssets(200)
	if err != nil {
		return nil, err
	}
	if assets == nil {
		assets = []Asset{}
	}
	return assets, nil
}

// GetRecentlyAddedAssets returns recently added assets.
func (a *App) GetRecentlyAddedAssets() ([]Asset, error) {
	assets, err := a.db.GetRecentlyAddedAssets(200)
	if err != nil {
		return nil, err
	}
	if assets == nil {
		assets = []Asset{}
	}
	return assets, nil
}

// --- Tag Methods ---

// GetAllTags returns every tag in the system.
func (a *App) GetAllTags() ([]Tag, error) {
	tags, err := a.db.ListTags()
	if err != nil {
		return nil, err
	}
	if tags == nil {
		tags = []Tag{}
	}
	return tags, nil
}

// GetTagsForAsset returns tags attached to a specific asset.
func (a *App) GetTagsForAsset(assetID int64) ([]Tag, error) {
	tags, err := a.db.GetTagsForAsset(assetID)
	if err != nil {
		return nil, err
	}
	if tags == nil {
		tags = []Tag{}
	}
	return tags, nil
}

// AddTagToAsset creates a tag if needed and attaches it to an asset.
func (a *App) AddTagToAsset(assetID int64, tagName string) ([]Tag, error) {
	tag, err := a.db.CreateTag(tagName)
	if err != nil {
		return nil, err
	}
	if err := a.db.TagAsset(assetID, tag.ID); err != nil {
		return nil, err
	}
	return a.GetTagsForAsset(assetID)
}

// RemoveTagFromAsset removes a tag from an asset.
func (a *App) RemoveTagFromAsset(assetID int64, tagID int64) ([]Tag, error) {
	if err := a.db.UntagAsset(assetID, tagID); err != nil {
		return nil, err
	}
	return a.GetTagsForAsset(assetID)
}

// --- Thumbnail Methods ---

// SaveThumbnail saves a base64-encoded PNG thumbnail for an asset.
// Called from the frontend after rendering with Three.js.
func (a *App) SaveThumbnail(assetID int64, base64PNG string) error {
	return a.db.SetThumbnail(assetID, base64PNG)
}

// SavePolyCount saves the triangle/polygon count for an asset.
// Called from the frontend after parsing with Three.js.
func (a *App) SavePolyCount(assetID int64, count int64) error {
	return a.db.SetPolyCount(assetID, count)
}

// GetThumbnail returns the base64 PNG data for an asset's thumbnail.
func (a *App) GetThumbnail(assetID int64) (string, error) {
	return a.db.GetThumbnail(assetID)
}

// ClearAllThumbnails wipes all cached thumbnails so they regenerate on next load.
func (a *App) ClearAllThumbnails() (int64, error) {
	return a.db.ClearAllThumbnails()
}

// --- Utility Methods ---

// OpenFileInFolder opens the system file manager with the file's directory.
func (a *App) OpenFileInFolder(absolutePath string) error {
	dir := filepath.Dir(absolutePath)
	return openFileManager(dir)
}

// --- Collection Methods ---

// CreateCollection makes a new collection shelf.
func (a *App) CreateCollection(name string, icon string) (*Collection, error) {
	return a.db.CreateCollection(name, icon)
}

// GetCollections returns all collections with their asset counts.
func (a *App) GetCollections() ([]Collection, error) {
	cols, err := a.db.ListCollections()
	if err != nil {
		return nil, err
	}
	if cols == nil {
		cols = []Collection{}
	}
	return cols, nil
}

// RenameCollection renames a collection.
func (a *App) RenameCollection(id int64, name string) error {
	return a.db.RenameCollection(id, name)
}

// DeleteCollection removes a collection (assets are not deleted).
func (a *App) DeleteCollection(id int64) error {
	return a.db.DeleteCollection(id)
}

// AddToCollection adds an asset to a collection.
func (a *App) AddToCollection(collectionID int64, assetID int64) error {
	return a.db.AddAssetToCollection(collectionID, assetID)
}

// RemoveFromCollection removes an asset from a collection.
func (a *App) RemoveFromCollection(collectionID int64, assetID int64) error {
	return a.db.RemoveAssetFromCollection(collectionID, assetID)
}

// GetCollectionAssets returns all assets in a collection.
func (a *App) GetCollectionAssets(collectionID int64) ([]Asset, error) {
	assets, err := a.db.GetAssetsInCollection(collectionID)
	if err != nil {
		return nil, err
	}
	if assets == nil {
		assets = []Asset{}
	}
	return assets, nil
}

// --- Bulk Operations ---

// BulkTagAssets applies a tag to multiple assets at once.
func (a *App) BulkTagAssets(assetIDs []int64, tagName string) error {
	return a.db.BulkTagAssets(assetIDs, tagName)
}

// BulkAddToCollection adds multiple assets to a collection at once.
func (a *App) BulkAddToCollection(collectionID int64, assetIDs []int64) error {
	return a.db.BulkAddToCollection(collectionID, assetIDs)
}

// GetTagsWithCounts returns all tags with usage counts, ordered by most common.
func (a *App) GetTagsWithCounts() ([]TagWithCount, error) {
	tags, err := a.db.ListTagsWithCounts()
	if err != nil {
		return nil, err
	}
	if tags == nil {
		tags = []TagWithCount{}
	}
	return tags, nil
}

// GetAssetsByTags returns assets that have ALL of the specified tags.
func (a *App) GetAssetsByTags(tagNames []string) ([]Asset, error) {
	assets, err := a.db.GetAssetsByTags(tagNames)
	if err != nil {
		return nil, err
	}
	if assets == nil {
		assets = []Asset{}
	}
	return assets, nil
}

// GetCollectionsForAsset returns which collections an asset belongs to.
func (a *App) GetCollectionsForAsset(assetID int64) ([]Collection, error) {
	cols, err := a.db.GetCollectionsForAsset(assetID)
	if err != nil {
		return nil, err
	}
	if cols == nil {
		cols = []Collection{}
	}
	return cols, nil
}
