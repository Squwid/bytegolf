import { Component, OnInit } from '@angular/core';
import 'brace';
import 'brace/mode/golang';
import 'brace/mode/java';
import 'brace/mode/php';
import 'brace/mode/python';
import 'brace/mode/ruby';
import 'brace/mode/rust';
import 'brace/mode/sql';
import 'brace/mode/text';
import 'brace/theme/dracula';
import { ToastrService } from 'ngx-toastr';
import { LANGUAGES } from '../consts/consts';
import { Language } from '../models/language';
import { Question } from '../models/question';

export interface Submission {
  id: string;
  correct: boolean;
  length: number;
}

@Component({
  selector: 'app-playview',
  templateUrl: './playview.component.html',
  styleUrls: ['./playview.component.scss']
})
export class PlayviewComponent implements OnInit {
  private defaultContent = 'print(\'Hello, World!)';

  public submitDisabled = true;

  public languages = LANGUAGES;
  public activeLanguage: Language = null;

  public braceActiveTheme = 'dracula';
  public braceContent = this.defaultContent;

  public question: Question = null;
  public questionLoading = false;

  // Stuff relating to past submissions
  public pastSubs: PastSubmission[] = null;
  public loadingPastSubs = true;

  constructor(private toastr: ToastrService) {
    if (this.languages.length !== 0) {
      this.activeLanguage = this.languages[0];
    } else {
      this.activeLanguage = {} as Language;
    }
  }

  ngOnInit() {
    this.getQuestion();
    this.getPastSubmissions();
  }

  // gets and sets the question using the id
  public getQuestion(): void {
    this.questionLoading = true;
    this.question = {
      id: '2020',
      name: 'Question Name Here',
      question: 'This is a question here but its going to be here when its done',
      live: true,
      difficulty: 'Hard'
    };
    this.questionLoading = false;
  }

  public onDelete() {
    this.braceContent = this.defaultContent;
    this.toastr.info('', 'Hello!', {tapToDismiss: true});
  }

  public setLanguage(lang: Language) {
    this.activeLanguage = lang;
  }

  public onSubmit() {
    console.log('Submission');
  }

  // getPastSubmissions gets the past submissions
  public getPastSubmissions(): void {
    // this.pastSubs = [{
    //   id: '100',
    //   correct: true,
    //   score: 20,
    //   language: 'Go',
    //   script: 'print(\"Hello, World!\")',
    //   date: '10/20/2020 10:35am'
    // },
    // {
    //   id: '100',
    //   correct: true,
    //   score: 20,
    //   language: 'Go',
    //   script: 'print(\"Hello, World!\")',
    //   date: '10/20/2020 10:35am'
    // },
    // {
    //   id: '100',
    //   correct: false,
    //   score: 20,
    //   language: 'Go',
    //   script: 'print(\"Hello, World!\")',
    //   date: '10/20/2020 10:35am'
    // },
    // {
    //   id: '100',
    //   correct: false,
    //   score: 20,
    //   language: 'Go',
    //   script: 'print(\"Hello, World!\")',
    //   date: '10/20/2020 10:35am'
    // }];
    this.pastSubs = [];
  }


}
