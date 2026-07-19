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

}

