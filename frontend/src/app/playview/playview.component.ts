import { HttpClient } from '@angular/common/http';
import { Component, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
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

// LeaderboardSpot is a spot on the leaderboard
export interface LeaderboardSpot {
  place: number;
  displayName: string;
  githubUrl: string;
  score: number;
  language: string;
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

  // Leaderboard stuff
  public leaders: LeaderboardSpot[] = null;
  public loadingLeaders = true;

  // the hole id for the entire page
  private holeId: string;

  constructor(
    private toastr: ToastrService,
    private http: HttpClient,
    private route: ActivatedRoute
    ) {
    if (this.languages.length !== 0) {
      this.activeLanguage = this.languages[0];
    } else {
      this.activeLanguage = {} as Language;
    }
  }

  ngOnInit() {
    this.holeId = this.route.snapshot.params.id;
    console.log('this id: ' + this.holeId);
    this.getQuestion();
    this.getPastSubmissions();
    this.getLeaders();

  }

  // get the leaders per hole from the backend
  public getLeaders(): void {
    this.loadingLeaders = true;
    this.leaders = [
      {
        place: 1,
        displayName: 'Squwid',
        githubUrl: 'https://github.com/Squwid',
        score: 25,
        language: 'Go'
      },
      {
        place: 2,
        displayName: 'iCollin',
        githubUrl: 'https://github.com/icollin',
        score: 30,
        language: 'Python 3'
      },
      {
        place: 3,
        displayName: 'Kraftcur',
        githubUrl: 'https://github.com/kraftcur',
        score: 35,
        language: 'Python 2'
      }
    ];
    this.loadingLeaders = false;
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
    // this.toastr.info('', 'Hello!', {tapToDismiss: true});
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
