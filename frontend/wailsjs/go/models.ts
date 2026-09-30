export namespace domain {
	
	export class Config {
	    token: string;
	    username: string;
	    client: string;
	    rate: string;
	    token_expiration: string;
	    cnpj: string;
	    company_name: string;
	
	    static createFrom(source: any = {}) {
	        return new Config(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.token = source["token"];
	        this.username = source["username"];
	        this.client = source["client"];
	        this.rate = source["rate"];
	        this.token_expiration = source["token_expiration"];
	        this.cnpj = source["cnpj"];
	        this.company_name = source["company_name"];
	    }
	}
	export class GithubActivity {
	    type: string;
	    repo: string;
	    title: string;
	    description: string;
	    url: string;
	    date: string;
	    ref: string;
	
	    static createFrom(source: any = {}) {
	        return new GithubActivity(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.type = source["type"];
	        this.repo = source["repo"];
	        this.title = source["title"];
	        this.description = source["description"];
	        this.url = source["url"];
	        this.date = source["date"];
	        this.ref = source["ref"];
	    }
	}
	export class Repository {
	    id: number;
	    name: string;
	    full_name: string;
	    private: boolean;
	    description: string;
	    html_url: string;
	
	    static createFrom(source: any = {}) {
	        return new Repository(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.full_name = source["full_name"];
	        this.private = source["private"];
	        this.description = source["description"];
	        this.html_url = source["html_url"];
	    }
	}

}

