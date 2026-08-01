export namespace steam {
	
	export class AppInfo {
	    appID: string;
	    name: string;
	    size: number;
	
	    static createFrom(source: any = {}) {
	        return new AppInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.appID = source["appID"];
	        this.name = source["name"];
	        this.size = source["size"];
	    }
	}
	export class CleanedItem {
	    app_id: string;
	    name: string;
	    size: number;
	
	    static createFrom(source: any = {}) {
	        return new CleanedItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.app_id = source["app_id"];
	        this.name = source["name"];
	        this.size = source["size"];
	    }
	}
	export class CleanRecord {
	    timestamp: string;
	    items_count: number;
	    bytes_freed: number;
	    items: CleanedItem[];
	
	    static createFrom(source: any = {}) {
	        return new CleanRecord(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.timestamp = source["timestamp"];
	        this.items_count = source["items_count"];
	        this.bytes_freed = source["bytes_freed"];
	        this.items = this.convertValues(source["items"], CleanedItem);
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

}

