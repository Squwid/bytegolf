import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { RouterModule, Routes } from '@angular/router';
import { ClarityModule } from '@clr/angular';
import { AceEditorModule } from 'ng2-ace-editor';
import { ToastrModule } from 'ngx-toastr';
import { AppComponent } from './app.component';
import { HolesComponent } from './holes/holes.component';
import { HomeComponent } from './home/home.component';
import { PlayviewComponent } from './playview/playview.component';

const routes: Routes = [
  {
    path: '',
    component: HomeComponent
  },
  {
    path: 'play/:id',
    component: PlayviewComponent,
  },
  {
    path: 'holes',
    component: HolesComponent
  }
];

@NgModule({
  declarations: [
    AppComponent,
    HolesComponent,
    PlayviewComponent,
    HomeComponent
  ],
  imports: [
    BrowserModule,
    ClarityModule,
    BrowserAnimationsModule,
    RouterModule.forRoot(
      routes,
    ),
    AceEditorModule,
    ToastrModule.forRoot()
  ],
  providers: [],
  bootstrap: [AppComponent]
})
export class AppModule { }
