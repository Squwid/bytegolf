import { HttpClient, HttpErrorResponse } from '@angular/common/http';
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

interface LoggedIn {
  logged_in: boolean;
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

  public question: Question = {} as Question;
  public questionLoading = false;

  // Stuff relating to past submissions
  public pastSubs: PastSubmission[] = [];
  public loadingPastSubs = true;

  // Leaderboard stuff
  public leaders: LeaderboardSpot[] = null;
  public loadingLeaders = true;
  public selfLeader: LeaderboardSpot = null;

  // the hole id for the entire page
  private holeId: string;

  public closedWarning = false;

  // hole not found
  public holeFound = true;

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

  public getSelfLeader(): void {
    this.selfLeader = {
      place: 22,
      displayName: 'Squwid',
      githubUrl: 'https://github.com/Squwid',
      score: 122,
      language: 'Go'
    };
  }

  // get the leaders per hole from the backend
  public getLeaders(): void {
    this.loadingLeaders = true;
    this.getSelfLeader();
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
    this.http.get('http://localhost:8080/hole?hole=' + this.holeId)
      .subscribe(
        (q: Question) => {
          this.holeFound = true;
          this.question = q;
          this.questionLoading = false;
          console.log('Found hole: ' + this.holeId);
          return;
        },
        (error: HttpErrorResponse) => {
          if (error.status === 404) {
            console.log('Did not find ' + this.holeId);
            this.holeFound = false;
            return;
          }
          console.log('Other error: ' + error.status);
        }
      );
    /*
    this.question = {
      id: '2020',
      name: 'Question Name Here',
      question: 'This is a question here but its going to be here when its done',
      live: true,
      difficulty: 'Hard'
    };
    */
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
    this.loadingPastSubs = true;
    this.http.get('http://localhost:8080/api/submissions?hole=' + this.holeId)
      .subscribe(
        (ps: PastSubmission[]) => {
          this.pastSubs = ps;
          this.loadingPastSubs = false;
          return;
        },
        (error: HttpErrorResponse) => {
          console.log('Error getting past submissions for ' + this.holeId);
          this.toastr.error(error.message, 'Error loading past submissions', {tapToDismiss: true});
          this.loadingPastSubs = false;
          return;
        }
      );
  }


}
