export class Submission {
    public language: string;
    public score: number;
    public user: string;
    public date: string;

    constructor(lang: string, score: number, date: string) {
        this.language = lang;
        this.score = score;
        this.date = date;
        this.user = '';
    }

}
