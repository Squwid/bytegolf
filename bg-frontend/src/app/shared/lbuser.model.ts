export class LbUser {
    public username: string;
    public score: number;
    public language: string;
    public codeImageUri: string;


    constructor(username: string, score: number, language: string) {
        this.username = username;
        this.score = score;
        this.language = language;
        this.codeImageUri = 'https://blog.golang.org/go-brand/Go-Logo/SVG/Go-Logo_LightBlue.svg';
    }
}
