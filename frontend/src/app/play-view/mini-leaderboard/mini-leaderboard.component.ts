import { animate, state, style, transition, trigger } from '@angular/animations';
import { Component, OnInit } from '@angular/core';
import { LeaderboardEntry } from '../../shared/leaderboardentry.model';


const LBDATA: LeaderboardEntry[] = [
  {
    position: 1,
    username: 'Squwid',
    language_url: 'akdjfasdf',
    language: 'golang',
    score: 20,
    submitted_date: '2019',
  },
  {
    position: 2,
    username: 'Squwid',
    language_url: 'akdjfasdf',
    language: 'golang',
    score: 20,
    submitted_date: '2019',
  },
  {
    position: 3,
    username: 'Squwid',
    language_url: 'akdjfasdf',
    language: 'golang',
    score: 20,
    submitted_date: '2019',
  }
];

@Component({
  selector: 'app-mini-leaderboard',
  templateUrl: './mini-leaderboard.component.html',
  styleUrls: ['./mini-leaderboard.component.css'],
  animations: [
    trigger('detailExpand', [
      state('collapsed', style({height: '0px', minHeight: '0'})),
      state('expanded', style({height: '*'})),
      transition('expanded <=> collapsed', animate('225ms cubic-bezier(0.4, 0.0, 0.2, 1)')),
    ]),
  ],
})
export class MiniLeaderboardComponent implements OnInit {
  public data = LBDATA;
  public columnsToDisplay = ['position', 'username', 'score', 'language'];
  public expandedElement: LeaderboardEntry | null;

  constructor() { }

  ngOnInit() {
  }

}
