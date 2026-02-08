export interface Asset {
  id: number;
  absolute_path: string;
  filename: string;
  file_size: number;
  folder_id: number;
  modified_at: string;
  thumbnail: string;
  favorited: number;
  last_used_at: string;
  created_at: string;
  updated_at: string;
}

export interface WatchFolder {
  id: number;
  path: string;
  created_at: string;
}

export interface Tag {
  id: number;
  name: string;
}

export interface TagWithCount {
  id: number;
  name: string;
  count: number;
}

export interface Collection {
  id: number;
  name: string;
  description: string;
  icon: string;
  asset_count: number;
  created_at: string;
}
