import { HttpClient } from '@angular/common/http';
import { Component, OnInit } from '@angular/core';
import { Globals } from '../../globals';
import { Submission } from '../../shared/submission.model';

export interface SubmissionResponse {
  items: PastSubmission[];
  total_count: number;
}

export interface PastSubmission {
  username: string;
  length: number;
  language: string;
  correct: boolean;
  submitted_date: string;
}

interface Alert {
  type: string;
  message: string;
}

// export const SUBMISSIONS: Submission[] = [
//   new Submission('python3', 25, )
// ];

const url = 'https://bytegolf.io';

@Component({
  selector: 'app-previous',
  templateUrl: './previous.component.html',
  styleUrls: ['./previous.component.css']
})
export class PreviousComponent implements OnInit {
  public panelOpenState = false;

  public submissions: Submission[] = [];

  constructor(
    private httpClient: HttpClient
  ) {}

  ngOnInit() {
    this.fetchData();
  }

  fetchData(): void {
    this.httpClient.get<Submission[]>(Globals.url, Globals.httpOptions)
      .subscribe(data => {
        this.submissions = data;
        console.log('got data: length = ' + this.submissions.length);
      });
  }


}
