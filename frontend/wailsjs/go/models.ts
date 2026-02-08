export namespace main {
	
	export class Asset {
	    id: number;
	    absolute_path: string;
	    filename: string;
	    folder_id: number;
	    file_size: number;
	    modified_at: string;
	    thumbnail: string;
	    favorited: number;
	    last_used_at: string;
	    poly_count: number;
	    created_at: string;
	    updated_at: string;
	
	    static createFrom(source: any = {}) {
	        return new Asset(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.absolute_path = source["absolute_path"];
	        this.filename = source["filename"];
	        this.folder_id = source["folder_id"];
	        this.file_size = source["file_size"];
	        this.modified_at = source["modified_at"];
	        this.thumbnail = source["thumbnail"];
	        this.favorited = source["favorited"];
	        this.last_used_at = source["last_used_at"];
	        this.poly_count = source["poly_count"];
	        this.created_at = source["created_at"];
	        this.updated_at = source["updated_at"];
	    }
	}
	export class BlenderStatus {
	    connected: boolean;
	    error?: string;
	
	    static createFrom(source: any = {}) {
	        return new BlenderStatus(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.connected = source["connected"];
	        this.error = source["error"];
	    }
	}
	export class Collection {
	    id: number;
	    name: string;
	    description: string;
	    icon: string;
	    asset_count: number;
	    created_at: string;
	
	    static createFrom(source: any = {}) {
	        return new Collection(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.description = source["description"];
	        this.icon = source["icon"];
	        this.asset_count = source["asset_count"];
	        this.created_at = source["created_at"];
	    }
	}
	export class Tag {
	    id: number;
	    name: string;
	
	    static createFrom(source: any = {}) {
	        return new Tag(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	    }
	}
	export class TagWithCount {
	    id: number;
	    name: string;
	    count: number;
	
	    static createFrom(source: any = {}) {
	        return new TagWithCount(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.count = source["count"];
	    }
	}
	export class WatchFolder {
	    id: number;
	    path: string;
	    created_at: string;
	
	    static createFrom(source: any = {}) {
	        return new WatchFolder(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.path = source["path"];
	        this.created_at = source["created_at"];
	    }
	}

}

