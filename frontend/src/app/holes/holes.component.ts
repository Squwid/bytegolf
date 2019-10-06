import { Component, OnInit } from '@angular/core';
import { QUESTIONS } from './holes';

@Component({
  selector: 'app-holes',
  templateUrl: './holes.component.html',
  styleUrls: ['./holes.component.css']
})
export class HolesComponent implements OnInit {
  holes = QUESTIONS;

  constructor() { }

  ngOnInit() {
  }
}
