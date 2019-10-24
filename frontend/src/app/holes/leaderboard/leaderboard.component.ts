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
  selector: 'app-leaderboard',
  templateUrl: './leaderboard.component.html',
  styleUrls: ['./leaderboard.component.css']
})
export class LeaderboardComponent implements OnInit {
  public data = LBDATA;

  constructor() { }

  ngOnInit() {
  }

}
