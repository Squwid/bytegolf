import { HttpClient, HttpErrorResponse } from '@angular/common/http';
import { Component } from '@angular/core';

export interface User {
  logged_in: boolean;
  username: string;
}

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent {
  title = 'frontend';

  public user = {
    logged_in: true,
    username: 'Squwid'
  };

  constructor(
    private http: HttpClient
  ) {
  }

  public checkLogin(): void {
    this.http.get('http://localhost:8080/check').subscribe(
      (u: User) => {
        this.user = u;
      },
      (error: HttpErrorResponse) => {
        console.log('Got error: ' + error);
      }
    );
  }
}
