export interface Submission {
    submission: Request;
    response: Response;
    bgid: string;
    correct: boolean;
    holeId: string;
    submitted_time: string;
    length: number;
}

export interface Response {
    output: string;
    statusCode: string;
    memory: string;
    cpuTime: string;
}

export interface Request {
    script: string;
    language: string;
    versionIndex: string;
    holeId: string;
}
