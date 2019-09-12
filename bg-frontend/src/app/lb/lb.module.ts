import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { DataViewModule } from 'primeng/dataview';
import { LbComponent } from './lb.component';



@NgModule({
  declarations: [
    LbComponent
  ],
  imports: [
    CommonModule,
    DataViewModule
  ],
  exports: [
    LbComponent,
    DataViewModule
  ],
  providers: [],
  bootstrap: [LbComponent]
})
export class LbModule { }
