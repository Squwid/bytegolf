export class Question {
    public difficulty: string;
    public title: string;
    public uuid: string;
    public createdDate: string;
    public question: string;

    constructor(diff: string, title: string, question: string) {
        this.difficulty = diff;
        this.title = title;
        this.question = question;
        this.uuid = '12345';
        this.createdDate = '10/24/2018';
    }
}
