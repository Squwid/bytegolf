export class Execute {
    public script: string;
    public language: string;
    public versionIndex: string;
    public holeId: string;

    constructor(script: string, lang: string, holeId: string) {
        this.script = script;
        this.language = lang;
        this.holeId = holeId;

        // todo: create a function to track the version indexes
        this.versionIndex = '2';
    }
}
