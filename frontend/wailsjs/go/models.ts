export namespace network {
	
	export class ARPEntry {
	    IPAddress: string;
	    MACAddress: string;
	    Interface: string;
	    Type: string;
	
	    static createFrom(source: any = {}) {
	        return new ARPEntry(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.IPAddress = source["IPAddress"];
	        this.MACAddress = source["MACAddress"];
	        this.Interface = source["Interface"];
	        this.Type = source["Type"];
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
	    FD: number;
	    Family: number;
	    Type: number;
	    LocalIP: string;
	    LocalPort: number;
	    RemoteIP: string;
	    RemotePort: number;
	    Status: string;
	    PID: number;
	    ProcessName: string;
	    Protocol: string;
	
	    static createFrom(source: any = {}) {
	        return new NetworkConnection(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.FD = source["FD"];
	        this.Family = source["Family"];
	        this.Type = source["Type"];
	        this.LocalIP = source["LocalIP"];
	        this.LocalPort = source["LocalPort"];
	        this.RemoteIP = source["RemoteIP"];
	        this.RemotePort = source["RemotePort"];
	        this.Status = source["Status"];
	        this.PID = source["PID"];
	        this.ProcessName = source["ProcessName"];
	        this.Protocol = source["Protocol"];
	    }
	}
	export class NetworkInterface {
	    Name: string;
	    DisplayName: string;
	    Description: string;
	    HardwareAddr: number[];
	    MTU: number;
	    Flags: string[];
	    Addrs: string[];
	    IsUp: boolean;
	    IsLoopback: boolean;
	    IsBroadcast: boolean;
	    IsPointToPoint: boolean;
	    IsMulticast: boolean;
	
	    static createFrom(source: any = {}) {
	        return new NetworkInterface(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Name = source["Name"];
	        this.DisplayName = source["DisplayName"];
	        this.Description = source["Description"];
	        this.HardwareAddr = source["HardwareAddr"];
	        this.MTU = source["MTU"];
	        this.Flags = source["Flags"];
	        this.Addrs = source["Addrs"];
	        this.IsUp = source["IsUp"];
	        this.IsLoopback = source["IsLoopback"];
	        this.IsBroadcast = source["IsBroadcast"];
	        this.IsPointToPoint = source["IsPointToPoint"];
	        this.IsMulticast = source["IsMulticast"];
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

