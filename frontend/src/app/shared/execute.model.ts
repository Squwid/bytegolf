export class Execute {
    public script: string;
    public language: string;
    public versionIndex: string;

    constructor(script: string, lang: string) {
        this.script = script;
        this.language = lang;

        // todo: create a function to track the version indexes
        this.versionIndex = '2';
    }
}
