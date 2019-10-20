import { HttpClient, HttpErrorResponse } from '@angular/common/http';
import { Component, OnInit } from '@angular/core';
import { Globals } from '../globals';

export interface User {
  logged_in: boolean;
  username: string;
}

@Component({
  selector: 'app-navbar',
  templateUrl: './navbar.component.html',
  styleUrls: ['./navbar.component.css']
})
export class NavbarComponent implements OnInit {
  public user: User = null;

  constructor(private http: HttpClient) { }

  ngOnInit() {
    this.http.get<User>(Globals.url + '/check', Globals.httpOptions)
      .subscribe(
        user => {
          this.user = user;
          console.log('got user. Logged in: ' + this.user.logged_in);
        },
        (error: HttpErrorResponse) => {console.log('error: ' + error); }
      );
  }

}
