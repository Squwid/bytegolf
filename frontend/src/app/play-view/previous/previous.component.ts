import { Component, OnInit } from '@angular/core';
import { Submission } from '../../shared/submission.model';

interface Alert {
  type: string;
  message: string;
}

@Component({
  selector: 'app-previous',
  templateUrl: './previous.component.html',
  styleUrls: ['./previous.component.css']
})
export class PreviousComponent implements OnInit {
  gameAlerts: Alert[] = [];
  previousSuccess = false;
  previousSubmission: Submission;
  showAlerts() {
    return (this.gameAlerts !== undefined && this.gameAlerts.length !== 0);
  }

  constructor() { }


  ngOnInit() {
    const sub = this.getBestScore();
    if (sub !== undefined) {
      this.gameAlerts.push({
        type: 'warning',
        message: 'You have not submitted a successful solution yet'
      });
      console.log('pushed');
      this.previousSuccess = true;
      this.previousSubmission = sub;
    }
  }

  closeGameAlert(alert: Alert) {
    this.gameAlerts.splice(this.gameAlerts.indexOf(alert), 1);
  }



  // getBestScore gets the players best score and sets it to the
  getBestScore() {
    // return undefined;
    return new Submission('golang', 10, '10/19/2019');
  }

}
