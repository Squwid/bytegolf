import { Component, OnInit } from '@angular/core';
import { HOLES } from './holes';

@Component({
  selector: 'app-holes',
  templateUrl: './holes.component.html',
  styleUrls: ['./holes.component.css']
})
export class HolesComponent implements OnInit {
  holes = HOLES;

  constructor() { }

  ngOnInit() {
  }
}
