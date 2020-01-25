import { HttpClient, HttpErrorResponse } from '@angular/common/http';
import { Component, OnInit } from '@angular/core';
import { ToastrService } from 'ngx-toastr';
import { Question } from '../models/question';

@Component({
  selector: 'app-holes',
  templateUrl: './holes.component.html',
  styleUrls: ['./holes.component.scss']
})
export class HolesComponent implements OnInit {
  public holes: Question[] = null;
  public loadingHoles = true;

  constructor(
    private http: HttpClient,
    private toastr: ToastrService
  ) { }

  public getHoles(): void {
    this.loadingHoles = true;
    this.http.get('http://localhost:8080/holes')
      .subscribe(
        (qs: Question[]) => {
          this.holes = qs;
          this.loadingHoles = false;
        },
        (error: HttpErrorResponse) => {
          this.holes = [];
          this.loadingHoles = false;
          this.toastr.error('Could not get holes, try reloading page...', 'Error', {tapToDismiss: true});
          console.log('Error getting holes' + error);
        }
      );
  }

  ngOnInit() {
    this.getHoles();
  }

}
