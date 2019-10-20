export class Language {
    public language: string;
    public icon: string;
    public tooltip: string;
    public mode: string;
    public versionIndex: string;

    constructor(language: string, icon: string, tip: string, mode: string) {
        this.language = language;
        this.icon = icon;
        this.tooltip = tip;
        this.mode = mode;

        // this exists because the version indexes are what the api has
        if (language === 'python') {
            this.versionIndex = '2';
        }
    }
}
