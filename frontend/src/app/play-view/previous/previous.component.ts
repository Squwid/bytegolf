import { Component } from '@angular/core';

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
export class PreviousComponent {
  public panelOpenState = false;
}
