import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Component, OnInit } from '@angular/core';
import { ToastrService } from 'ngx-toastr';
import { Submission } from '../../shared/submission.model';

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
  gameAlerts: Alert[] = [];
  previousSuccess = false;
  previousSubmission: Submission;
  showAlerts() {
    return (this.gameAlerts !== undefined && this.gameAlerts.length !== 0);
  }

  constructor(private http: HttpClient, private toastr: ToastrService) { }


  retreivePastExecutes() {
    console.log('past executes');
    this.http.get(url + '/compile', {
      headers: new HttpHeaders({
        'Content-Type': 'application/json',
      })
    }).subscribe(responseData => {
      console.log('compiles: ' + responseData);
    });
  }

  ngOnInit() {
    this.retreivePastExecutes();
    const sub = this.getBestScore();
    if (sub !== undefined) {
      this.gameAlerts.push({
        type: 'warning',
        message: 'You have not submitted a successful solution yet'
      });
      // console.log('pushed');
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
