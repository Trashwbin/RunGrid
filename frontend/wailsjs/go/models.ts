export namespace domain {
	
	export class Group {
	    id: string;
	    name: string;
	    order: number;
	    color: string;
	    category: string;
	    icon: string;
	
	    static createFrom(source: any = {}) {
	        return new Group(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.order = source["order"];
	        this.color = source["color"];
	        this.category = source["category"];
	        this.icon = source["icon"];
	    }
	}
	export class GroupInput {
	    name: string;
	    order: number;
	    color: string;
	    category: string;
	    icon: string;
	
	    static createFrom(source: any = {}) {
	        return new GroupInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.order = source["order"];
	        this.color = source["color"];
	        this.category = source["category"];
	        this.icon = source["icon"];
	    }
	}
	export class HotkeyIssue {
	    id: string;
	    keys: string;
	    reason: string;
	
	    static createFrom(source: any = {}) {
	        return new HotkeyIssue(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.keys = source["keys"];
	        this.reason = source["reason"];
	    }
	}
	export class HotkeyApplyResult {
	    issues: HotkeyIssue[];
	
	    static createFrom(source: any = {}) {
	        return new HotkeyApplyResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.issues = this.convertValues(source["issues"], HotkeyIssue);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class HotkeyBinding {
	    id: string;
	    keys: string;
	
	    static createFrom(source: any = {}) {
	        return new HotkeyBinding(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.keys = source["keys"];
	    }
	}
	
	export class Item {
	    id: string;
	    name: string;
	    path: string;
	    target_name: string;
	    type: string;
	    icon_path: string;
	    group_id: string;
	    tags: string[];
	    favorite: boolean;
	    launch_count: number;
	    // Go type: time
	    last_used_at?: any;
	    hidden: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Item(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.path = source["path"];
	        this.target_name = source["target_name"];
	        this.type = source["type"];
	        this.icon_path = source["icon_path"];
	        this.group_id = source["group_id"];
	        this.tags = source["tags"];
	        this.favorite = source["favorite"];
	        this.launch_count = source["launch_count"];
	        this.last_used_at = this.convertValues(source["last_used_at"], null);
	        this.hidden = source["hidden"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ItemInput {
	    name: string;
	    path: string;
	    target_name: string;
	    type: string;
	    icon_path: string;
	    group_id: string;
	    tags: string[];
	    favorite: boolean;
	    hidden: boolean;
	
	    static createFrom(source: any = {}) {
	        return new ItemInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.path = source["path"];
	        this.target_name = source["target_name"];
	        this.type = source["type"];
	        this.icon_path = source["icon_path"];
	        this.group_id = source["group_id"];
	        this.tags = source["tags"];
	        this.favorite = source["favorite"];
	        this.hidden = source["hidden"];
	    }
	}
	export class ItemUpdate {
	    id: string;
	    name: string;
	    path: string;
	    target_name: string;
	    type: string;
	    icon_path: string;
	    group_id: string;
	    tags: string[];
	    favorite: boolean;
	    hidden: boolean;
	
	    static createFrom(source: any = {}) {
	        return new ItemUpdate(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.path = source["path"];
	        this.target_name = source["target_name"];
	        this.type = source["type"];
	        this.icon_path = source["icon_path"];
	        this.group_id = source["group_id"];
	        this.tags = source["tags"];
	        this.favorite = source["favorite"];
	        this.hidden = source["hidden"];
	    }
	}
	export class Point {
	    x: number;
	    y: number;
	
	    static createFrom(source: any = {}) {
	        return new Point(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.x = source["x"];
	        this.y = source["y"];
	    }
	}
	export class ScanResult {
	    total: number;
	    inserted: number;
	    skipped: number;
	
	    static createFrom(source: any = {}) {
	        return new ScanResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.total = source["total"];
	        this.inserted = source["inserted"];
	        this.skipped = source["skipped"];
	    }
	}

}

