export namespace network {
	
	export class ARPEntry {
	    ipAddress: string;
	    macAddress: string;
	    interface: string;
	    type: string;
	
	    static createFrom(source: any = {}) {
	        return new ARPEntry(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ipAddress = source["ipAddress"];
	        this.macAddress = source["macAddress"];
	        this.interface = source["interface"];
	        this.type = source["type"];
	    }
	}
	export class CaptureTemplate {
	    name: string;
	    description: string;
	    bpfFilter: string;
	    duration: number;
	
	    static createFrom(source: any = {}) {
	        return new CaptureTemplate(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.description = source["description"];
	        this.bpfFilter = source["bpfFilter"];
	        this.duration = source["duration"];
	    }
	}
	export class NetworkConnection {
	    fd: number;
	    family: number;
	    type: number;
	    localIP: string;
	    localPort: number;
	    remoteIP: string;
	    remotePort: number;
	    status: string;
	    pid: number;
	    processName: string;
	    protocol: string;
	
	    static createFrom(source: any = {}) {
	        return new NetworkConnection(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.fd = source["fd"];
	        this.family = source["family"];
	        this.type = source["type"];
	        this.localIP = source["localIP"];
	        this.localPort = source["localPort"];
	        this.remoteIP = source["remoteIP"];
	        this.remotePort = source["remotePort"];
	        this.status = source["status"];
	        this.pid = source["pid"];
	        this.processName = source["processName"];
	        this.protocol = source["protocol"];
	    }
	}
	export class NetworkInterface {
	    name: string;
	    displayName: string;
	    description: string;
	    hardwareAddr: string;
	    mtu: number;
	    flags: string[];
	    addrs: string[];
	    isUp: boolean;
	    isLoopback: boolean;
	    isBroadcast: boolean;
	    isPointToPoint: boolean;
	    isMulticast: boolean;
	
	    static createFrom(source: any = {}) {
	        return new NetworkInterface(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.displayName = source["displayName"];
	        this.description = source["description"];
	        this.hardwareAddr = source["hardwareAddr"];
	        this.mtu = source["mtu"];
	        this.flags = source["flags"];
	        this.addrs = source["addrs"];
	        this.isUp = source["isUp"];
	        this.isLoopback = source["isLoopback"];
	        this.isBroadcast = source["isBroadcast"];
	        this.isPointToPoint = source["isPointToPoint"];
	        this.isMulticast = source["isMulticast"];
	    }
	}
	export class PublicIPInfo {
	    ipAddress: string;
	    interfaceName: string;
	
	    static createFrom(source: any = {}) {
	        return new PublicIPInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ipAddress = source["ipAddress"];
	        this.interfaceName = source["interfaceName"];
	    }
	}

}

export namespace system {
	
	export class SystemInfo {
	    username: string;
	    os: string;
	    arch: string;
	    hostname: string;
	    workingDir: string;
	
	    static createFrom(source: any = {}) {
	        return new SystemInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.username = source["username"];
	        this.os = source["os"];
	        this.arch = source["arch"];
	        this.hostname = source["hostname"];
	        this.workingDir = source["workingDir"];
	    }
	}

}

export namespace tools {
	
	export class PingResult {
	    host: string;
	    ip: string;
	    packets: number;
	    loss: string;
	    minRtt: string;
	    avgRtt: string;
	    maxRtt: string;
	    error?: string;
	
	    static createFrom(source: any = {}) {
	        return new PingResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.host = source["host"];
	        this.ip = source["ip"];
	        this.packets = source["packets"];
	        this.loss = source["loss"];
	        this.minRtt = source["minRtt"];
	        this.avgRtt = source["avgRtt"];
	        this.maxRtt = source["maxRtt"];
	        this.error = source["error"];
	    }
	}
	export class TracerouteHop {
	    n: number;
	    host: string;
	    address: string;
	    rtt: string;
	
	    static createFrom(source: any = {}) {
	        return new TracerouteHop(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.n = source["n"];
	        this.host = source["host"];
	        this.address = source["address"];
	        this.rtt = source["rtt"];
	    }
	}

}

