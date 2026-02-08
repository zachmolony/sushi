package main

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"
)

// WatchFolder represents a root directory being watched for assets.
type WatchFolder struct {
	ID        int64  `json:"id"`
	Path      string `json:"path"`
	CreatedAt string `json:"created_at"`
}

// Asset represents a single .glb/.gltf file found on disk.
type Asset struct {
	ID           int64  `json:"id"`
	AbsolutePath string `json:"absolute_path"`
	Filename     string `json:"filename"`
	FolderID     int64  `json:"folder_id"`
	FileSize     int64  `json:"file_size"`
	ModifiedAt   string `json:"modified_at"`
	Thumbnail    string `json:"thumbnail"`
	Favorited    int64  `json:"favorited"`
	LastUsedAt   string `json:"last_used_at"`
	PolyCount    int64  `json:"poly_count"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

// Tag represents a user-defined label.
type Tag struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// Database wraps the SQLite connection and provides all data methods.
type Database struct {
	db *sql.DB
}

// NewDatabase opens (or creates) the SQLite database and runs the schema.
func NewDatabase() (*Database, error) {
	dataHome := os.Getenv("XDG_DATA_HOME")
	if dataHome == "" {
		home, _ := os.UserHomeDir()
		dataHome = filepath.Join(home, ".local", "share")
	}
	dbDir := filepath.Join(dataHome, "sushi")
	os.MkdirAll(dbDir, 0755)
	dbPath := filepath.Join(dbDir, "sushi.db")

	db, err := sql.Open("sqlite", dbPath+"?_pragma=journal_mode(wal)&_pragma=foreign_keys(on)")
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}
	// SQLite only supports one writer — limit pool to avoid SQLITE_BUSY errors
	db.SetMaxOpenConns(1)

	d := &Database{db: db}
	if err := d.migrate(); err != nil {
		return nil, fmt.Errorf("migrate: %w", err)
	}
	return d, nil
}

func (d *Database) migrate() error {
	schema := `
	CREATE TABLE IF NOT EXISTS watch_folders (
		id         INTEGER PRIMARY KEY AUTOINCREMENT,
		path       TEXT    NOT NULL UNIQUE,
		created_at TEXT    NOT NULL DEFAULT (datetime('now'))
	);

	CREATE TABLE IF NOT EXISTS assets (
		id            INTEGER PRIMARY KEY AUTOINCREMENT,
		absolute_path TEXT    NOT NULL UNIQUE,
		filename      TEXT    NOT NULL,
		folder_id     INTEGER NOT NULL REFERENCES watch_folders(id) ON DELETE CASCADE,
		file_size     INTEGER NOT NULL,
		modified_at   TEXT    NOT NULL,
		thumbnail     TEXT    DEFAULT '',
		favorited     INTEGER NOT NULL DEFAULT 0,
		last_used_at  TEXT    DEFAULT '',
		poly_count    INTEGER NOT NULL DEFAULT 0,
		created_at    TEXT    NOT NULL DEFAULT (datetime('now')),
		updated_at    TEXT    NOT NULL DEFAULT (datetime('now'))
	);

	CREATE INDEX IF NOT EXISTS idx_assets_folder   ON assets(folder_id);
	CREATE INDEX IF NOT EXISTS idx_assets_filename  ON assets(filename);

	CREATE TABLE IF NOT EXISTS tags (
		id   INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT    NOT NULL UNIQUE
	);

	CREATE TABLE IF NOT EXISTS asset_tags (
		asset_id INTEGER NOT NULL REFERENCES assets(id) ON DELETE CASCADE,
		tag_id   INTEGER NOT NULL REFERENCES tags(id)   ON DELETE CASCADE,
		PRIMARY KEY (asset_id, tag_id)
	);

	CREATE TABLE IF NOT EXISTS collections (
		id          INTEGER PRIMARY KEY AUTOINCREMENT,
		name        TEXT    NOT NULL UNIQUE,
		description TEXT    NOT NULL DEFAULT '',
		icon        TEXT    NOT NULL DEFAULT '📁',
		created_at  TEXT    NOT NULL DEFAULT (datetime('now'))
	);

	CREATE TABLE IF NOT EXISTS collection_assets (
		collection_id INTEGER NOT NULL REFERENCES collections(id) ON DELETE CASCADE,
		asset_id      INTEGER NOT NULL REFERENCES assets(id)      ON DELETE CASCADE,
		added_at      TEXT    NOT NULL DEFAULT (datetime('now')),
		PRIMARY KEY (collection_id, asset_id)
	);
	`
	_, err := d.db.Exec(schema)
	if err != nil {
		return err
	}

	// Add columns if they don't exist (for existing databases)
	d.db.Exec("ALTER TABLE assets ADD COLUMN favorited INTEGER NOT NULL DEFAULT 0")
	d.db.Exec("ALTER TABLE assets ADD COLUMN last_used_at TEXT DEFAULT ''")
	d.db.Exec("ALTER TABLE assets ADD COLUMN poly_count INTEGER NOT NULL DEFAULT 0")

	return nil
}

// ClearAllThumbnails resets thumbnail and poly_count for all assets so they regenerate.
func (d *Database) ClearAllThumbnails() (int64, error) {
	res, err := d.db.Exec("UPDATE assets SET thumbnail = '', poly_count = 0")
	if err != nil {
		return 0, err
	}
	n, _ := res.RowsAffected()
	return n, nil
}

// --- Watch Folders ---

func (d *Database) AddWatchFolder(path string) (*WatchFolder, error) {
	res, err := d.db.Exec("INSERT INTO watch_folders (path) VALUES (?)", path)
	if err != nil {
		return nil, err
	}
	id, _ := res.LastInsertId()
	return d.GetWatchFolder(id)
}

func (d *Database) GetWatchFolder(id int64) (*WatchFolder, error) {
	row := d.db.QueryRow("SELECT id, path, created_at FROM watch_folders WHERE id = ?", id)
	f := &WatchFolder{}
	err := row.Scan(&f.ID, &f.Path, &f.CreatedAt)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (d *Database) ListWatchFolders() ([]WatchFolder, error) {
	rows, err := d.db.Query("SELECT id, path, created_at FROM watch_folders ORDER BY path")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var folders []WatchFolder
	for rows.Next() {
		var f WatchFolder
		if err := rows.Scan(&f.ID, &f.Path, &f.CreatedAt); err != nil {
			return nil, err
		}
		folders = append(folders, f)
	}
	return folders, nil
}

func (d *Database) RemoveWatchFolder(id int64) error {
	_, err := d.db.Exec("DELETE FROM watch_folders WHERE id = ?", id)
	return err
}

// --- Assets ---

func (d *Database) UpsertAsset(absolutePath string, folderID int64, fileSize int64, modifiedAt time.Time) (*Asset, error) {
	filename := filepath.Base(absolutePath)
	modStr := modifiedAt.UTC().Format(time.RFC3339)
	nowStr := time.Now().UTC().Format(time.RFC3339)

	_, err := d.db.Exec(`
		INSERT INTO assets (absolute_path, filename, folder_id, file_size, modified_at, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(absolute_path) DO UPDATE SET
			file_size   = excluded.file_size,
			modified_at = excluded.modified_at,
			updated_at  = ?
	`, absolutePath, filename, folderID, fileSize, modStr, nowStr, nowStr, nowStr)
	if err != nil {
		return nil, err
	}

	row := d.db.QueryRow("SELECT id, absolute_path, filename, folder_id, file_size, modified_at, thumbnail, favorited, last_used_at, poly_count, created_at, updated_at FROM assets WHERE absolute_path = ?", absolutePath)
	a := &Asset{}
	err = row.Scan(&a.ID, &a.AbsolutePath, &a.Filename, &a.FolderID, &a.FileSize, &a.ModifiedAt, &a.Thumbnail, &a.Favorited, &a.LastUsedAt, &a.PolyCount, &a.CreatedAt, &a.UpdatedAt)
	return a, err
}

func (d *Database) ListAssets() ([]Asset, error) {
	rows, err := d.db.Query("SELECT id, absolute_path, filename, folder_id, file_size, modified_at, thumbnail, favorited, last_used_at, poly_count, created_at, updated_at FROM assets ORDER BY filename")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var assets []Asset
	for rows.Next() {
		var a Asset
		if err := rows.Scan(&a.ID, &a.AbsolutePath, &a.Filename, &a.FolderID, &a.FileSize, &a.ModifiedAt, &a.Thumbnail, &a.Favorited, &a.LastUsedAt, &a.PolyCount, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, err
		}
		assets = append(assets, a)
	}
	return assets, nil
}

func (d *Database) DeleteAssetByPath(absolutePath string) error {
	_, err := d.db.Exec("DELETE FROM assets WHERE absolute_path = ?", absolutePath)
	return err
}

func (d *Database) GetAssetByID(id int64) (*Asset, error) {
	row := d.db.QueryRow(`
		SELECT id, absolute_path, filename, folder_id, file_size, modified_at,
		       thumbnail, favorited, last_used_at, poly_count, created_at, updated_at
		FROM assets WHERE id = ?
	`, id)
	a := &Asset{}
	err := row.Scan(&a.ID, &a.AbsolutePath, &a.Filename, &a.FolderID, &a.FileSize,
		&a.ModifiedAt, &a.Thumbnail, &a.Favorited, &a.LastUsedAt, &a.PolyCount,
		&a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func (d *Database) DeleteAssetByID(id int64) error {
	_, err := d.db.Exec("DELETE FROM assets WHERE id = ?", id)
	return err
}

func (d *Database) DeleteAssetsByIDs(ids []int64) (int, error) {
	if len(ids) == 0 {
		return 0, nil
	}
	tx, err := d.db.Begin()
	if err != nil {
		return 0, err
	}
	count := 0
	for _, id := range ids {
		res, err := tx.Exec("DELETE FROM assets WHERE id = ?", id)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
		affected, _ := res.RowsAffected()
		count += int(affected)
	}
	return count, tx.Commit()
}

func (d *Database) DeleteAssetsByFolder(folderID int64) error {
	_, err := d.db.Exec("DELETE FROM assets WHERE folder_id = ?", folderID)
	return err
}

// PruneAssetsForFolder removes DB rows for files that no longer exist on disk.
func (d *Database) PruneAssetsForFolder(folderID int64) (int, error) {
	rows, err := d.db.Query("SELECT id, absolute_path FROM assets WHERE folder_id = ?", folderID)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	var toDelete []int64
	for rows.Next() {
		var id int64
		var path string
		if err := rows.Scan(&id, &path); err != nil {
			return 0, err
		}
		if _, err := os.Stat(path); os.IsNotExist(err) {
			toDelete = append(toDelete, id)
		}
	}

	for _, id := range toDelete {
		d.db.Exec("DELETE FROM assets WHERE id = ?", id)
	}
	return len(toDelete), nil
}

// --- Untagged / Favorites / Recently Used ---

func (d *Database) GetUntaggedAssets() ([]Asset, error) {
	rows, err := d.db.Query(`
		SELECT a.id, a.absolute_path, a.filename, a.folder_id, a.file_size, a.modified_at, a.thumbnail, a.favorited, a.last_used_at, a.poly_count, a.created_at, a.updated_at
		FROM assets a
		LEFT JOIN asset_tags at ON at.asset_id = a.id
		WHERE at.asset_id IS NULL
		ORDER BY a.filename
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var assets []Asset
	for rows.Next() {
		var a Asset
		if err := rows.Scan(&a.ID, &a.AbsolutePath, &a.Filename, &a.FolderID, &a.FileSize, &a.ModifiedAt, &a.Thumbnail, &a.Favorited, &a.LastUsedAt, &a.PolyCount, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, err
		}
		assets = append(assets, a)
	}
	return assets, nil
}

func (d *Database) GetFavoritedAssets() ([]Asset, error) {
	rows, err := d.db.Query(`
		SELECT id, absolute_path, filename, folder_id, file_size, modified_at, thumbnail, favorited, last_used_at, poly_count, created_at, updated_at
		FROM assets WHERE favorited = 1 ORDER BY filename
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var assets []Asset
	for rows.Next() {
		var a Asset
		if err := rows.Scan(&a.ID, &a.AbsolutePath, &a.Filename, &a.FolderID, &a.FileSize, &a.ModifiedAt, &a.Thumbnail, &a.Favorited, &a.LastUsedAt, &a.PolyCount, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, err
		}
		assets = append(assets, a)
	}
	return assets, nil
}

func (d *Database) ToggleFavorite(assetID int64) (bool, error) {
	var current int
	err := d.db.QueryRow("SELECT favorited FROM assets WHERE id = ?", assetID).Scan(&current)
	if err != nil {
		return false, err
	}
	newVal := 1 - current
	_, err = d.db.Exec("UPDATE assets SET favorited = ? WHERE id = ?", newVal, assetID)
	return newVal == 1, err
}

func (d *Database) SetAssetUsed(assetID int64) error {
	nowStr := time.Now().UTC().Format(time.RFC3339)
	_, err := d.db.Exec("UPDATE assets SET last_used_at = ? WHERE id = ?", nowStr, assetID)
	return err
}

func (d *Database) BulkSetFavorite(assetIDs []int64, favorited bool) error {
	val := 0
	if favorited {
		val = 1
	}
	for _, aid := range assetIDs {
		d.db.Exec("UPDATE assets SET favorited = ? WHERE id = ?", val, aid)
	}
	return nil
}

func (d *Database) GetRecentlyUsedAssets(limit int) ([]Asset, error) {
	rows, err := d.db.Query(`
		SELECT id, absolute_path, filename, folder_id, file_size, modified_at, thumbnail, favorited, last_used_at, poly_count, created_at, updated_at
		FROM assets WHERE last_used_at != '' ORDER BY last_used_at DESC LIMIT ?
	`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var assets []Asset
	for rows.Next() {
		var a Asset
		if err := rows.Scan(&a.ID, &a.AbsolutePath, &a.Filename, &a.FolderID, &a.FileSize, &a.ModifiedAt, &a.Thumbnail, &a.Favorited, &a.LastUsedAt, &a.PolyCount, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, err
		}
		assets = append(assets, a)
	}
	return assets, nil
}

func (d *Database) GetRecentlyAddedAssets(limit int) ([]Asset, error) {
	rows, err := d.db.Query(`
		SELECT id, absolute_path, filename, folder_id, file_size, modified_at, thumbnail, favorited, last_used_at, poly_count, created_at, updated_at
		FROM assets ORDER BY created_at DESC LIMIT ?
	`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var assets []Asset
	for rows.Next() {
		var a Asset
		if err := rows.Scan(&a.ID, &a.AbsolutePath, &a.Filename, &a.FolderID, &a.FileSize, &a.ModifiedAt, &a.Thumbnail, &a.Favorited, &a.LastUsedAt, &a.PolyCount, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, err
		}
		assets = append(assets, a)
	}
	return assets, nil
}

// --- Tags ---

func (d *Database) CreateTag(name string) (*Tag, error) {
	_, err := d.db.Exec("INSERT OR IGNORE INTO tags (name) VALUES (?)", name)
	if err != nil {
		return nil, err
	}
	// Always look up by name to get the correct ID
	row := d.db.QueryRow("SELECT id, name FROM tags WHERE name = ?", name)
	t := &Tag{}
	err = row.Scan(&t.ID, &t.Name)
	return t, err
}

func (d *Database) ListTags() ([]Tag, error) {
	rows, err := d.db.Query("SELECT id, name FROM tags ORDER BY name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []Tag
	for rows.Next() {
		var t Tag
		if err := rows.Scan(&t.ID, &t.Name); err != nil {
			return nil, err
		}
		tags = append(tags, t)
	}
	return tags, nil
}

// TagWithCount includes usage count for sorting by popularity.
type TagWithCount struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Count int    `json:"count"`
}

func (d *Database) ListTagsWithCounts() ([]TagWithCount, error) {
	rows, err := d.db.Query(`
		SELECT t.id, t.name, COUNT(at.asset_id) as cnt
		FROM tags t
		LEFT JOIN asset_tags at ON at.tag_id = t.id
		GROUP BY t.id, t.name
		ORDER BY cnt DESC, t.name ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []TagWithCount
	for rows.Next() {
		var t TagWithCount
		if err := rows.Scan(&t.ID, &t.Name, &t.Count); err != nil {
			return nil, err
		}
		tags = append(tags, t)
	}
	return tags, nil
}

func (d *Database) TagAsset(assetID, tagID int64) error {
	_, err := d.db.Exec("INSERT OR IGNORE INTO asset_tags (asset_id, tag_id) VALUES (?, ?)", assetID, tagID)
	return err
}

func (d *Database) UntagAsset(assetID, tagID int64) error {
	_, err := d.db.Exec("DELETE FROM asset_tags WHERE asset_id = ? AND tag_id = ?", assetID, tagID)
	return err
}

func (d *Database) GetTagsForAsset(assetID int64) ([]Tag, error) {
	rows, err := d.db.Query(`
		SELECT t.id, t.name FROM tags t
		JOIN asset_tags at ON at.tag_id = t.id
		WHERE at.asset_id = ?
		ORDER BY t.name
	`, assetID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []Tag
	for rows.Next() {
		var t Tag
		if err := rows.Scan(&t.ID, &t.Name); err != nil {
			return nil, err
		}
		tags = append(tags, t)
	}
	return tags, nil
}

func (d *Database) BulkTagAssets(assetIDs []int64, tagName string) error {
	tag, err := d.CreateTag(tagName)
	if err != nil {
		return err
	}
	for _, aid := range assetIDs {
		d.db.Exec("INSERT OR IGNORE INTO asset_tags (asset_id, tag_id) VALUES (?, ?)", aid, tag.ID)
	}
	return nil
}

func (d *Database) BulkAddToCollection(collectionID int64, assetIDs []int64) error {
	for _, aid := range assetIDs {
		d.db.Exec("INSERT OR IGNORE INTO collection_assets (collection_id, asset_id) VALUES (?, ?)", collectionID, aid)
	}
	return nil
}

func (d *Database) Close() error {
	return d.db.Close()
}

// --- Assets by Tag ---

func (d *Database) GetAssetsByTag(tagName string) ([]Asset, error) {
	rows, err := d.db.Query(`
		SELECT a.id, a.absolute_path, a.filename, a.folder_id, a.file_size, a.modified_at, a.thumbnail, a.favorited, a.last_used_at, a.poly_count, a.created_at, a.updated_at
		FROM assets a
		JOIN asset_tags at ON at.asset_id = a.id
		JOIN tags t ON t.id = at.tag_id
		WHERE t.name = ?
		ORDER BY a.filename
	`, tagName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var assets []Asset
	for rows.Next() {
		var a Asset
		if err := rows.Scan(&a.ID, &a.AbsolutePath, &a.Filename, &a.FolderID, &a.FileSize, &a.ModifiedAt, &a.Thumbnail, &a.Favorited, &a.LastUsedAt, &a.PolyCount, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, err
		}
		assets = append(assets, a)
	}
	return assets, nil
}

func (d *Database) GetAssetsByTags(tagNames []string) ([]Asset, error) {
	if len(tagNames) == 0 {
		return d.ListAssets()
	}
	// Build placeholders
	placeholders := ""
	args := make([]interface{}, len(tagNames))
	for i, name := range tagNames {
		if i > 0 {
			placeholders += ","
		}
		placeholders += "?"
		args[i] = name
	}
	args = append(args, len(tagNames))

	query := fmt.Sprintf(`
		SELECT a.id, a.absolute_path, a.filename, a.folder_id, a.file_size, a.modified_at, a.thumbnail, a.favorited, a.last_used_at, a.poly_count, a.created_at, a.updated_at
		FROM assets a
		JOIN asset_tags at ON at.asset_id = a.id
		JOIN tags t ON t.id = at.tag_id
		WHERE t.name IN (%s)
		GROUP BY a.id
		HAVING COUNT(DISTINCT t.id) = ?
		ORDER BY a.filename
	`, placeholders)

	rows, err := d.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var assets []Asset
	for rows.Next() {
		var a Asset
		if err := rows.Scan(&a.ID, &a.AbsolutePath, &a.Filename, &a.FolderID, &a.FileSize, &a.ModifiedAt, &a.Thumbnail, &a.Favorited, &a.LastUsedAt, &a.PolyCount, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, err
		}
		assets = append(assets, a)
	}
	return assets, nil
}

// GetAssetIDsByTags returns IDs of assets that have ANY of the given tags.
func (d *Database) GetAssetIDsByTags(tagNames []string) ([]int64, error) {
	if len(tagNames) == 0 {
		return nil, nil
	}
	placeholders := ""
	args := make([]interface{}, len(tagNames))
	for i, name := range tagNames {
		if i > 0 {
			placeholders += ","
		}
		placeholders += "?"
		args[i] = name
	}
	query := fmt.Sprintf(`
		SELECT DISTINCT a.id
		FROM assets a
		JOIN asset_tags at ON at.asset_id = a.id
		JOIN tags t ON t.id = at.tag_id
		WHERE t.name IN (%s)
	`, placeholders)

	rows, err := d.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

// --- Collections ---

// Collection represents a user-curated group of assets.
type Collection struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	AssetCount  int    `json:"asset_count"`
	CreatedAt   string `json:"created_at"`
}

func (d *Database) CreateCollection(name string, icon string) (*Collection, error) {
	if icon == "" {
		icon = "📁"
	}
	res, err := d.db.Exec("INSERT INTO collections (name, icon) VALUES (?, ?)", name, icon)
	if err != nil {
		return nil, err
	}
	id, _ := res.LastInsertId()
	return d.GetCollection(id)
}

func (d *Database) GetCollection(id int64) (*Collection, error) {
	row := d.db.QueryRow(`
		SELECT c.id, c.name, c.description, c.icon, c.created_at,
			(SELECT COUNT(*) FROM collection_assets ca WHERE ca.collection_id = c.id) as asset_count
		FROM collections c WHERE c.id = ?
	`, id)
	c := &Collection{}
	err := row.Scan(&c.ID, &c.Name, &c.Description, &c.Icon, &c.CreatedAt, &c.AssetCount)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (d *Database) ListCollections() ([]Collection, error) {
	rows, err := d.db.Query(`
		SELECT c.id, c.name, c.description, c.icon, c.created_at,
			(SELECT COUNT(*) FROM collection_assets ca WHERE ca.collection_id = c.id) as asset_count
		FROM collections c
		ORDER BY c.name
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var collections []Collection
	for rows.Next() {
		var c Collection
		if err := rows.Scan(&c.ID, &c.Name, &c.Description, &c.Icon, &c.CreatedAt, &c.AssetCount); err != nil {
			return nil, err
		}
		collections = append(collections, c)
	}
	return collections, nil
}

func (d *Database) RenameCollection(id int64, name string) error {
	_, err := d.db.Exec("UPDATE collections SET name = ? WHERE id = ?", name, id)
	return err
}

func (d *Database) DeleteCollection(id int64) error {
	_, err := d.db.Exec("DELETE FROM collections WHERE id = ?", id)
	return err
}

func (d *Database) AddAssetToCollection(collectionID, assetID int64) error {
	_, err := d.db.Exec("INSERT OR IGNORE INTO collection_assets (collection_id, asset_id) VALUES (?, ?)", collectionID, assetID)
	return err
}

func (d *Database) RemoveAssetFromCollection(collectionID, assetID int64) error {
	_, err := d.db.Exec("DELETE FROM collection_assets WHERE collection_id = ? AND asset_id = ?", collectionID, assetID)
	return err
}

func (d *Database) GetAssetsInCollection(collectionID int64) ([]Asset, error) {
	rows, err := d.db.Query(`
		SELECT a.id, a.absolute_path, a.filename, a.folder_id, a.file_size, a.modified_at, a.thumbnail, a.favorited, a.last_used_at, a.poly_count, a.created_at, a.updated_at
		FROM assets a
		JOIN collection_assets ca ON ca.asset_id = a.id
		WHERE ca.collection_id = ?
		ORDER BY a.filename
	`, collectionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var assets []Asset
	for rows.Next() {
		var a Asset
		if err := rows.Scan(&a.ID, &a.AbsolutePath, &a.Filename, &a.FolderID, &a.FileSize, &a.ModifiedAt, &a.Thumbnail, &a.Favorited, &a.LastUsedAt, &a.PolyCount, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, err
		}
		assets = append(assets, a)
	}
	return assets, nil
}

func (d *Database) GetCollectionsForAsset(assetID int64) ([]Collection, error) {
	rows, err := d.db.Query(`
		SELECT c.id, c.name, c.description, c.icon, c.created_at,
			(SELECT COUNT(*) FROM collection_assets ca2 WHERE ca2.collection_id = c.id) as asset_count
		FROM collections c
		JOIN collection_assets ca ON ca.collection_id = c.id
		WHERE ca.asset_id = ?
		ORDER BY c.name
	`, assetID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var collections []Collection
	for rows.Next() {
		var c Collection
		if err := rows.Scan(&c.ID, &c.Name, &c.Description, &c.Icon, &c.CreatedAt, &c.AssetCount); err != nil {
			return nil, err
		}
		collections = append(collections, c)
	}
	return collections, nil
}

// --- Thumbnails ---

func (d *Database) SetThumbnail(assetID int64, base64PNG string) error {
	_, err := d.db.Exec("UPDATE assets SET thumbnail = ? WHERE id = ?", base64PNG, assetID)
	return err
}

func (d *Database) SetPolyCount(assetID int64, count int64) error {
	_, err := d.db.Exec("UPDATE assets SET poly_count = ? WHERE id = ?", count, assetID)
	return err
}

func (d *Database) GetThumbnail(assetID int64) (string, error) {
	var thumb string
	err := d.db.QueryRow("SELECT thumbnail FROM assets WHERE id = ?", assetID).Scan(&thumb)
	return thumb, err
}
